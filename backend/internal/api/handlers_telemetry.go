package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleGetTelemetry(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "node_id tidak valid")
		return
	}

	now := time.Now().UTC()
	from := now.Add(-24 * time.Hour)
	to := now

	queryFrom := r.URL.Query().Get("from")
	if queryFrom != "" {
		t, err := time.Parse(time.RFC3339, queryFrom)
		if err != nil {
			writeError(w, http.StatusBadRequest, "parameter 'from' harus berformat ISO 8601 UTC (RFC3339)")
			return
		}
		from = t.UTC()
	}

	queryTo := r.URL.Query().Get("to")
	if queryTo != "" {
		t, err := time.Parse(time.RFC3339, queryTo)
		if err != nil {
			writeError(w, http.StatusBadRequest, "parameter 'to' harus berformat ISO 8601 UTC (RFC3339)")
			return
		}
		to = t.UTC()
	}

	if from.After(to) {
		writeError(w, http.StatusBadRequest, "parameter 'from' tidak boleh lebih baru dari 'to'")
		return
	}

	if s.store == nil {
		writeError(w, http.StatusInternalServerError, "basis data belum diinisialisasi")
		return
	}

	// Pastikan node ada
	if _, err := s.store.GetNodeDetail(r.Context(), nodeID); err != nil {
		if errors.Is(err, db.ErrNodeNotFound) {
			writeError(w, http.StatusNotFound, "akuarium tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memeriksa keberadaan node")
		return
	}

	points, err := s.store.GetTelemetryHistory(r.Context(), nodeID, from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mengambil data riwayat telemetry")
		return
	}

	writeJSON(w, http.StatusOK, points)
}
