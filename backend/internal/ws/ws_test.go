package ws_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/ws"
	"github.com/coder/websocket"
)

const testSecret = "super-secret-key-that-is-at-least-32-characters-long!"

func TestHubImplementsBroadcaster(t *testing.T) {
	var _ core.Broadcaster = (*ws.Hub)(nil)
}

func TestWSHandlerAuth(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	hub := ws.NewHub()
	handler := ws.Handler(hub, authSvc)

	// Test 1: Tanpa query token -> 401
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", rec.Code)
	}

	// Test 2: Token tidak valid -> 401
	req = httptest.NewRequest(http.MethodGet, "/ws?token=invalid-token", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", rec.Code)
	}
}

func TestWSConnectionAndBroadcast(t *testing.T) {
	authSvc, _ := auth.NewService(testSecret)
	hub := ws.NewHub()
	token, _ := authSvc.Sign(1, "admin")

	mux := http.NewServeMux()
	mux.Handle("/ws", ws.Handler(hub, authSvc))

	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?token=" + token

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("websocket.Dial gagal: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Beri jeda singkat agar klien terdaftar di Hub
	time.Sleep(50 * time.Millisecond)
	if hub.ClientCount() != 1 {
		t.Fatalf("jumlah klien di hub: got %d, want 1", hub.ClientCount())
	}

	// Publish telemetry event dari hub
	event := core.TelemetryEvent{
		Type:   "telemetry",
		NodeID: "aq-1a2b3c",
		TS:     time.Date(2026, 11, 8, 3, 0, 0, 0, time.UTC),
		Temp:   27.3,
		NTU:    12.4,
		Classes: core.Classes{
			Temp:      core.ClassNormal,
			Turbidity: core.ClassNormal,
			Overall:   core.ClassNormal,
			Status:    core.ClassNormal,
		},
	}
	hub.Publish(event)

	// Baca pesan di sisi klien
	readCtx, readCancel := context.WithTimeout(ctx, 2*time.Second)
	defer readCancel()

	msgType, data, err := conn.Read(readCtx)
	if err != nil {
		t.Fatalf("conn.Read gagal: %v", err)
	}
	if msgType != websocket.MessageText {
		t.Errorf("got msgType %v, want MessageText", msgType)
	}

	var received core.TelemetryEvent
	if err := json.Unmarshal(data, &received); err != nil {
		t.Fatalf("gagal unmarshal event: %v", err)
	}

	if received.Type != "telemetry" || received.NodeID != "aq-1a2b3c" || received.Temp != 27.3 {
		t.Errorf("data event tidak cocok: %+v", received)
	}
}

func TestHubSlowConsumerEviction(t *testing.T) {
	hub := ws.NewHub()

	// Buat client langsung dan daftarkan ke hub tanpa menjalankan writePump
	// agar channel tidak dikosongkan.
	client := ws.NewTestClient(hub)
	hub.RegisterForTest(client)

	if hub.ClientCount() != 1 {
		t.Fatalf("jumlah klien di hub: got %d, want 1", hub.ClientCount())
	}

	// Isi antrean 32 pesan hingga penuh
	for i := 0; i < 32; i++ {
		hub.Publish(map[string]any{"idx": i})
	}

	if hub.ClientCount() != 1 {
		t.Fatalf("klien masih harus ada saat buffer baru pas 32, sisa: %d", hub.ClientCount())
	}

	// Kirim pesan ke-33 untuk memicu proteksi slow consumer
	hub.Publish(map[string]any{"idx": 33})

	// Klien harus segera dikeluarkan dari Hub
	if hub.ClientCount() != 0 {
		t.Errorf("klien lambat harus dikeluarkan dari hub, sisa: %d", hub.ClientCount())
	}
}

type fakeOwnerResolver struct {
	owners map[string]int64
}

func (f *fakeOwnerResolver) NodeOwner(_ context.Context, nodeID string) (int64, bool, error) {
	ownerID, ok := f.owners[nodeID]
	return ownerID, ok, nil
}

func TestHubPublishOwnershipAndAdminFiltering(t *testing.T) {
	resolver := &fakeOwnerResolver{
		owners: map[string]int64{
			"aq-milik-a": 10,
			// "aq-tanpa-pemilik" belum diklaim
		},
	}
	hub := ws.NewHub(resolver)

	clientA := ws.NewTestClientWithIdentity(hub, 10, false)
	clientB := ws.NewTestClientWithIdentity(hub, 20, false)
	clientAdmin := ws.NewTestClientWithIdentity(hub, 1, true)

	hub.RegisterForTest(clientA)
	hub.RegisterForTest(clientB)
	hub.RegisterForTest(clientAdmin)

	// 1. Telemetry untuk aq-milik-a -> hanya diterima oleh clientA (pemilik), B dan Admin tidak menerima
	hub.Publish(core.TelemetryEvent{
		Type:   "telemetry",
		NodeID: "aq-milik-a",
		TS:     time.Now().UTC(),
		Temp:   27.5,
		NTU:    10.0,
	})

	if len(ws.SendChanForTest(clientA)) != 1 {
		t.Fatalf("clientA (pemilik) seharusnya menerima telemetry, got %d", len(ws.SendChanForTest(clientA)))
	}
	<-ws.SendChanForTest(clientA)

	if len(ws.SendChanForTest(clientB)) != 0 {
		t.Fatalf("clientB (bukan pemilik) tidak boleh menerima telemetry node A")
	}
	if len(ws.SendChanForTest(clientAdmin)) != 0 {
		t.Fatalf("admin tidak boleh menerima telemetry node klien")
	}

	// 2. node_status untuk aq-milik-a -> diterima oleh clientA (pemilik) dan clientAdmin, tetapi bukan clientB
	hub.Publish(core.NodeStatusEvent{
		Type:   "node_status",
		NodeID: "aq-milik-a",
		Online: true,
	})

	if len(ws.SendChanForTest(clientA)) != 1 {
		t.Fatalf("clientA seharusnya menerima node_status miliknya")
	}
	<-ws.SendChanForTest(clientA)

	if len(ws.SendChanForTest(clientAdmin)) != 1 {
		t.Fatalf("admin seharusnya menerima node_status")
	}
	<-ws.SendChanForTest(clientAdmin)

	if len(ws.SendChanForTest(clientB)) != 0 {
		t.Fatalf("clientB tidak boleh menerima node_status node A")
	}

	// 3. Event untuk node tanpa pemilik (aq-tanpa-pemilik):
	// telemetry tidak diterima siapa pun; node_status hanya diterima oleh admin.
	hub.Publish(core.TelemetryEvent{
		Type:   "telemetry",
		NodeID: "aq-tanpa-pemilik",
		TS:     time.Now().UTC(),
		Temp:   28.0,
		NTU:    15.0,
	})
	hub.Publish(core.NodeStatusEvent{
		Type:   "node_status",
		NodeID: "aq-tanpa-pemilik",
		Online: true,
	})

	if len(ws.SendChanForTest(clientA)) != 0 || len(ws.SendChanForTest(clientB)) != 0 {
		t.Fatalf("klien biasa tidak boleh menerima event dari node tanpa pemilik")
	}
	if len(ws.SendChanForTest(clientAdmin)) != 1 {
		t.Fatalf("admin seharusnya hanya menerima 1 event (node_status) dari node tanpa pemilik, got %d", len(ws.SendChanForTest(clientAdmin)))
	}
}
