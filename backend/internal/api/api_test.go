package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/api"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
)

const testSecret = "super-secret-key-that-is-at-least-32-characters-long!"

type fakeNodeRecord struct {
	detail    db.NodeDetail
	ownerID   *int64
	claimCode string
	createdAt time.Time
}

// SEMENTARA: implementasi DataStore in-memory untuk pengujian unit API.
type fakeAPIStore struct {
	mu      sync.Mutex
	nextUID int64
	users   map[int64]*db.User
	nodes   map[string]*fakeNodeRecord
}

func newFakeAPIStore() *fakeAPIStore {
	return &fakeAPIStore{
		nextUID: 1,
		users:   make(map[int64]*db.User),
		nodes:   make(map[string]*fakeNodeRecord),
	}
}

func (f *fakeAPIStore) EnsureAdminUser(_ context.Context, username, passwordHash string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.users {
		if u.Username == username {
			u.IsAdmin = true
			return nil
		}
	}
	id := f.nextUID
	f.nextUID++
	f.users[id] = &db.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		IsAdmin:      true,
	}
	return nil
}

func (f *fakeAPIStore) CreateUser(_ context.Context, username, passwordHash string) (*db.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.users {
		if u.Username == username {
			return nil, db.ErrUsernameTaken
		}
	}
	id := f.nextUID
	f.nextUID++
	u := &db.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		IsAdmin:      false,
	}
	f.users[id] = u
	copyUser := *u
	return &copyUser, nil
}

func (f *fakeAPIStore) GetUserByUsername(_ context.Context, username string) (*db.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.users {
		if u.Username == username {
			copyUser := *u
			return &copyUser, nil
		}
	}
	return nil, db.ErrUserNotFound
}

func (f *fakeAPIStore) GetUserByID(_ context.Context, userID int64) (*db.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	u, ok := f.users[userID]
	if !ok {
		return nil, db.ErrUserNotFound
	}
	copyUser := *u
	return &copyUser, nil
}

func (f *fakeAPIStore) UpdateUserTelegramChatID(_ context.Context, userID int64, chatID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	u, ok := f.users[userID]
	if !ok {
		return db.ErrUserNotFound
	}
	if chatID == "" {
		u.TelegramChatID = nil
	} else {
		c := chatID
		u.TelegramChatID = &c
	}
	return nil
}

func (f *fakeAPIStore) ListNodeSummaries(_ context.Context, ownerID int64) ([]db.NodeSummary, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	res := make([]db.NodeSummary, 0)
	for _, n := range f.nodes {
		if n.ownerID != nil && *n.ownerID == ownerID {
			res = append(res, db.NodeSummary{
				NodeID:   n.detail.NodeID,
				Name:     n.detail.Name,
				Online:   n.detail.Online,
				Mode:     n.detail.Mode,
				LastSeen: n.detail.LastSeen,
				Temp:     n.detail.Temp,
				NTU:      n.detail.NTU,
				Status:   n.detail.Status,
			})
		}
	}
	return res, nil
}

func (f *fakeAPIStore) GetNodeDetail(_ context.Context, ownerID int64, nodeID string) (*db.NodeDetail, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok || n.ownerID == nil || *n.ownerID != ownerID {
		return nil, db.ErrNodeNotFound
	}
	d := n.detail
	return &d, nil
}

func (f *fakeAPIStore) UpdateNodeName(_ context.Context, ownerID int64, nodeID, name string) (*db.NodeDetail, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok || n.ownerID == nil || *n.ownerID != ownerID {
		return nil, db.ErrNodeNotFound
	}
	n.detail.Name = name
	d := n.detail
	return &d, nil
}

func (f *fakeAPIStore) UpdateNodeMode(_ context.Context, ownerID int64, nodeID string, mode core.Mode) (*db.NodeDetail, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok || n.ownerID == nil || *n.ownerID != ownerID {
		return nil, db.ErrNodeNotFound
	}
	n.detail.Mode = mode
	d := n.detail
	return &d, nil
}

func (f *fakeAPIStore) ClaimNode(_ context.Context, userID int64, nodeID, claimCode string) (*db.NodeDetail, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok || !strings.EqualFold(n.claimCode, claimCode) {
		return nil, db.ErrInvalidClaim
	}
	if n.ownerID != nil && *n.ownerID != userID {
		return nil, db.ErrNodeAlreadyClaimed
	}
	uid := userID
	n.ownerID = &uid
	d := n.detail
	return &d, nil
}

