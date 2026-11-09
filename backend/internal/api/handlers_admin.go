package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

var validAdminNodeIDRegexp = regexp.MustCompile(`^(aq-[0-9a-f]{6}|aq-sim[0-9]{3})$`)

type createAdminNodeRequest struct {
	NodeID string `json:"node_id"`
}

func (s *Server) handleAdminListNodes(w http.ResponseWriter, r *http.Request) {
	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	nodes, err := s.store.ListAdminNodes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mengambil daftar node admin")
		return
	}

	writeJSON(w, http.StatusOK, nodes)
}

func (s *Server) handleAdminCreateNode(w http.ResponseWriter, r *http.Request) {
	var req createAdminNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "format JSON tidak valid")
		return
	}

	nodeID := strings.TrimSpace(req.NodeID)
	if !validAdminNodeIDRegexp.MatchString(nodeID) {
		writeError(w, http.StatusBadRequest, "format node_id tidak valid")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	node, err := s.store.CreateAdminNode(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, db.ErrNodeAlreadyExists) {
			writeError(w, http.StatusConflict, "Node sudah terdaftar")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mendaftarkan node")
		return
	}

	// Kirim config retained manual agar sudah tersedia di broker sebelum alat pertama kali menyala.
	if s.commander != nil {
		if err := s.commander.SetMode(r.Context(), nodeID, core.ModeManual); err != nil {
			writeError(w, http.StatusInternalServerError, "gagal mengirim konfigurasi awal ke broker MQTT")
			return
		}
	}

	writeJSON(w, http.StatusCreated, node)
}

func (s *Server) handleAdminRegenerateClaimCode(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	node, err := s.store.RegenerateClaimCode(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memperbarui kode klaim")
		return
	}

	writeJSON(w, http.StatusOK, node)
}

func (s *Server) handleAdminDeleteNodeOwner(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	if err := s.store.AdminReleaseNodeOwner(r.Context(), nodeID); err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal melepas kepemilikan node")
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

// handleAdminDeleteNode menghapus node beserta seluruh datanya secara cascade.
// Catatan: node fisik atau simulator yang masih menyala dan mengirim pesan MQTT
// akan terdaftar lagi secara otomatis melalui mekanisme EnsureNode pada ingest.
func (s *Server) handleAdminDeleteNode(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	if !s.hasStore() {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	if err := s.store.AdminDeleteNode(r.Context(), nodeID); err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal menghapus node")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
