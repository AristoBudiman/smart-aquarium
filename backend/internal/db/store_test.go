package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/db"
)

func TestStoreImplementsCoreStore(t *testing.T) {
	var _ core.Store = (*db.Store)(nil)
}

func TestConnectInvalidURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := db.Connect(ctx, "invalid-url://:bad-format")
	if err == nil {
		t.Fatal("diharapkan error saat format database URL tidak valid, tetapi tidak ada error")
	}
}

func TestConnectCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // langsung batalkan

	_, err := db.Connect(ctx, "postgres://user:pass@localhost:5432/dbname")
	if err == nil {
		t.Fatal("diharapkan error saat context dibatalkan, tetapi tidak ada error")
	}
}

func TestNewStore(t *testing.T) {
	store := db.New(nil)
	if store == nil {
		t.Fatal("diharapkan store tidak nil")
	}
	if store.Pool() != nil {
		t.Fatal("diharapkan pool bernilai nil")
	}
	store.Close()
}
