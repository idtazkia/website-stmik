package model

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

// Connect establishes a connection pool to the database
func Connect(ctx context.Context, dsn string) error {
	var err error
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection pool
func Close() {
	if pool != nil {
		pool.Close()
	}
}

// Pool returns the database connection pool
func Pool() *pgxpool.Pool {
	return pool
}
