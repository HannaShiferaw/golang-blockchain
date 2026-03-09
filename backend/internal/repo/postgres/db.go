package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func ConnectFromEnv(ctx context.Context) (*DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return nil, fmt.Errorf("POSTGRES_HOST not set")
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	db := os.Getenv("POSTGRES_DB")
	if db == "" {
		db = "coffee_demo"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "coffee"
	}
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = "coffee"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 4
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db != nil && db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *DB) EnsureSchema(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS audit_tx (
			tx_id TEXT PRIMARY KEY,
			tx_type TEXT NOT NULL,
			actor_id TEXT NOT NULL,
			actor_role TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			tx_hash TEXT NOT NULL,
			block_index INT NOT NULL,
			block_hash TEXT NOT NULL,
			payload JSONB NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS order_index (
			order_id TEXT PRIMARY KEY,
			exporter_id TEXT NOT NULL,
			buyer_id TEXT NOT NULL,
			status TEXT NOT NULL,
			total_usd NUMERIC NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		);`,
	}
	for _, s := range stmts {
		if _, err := db.Pool.Exec(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

