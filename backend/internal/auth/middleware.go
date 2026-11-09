package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey struct{}

var userContextKey = contextKey{}

// WithClaims menyisipkan data Claims pengguna ke dalam context.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, userContextKey, claims)
}

// FromContext mengambil data Claims pengguna dari context.
func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(userContextKey).(*Claims)
	return claims, ok && claims != nil
}

// Middleware memeriksa header Authorization: Bearer <token>, memvalidasinya,
// dan menyisipkan identitas pengguna ke dalam context request.
func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeUnauthorized(w, "header Authorization diperlukan")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeUnauthorized(w, "format header Authorization harus 'Bearer <token>'")
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			writeUnauthorized(w, "token tidak boleh kosong")
			return
		}

		claims, err := s.Verify(tokenStr)
		if err != nil {
			writeUnauthorized(w, "token tidak valid atau telah kedaluwarsa")
			return
		}

		ctx := WithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin memastikan pengguna yang terautentikasi memiliki hak akses admin.
// Mengembalikan 403 {"error":"Khusus admin"} jika bukan admin.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := FromContext(r.Context())
		if !ok || claims == nil {
			writeUnauthorized(w, "tidak terautentikasi")
			return
		}
		if !claims.Admin {
			writeForbidden(w, "Khusus admin")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin adalah method pembantu pada Service yang memanggil RequireAdmin.
func (s *Service) RequireAdmin(next http.Handler) http.Handler {
	return RequireAdmin(next)
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func writeForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
