package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config centralizes runtime settings for the core ledger service.
type Config struct {
	AppEnv                string
	HTTPAddr              string
	HTTPReadTimeout       time.Duration
	HTTPWriteTimeout      time.Duration
	HTTPIdleTimeout       time.Duration
	HTTPShutdownTimeout   time.Duration
	DatabaseDSN           string
	RedisAddr             string
	RedisUsername         string
	RedisPassword         string
	RedisDB               int
	AppRoleHeader         string
	IdempotencyHeader     string
	RedisStreamKey        string
	AccountCacheTTL       time.Duration
	OutboxRelayInterval   time.Duration
	OutboxRelayRetryDelay time.Duration
	OutboxRelayClaimTTL   time.Duration
	OutboxRelayBatchSize  int
}

// Load reads environment-driven configuration with safe defaults for local development.
func Load() (Config, error) {
	var cfg Config

	cfg.AppEnv = envOrDefault("APP_ENV", "development")
	cfg.HTTPAddr = envOrDefault("HTTP_ADDR", ":8080")
	cfg.DatabaseDSN = envOrDefault("DB_DSN", "postgres://postgres:postgres@localhost:5432/vct_ledger?sslmode=disable")
	cfg.RedisAddr = envOrDefault("REDIS_ADDR", "localhost:6379")
	cfg.RedisUsername = envOrDefault("REDIS_USERNAME", "")
	cfg.RedisPassword = envOrDefault("REDIS_PASSWORD", "")
	cfg.AppRoleHeader = envOrDefault("APP_ROLE_HEADER", "X-App-Role")
	cfg.IdempotencyHeader = envOrDefault("IDEMPOTENCY_HEADER", "Idempotency-Key")
	cfg.RedisStreamKey = envOrDefault("REDIS_STREAM_KEY", "ledger.events")

	var err error
	if cfg.HTTPReadTimeout, err = parseDurationEnv("HTTP_READ_TIMEOUT", 15*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.HTTPWriteTimeout, err = parseDurationEnv("HTTP_WRITE_TIMEOUT", 15*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.HTTPIdleTimeout, err = parseDurationEnv("HTTP_IDLE_TIMEOUT", time.Minute); err != nil {
		return Config{}, err
	}
	if cfg.HTTPShutdownTimeout, err = parseDurationEnv("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.AccountCacheTTL, err = parseDurationEnv("ACCOUNT_CACHE_TTL", 15*time.Minute); err != nil {
		return Config{}, err
	}
	if cfg.OutboxRelayInterval, err = parseDurationEnv("OUTBOX_RELAY_INTERVAL", 5*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.OutboxRelayRetryDelay, err = parseDurationEnv("OUTBOX_RELAY_RETRY_DELAY", 30*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.OutboxRelayClaimTTL, err = parseDurationEnv("OUTBOX_RELAY_CLAIM_TTL", 2*time.Minute); err != nil {
		return Config{}, err
	}
	if cfg.RedisDB, err = parseIntEnv("REDIS_DB", 0); err != nil {
		return Config{}, err
	}
	if cfg.OutboxRelayBatchSize, err = parseIntEnv("OUTBOX_RELAY_BATCH_SIZE", 100); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func envOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func parseDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	raw := envOrDefault(key, fallback.String())
	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return value, nil
}

func parseIntEnv(key string, fallback int) (int, error) {
	raw := envOrDefault(key, strconv.Itoa(fallback))
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return value, nil
}
