package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type updateTelegramRequest struct {
	TelegramChatID string `json:"telegram_chat_id"`
}

// BootstrapAdmin menginisialisasi akun administrator awal jika belum ada di database.
func BootstrapAdmin(ctx context.Context, store *db.Store, username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return nil
	}

	hashed, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("bootstrap admin hash password gagal: %w", err)
	}

	if err := store.EnsureAdminUser(ctx, username, hashed); err != nil {
		return fmt.Errorf("bootstrap admin ensure user gagal: %w", err)
	}

	log.Printf("api: akun admin %q terverifikasi/diinisialisasi", username)
	return nil
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusUnauthorized, "username atau password salah")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusUnauthorized, "username atau password salah")
		return
	}

	user, err := s.store.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			writeError(w, http.StatusUnauthorized, "username atau password salah")
			return
		}
		writeError(w, http.StatusInternalServerError, "kesalahan internal basis data")
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "username atau password salah")
		return
	}

	token, err := s.authSvc.Sign(user.ID, user.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal membuat token autentikasi")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (s *Server) handleUpdateTelegram(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	var req updateTelegramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	chatID := strings.TrimSpace(req.TelegramChatID)
	if err := s.store.UpdateUserTelegramChatID(r.Context(), claims.UserID, chatID); err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memperbarui chat ID telegram")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
