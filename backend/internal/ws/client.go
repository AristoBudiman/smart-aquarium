package ws

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	queueSize    = 32
	writeTimeout = 5 * time.Second
)

// Client merepresentasikan satu koneksi WebSocket aktif dengan antrean pesan keluar terisolasi.
type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	userID    int64
	isAdmin   bool
	send      chan []byte
	closeOnce sync.Once
}

func newClient(hub *Hub, conn *websocket.Conn, userID int64, isAdmin bool) *Client {
	return &Client{
		hub:     hub,
		conn:    conn,
		userID:  userID,
		isAdmin: isAdmin,
		send:    make(chan []byte, queueSize),
	}
}

func (c *Client) closeSend() {
	c.closeOnce.Do(func() {
		close(c.send)
	})
}

func (c *Client) writePump(ctx context.Context) {
	defer func() {
		c.hub.unregister(c)
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister(c)
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, _, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
	}
}
