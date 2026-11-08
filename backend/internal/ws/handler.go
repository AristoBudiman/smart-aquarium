package ws

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/coder/websocket"
)

// Handler mengembalikan http.HandlerFunc untuk endpoint GET /ws?token=JWT.
func Handler(hub *Hub, authSvc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(r.URL.Query().Get("token"))
		if token == "" {
			writeUnauthorized(w, "parameter query token diperlukan")
			return
		}

		if authSvc != nil {
			if _, err := authSvc.Verify(token); err != nil {
				writeUnauthorized(w, "token tidak valid atau telah kedaluwarsa")
				return
			}
		}

		if hub == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "websocket hub belum diinisialisasi",
			})
			return
		}

		opts := &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		}
		conn, err := websocket.Accept(w, r, opts)
		if err != nil {
			return
		}

		client := newClient(hub, conn)
		hub.register(client)

		ctx := r.Context()
		go client.writePump(ctx)
		client.readPump(ctx)
	}
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
