// Package httputil provides HTTP helpers for the VCT Platform.
// This is the canonical location for HTTP utilities in the Hybrid Architecture.
//
// These helpers can be used by any module's handler without depending on
// the httpapi package, breaking the circular dependency chain.
package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"vct-platform/backend/internal/apierror"
	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/shared/auth"
)

// ═══════════════════════════════════════════════════════════════
// JSON Decode Helpers
// ═══════════════════════════════════════════════════════════════

// DecodeJSON decodes the request body into the given output struct.
func DecodeJSON(r *http.Request, output any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(output); err != nil {
		return fmt.Errorf("json không hợp lệ: %w", err)
	}
	return nil
}

// DecodeObject decodes the request body into a map[string]any.
func DecodeObject(r *http.Request) (map[string]any, error) {
	var payload map[string]any
	if err := DecodeJSON(r, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// ═══════════════════════════════════════════════════════════════
// Response Helpers — Unified API Envelope
// ═══════════════════════════════════════════════════════════════

// JSON writes a successful JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// JSONBytes writes raw JSON bytes as a response.
func JSONBytes(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}

// ── Standard Envelope ────────────────────────────────────────

// Envelope is the standard API response wrapper.
type Envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
}

// ErrorBody is the error section of the API envelope.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Success writes a success envelope response.
func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, Envelope{Success: true, Data: data})
}

// SuccessWithMeta writes a success envelope response with metadata.
func SuccessWithMeta(w http.ResponseWriter, status int, data, meta any) {
	JSON(w, status, Envelope{Success: true, Data: data, Meta: meta})
}

// Error writes an error envelope response.
func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, Envelope{
		Success: false,
		Error:   &ErrorBody{Code: code, Message: message},
	})
}

// ErrorWithDetails writes an error response with additional details.
func ErrorWithDetails(w http.ResponseWriter, status int, code, message string, details any) {
	JSON(w, status, Envelope{
		Success: false,
		Error:   &ErrorBody{Code: code, Message: message, Details: details},
	})
}

// InternalError writes a generic 500 error (hides internal details from client).
func InternalError(w http.ResponseWriter, err error) {
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Lỗi hệ thống, vui lòng thử lại sau")
}

// ═══════════════════════════════════════════════════════════════
// WriteError — Auto-maps apierror.Error to HTTP status
// ═══════════════════════════════════════════════════════════════

// WriteError inspects the error type and writes the appropriate HTTP response.
// Supports both structured apierror.Error and plain errors.
func WriteError(w http.ResponseWriter, err error) {
	var apiErr *apierror.Error
	if errors.As(err, &apiErr) {
		status := mapErrorCodeToStatus(apiErr.Code)
		Error(w, status, apiErr.Code, apiErr.Message)
		return
	}
	// Fallback for non-structured errors
	InternalError(w, err)
}

// mapErrorCodeToStatus maps error code prefixes to HTTP status codes.
func mapErrorCodeToStatus(code string) int {
	switch {
	case strings.HasPrefix(code, "AUTH_401"), strings.HasPrefix(code, "AUTH_ERR"):
		return http.StatusUnauthorized
	case strings.HasPrefix(code, "AUTH_403"), strings.HasPrefix(code, "FORBIDDEN"):
		return http.StatusForbidden
	case strings.HasPrefix(code, "STORE_404"), strings.HasPrefix(code, "TENANT_404"), strings.HasPrefix(code, "TEMPLATE_404"):
		return http.StatusNotFound
	case strings.HasPrefix(code, "STORE_409"), strings.HasPrefix(code, "BIZ_409"), strings.HasPrefix(code, "AUTH_409"):
		return http.StatusConflict
	case strings.HasPrefix(code, "VAL_400"), strings.HasPrefix(code, "STORE_400"), strings.HasPrefix(code, "BIZ_400"), strings.HasPrefix(code, "WORKER_ERR"), strings.HasPrefix(code, "AUTH_400"):
		return http.StatusBadRequest
	case strings.HasPrefix(code, "BIZ_429"), strings.HasPrefix(code, "AUTH_429"):
		return http.StatusTooManyRequests
	case strings.HasPrefix(code, "WEBHOOK_ERR"):
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// ═══════════════════════════════════════════════════════════════
// Module Registration Interface
// ═══════════════════════════════════════════════════════════════

// Module is the interface that domain modules implement to self-register
// their HTTP routes. This replaces the centralized router.go approach.
type Module interface {
	// RegisterRoutes registers the module's HTTP routes on the given mux.
	RegisterRoutes(mux *http.ServeMux)
}

// ═══════════════════════════════════════════════════════════════
// Real-time Event Broadcasting
// ═══════════════════════════════════════════════════════════════

// EventBroadcaster is the interface for triggering real-time entity change events.
// Modules use this to notify other parts of the system (and connected clients)
// of data changes without direct coupling to the core event bus or WebSocket server.
type EventBroadcaster interface {
	BroadcastEntityChange(entity, action, id string, data map[string]any, metadata map[string]any)
}

// DummyBroadcaster is a no-op implementation for testing or simple modules.
type DummyBroadcaster struct{}

func (d DummyBroadcaster) BroadcastEntityChange(entity, action, id string, data map[string]any, metadata map[string]any) {
}

// ToMap converts an object to map[string]any for broadcasting.
func ToMap(obj any) map[string]any {
	if obj == nil {
		return nil
	}
	b, _ := json.Marshal(obj)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	return m
}

// ═══════════════════════════════════════════════════════════════
// Auth & Context Helpers
// ═══════════════════════════════════════════════════════════════

type contextKey string

const principalKey contextKey = "principal"

// WithPrincipal returns a copy of r with the principal in its context.
func WithPrincipal(r *http.Request, p auth.Principal) *http.Request {
	ctx := context.WithValue(r.Context(), principalKey, p)
	return r.WithContext(ctx)
}

// GetPrincipal retrieves the principal from the request context.
func GetPrincipal(r *http.Request) (auth.Principal, bool) {
	p, ok := r.Context().Value(principalKey).(auth.Principal)
	return p, ok
}

// ═══════════════════════════════════════════════════════════════
// Middleware
// ═══════════════════════════════════════════════════════════════

// Middleware is a standard function signature for HTTP middleware.
type Middleware func(http.Handler) http.Handler

// RequirePermission returns a middleware that checks if the current user
// has the required permission for the given entity and action.
//
// It assumes the Principal is already in the context (placed by a top-level auth middleware).
func RequirePermission(entity string, action authz.EntityAction) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := GetPrincipal(r)
			if !ok {
				Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu đăng nhập")
				return
			}

			if authz.CanEntityAction(p.User.Role, entity, action) {
				next.ServeHTTP(w, r)
				return
			}

			// Forbidden if role doesn't have permission
			Error(w, http.StatusForbidden, "AUTH_403", fmt.Sprintf("Vai trò %s không có quyền %s trên %s", p.User.Role, action, entity))
		})
	}
}

// Chain applies a list of middleware to an http.Handler.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
