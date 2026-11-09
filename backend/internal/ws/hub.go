package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/coder/websocket"
)

var _ core.Broadcaster = (*Hub)(nil)

// OwnerResolver mencari informasi kepemilikan node untuk penyaringan event WebSocket.
type OwnerResolver interface {
	NodeOwner(ctx context.Context, nodeID string) (ownerID int64, owned bool, err error)
}

// Hub mengelola koneksi seluruh klien WebSocket yang aktif dan menyiarkan event secara thread-safe.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
	owners  OwnerResolver
}

// NewHub membuat instance Hub baru.
func NewHub(owners ...OwnerResolver) *Hub {
	var resolver OwnerResolver
	if len(owners) > 0 {
		resolver = owners[0]
	}
	return &Hub{
		clients: make(map[*Client]bool),
		owners:  resolver,
	}
}

// SetOwnerResolver menetapkan penyedia informasi kepemilikan node pada Hub.
func (h *Hub) SetOwnerResolver(owners OwnerResolver) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.owners = owners
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		c.closeSend()
	}
}

// ClientCount mengembalikan jumlah klien yang sedang terhubung ke Hub.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Publish menyiarkan event ke klien yang berhak sesuai kontrak:
// - Klien biasa hanya menerima event dari node miliknya (UserID == ownerID).
// - Admin hanya menerima event bertipe "node_status" (baik node punya pemilik maupun belum).
// - Node tanpa pemilik hanya menyiarkan "node_status" ke admin.
// Jika buffer pesan klien penuh (konsumen lambat), Hub segera mengeluarkan dan menutup koneksi klien tersebut untuk mencegah pemblokiran sistem.
func (h *Hub) Publish(event any) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("ws: gagal serialisasi event ke JSON: %v", err)
		return
	}

	var meta struct {
		Type   string `json:"type"`
		NodeID string `json:"node_id"`
	}
	_ = json.Unmarshal(data, &meta)

	h.mu.RLock()
	owners := h.owners
	h.mu.RUnlock()

	var (
		ownerID   int64
		owned     bool
		hasFilter bool
	)
	if owners != nil && meta.NodeID != "" {
		hasFilter = true
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		ownerID, owned, err = owners.NodeOwner(ctx, meta.NodeID)
		cancel()
		if err != nil {
			log.Printf("ws: gagal memeriksa pemilik node %s: %v", meta.NodeID, err)
			return
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.clients {
		if c.isAdmin {
			if meta.Type != "node_status" {
				continue
			}
		} else if hasFilter {
			if !owned || c.userID != ownerID {
				continue
			}
		}

		select {
		case c.send <- data:
		default:
			// Buffer penuh: keluarkan dan tutup koneksi klien lambat segera
			delete(h.clients, c)
			c.closeSend()
			if c.conn != nil {
				_ = c.conn.Close(websocket.StatusPolicyViolation, "antrean pesan penuh (slow consumer)")
			}
		}
	}
}
