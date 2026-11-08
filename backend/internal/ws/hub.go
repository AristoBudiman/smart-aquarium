package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/coder/websocket"
)

var _ core.Broadcaster = (*Hub)(nil)

// Hub mengelola koneksi seluruh klien WebSocket yang aktif dan menyiarkan event secara thread-safe.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

// NewHub membuat instance Hub baru.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
	}
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

// Publish menyiarkan event ke seluruh klien yang terhubung sesuai kontrak.
// Jika buffer pesan klien penuh (konsumen lambat), Hub segera mengeluarkan dan menutup koneksi klien tersebut untuk mencegah pemblokiran sistem.
func (h *Hub) Publish(event any) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("ws: gagal serialisasi event ke JSON: %v", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.clients {
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
