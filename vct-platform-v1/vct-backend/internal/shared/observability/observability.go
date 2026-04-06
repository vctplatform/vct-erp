// Package observability provides unified logging, metrics, and tracing
// for the VCT Platform. This is the canonical location for observability
// infrastructure in the Hybrid Architecture.
//
// It re-exports logging primitives and provides a unified entry point
// for setting up the full observability stack.
package observability

import (
	"log/slog"

	"vct-platform/backend/internal/logging"
)

// ── Logging — Re-exports from internal/logging ──────────────

// LogConfig holds logger configuration.
type LogConfig = logging.Config

// DefaultLogConfig returns production-ready defaults.
var DefaultLogConfig = logging.DefaultConfig

// NewLogger creates a configured structured logger.
var NewLogger = logging.NewLogger

// NewLoggerFromEnv creates a logger using environment-based defaults.
var NewLoggerFromEnv = logging.New

// WithContext stores a logger in the context.
var WithContext = logging.WithContext

// FromContext retrieves the logger from context, or returns the default.
var FromContext = logging.FromContext

// ── Component Logger Factory ─────────────────────────────────

// Logger creates a child logger tagged with a component name.
// Use this in domain modules to identify log source.
//
//	log := observability.Logger(baseLogger, "scoring")
//	log.Info("match started", slog.String("match_id", id))
func Logger(base *slog.Logger, component string) *slog.Logger {
	if base == nil {
		base = slog.Default()
	}
	return base.With(slog.String("component", component))
}
