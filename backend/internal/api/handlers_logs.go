package api

import (
	"errors"
	"net/http"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

type nodeLogsResponse struct {
	Actuators []db.ActuatorLog `json:"actuators"`
	Alerts    []db.AlertLog    `json:"alerts"`
}

func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
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

	// Pastikan node ada dan dimiliki oleh user yang login (disaring di SQL).
	if _, err := s.store.GetNodeDetail(r.Context(), claims.UserID, nodeID); err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "Node tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memeriksa keberadaan node")
		return
	}

	actuators, alerts, err := s.store.GetNodeLogs(r.Context(), nodeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mengambil log aktivitas node")
		return
	}

	writeJSON(w, http.StatusOK, nodeLogsResponse{
		Actuators: actuators,
		Alerts:    alerts,
	})
}
