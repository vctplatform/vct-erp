// Package database provides database connection pool management for the VCT Platform.
// This is the canonical location for DB connectivity in the Hybrid Architecture.
//
// Individual modules access the database through this package rather than
// importing store/ directly. This allows modules to own their own queries
// while sharing a common connection pool.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // registers "pgx" driver for database/sql
)

// Pool wraps a *sql.DB with convenience methods and health check support.
type Pool struct {
	DB     *sql.DB
	logger *slog.Logger
}

// Config holds database connection configuration.
type Config struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns sensible default connection pool settings.
func DefaultConfig(url string) Config {
	return Config{
		URL:             url,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// NewPool creates a new database connection pool.
func NewPool(cfg Config, logger *slog.Logger) (*Pool, error) {
	if logger == nil {
		logger = slog.Default()
	}

	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("database.NewPool: sql.Open failed: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database.NewPool: ping failed: %w", err)
	}

	logger.Info("database pool connected",
		slog.String("driver", "pgx"),
		slog.Int("max_open", cfg.MaxOpenConns),
		slog.Int("max_idle", cfg.MaxIdleConns),
	)

	return &Pool{DB: db, logger: logger}, nil
}

// Ping checks the database connection health.
func (p *Pool) Ping(ctx context.Context) error {
	return p.DB.PingContext(ctx)
}

// Close closes the database connection pool.
func (p *Pool) Close() error {
	p.logger.Info("database pool closing")
	return p.DB.Close()
}

// Stats returns database connection pool statistics.
func (p *Pool) Stats() sql.DBStats {
	return p.DB.Stats()
}
