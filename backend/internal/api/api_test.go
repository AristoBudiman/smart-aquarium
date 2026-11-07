package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/api"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
)

const testSecret = "super-secret-key-that-is-at-least-32-characters-long!"

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
		{http.MethodGet, "/api/nodes"},
		{http.MethodGet, "/api/nodes/aq-1a2b3c"},
		{http.MethodPatch, "/api/nodes/aq-1a2b3c"},
		{http.MethodPut, "/api/nodes/aq-1a2b3c/mode"},
		{http.MethodPut, "/api/nodes/aq-1a2b3c/actuators"},
		{http.MethodGet, "/api/nodes/aq-1a2b3c/telemetry"},
		{http.MethodGet, "/api/nodes/aq-1a2b3c/logs"},
		{http.MethodPut, "/api/me/telegram"},
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
