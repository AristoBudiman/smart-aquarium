package ws

// NewTestClient membuat instance Client untuk keperluan pengujian unit.
func NewTestClient(hub *Hub) *Client {
	return newClient(hub, nil, 0, false)
}

// NewTestClientWithIdentity membuat instance Client dengan identitas user/admin untuk pengujian unit.
func NewTestClientWithIdentity(hub *Hub, userID int64, isAdmin bool) *Client {
	return newClient(hub, nil, userID, isAdmin)
}

// SendChanForTest mengembalikan channel pesan keluar milik Client untuk pengujian unit.
func SendChanForTest(c *Client) <-chan []byte {
	return c.send
}

// RegisterForTest mendaftarkan klien ke hub untuk keperluan pengujian unit.
func (h *Hub) RegisterForTest(c *Client) {
	h.register(c)
}
