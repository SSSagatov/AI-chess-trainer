package postgres

import (
	"context"
	"strings"
	"testing"

	"github.com/SSSagatov/AI-chess-trainer/internal/config"
)

func TestConnectNilConfig(t *testing.T) {
	pool, err := Connect(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "database config is nil") {
		t.Fatalf("Connect(nil) error = %v, want nil config error", err)
	}
	if pool != nil {
		t.Fatal("Connect(nil) returned a non-nil pool")
	}
}

func TestConnectInvalidDSN(t *testing.T) {
	cfg := &config.DBConfig{
		Name: "chess", Username: "user", Password: "secret",
		Host: "localhost", Port: "invalid", SSLMode: "disable",
	}
	pool, err := Connect(context.Background(), cfg)
	if err == nil || !strings.Contains(err.Error(), "create PostgreSQL pool") {
		t.Fatalf("Connect(invalid DSN) error = %v, want pool creation error", err)
	}
	if pool != nil {
		t.Fatal("Connect(invalid DSN) returned a non-nil pool")
	}
}
