// Package requestid provides request ID generation, context propagation,
// and HTTP middleware for end-to-end request correlation tracking.
package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

// ═══════════════════════════════════════════════════════════════
// Headers
// ═══════════════════════════════════════════════════════════════

const (
	// HeaderRequestID is the primary request identifier.
	HeaderRequestID = "X-Request-ID"
	// HeaderCorrelationID is propagated across service boundaries.
	HeaderCorrelationID = "X-Correlation-ID"
)

// ═══════════════════════════════════════════════════════════════
// ID Generation
// ═══════════════════════════════════════════════════════════════

// NewID generates a 16-character hex request ID (8 random bytes).
func NewID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// NewLongID generates a 32-character hex ID (16 random bytes, UUID-length).
func NewLongID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ═══════════════════════════════════════════════════════════════
// Context Propagation
// ═══════════════════════════════════════════════════════════════

type ctxKey string

const (
	requestIDKey     ctxKey = "request_id"
	correlationIDKey ctxKey = "correlation_id"
)

// WithRequestID stores a request ID in context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID extracts the request ID from context.
func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// WithCorrelationID stores a correlation ID in context.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

// GetCorrelationID extracts the correlation ID from context.
func GetCorrelationID(ctx context.Context) string {
	if v, ok := ctx.Value(correlationIDKey).(string); ok {
		return v
	}
	return ""
}

// ═══════════════════════════════════════════════════════════════
// HTTP Middleware
// ═══════════════════════════════════════════════════════════════

// Middleware injects Request-ID and Correlation-ID into the request
// context and response headers. If a client sends X-Request-ID or
// X-Correlation-ID, those values are reused (for cross-service tracing).
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request ID: use incoming or generate new
		reqID := sanitize(r.Header.Get(HeaderRequestID))
		if reqID == "" {
			reqID = NewID()
		}

		// Correlation ID: use incoming or generate new
		corrID := sanitize(r.Header.Get(HeaderCorrelationID))
		if corrID == "" {
			corrID = NewLongID()
		}

		// Store in context
		ctx := WithRequestID(r.Context(), reqID)
		ctx = WithCorrelationID(ctx, corrID)

		// Echo in response headers
		w.Header().Set(HeaderRequestID, reqID)
		w.Header().Set(HeaderCorrelationID, corrID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// InjectOutgoing adds request/correlation IDs from context to outgoing HTTP requests.
// Use this when calling downstream services.
func InjectOutgoing(ctx context.Context, req *http.Request) {
	if id := GetRequestID(ctx); id != "" {
		req.Header.Set(HeaderRequestID, id)
	}
	if id := GetCorrelationID(ctx); id != "" {
		req.Header.Set(HeaderCorrelationID, id)
	}
}

// sanitize strips spaces and limits length to prevent header injection.
func sanitize(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 64 {
		return s[:64]
	}
	return s
}
