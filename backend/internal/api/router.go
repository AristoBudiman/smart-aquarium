package api

import (
	"net/http"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server mengelola routing REST API dan integrasi dependensi backend.
type Server struct {
	router      chi.Router
	store       *db.Store
	authSvc     *auth.Service
	commander   core.Commander
	broadcaster core.Broadcaster
	wsHub       *ws.Hub
}

// NewServer membuat instance Server baru dan mendaftarkan seluruh endpoint kontrak.
func NewServer(
	store *db.Store,
	authSvc *auth.Service,
	commander core.Commander,
	broadcaster core.Broadcaster,
) *Server {
	s := &Server{
		store:       store,
		authSvc:     authSvc,
		commander:   commander,
		broadcaster: broadcaster,
	}

	if hub, ok := broadcaster.(*ws.Hub); ok {
		s.wsHub = hub
	}

	r := chi.NewRouter()

	// Middleware standar
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Endpoint publik
	r.Get("/api/health", s.handleHealth)
	r.Post("/api/auth/login", s.handleLogin)

	// Endpoint WebSocket (verifikasi token dari query param ?token=JWT)
	r.Get("/ws", ws.Handler(s.wsHub, s.authSvc))

	// Endpoint terproteksi JWT
	r.Group(func(protected chi.Router) {
		if s.authSvc != nil {
			protected.Use(s.authSvc.Middleware)
		}

		// Node endpoints
		protected.Get("/api/nodes", s.handleListNodes)
		protected.Get("/api/nodes/{id}", s.handleGetNode)
		protected.Patch("/api/nodes/{id}", s.handlePatchNode)
		protected.Put("/api/nodes/{id}/mode", s.handleUpdateMode)
		protected.Put("/api/nodes/{id}/actuators", s.handleSetActuators)
		protected.Get("/api/nodes/{id}/telemetry", s.handleGetTelemetry)
		protected.Get("/api/nodes/{id}/logs", s.handleGetLogs)

		// Profile / Telegram endpoint
		protected.Put("/api/me/telegram", s.handleUpdateTelegram)
	})

	s.router = r
	return s
}

// SetWSHub menetapkan WebSocket hub jika dibuat secara terpisah.
func (s *Server) SetWSHub(hub *ws.Hub) {
	s.wsHub = hub
	s.broadcaster = hub
}

// Router mengembalikan router chi yang telah dikonfigurasi.
func (s *Server) Router() chi.Router {
	return s.router
}

// ServeHTTP mengimplementasikan http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
