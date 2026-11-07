package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
)

// SEMENTARA: implementasi Store palsu untuk unit test Peran 3.
type fakeNotifyStore struct {
	core.Store
	chatIDs []string
	marked  chan int64
}

func (f *fakeNotifyStore) TelegramChatIDs(_ context.Context) ([]string, error) {
	return f.chatIDs, nil
}

func (f *fakeNotifyStore) MarkAlertSent(_ context.Context, id int64) error {
	if f.marked != nil {
		f.marked <- id
	}
	return nil
}

func TestFormatMessage(t *testing.T) {
	w := FormatMessage(core.Alert{
		Level:   core.ClassWarning,
		Message: "Akuarium 1: suhu 33,1 °C (WARNING), kekeruhan 18 NTU (NORMAL). Kipas dinyalakan otomatis.",
	})
	wantW := "⚠️ WARNING – Akuarium 1: suhu 33,1 °C (WARNING), kekeruhan 18 NTU (NORMAL). Kipas dinyalakan otomatis."
	if w != wantW {
		t.Fatalf("FormatMessage WARNING = %q, ingin %q", w, wantW)
	}

	s := FormatMessage(core.Alert{
		Level:   core.ClassSevere,
		Message: "Akuarium 1: suhu 36,0 °C (SEVERE), kekeruhan 250 NTU (SEVERE).",
	})
	if !strings.HasPrefix(s, "🚨 DARURAT") {
		t.Fatalf("FormatMessage SEVERE = %q, ingin awalan 🚨 DARURAT", s)
	}
}

func TestNotifySendMessageAndMarkAlertSent(t *testing.T) {
	var (
		mu       sync.Mutex
		requests []sendMessageReq
		reqPaths []string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req sendMessageReq
		_ = json.NewDecoder(r.Body).Decode(&req)
		mu.Lock()
		requests = append(requests, req)
		reqPaths = append(reqPaths, r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	store := &fakeNotifyStore{
		chatIDs: []string{"123456", "789012"},
		marked:  make(chan int64, 1),
	}

	svc := New(store, "TOKEN_RAHASIA_UJI")
	svc.baseURL = srv.URL

	svc.Notify(context.Background(), core.Alert{
		ID:      42,
		NodeID:  "aq-1a2b3c",
		Level:   core.ClassWarning,
		Message: "Akuarium 1: suhu 33,1 °C (WARNING), kekeruhan 18 NTU (NORMAL). Kipas dinyalakan otomatis.",
	})

	select {
	case id := <-store.marked:
		if id != 42 {
			t.Fatalf("MarkAlertSent id = %d, ingin 42", id)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("MarkAlertSent tidak dipanggil tepat waktu")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(requests) != 2 {
		t.Fatalf("jumlah request sendMessage = %d, ingin 2", len(requests))
	}
	if reqPaths[0] != "/botTOKEN_RAHASIA_UJI/sendMessage" {
		t.Fatalf("path = %s, ingin /botTOKEN_RAHASIA_UJI/sendMessage", reqPaths[0])
	}
	if !strings.HasPrefix(requests[0].Text, "⚠️ WARNING – ") {
		t.Fatalf("teks pesan = %q", requests[0].Text)
	}
}

func TestNotifyEmptyTokenDoesNotLeakOrSend(t *testing.T) {
	store := &fakeNotifyStore{
		chatIDs: []string{"123456"},
		marked:  make(chan int64, 1),
	}
	svc := New(store, "")

	var buf bytes.Buffer
	origOut := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(origOut)

	svc.Notify(context.Background(), core.Alert{
		ID:      7,
		NodeID:  "aq-1a2b3c",
		Level:   core.ClassSevere,
		Message: "Akuarium 1: suhu 37,0 °C (SEVERE).",
	})

	select {
	case id := <-store.marked:
		t.Fatalf("MarkAlertSent seharusnya tidak dipanggil saat token kosong, dapat id=%d", id)
	case <-time.After(50 * time.Millisecond):
	}

	if !strings.Contains(buf.String(), "🚨 DARURAT") {
		t.Fatalf("log tidak mencatat pesan darurat: %q", buf.String())
	}
}

func TestRedactTokenOnError(t *testing.T) {
	token := "SECRET_BOT_TOKEN_999"
	svc := New(nil, token)
	terfilter := svc.redactToken("Post \"https://api.telegram.org/bot" + token + "/sendMessage\": dial tcp: connection refused")
	if strings.Contains(terfilter, token) {
		t.Fatalf("token bocor pada pesan error: %s", terfilter)
	}
}

