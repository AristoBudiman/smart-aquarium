package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/api"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/auth"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/ingest"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/mqtt"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/notify"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/ws"
)

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	mqttURL := strings.TrimSpace(os.Getenv("MQTT_URL"))
	mqttUsername := strings.TrimSpace(os.Getenv("MQTT_USERNAME"))
	mqttPassword := os.Getenv("MQTT_PASSWORD")
	jwtSecret := os.Getenv("JWT_SECRET")
	adminUsername := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	telegramBotToken := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))

	if len(jwtSecret) < 32 {
		log.Fatal("konfigurasi tidak valid: JWT_SECRET wajib diisi minimal 32 karakter")
	}
	if databaseURL == "" {
		log.Fatal("konfigurasi tidak valid: DATABASE_URL wajib diisi")
	}
	if mqttURL == "" {
		mqttURL = "tcp://localhost:1883"
	}

	authSvc, err := auth.NewService(jwtSecret)
	if err != nil {
		log.Fatalf("gagal menginisialisasi layanan auth: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}
	defer store.Close()

	if err := api.BootstrapAdmin(ctx, store, adminUsername, adminPassword); err != nil {
		log.Fatalf("gagal menginisialisasi akun admin: %v", err)
	}

	hub := ws.NewHub()
	notifier := notify.New(store, telegramBotToken)

	mqttCfg := mqtt.Config{
		URL:      mqttURL,
		Username: mqttUsername,
		Password: mqttPassword,
		ClientID: "backend",
	}
	mqttClient := mqtt.NewWithConfig(ctx, mqttCfg, nil)

	ingester := ingest.New(store, mqttClient, hub, notifier)
	mqttClient.SetHandler(ingester.HandleMessage)

	if err := mqttClient.Connect(ctx); err != nil {
		log.Fatalf("gagal terhubung ke broker MQTT: %v", err)
	}
	defer mqttClient.Disconnect()

	go mqttClient.StartHeartbeat(ctx)

	apiServer := api.NewServer(store, authSvc, mqttClient, hub)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: apiServer,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Println("backend berjalan di :8080")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		log.Println("menerima sinyal berhenti, menutup server...")
	case err := <-errCh:
		if err != nil {
			log.Fatalf("server HTTP berhenti dengan error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("gagal menutup server HTTP dengan rapi: %v", err)
	}
}
