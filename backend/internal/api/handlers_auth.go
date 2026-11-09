package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
)

var validUsernameRegexp = regexp.MustCompile(`^[a-z0-9_.-]{3,32}$`)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type meResponse struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	IsAdmin        bool   `json:"is_admin"`
	TelegramLinked bool   `json:"telegram_linked"`
}

type updateTelegramRequest struct {
	TelegramChatID string `json:"telegram_chat_id"`
}

// BootstrapAdmin menginisialisasi akun administrator awal jika belum ada di database,
// atau memastikan is_admin = TRUE jika username sudah ada tanpa mengubah password lama.
func BootstrapAdmin(ctx context.Context, store DataStore, username, password string) error {
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

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	user, err := s.store.GetUserByUsername(r.Context(), req.Username)
	if errors.Is(err, db.ErrUserNotFound) && strings.ToLower(req.Username) != req.Username {
		user, err = s.store.GetUserByUsername(r.Context(), strings.ToLower(req.Username))
	}
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

	token, err := s.authSvc.Sign(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal membuat token autentikasi")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	username := strings.ToLower(strings.TrimSpace(req.Username))
	if !validUsernameRegexp.MatchString(username) {
		writeError(w, http.StatusBadRequest, "username harus 3-32 karakter dan hanya boleh berisi a-z, 0-9, _, ., atau -")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password minimal 8 karakter")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memproses kata sandi")
		return
	}

	user, err := s.store.CreateUser(r.Context(), username, hashed)
	if err != nil {
		if errors.Is(err, db.ErrUsernameTaken) {
			writeError(w, http.StatusConflict, "Username sudah dipakai")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mendaftarkan pengguna")
		return
	}

	token, err := s.authSvc.Sign(user.ID, user.Username, false)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal membuat token autentikasi")
		return
	}

	writeJSON(w, http.StatusCreated, loginResponse{Token: token})
}

func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	user, err := s.store.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			writeError(w, http.StatusUnauthorized, "pengguna tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mengambil profil pengguna")
		return
	}

	linked := user.TelegramChatID != nil && strings.TrimSpace(*user.TelegramChatID) != ""
	writeJSON(w, http.StatusOK, meResponse{
		ID:             user.ID,
		Username:       user.Username,
		IsAdmin:        user.IsAdmin,
		TelegramLinked: linked,
	})
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

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	chatID := strings.TrimSpace(req.TelegramChatID)
	if err := s.store.UpdateUserTelegramChatID(r.Context(), claims.UserID, chatID); err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memperbarui chat ID telegram")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