func (f *fakeAPIStore) ReleaseNodeOwnership(_ context.Context, ownerID int64, nodeID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok || n.ownerID == nil || *n.ownerID != ownerID {
		return db.ErrNodeNotFound
	}
	n.ownerID = nil
	n.detail.Mode = core.ModeManual
	return nil
}

func (f *fakeAPIStore) toAdminNodeLocked(n *fakeNodeRecord) db.AdminNode {
	var ownerUsername *string
	if n.ownerID != nil {
		if u, ok := f.users[*n.ownerID]; ok {
			uname := u.Username
			ownerUsername = &uname
		}
	}
	return db.AdminNode{
		NodeID:        n.detail.NodeID,
		Name:          n.detail.Name,
		Online:        n.detail.Online,
		LastSeen:      n.detail.LastSeen,
		CreatedAt:     n.createdAt,
		ClaimCode:     n.claimCode,
		OwnerUsername: ownerUsername,
	}
}

func (f *fakeAPIStore) ListAdminNodes(_ context.Context) ([]db.AdminNode, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	res := make([]db.AdminNode, 0, len(f.nodes))
	for _, n := range f.nodes {
		res = append(res, f.toAdminNodeLocked(n))
	}
	return res, nil
}

func (f *fakeAPIStore) CreateAdminNode(_ context.Context, nodeID string) (*db.AdminNode, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, exists := f.nodes[nodeID]; exists {
		return nil, db.ErrNodeAlreadyExists
	}
	rec := &fakeNodeRecord{
		detail: db.NodeDetail{
			NodeID: nodeID,
			Name:   "Akuarium " + nodeID,
			Mode:   core.ModeManual,
		},
		claimCode: "ABCD1234",
		createdAt: time.Now().UTC(),
	}
	f.nodes[nodeID] = rec
	adminNode := f.toAdminNodeLocked(rec)
	return &adminNode, nil
}

func (f *fakeAPIStore) RegenerateClaimCode(_ context.Context, nodeID string) (*db.AdminNode, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok {
		return nil, db.ErrNodeNotFound
	}
	n.claimCode = fmt.Sprintf("NEW%05d", len(n.claimCode))
	adminNode := f.toAdminNodeLocked(n)
	return &adminNode, nil
}

func (f *fakeAPIStore) AdminReleaseNodeOwner(_ context.Context, nodeID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok {
		return db.ErrNodeNotFound
	}
	n.ownerID = nil
	n.detail.Mode = core.ModeManual
	return nil
}

func (f *fakeAPIStore) AdminDeleteNode(_ context.Context, nodeID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.nodes[nodeID]; !ok {
		return db.ErrNodeNotFound
	}
	delete(f.nodes, nodeID)
	return nil
}

func (f *fakeAPIStore) GetTelemetryHistory(_ context.Context, _ string, _, _ time.Time) ([]db.TelemetryPoint, error) {
	return []db.TelemetryPoint{}, nil
}

func (f *fakeAPIStore) GetNodeLogs(_ context.Context, _ string) ([]db.ActuatorLog, []db.AlertLog, error) {
	return []db.ActuatorLog{}, []db.AlertLog{}, nil
}

func (f *fakeAPIStore) LogActuators(_ context.Context, _ string, _ core.Actuators, _ core.Source) error {
	return nil
}

func TestHealthEndpoint(t *testing.T) {
	server := api.NewServer(nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var res map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("gagal decode response: %v", err)
	}
	if res["status"] != "ok" {
		t.Errorf("got status %q, want 'ok'", res["status"])
	}
}

func TestProtectedEndpointsRequireAuth(t *testing.T) {
	authSvc, err := auth.NewService(testSecret)
	if err != nil {
		t.Fatalf("gagal inisialisasi auth service: %v", err)
	}

	server := api.NewServer(nil, authSvc, nil, nil)

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/me"},
		{http.MethodPut, "/api/me/telegram"},
		{http.MethodPost, "/api/nodes/claim"},
		{http.MethodGet, "/api/nodes"},
		{http.MethodGet, "/api/nodes/aq-1a2b3c"},
		{http.MethodPatch, "/api/nodes/aq-1a2b3c"},
		{http.MethodDelete, "/api/nodes/aq-1a2b3c"},
		{http.MethodPut, "/api/nodes/aq-1a2b3c/mode"},
		{http.MethodPut, "/api/nodes/aq-1a2b3c/actuators"},
		{http.MethodGet, "/api/nodes/aq-1a2b3c/telemetry"},
		{http.MethodGet, "/api/nodes/aq-1a2b3c/logs"},
		{http.MethodGet, "/api/admin/nodes"},
		{http.MethodPost, "/api/admin/nodes"},
		{http.MethodPost, "/api/admin/nodes/aq-1a2b3c/claim-code"},
		{http.MethodDelete, "/api/admin/nodes/aq-1a2b3c/owner"},
		{http.MethodDelete, "/api/admin/nodes/aq-1a2b3c"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			req := httptest.NewRequest(ep.method, ep.path, nil)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("%s %s: got status %d, want 401", ep.method, ep.path, rec.Code)
			}
		})
	}
}

