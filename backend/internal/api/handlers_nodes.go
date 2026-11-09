package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

type claimNodeRequest struct {
	NodeID    string `json:"node_id"`
	ClaimCode string `json:"claim_code"`
}

type updateNameRequest struct {
	Name string `json:"name"`
}

type updateModeRequest struct {
	Mode string `json:"mode"`
}

type setActuatorsRequest struct {
	Heater *int `json:"heater"`
	Fan    *int `json:"fan"`
	Pump   *int `json:"pump"`
}

func (s *Server) handleClaimNode(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	if claims.Admin {
		writeError(w, http.StatusForbidden, "Admin tidak dapat mengklaim akuarium")
		return
	}

	now := time.Now().UTC()
	if s.isClaimRateLimited(claims.UserID, now) {
		writeError(w, http.StatusTooManyRequests, "Terlalu banyak percobaan, coba lagi nanti")
		return
	}

	var req claimNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	nodeID := strings.TrimSpace(req.NodeID)
	claimCode := strings.TrimSpace(req.ClaimCode)
	if nodeID == "" || claimCode == "" {
		s.recordClaimFailure(claims.UserID, now)
		writeError(w, http.StatusNotFound, "ID node atau kode klaim salah")
		return
	}

	detail, err := s.store.ClaimNode(r.Context(), claims.UserID, nodeID, claimCode)
	if err != nil {
		if errors.Is(err, db.ErrInvalidClaim) {
			s.recordClaimFailure(claims.UserID, now)
			writeError(w, http.StatusNotFound, "ID node atau kode klaim salah")
			return
		}
		if errors.Is(err, db.ErrNodeAlreadyClaimed) {
			s.recordClaimFailure(claims.UserID, now)
			writeError(w, http.StatusConflict, "Akuarium sudah terdaftar di akun lain")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mengklaim akuarium")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) handleDeleteNode(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	if err := s.store.ReleaseNodeOwnership(r.Context(), claims.UserID, nodeID); err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal melepas kepemilikan akuarium")
		return
	}

	if s.commander != nil {
		if err := s.commander.SetMode(r.Context(), nodeID, core.ModeManual); err != nil {
			writeError(w, http.StatusInternalServerError, "gagal mengirim perintah mode ke node")
			return
		}
	}

	if s.broadcaster != nil {
		s.broadcaster.Publish(core.ModeEvent{
			Type:   "mode",
			NodeID: nodeID,
			Mode:   core.ModeManual,
		})
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListNodes(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	nodes, err := s.store.ListNodeSummaries(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mengambil daftar akuarium")
		return
	}
	writeJSON(w, http.StatusOK, nodes)
}

func (s *Server) handleGetNode(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	detail, err := s.store.GetNodeDetail(r.Context(), claims.UserID, nodeID)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mengambil detail akuarium")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) handlePatchNode(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	var req updateNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "nama akuarium tidak boleh kosong")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	detail, err := s.store.UpdateNodeName(r.Context(), claims.UserID, nodeID, name)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memperbarui nama akuarium")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) handleUpdateMode(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	var req updateModeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	modeStr := strings.ToLower(strings.TrimSpace(req.Mode))
	if modeStr != string(core.ModeAuto) && modeStr != string(core.ModeManual) {
		writeError(w, http.StatusBadRequest, "mode harus berupa 'auto' atau 'manual'")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	mode := core.Mode(modeStr)
	detail, err := s.store.UpdateNodeMode(r.Context(), claims.UserID, nodeID, mode)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memperbarui mode akuarium")
		return
	}

	if s.commander != nil {
		if err := s.commander.SetMode(r.Context(), nodeID, mode); err != nil {
			writeError(w, http.StatusInternalServerError, "gagal mengirim perintah mode ke node")
			return
		}
	}

	if s.broadcaster != nil {
		s.broadcaster.Publish(core.ModeEvent{
			Type:   "mode",
			NodeID: nodeID,
			Mode:   mode,
		})
	}

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) handleSetActuators(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.FromContext(r.Context())
	if !ok || claims == nil {
		writeError(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}

	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	var req setActuatorsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	if req.Heater == nil || req.Fan == nil || req.Pump == nil {
		writeError(w, http.StatusBadRequest, "heater, fan, dan pump harus disediakan")
		return
	}

	if (*req.Heater != 0 && *req.Heater != 1) ||
		(*req.Fan != 0 && *req.Fan != 1) ||
		(*req.Pump != 0 && *req.Pump != 1) {
		writeError(w, http.StatusBadRequest, "nilai aktuator harus 0 atau 1")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	detail, err := s.store.GetNodeDetail(r.Context(), claims.UserID, nodeID)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memeriksa status node")
		return
	}

	// Kontrak: Tolak 409 Conflict saat mode Auto
	if detail.Mode == core.ModeAuto {
		writeError(w, http.StatusConflict, "perintah aktuator manual ditolak karena node berada dalam mode auto")
		return
	}

	acts := core.Actuators{
		Heater: *req.Heater,
		Fan:    *req.Fan,
		Pump:   *req.Pump,
	}

	if s.commander != nil {
		if err := s.commander.SendActuators(r.Context(), nodeID, acts); err != nil {
			writeError(w, http.StatusInternalServerError, "gagal mengirim perintah aktuator")
			return
		}
	}

	if err := s.store.LogActuators(r.Context(), nodeID, acts, core.SourceManual); err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mencatat log aktuator")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
