package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "super-secret-key-that-is-at-least-32-characters-long!"

func TestNewServiceSecretLength(t *testing.T) {
	_, err := auth.NewService("too-short")
	if err == nil {
		t.Fatal("diharapkan error jika panjang secret < 32 karakter")
	}

	svc, err := auth.NewService(testSecret)
	if err != nil {
		t.Fatalf("gagal membuat service: %v", err)
	}
	if svc == nil {
		t.Fatal("service bernilai nil")
	}
}

func TestSignAndVerify(t *testing.T) {
	svc, err := auth.NewService(testSecret)
	if err != nil {
		t.Fatalf("gagal membuat service: %v", err)
	}

	token, err := svc.Sign(42, "admin")
	if err != nil {
		t.Fatalf("gagal sign token: %v", err)
	}

	claims, err := svc.Verify(token)
	if err != nil {
		t.Fatalf("gagal verify token: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID tidak sesuai: got %d, want 42", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("Username tidak sesuai: got %q, want %q", claims.Username, "admin")
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	svc, err := auth.NewService(testSecret)
	if err != nil {
		t.Fatalf("gagal membuat service: %v", err)
	}

	_, err = svc.Verify("invalid.token.string")
	if err == nil {
		t.Fatal("diharapkan error verifikasi token palsu")
	}

	otherSvc, _ := auth.NewService("another-different-secret-key-32-chars-long!!")
	token, _ := otherSvc.Sign(1, "user")
	_, err = svc.Verify(token)
	if err == nil {
		t.Fatal("diharapkan error saat token diverifikasi dengan secret berbeda")
	}

	// Test expired token
	past := time.Now().Add(-1 * time.Hour)
	expiredClaims := auth.Claims{
		UserID:   1,
		Username: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user",
			IssuedAt:  jwt.NewNumericDate(past.Add(-24 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(past),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	signedExpired, _ := expiredToken.SignedString([]byte(testSecret))
	_, err = svc.Verify(signedExpired)
	if err == nil {
		t.Fatal("diharapkan error saat token telah kedaluwarsa")
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "rahasia123"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("gagal hash password: %v", err)
	}

	if !auth.CheckPassword(password, hash) {
		t.Fatal("password yang benar gagal diverifikasi")
	}

	if auth.CheckPassword("salah123", hash) {
		t.Fatal("password yang salah terverifikasi")
	}
}

func TestMiddleware(t *testing.T) {
	svc, _ := auth.NewService(testSecret)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := auth.FromContext(r.Context())
		if !ok || claims == nil {
			http.Error(w, "no claims", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"user_id":  claims.UserID,
			"username": claims.Username,
		})
	})

	handler := svc.Middleware(nextHandler)

	// Test 1: Tanpa header Authorization -> 401
	req := httptest.NewRequest(http.MethodGet, "/api/nodes", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", rec.Code)
	}

	// Test 2: Format header salah -> 401
	req = httptest.NewRequest(http.MethodGet, "/api/nodes", nil)
	req.Header.Set("Authorization", "Basic 12345")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", rec.Code)
	}

	// Test 3: Token valid -> 200
	token, _ := svc.Sign(99, "john_doe")
	req = httptest.NewRequest(http.MethodGet, "/api/nodes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want 200", rec.Code)
	}

	var res map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("gagal decode response: %v", err)
	}
	if res["username"] != "john_doe" {
		t.Errorf("got username %v, want john_doe", res["username"])
	}
}
