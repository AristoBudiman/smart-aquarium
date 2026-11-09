package api

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/ws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Pastikan *db.Store mengimplementasikan DataStore.
var _ DataStore = (*db.Store)(nil)

// DataStore mendefinisikan operasi basis data yang dibutuhkan oleh REST API.
type DataStore interface {
	EnsureAdminUser(ctx context.Context, username, passwordHash string) error
	CreateUser(ctx context.Context, username, passwordHash string) (*db.User, error)
	GetUserByUsername(ctx context.Context, username string) (*db.User, error)
	GetUserByID(ctx context.Context, userID int64) (*db.User, error)
	UpdateUserTelegramChatID(ctx context.Context, userID int64, chatID string) error
	ListNodeSummaries(ctx context.Context, ownerID int64) ([]db.NodeSummary, error)
	GetNodeDetail(ctx context.Context, ownerID int64, nodeID string) (*db.NodeDetail, error)
	UpdateNodeName(ctx context.Context, ownerID int64, nodeID, name string) (*db.NodeDetail, error)
	UpdateNodeMode(ctx context.Context, ownerID int64, nodeID string, mode core.Mode) (*db.NodeDetail, error)
	ClaimNode(ctx context.Context, userID int64, nodeID, claimCode string) (*db.NodeDetail, error)
	ReleaseNodeOwnership(ctx context.Context, ownerID int64, nodeID string) error
	ListAdminNodes(ctx context.Context) ([]db.AdminNode, error)
	CreateAdminNode(ctx context.Context, nodeID string) (*db.AdminNode, error)
	RegenerateClaimCode(ctx context.Context, nodeID string) (*db.AdminNode, error)
	AdminReleaseNodeOwner(ctx context.Context, nodeID string) error
	AdminDeleteNode(ctx context.Context, nodeID string) error
	GetTelemetryHistory(ctx context.Context, nodeID string, from, to time.Time) ([]db.TelemetryPoint, error)
	GetNodeLogs(ctx context.Context, nodeID string) ([]db.ActuatorLog, []db.AlertLog, error)
	LogActuators(ctx context.Context, nodeID string, a core.Actuators, source core.Source) error
}

const (
	claimRateLimitMax    = 10
	claimRateLimitWindow = 10 * time.Minute
)

// Server mengelola routing REST API dan integrasi dependensi backend.
type Server struct {
	router      chi.Router
	store       DataStore
	authSvc     *auth.Service
	commander   core.Commander
	broadcaster core.Broadcaster
	wsHub       *ws.Hub

	claimMu       sync.Mutex
	claimFailures map[int64][]time.Time
}

// NewServer membuat instance Server baru dan mendaftarkan seluruh endpoint kontrak.
func NewServer(
	store DataStore,
	authSvc *auth.Service,
	commander core.Commander,
	broadcaster core.Broadcaster,
) *Server {
	s := &Server{
		store:         store,
		authSvc:       authSvc,
		commander:     commander,
		broadcaster:   broadcaster,
		claimFailures: make(map[int64][]time.Time),
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
	r.Post("/api/auth/register", s.handleRegister)

	// Endpoint WebSocket (verifikasi token dari query param ?token=JWT)
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub := s.wsHub
		if hub == nil {
			if h, ok := s.broadcaster.(*ws.Hub); ok {
				hub = h
			}
		}
		ws.Handler(hub, s.authSvc)(w, r)
	})

	// Endpoint terproteksi JWT
	r.Group(func(protected chi.Router) {
		if s.authSvc != nil {
			protected.Use(s.authSvc.Middleware)
		}

		// Profile & Telegram endpoints
		protected.Get("/api/me", s.handleGetMe)
		protected.Put("/api/me/telegram", s.handleUpdateTelegram)

		// Klaim & kepemilikan node klien
		protected.Post("/api/nodes/claim", s.handleClaimNode)
		protected.Get("/api/nodes", s.handleListNodes)
		protected.Get("/api/nodes/{id}", s.handleGetNode)
		protected.Patch("/api/nodes/{id}", s.handlePatchNode)
		protected.Delete("/api/nodes/{id}", s.handleDeleteNode)
		protected.Put("/api/nodes/{id}/mode", s.handleUpdateMode)
		protected.Put("/api/nodes/{id}/actuators", s.handleSetActuators)
		protected.Get("/api/nodes/{id}/telemetry", s.handleGetTelemetry)
		protected.Get("/api/nodes/{id}/logs", s.handleGetLogs)

		// Endpoint khusus admin
		protected.Route("/api/admin", func(admin chi.Router) {
			admin.Use(auth.RequireAdmin)
			admin.Get("/nodes", s.handleAdminListNodes)
			admin.Post("/nodes", s.handleAdminCreateNode)
			admin.Post("/nodes/{id}/claim-code", s.handleAdminRegenerateClaimCode)
			admin.Delete("/nodes/{id}/owner", s.handleAdminDeleteNodeOwner)
			admin.Delete("/nodes/{id}", s.handleAdminDeleteNode)
		})
	})

	s.router = r
	return s
}

func (s *Server) hasStore() bool {
	if s.store == nil {
		return false
	}
	v := reflect.ValueOf(s.store)
	return v.Kind() != reflect.Ptr || !v.IsNil()
}

func (s *Server) isClaimRateLimited(userID int64, now time.Time) bool {
	s.claimMu.Lock()
	defer s.claimMu.Unlock()

	cutoff := now.Add(-claimRateLimitWindow)
	attempts := s.claimFailures[userID]
	kept := attempts[:0]
	for _, ts := range attempts {
		if ts.After(cutoff) {
			kept = append(kept, ts)
		}
	}
	s.claimFailures[userID] = kept
	return len(kept) >= claimRateLimitMax
}

func (s *Server) recordClaimFailure(userID int64, now time.Time) {
	s.claimMu.Lock()
	defer s.claimMu.Unlock()

	cutoff := now.Add(-claimRateLimitWindow)
	attempts := s.claimFailures[userID]
	kept := attempts[:0]
	for _, ts := range attempts {
		if ts.After(cutoff) {
			kept = append(kept, ts)
		}
	}
	s.claimFailures[userID] = append(kept, now)
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
