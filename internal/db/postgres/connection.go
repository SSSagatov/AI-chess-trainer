package postgres

import (
	"context"
	"fmt"

	"github.com/SSSagatov/AI-chess-trainer/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a connection pool and checks that PostgreSQL is reachable.
// The caller must close the returned pool when it is no longer needed.
func Connect(ctx context.Context, cfg *config.DBConfig) (*pgxpool.Pool, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config is nil")
	}

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	return pool, nil
}
