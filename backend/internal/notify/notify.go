package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
)

const defaultBaseURL = "https://api.telegram.org"

// Service mengelola pengiriman notifikasi alert ke Telegram.
type Service struct {
	store   core.Store
	token   string
	baseURL string
	client  *http.Client
}

// NewFromEnv membuat Service notifikasi dengan membaca TELEGRAM_BOT_TOKEN dari environment.
func NewFromEnv(store core.Store) *Service {
	return New(store, os.Getenv("TELEGRAM_BOT_TOKEN"))
}

// New membuat Service notifikasi dengan token tertentu.
func New(store core.Store, token string) *Service {
	return &Service{
		store:   store,
		token:   strings.TrimSpace(token),
		baseURL: defaultBaseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FormatMessage menyusun teks pesan Telegram sesuai level keparahan kontrak.
func FormatMessage(alert core.Alert) string {
	msg := strings.TrimSpace(alert.Message)
	switch alert.Level {
	case core.ClassSevere:
		if strings.HasPrefix(msg, "🚨") {
			return msg
		}
		return "🚨 DARURAT – " + msg
	case core.ClassWarning:
		if strings.HasPrefix(msg, "⚠️") {
			return msg
		}
		return "⚠️ WARNING – " + msg
	default:
		if strings.HasPrefix(msg, "⚠️") || strings.HasPrefix(msg, "🚨") {
			return msg
		}
		return fmt.Sprintf("⚠️ %s – %s", alert.Level, msg)
	}
}

// Notify mengirim pesan alert ke Telegram di goroutine terpisah.
// Jika TELEGRAM_BOT_TOKEN kosong, hanya mencatat ke log dan langsung kembali.
func (s *Service) Notify(_ context.Context, alert core.Alert) {
	text := FormatMessage(alert)
	if s.token == "" {
		log.Printf("notify: TELEGRAM_BOT_TOKEN kosong, notifikasi dicatat di log: %s", text)
		return
	}

	go s.sendAsync(alert, text)
}

type sendMessageReq struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (s *Service) sendAsync(alert core.Alert, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if s.store == nil {
		return
	}

	chatIDs, err := s.store.TelegramChatIDs(ctx)
	if err != nil {
		log.Printf("notify: gagal ambil daftar chat_id: %v", err)
		return
	}

	berhasil := false
	for _, rawID := range chatIDs {
		chatID := strings.TrimSpace(rawID)
		if chatID == "" {
			continue
		}
		if err := s.sendOne(ctx, chatID, text); err != nil {
			log.Printf("notify: gagal kirim ke chat_id %s: %s", chatID, s.redactToken(err.Error()))
			continue
		}
		berhasil = true
	}

	if berhasil && alert.ID > 0 {
		if err := s.store.MarkAlertSent(ctx, alert.ID); err != nil {
			log.Printf("notify: gagal MarkAlertSent id=%d: %v", alert.ID, err)
		}
	}
}

func (s *Service) sendOne(ctx context.Context, chatID, text string) error {
	body, err := json.Marshal(sendMessageReq{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/bot%s/sendMessage", strings.TrimRight(s.baseURL, "/"), s.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("buat request gagal")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status HTTP %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) redactToken(msg string) string {
	if s.token == "" {
		return msg
	}
	return strings.ReplaceAll(msg, s.token, "[REDACTED]")
}
