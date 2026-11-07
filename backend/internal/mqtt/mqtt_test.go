package mqtt

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	paho "github.com/eclipse/paho.mqtt.golang"
)

type fakeToken struct {
	done chan struct{}
	err  error
}

func newClosedToken(err error) *fakeToken {
	ch := make(chan struct{})
	close(ch)
	return &fakeToken{done: ch, err: err}
}

func (t *fakeToken) Wait() bool                       { <-t.done; return true }
func (t *fakeToken) WaitTimeout(d time.Duration) bool { <-t.done; return true }
func (t *fakeToken) Done() <-chan struct{}            { return t.done }
func (t *fakeToken) Error() error                     { return t.err }

type pubCall struct {
	topic    string
	qos      byte
	retained bool
	payload  []byte
}

type fakePahoClient struct {
	paho.Client
	mu       sync.Mutex
	pubs     []pubCall
	subs     map[string]byte
	callback paho.MessageHandler
}

func (f *fakePahoClient) Publish(topic string, qos byte, retained bool, payload interface{}) paho.Token {
	f.mu.Lock()
	defer f.mu.Unlock()
	var b []byte
	switch v := payload.(type) {
	case []byte:
		b = append([]byte(nil), v...)
	case string:
		b = []byte(v)
	}
	f.pubs = append(f.pubs, pubCall{
		topic:    topic,
		qos:      qos,
		retained: retained,
		payload:  b,
	})
	return newClosedToken(nil)
}

func (f *fakePahoClient) SubscribeMultiple(filters map[string]byte, callback paho.MessageHandler) paho.Token {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.subs = make(map[string]byte, len(filters))
	for k, v := range filters {
		f.subs[k] = v
	}
	f.callback = callback
	return newClosedToken(nil)
}

type fakeMessage struct {
	paho.Message
	topic   string
	payload []byte
}

func (m fakeMessage) Topic() string   { return m.topic }
func (m fakeMessage) Payload() []byte { return m.payload }

func TestSetModeAndSendActuators(t *testing.T) {
	fp := &fakePahoClient{}
	c := New(fp, nil)
	ctx := context.Background()

	if err := c.SetMode(ctx, "aq-1a2b3c", core.ModeAuto); err != nil {
		t.Fatalf("SetMode gagal: %v", err)
	}
	if err := c.SendActuators(ctx, "aq-1a2b3c", core.Actuators{Heater: 0, Fan: 1, Pump: 0}); err != nil {
		t.Fatalf("SendActuators gagal: %v", err)
	}

	fp.mu.Lock()
	defer fp.mu.Unlock()
	if len(fp.pubs) != 2 {
		t.Fatalf("jumlah publish = %d, ingin 2", len(fp.pubs))
	}

	// Cek SetMode: QoS 1, retained true, topik aquarium/aq-1a2b3c/config
	p0 := fp.pubs[0]
	if p0.topic != "aquarium/aq-1a2b3c/config" || p0.qos != 1 || !p0.retained {
		t.Fatalf("SetMode publish salah: %+v", p0)
	}
	if string(p0.payload) != `{"mode":"auto"}` {
		t.Fatalf("SetMode payload = %s, ingin {\"mode\":\"auto\"}", string(p0.payload))
	}

	// Cek SendActuators: QoS 1, retained false, topik aquarium/aq-1a2b3c/cmd
	p1 := fp.pubs[1]
	if p1.topic != "aquarium/aq-1a2b3c/cmd" || p1.qos != 1 || p1.retained {
		t.Fatalf("SendActuators publish salah: %+v", p1)
	}
	if string(p1.payload) != `{"heater":0,"fan":1,"pump":0}` {
		t.Fatalf("SendActuators payload = %s, ingin {\"heater\":0,\"fan\":1,\"pump\":0}", string(p1.payload))
	}
}

func TestSubscribeAndNonBlockingCallback(t *testing.T) {
	fp := &fakePahoClient{}
	diterima := make(chan string, 1)
	blokir := make(chan struct{})

	c := New(fp, func(_ context.Context, topic string, payload []byte) {
		diterima <- topic + ":" + string(payload)
		<-blokir
	})

	c.subscribeAll(fp)

	fp.mu.Lock()
	if len(fp.subs) != 3 {
		t.Fatalf("jumlah topik subscribe = %d, ingin 3", len(fp.subs))
	}
	cb := fp.callback
	fp.mu.Unlock()

	// Callback Paho harus segera kembali meskipun handler sedang memproses
	selesaiCallback := make(chan struct{})
	go func() {
		cb(fp, fakeMessage{topic: "aquarium/aq-test01/telemetry", payload: []byte(`{"temp":27.3,"ntu":12.4}`)})
		close(selesaiCallback)
	}()

	select {
	case <-selesaiCallback:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("callback Paho terblokir oleh handler")
	}

	select {
	case msg := <-diterima:
		if msg != `aquarium/aq-test01/telemetry:{"temp":27.3,"ntu":12.4}` {
			t.Fatalf("pesan diterima = %s", msg)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("handler tidak terpanggil")
	}
	close(blokir)
}

func TestStartHeartbeat(t *testing.T) {
	fp := &fakePahoClient{}
	c := New(fp, nil)
	c.heartbeatInterval = 15 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		c.StartHeartbeat(ctx)
		close(done)
	}()

	time.Sleep(40 * time.Millisecond)
	cancel()
	<-done

	fp.mu.Lock()
	defer fp.mu.Unlock()
	if len(fp.pubs) == 0 {
		t.Fatal("heartbeat tidak terkirim")
	}
	hb := fp.pubs[0]
	if hb.topic != "server/heartbeat" || hb.qos != 0 || hb.retained {
		t.Fatalf("heartbeat publish salah: %+v", hb)
	}
	var isi struct {
		TS int64 `json:"ts"`
	}
	if err := json.Unmarshal(hb.payload, &isi); err != nil || isi.TS <= 0 {
		t.Fatalf("payload heartbeat tidak valid: %s (%v)", string(hb.payload), err)
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("MQTT_URL", "tcp://mosquitto:1883")
	t.Setenv("MQTT_USERNAME", "backend")
	t.Setenv("MQTT_PASSWORD", "rahasia")

	cfg := ConfigFromEnv()
	if cfg.URL != "tcp://mosquitto:1883" || cfg.Username != "backend" || cfg.Password != "rahasia" {
		t.Fatalf("ConfigFromEnv = %+v", cfg)
	}
}

