package ws

// NewTestClient membuat instance Client untuk keperluan pengujian unit.
func NewTestClient(hub *Hub) *Client {
	return newClient(hub, nil)
}

// RegisterForTest mendaftarkan klien ke hub untuk keperluan pengujian unit.
func (h *Hub) RegisterForTest(c *Client) {
	h.register(c)
}