func TestTempCommander(t *testing.T) {
	cmd := api.NewTempCommander()
	ctx := context.Background()

	if err := cmd.SetMode(ctx, "aq-001", core.ModeAuto); err != nil {
		t.Errorf("TempCommander.SetMode gagal: %v", err)
	}

	if err := cmd.SendActuators(ctx, "aq-001", core.Actuators{Heater: 1, Fan: 0, Pump: 1}); err != nil {
		t.Errorf("TempCommander.SendActuators gagal: %v", err)
	}
}

func TestLoginValidation(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	server := api.NewServer(nil, authSvc, nil, nil)

	// Test 1: Payload kosong/rusak -> 401
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{invalid json}`))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", rec.Code)
	}

	// Test 2: Username kosong -> 401
	body, _ := json.Marshal(map[string]string{"username": "", "password": "password"})
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", rec.Code)
	}
}

func TestRegisterAndMeFlow(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	store := newFakeAPIStore()
	server := api.NewServer(store, authSvc, api.NewTempCommander(), nil)

	// 1. Validasi gagal -> 400
	invalidCases := []struct {
		name     string
		username string
		password string
	}{
		{"username terlalu pendek", "ab", "password123"},
		{"username terlalu panjang", strings.Repeat("a", 33), "password123"},
		{"username karakter terlarang", "budi@uns", "password123"},
		{"password kurang dari 8 karakter", "budi_01", "1234567"},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]string{
				"username": tc.username,
				"password": tc.password,
			})
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("%s: got status %d, want 400", tc.name, rec.Code)
			}
		})
	}

	// 2. Register sukses (huruf kapital diubah ke huruf kecil) -> 201 {"token"}
	body, _ := json.Marshal(map[string]string{
		"username": "Budi_Santoso.01",
		"password": "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("register got status %d, want 201: %s", rec.Code, rec.Body.String())
	}

	var regRes map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&regRes)
	token := regRes["token"]
	if token == "" {
		t.Fatal("token kosong saat register")
	}

	claims, err := authSvc.Verify(token)
	if err != nil {
		t.Fatalf("verify token gagal: %v", err)
	}
	if claims.Username != "budi_santoso.01" || claims.Admin {
		t.Fatalf("claims tidak sesuai: %+v", claims)
	}

	// 3. Username sudah dipakai -> 409 {"error":"Username sudah dipakai"}
	req = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("register duplikat got status %d, want 409", rec.Code)
	}
	var errRes map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&errRes)
	if errRes["error"] != "Username sudah dipakai" {
		t.Errorf("got error %q, want %q", errRes["error"], "Username sudah dipakai")
	}

	// 4. GET /api/me sebelum & sesudah mengisi telegram
	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/me got status %d, want 200", rec.Code)
	}
	rawMe := rec.Body.String()
	if strings.Contains(rawMe, "telegram_chat_id") {
		t.Fatalf("GET /api/me tidak boleh mengembalikan telegram_chat_id: %s", rawMe)
	}
	var me map[string]any
	_ = json.Unmarshal([]byte(rawMe), &me)
	if me["username"] != "budi_santoso.01" || me["is_admin"] != false || me["telegram_linked"] != false {
		t.Fatalf("isi /api/me tidak sesuai: %+v", me)
	}

	// Isi telegram chat ID -> telegram_linked menjadi true
	tgBody, _ := json.Marshal(map[string]string{"telegram_chat_id": "99887766"})
	req = httptest.NewRequest(http.MethodPut, "/api/me/telegram", bytes.NewReader(tgBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("PUT /api/me/telegram got %d, want 204", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	_ = json.NewDecoder(rec.Body).Decode(&me)
	if me["telegram_linked"] != true {
		t.Fatalf("telegram_linked seharusnya true, got %+v", me)
	}
}

func TestClaimNodeAndOwnershipIsolation(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	store := newFakeAPIStore()
	server := api.NewServer(store, authSvc, api.NewTempCommander(), nil)

	userA, _ := store.CreateUser(context.Background(), "usera", "hash")
	userB, _ := store.CreateUser(context.Background(), "userb", "hash")
	_ = store.EnsureAdminUser(context.Background(), "admin", "hash")
	adminUser, _ := store.GetUserByUsername(context.Background(), "admin")

	tokenA, _ := authSvc.Sign(userA.ID, userA.Username, false)
	tokenB, _ := authSvc.Sign(userB.ID, userB.Username, false)
	tokenAdmin, _ := authSvc.Sign(adminUser.ID, adminUser.Username, true)

	// Siapkan node aq-1a2b3c dengan claim code ABCD1234
	_, _ = store.CreateAdminNode(context.Background(), "aq-1a2b3c")

	// 1. Admin memanggil POST /api/nodes/claim -> 403
	claimBody, _ := json.Marshal(map[string]string{"node_id": "aq-1a2b3c", "claim_code": "ABCD1234"})
	req := httptest.NewRequest(http.MethodPost, "/api/nodes/claim", bytes.NewReader(claimBody))
	req.Header.Set("Authorization", "Bearer "+tokenAdmin)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("admin claim got status %d, want 403", rec.Code)
	}

	// 2. Sebelum diklaim, User A mengakses GET /api/nodes/aq-1a2b3c -> 404 {"error":"Node tidak ditemukan"}
	req = httptest.NewRequest(http.MethodGet, "/api/nodes/aq-1a2b3c", nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("akses node belum diklaim got %d, want 404", rec.Code)
	}

	// 3. Kode klaim salah atau node tidak ada -> 404 {"error":"ID node atau kode klaim salah"}
	wrongBody, _ := json.Marshal(map[string]string{"node_id": "aq-1a2b3c", "claim_code": "SALAH999"})
	req = httptest.NewRequest(http.MethodPost, "/api/nodes/claim", bytes.NewReader(wrongBody))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("kode salah got %d, want 404", rec.Code)
	}
	var errRes map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&errRes)
	if errRes["error"] != "ID node atau kode klaim salah" {
		t.Errorf("got error %q, want %q", errRes["error"], "ID node atau kode klaim salah")
	}

	// 4. User A klaim dengan huruf kecil & spasi -> 200 NodeDetail
	okBody, _ := json.Marshal(map[string]string{"node_id": " aq-1a2b3c ", "claim_code": " abcd1234 "})
	req = httptest.NewRequest(http.MethodPost, "/api/nodes/claim", bytes.NewReader(okBody))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("User A klaim got %d, want 200: %s", rec.Code, rec.Body.String())
	}

	// 5. User A klaim ulang node miliknya sendiri -> 200
	req = httptest.NewRequest(http.MethodPost, "/api/nodes/claim", bytes.NewReader(okBody))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("User A klaim ulang milik sendiri got %d, want 200", rec.Code)
	}

	// 6. User B klaim node milik User A dengan kode benar -> 409 {"error":"Akuarium sudah terdaftar di akun lain"}
	req = httptest.NewRequest(http.MethodPost, "/api/nodes/claim", bytes.NewReader(okBody))
	req.Header.Set("Authorization", "Bearer "+tokenB)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("User B klaim milik orang lain got %d, want 409", rec.Code)
	}
	_ = json.NewDecoder(rec.Body).Decode(&errRes)
	if errRes["error"] != "Akuarium sudah terdaftar di akun lain" {
		t.Errorf("got error %q, want %q", errRes["error"], "Akuarium sudah terdaftar di akun lain")
	}

	// 7. Semua endpoint node lama menolak User B (bukan pemilik) dengan 404 {"error":"Node tidak ditemukan"}
	nonOwnerRequests := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/api/nodes/aq-1a2b3c", ""},
		{http.MethodPatch, "/api/nodes/aq-1a2b3c", `{"name":"Hack"}`},
		{http.MethodPut, "/api/nodes/aq-1a2b3c/mode", `{"mode":"auto"}`},
		{http.MethodPut, "/api/nodes/aq-1a2b3c/actuators", `{"heater":1,"fan":0,"pump":0}`},
		{http.MethodGet, "/api/nodes/aq-1a2b3c/telemetry", ""},
		{http.MethodGet, "/api/nodes/aq-1a2b3c/logs", ""},
		{http.MethodDelete, "/api/nodes/aq-1a2b3c", ""},
	}
	for _, nr := range nonOwnerRequests {
		req = httptest.NewRequest(nr.method, nr.path, bytes.NewBufferString(nr.body))
		req.Header.Set("Authorization", "Bearer "+tokenB)
		rec = httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("User B %s %s got %d, want 404", nr.method, nr.path, rec.Code)
		}
		var e map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&e)
		if e["error"] != "Node tidak ditemukan" {
			t.Errorf("User B %s %s error = %q, want 'Node tidak ditemukan'", nr.method, nr.path, e["error"])
		}
	}

	// 8. Admin GET /api/nodes harus mengembalikan array kosong []
	req = httptest.NewRequest(http.MethodGet, "/api/nodes", nil)
	req.Header.Set("Authorization", "Bearer "+tokenAdmin)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin GET /api/nodes got %d, want 200", rec.Code)
	}
	var adminNodeList []db.NodeSummary
	_ = json.NewDecoder(rec.Body).Decode(&adminNodeList)
	if len(adminNodeList) != 0 {
		t.Fatalf("admin GET /api/nodes harus kosong, got %+v", adminNodeList)
	}

	// 9. User A bisa melihat dan mengubah mode ke auto, lalu DELETE /api/nodes/aq-1a2b3c melepas kepemilikan dan mereset ke manual
	req = httptest.NewRequest(http.MethodPut, "/api/nodes/aq-1a2b3c/mode", bytes.NewBufferString(`{"mode":"auto"}`))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("User A PUT mode got %d, want 200", rec.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/nodes/aq-1a2b3c", nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("User A DELETE node got %d, want 204", rec.Code)
	}

	// Setelah dilepas, User A tidak bisa mengakses detailnya lagi (404)
	req = httptest.NewRequest(http.MethodGet, "/api/nodes/aq-1a2b3c", nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("setelah unclaim, User A GET detail got %d, want 404", rec.Code)
	}
}

func TestClaimRateLimit(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	store := newFakeAPIStore()
	server := api.NewServer(store, authSvc, nil, nil)

	user, _ := store.CreateUser(context.Background(), "tester", "hash")
	token, _ := authSvc.Sign(user.ID, user.Username, false)

	badBody, _ := json.Marshal(map[string]string{"node_id": "aq-tidak-ada", "claim_code": "SALAH123"})

	// 10 percobaan gagal pertama mengembalikan 404
	for i := 1; i <= 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/nodes/claim", bytes.NewReader(badBody))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("percobaan ke-%d got status %d, want 404", i, rec.Code)
		}
	}

	// Percobaan ke-11 harus ditolak dengan 429 {"error":"Terlalu banyak percobaan, coba lagi nanti"}
	req := httptest.NewRequest(http.MethodPost, "/api/nodes/claim", bytes.NewReader(badBody))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("percobaan ke-11 got status %d, want 429", rec.Code)
	}
	var errRes map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&errRes)
	if errRes["error"] != "Terlalu banyak percobaan, coba lagi nanti" {
		t.Errorf("got error %q, want %q", errRes["error"], "Terlalu banyak percobaan, coba lagi nanti")
	}
}

func TestAdminEndpointsAuthorizationAndFlow(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	store := newFakeAPIStore()
	server := api.NewServer(store, authSvc, api.NewTempCommander(), nil)

	clientToken, _ := authSvc.Sign(2, "budi", false)
	adminToken, _ := authSvc.Sign(1, "admin", true)

	adminEndpoints := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/api/admin/nodes", ""},
		{http.MethodPost, "/api/admin/nodes", `{"node_id":"aq-1a2b3c"}`},
		{http.MethodPost, "/api/admin/nodes/aq-1a2b3c/claim-code", ""},
		{http.MethodDelete, "/api/admin/nodes/aq-1a2b3c/owner", ""},
		{http.MethodDelete, "/api/admin/nodes/aq-1a2b3c", ""},
	}

	// 1. Klien non-admin ditolak 403 {"error":"Khusus admin"} pada semua endpoint admin
	for _, ep := range adminEndpoints {
		req := httptest.NewRequest(ep.method, ep.path, bytes.NewBufferString(ep.body))
		req.Header.Set("Authorization", "Bearer "+clientToken)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("non-admin %s %s got %d, want 403", ep.method, ep.path, rec.Code)
		}
		var errRes map[string]string
		_ = json.NewDecoder(rec.Body).Decode(&errRes)
		if errRes["error"] != "Khusus admin" {
			t.Errorf("error = %q, want 'Khusus admin'", errRes["error"])
		}
	}

	// 2. Admin POST /api/admin/nodes validasi format node_id -> 400 kalau salah
	req := httptest.NewRequest(http.MethodPost, "/api/admin/nodes", bytes.NewBufferString(`{"node_id":"node-salah"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("admin create node invalid format got %d, want 400", rec.Code)
	}

	// 3. Admin POST /api/admin/nodes format benar (aq-1a2b3c & aq-sim001) -> 201, duplikat -> 409
	req = httptest.NewRequest(http.MethodPost, "/api/admin/nodes", bytes.NewBufferString(`{"node_id":"aq-1a2b3c"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("admin create node got %d, want 201", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/admin/nodes", bytes.NewBufferString(`{"node_id":"aq-1a2b3c"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("admin create duplicate node got %d, want 409", rec.Code)
	}

	// 4. Admin GET /api/admin/nodes -> tidak boleh memuat field sensor/aktuator (temp, ntu, heater, dll.)
	req = httptest.NewRequest(http.MethodGet, "/api/admin/nodes", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin list nodes got %d, want 200", rec.Code)
	}
	rawAdminList := rec.Body.String()
	for _, forbidden := range []string{`"temp"`, `"ntu"`, `"class_temp"`, `"heater"`, `"fan"`, `"pump"`} {
		if strings.Contains(rawAdminList, forbidden) {
			t.Fatalf("response admin memuat field sensor/aktuator terlarang %s: %s", forbidden, rawAdminList)
		}
	}

	// 5. Admin regenerate claim code -> 200
	req = httptest.NewRequest(http.MethodPost, "/api/admin/nodes/aq-1a2b3c/claim-code", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin regenerate claim code got %d, want 200", rec.Code)
	}

	// 6. Admin lepas kepemilikan -> 204
	req = httptest.NewRequest(http.MethodDelete, "/api/admin/nodes/aq-1a2b3c/owner", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("admin delete owner got %d, want 204", rec.Code)
	}

	// 7. Admin hapus node -> 204
	req = httptest.NewRequest(http.MethodDelete, "/api/admin/nodes/aq-1a2b3c", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("admin delete node got %d, want 204", rec.Code)
	}
}

func TestSetActuatorsInputValidation(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	server := api.NewServer(nil, authSvc, nil, nil)
	token, _ := authSvc.Sign(1, "admin")

	tests := []struct {
		name       string
		payload    string
		wantStatus int
	}{
		{
			name:       "field hilang",
			payload:    `{"heater": 1}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nilai bukan 0 atau 1",
			payload:    `{"heater": 2, "fan": 0, "pump": 0}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nilai negatif",
			payload:    `{"heater": 0, "fan": -1, "pump": 0}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/nodes/aq-001/actuators", bytes.NewBufferString(tt.payload))
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("%s: got status %d, want %d", tt.name, rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestUpdateModeInputValidation(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	server := api.NewServer(nil, authSvc, nil, nil)
	token, _ := authSvc.Sign(1, "admin")

	tests := []struct {
		name       string
		mode       string
		wantStatus int
	}{
		{"mode invalid", "semi-auto", http.StatusBadRequest},
		{"mode kosong", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]string{"mode": tt.mode})
			req := httptest.NewRequest(http.MethodPut, "/api/nodes/aq-001/mode", bytes.NewReader(body))
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()

			server.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("%s: got status %d, want %d", tt.name, rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestTelemetryQueryTimeValidation(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	server := api.NewServer(nil, authSvc, nil, nil)
	token, _ := authSvc.Sign(1, "admin")

	// Waktu 'from' salah format -> 400
	req := httptest.NewRequest(http.MethodGet, "/api/nodes/aq-001/telemetry?from=2026-11-08", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want 400", rec.Code)
	}

	// Waktu 'from' lebih baru dari 'to' -> 400
	req = httptest.NewRequest(http.MethodGet, "/api/nodes/aq-001/telemetry?from=2026-11-08T10:00:00Z&to=2026-11-08T08:00:00Z", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want 400", rec.Code)
	}
}

func TestTelemetryDeltaCalculation(t *testing.T) {
	now := time.Now().UTC()
	lessThan48 := now.Add(24 * time.Hour).Sub(now)
	if lessThan48 > 48*time.Hour {
		t.Fatalf("delta 24 jam tidak boleh > 48 jam")
	}

	moreThan48 := now.Add(72 * time.Hour).Sub(now)
	if moreThan48 <= 48*time.Hour {
		t.Fatalf("delta 72 jam harus > 48 jam")
	}
}
