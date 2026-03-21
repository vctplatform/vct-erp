package config

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	CORSAllowedOrigins    []string
	DatabaseDSN           string
	DBMaxOpenConns        int
	DBMaxIdleConns        int
	DBConnMaxLifetime     time.Duration
	RedisAddr             string
	RedisUsername         string
	RedisPassword         string
	RedisDB               int
	AppRoleHeader         string
	AppActorHeader        string
	IdempotencyHeader     string
	RedisStreamKey        string
	DashboardCacheTTL     time.Duration
	FinanceEventsChannel  string
	AccountCacheTTL       time.Duration
	MongoURI              string
	MongoDatabase         string
	MongoAuditCollection  string
	OutboxRelayInterval   time.Duration
	OutboxRelayRetryDelay time.Duration
	OutboxRelayClaimTTL   time.Duration
	OutboxRelayBatchSize  int
	OutboxRelayRunTimeout time.Duration
}

// Load reads environment-driven configuration with safe defaults for local development.
func Load() (Config, error) {
	var cfg Config

	if err := loadDotEnv(".env", filepath.Join("backend", ".env")); err != nil {
		return Config{}, err
	}

	cfg.AppEnv = envOrDefault("APP_ENV", "development")
	cfg.HTTPAddr = envOrDefault("HTTP_ADDR", ":8080")
	cfg.CORSAllowedOrigins = parseCSVEnv("CORS_ALLOWED_ORIGINS", []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	})
	cfg.DatabaseDSN = envOrDefault("DB_DSN", "postgres://postgres:postgres@localhost:5432/vct_ledger?sslmode=disable")
	cfg.AppRoleHeader = envOrDefault("APP_ROLE_HEADER", "X-App-Role")
	cfg.AppActorHeader = envOrDefault("APP_ACTOR_HEADER", "X-Actor-ID")
	cfg.IdempotencyHeader = envOrDefault("IDEMPOTENCY_HEADER", "Idempotency-Key")
	cfg.RedisStreamKey = envOrDefault("REDIS_STREAM_KEY", "ledger.events")
	cfg.FinanceEventsChannel = envOrDefault("FINANCE_EVENTS_CHANNEL", "finance.dashboard.events")
	cfg.MongoURI = envOrDefault("MONGO_URI", "mongodb://localhost:27017")
	cfg.MongoDatabase = envOrDefault("MONGO_DATABASE", "vct_finance")
	cfg.MongoAuditCollection = envOrDefault("MONGO_AUDIT_COLLECTION", "ledger_void_audit")

	cfg.RedisAddr = "localhost:6379"
	cfg.RedisUsername = ""
	cfg.RedisPassword = ""
	cfg.RedisDB = 0
	if redisURL := strings.TrimSpace(envOrDefault("REDIS_URL", "")); redisURL != "" {
		if err := applyRedisURL(&cfg, redisURL); err != nil {
			return Config{}, err
		}
	}
	cfg.RedisAddr = envOrDefault("REDIS_ADDR", cfg.RedisAddr)
	cfg.RedisUsername = envOrDefault("REDIS_USERNAME", cfg.RedisUsername)
	cfg.RedisPassword = envOrDefault("REDIS_PASSWORD", cfg.RedisPassword)

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
	if cfg.DBConnMaxLifetime, err = parseDurationEnv("DB_CONN_MAX_LIFETIME", 30*time.Minute); err != nil {
		return Config{}, err
	}
	if cfg.AccountCacheTTL, err = parseDurationEnv("ACCOUNT_CACHE_TTL", 15*time.Minute); err != nil {
		return Config{}, err
	}
	if cfg.DashboardCacheTTL, err = parseDurationEnv("DASHBOARD_CACHE_TTL", 60*time.Second); err != nil {
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
	if cfg.OutboxRelayRunTimeout, err = parseDurationEnv("OUTBOX_RELAY_RUN_TIMEOUT", 20*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.RedisDB, err = parseIntEnv("REDIS_DB", cfg.RedisDB); err != nil {
		return Config{}, err
	}
	if cfg.DBMaxOpenConns, err = parseIntEnv("DB_MAX_OPEN_CONNS", 20); err != nil {
		return Config{}, err
	}
	if cfg.DBMaxIdleConns, err = parseIntEnv("DB_MAX_IDLE_CONNS", 10); err != nil {
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

func loadDotEnv(paths ...string) error {
	for _, candidate := range paths {
		if candidate == "" {
			continue
		}
		if err := loadDotEnvFile(candidate); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		return nil
	}
	return nil
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("parse %s: invalid line %q", path, line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return fmt.Errorf("parse %s: empty environment key", path)
		}
		if len(value) >= 2 {
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}
			if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
				value = strings.Trim(value, "'")
			}
		}

		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("set %s from %s: %w", key, path, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan %s: %w", path, err)
	}

	return nil
}

func applyRedisURL(cfg *Config, raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse REDIS_URL: %w", err)
	}
	if parsed.Host == "" {
		return fmt.Errorf("parse REDIS_URL: missing host")
	}

	cfg.RedisAddr = parsed.Host
	if parsed.User != nil {
		cfg.RedisUsername = parsed.User.Username()
		if password, ok := parsed.User.Password(); ok {
			cfg.RedisPassword = password
		}
	}

	trimmedDB := strings.TrimPrefix(parsed.Path, "/")
	if trimmedDB != "" {
		db, err := strconv.Atoi(trimmedDB)
		if err != nil {
			return fmt.Errorf("parse REDIS_URL database: %w", err)
		}
		cfg.RedisDB = db
	}

	return nil
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

func parseCSVEnv(key string, fallback []string) []string {
	raw := strings.TrimSpace(envOrDefault(key, ""))
	if raw == "" {
		return append([]string(nil), fallback...)
	}

	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	if len(values) == 0 {
		return append([]string(nil), fallback...)
	}
	return values
}
