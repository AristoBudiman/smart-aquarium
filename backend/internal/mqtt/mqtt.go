package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	paho "github.com/eclipse/paho.mqtt.golang"
)

const (
	TopicSubTelemetry = "aquarium/+/telemetry"
	TopicSubState     = "aquarium/+/state"
	TopicSubStatus    = "aquarium/+/status"
	TopicHeartbeat    = "server/heartbeat"

	HeartbeatInterval = 30 * time.Second
)

// MessageHandler memproses pesan MQTT yang masuk dari topik langganan.
type MessageHandler func(ctx context.Context, topic string, payload []byte)

// Config menyimpan konfigurasi koneksi broker MQTT.
type Config struct {
	URL      string
	Username string
	Password string
	ClientID string
}

// ConfigFromEnv membaca MQTT_URL, MQTT_USERNAME, dan MQTT_PASSWORD dari environment.
func ConfigFromEnv() Config {
	url := os.Getenv("MQTT_URL")
	if url == "" {
		url = "tcp://localhost:1883"
	}
	return Config{
		URL:      url,
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
		ClientID: "backend",
	}
}

// Client mengelola koneksi MQTT, implementasi core.Commander, langganan topik, dan heartbeat.
type Client struct {
	raw               paho.Client
	mu                sync.RWMutex
	ctx               context.Context
	handler           MessageHandler
	heartbeatInterval time.Duration
}

var _ core.Commander = (*Client)(nil)

// New membuat instance Client dari paho.Client yang sudah ada (berguna untuk pengujian).
func New(raw paho.Client, handler MessageHandler) *Client {
	return &Client{
		raw:               raw,
		ctx:               context.Background(),
		handler:           handler,
		heartbeatInterval: HeartbeatInterval,
	}
}

// NewFromEnv membuat dan mengonfigurasi Client MQTT dari variabel lingkungan.
// Auto-reconnect diaktifkan dan langganan topik didaftarkan ulang pada OnConnect.
func NewFromEnv(ctx context.Context, handler MessageHandler) *Client {
	return NewWithConfig(ctx, ConfigFromEnv(), handler)
}

// NewWithConfig membuat Client MQTT menggunakan Config yang diberikan.
func NewWithConfig(ctx context.Context, cfg Config, handler MessageHandler) *Client {
	if ctx == nil {
		ctx = context.Background()
	}
	c := &Client{
		ctx:               ctx,
		handler:           handler,
		heartbeatInterval: HeartbeatInterval,
	}

	opts := paho.NewClientOptions().
		AddBroker(cfg.URL).
		SetClientID(cfg.ClientID).
		SetAutoReconnect(true).
		SetConnectRetry(true)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}

	opts.SetOnConnectHandler(func(pc paho.Client) {
		c.subscribeAll(pc)
	})

	c.raw = paho.NewClient(opts)
	return c
}

// SetHandler memperbarui handler pesan masuk.
func (c *Client) SetHandler(handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handler = handler
}

// Connect membuka koneksi ke broker MQTT.
func (c *Client) Connect(ctx context.Context) error {
	tok := c.raw.Connect()
	return waitToken(ctx, tok)
}

// Disconnect menutup koneksi ke broker MQTT.
func (c *Client) Disconnect() {
	c.raw.Disconnect(250)
}

func (c *Client) subscribeAll(pc paho.Client) {
	filters := map[string]byte{
		TopicSubTelemetry: 0,
		TopicSubState:     0,
		TopicSubStatus:    1,
	}
	tok := pc.SubscribeMultiple(filters, c.onMessage)
	go func() {
		if tok.Wait() && tok.Error() != nil {
			log.Printf("mqtt: gagal subscribe topik: %v", tok.Error())
		}
	}()
}

// onMessage dipanggil oleh Paho saat pesan tiba. Tidak memblokir callback Paho.
func (c *Client) onMessage(_ paho.Client, msg paho.Message) {
	topic := msg.Topic()
	payload := append([]byte(nil), msg.Payload()...)

	c.mu.RLock()
	h := c.handler
	ctx := c.ctx
	c.mu.RUnlock()

	if h == nil {
		return
	}
	go h(ctx, topic, payload)
}

type modePayload struct {
	Mode core.Mode `json:"mode"`
}

// SetMode mengirim konfigurasi mode ke topik aquarium/{node_id}/config dengan QoS 1 dan retained true.
func (c *Client) SetMode(ctx context.Context, nodeID string, mode core.Mode) error {
	data, err := json.Marshal(modePayload{Mode: mode})
	if err != nil {
		return err
	}
	topic := fmt.Sprintf("aquarium/%s/config", nodeID)
	tok := c.raw.Publish(topic, 1, true, data)
	return waitToken(ctx, tok)
}

// SendActuators mengirim perintah aktuator ke topik aquarium/{node_id}/cmd dengan QoS 1 dan retained false.
func (c *Client) SendActuators(ctx context.Context, nodeID string, a core.Actuators) error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	topic := fmt.Sprintf("aquarium/%s/cmd", nodeID)
	tok := c.raw.Publish(topic, 1, false, data)
	return waitToken(ctx, tok)
}

type heartbeatPayload struct {
	TS int64 `json:"ts"`
}

// PublishHeartbeat mengirim satu pesan heartbeat ke topik server/heartbeat dengan QoS 0 dan retained false.
func (c *Client) PublishHeartbeat(ctx context.Context) error {
	data, err := json.Marshal(heartbeatPayload{TS: time.Now().Unix()})
	if err != nil {
		return err
	}
	tok := c.raw.Publish(TopicHeartbeat, 0, false, data)
	return waitToken(ctx, tok)
}

// StartHeartbeat mengirim heartbeat setiap 30 detik sampai context selesai.
func (c *Client) StartHeartbeat(ctx context.Context) {
	interval := c.heartbeatInterval
	if interval <= 0 {
		interval = HeartbeatInterval
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.PublishHeartbeat(ctx); err != nil && ctx.Err() == nil {
				log.Printf("mqtt: gagal kirim heartbeat: %v", err)
			}
		}
	}
}

func waitToken(ctx context.Context, tok paho.Token) error {
	if ctx == nil {
		tok.Wait()
		return tok.Error()
	}
	select {
	case <-tok.Done():
		return tok.Error()
	case <-ctx.Done():
		return ctx.Err()
	}
}
