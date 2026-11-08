package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

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

func (s *Server) handleListNodes(w http.ResponseWriter, r *http.Request) {
	if s.store == nil {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	nodes, err := s.store.ListNodeSummaries(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mengambil daftar akuarium")
		return
	}
	writeJSON(w, http.StatusOK, nodes)
}

func (s *Server) handleGetNode(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	if s.store == nil {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	detail, err := s.store.GetNodeDetail(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "akuarium tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mengambil detail akuarium")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) handlePatchNode(w http.ResponseWriter, r *http.Request) {
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

	if s.store == nil {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	detail, err := s.store.UpdateNodeName(r.Context(), nodeID, name)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "akuarium tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memperbarui nama akuarium")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) handleUpdateMode(w http.ResponseWriter, r *http.Request) {
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

	if s.store == nil {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	mode := core.Mode(modeStr)
	detail, err := s.store.UpdateNodeMode(r.Context(), nodeID, mode)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "akuarium tidak ditemukan")
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

	if s.store == nil {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	detail, err := s.store.GetNodeDetail(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "akuarium tidak ditemukan")
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
