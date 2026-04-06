package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"vct-platform/backend/internal/auth"
	"vct-platform/backend/internal/authz"
)

// ── Request ID ───────────────────────────────────────────────

type contextKey string

const requestIDKey contextKey = "requestID"

func generateRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// withRequestID injects a unique request ID into the context and response header.
func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = generateRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// ── Authentication Middleware ─────────────────────────────────

func (s *Server) withAuth(next func(http.ResponseWriter, *http.Request, auth.Principal)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		principal, err := s.principalFromRequest(r)
		if err != nil {
			writeAuthError(w, err)
			return
		}
		next(w, r, principal)
	}
}

func (s *Server) principalFromRequest(r *http.Request) (auth.Principal, error) {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		return auth.Principal{}, fmt.Errorf("%w: thiếu bearer token", auth.ErrUnauthorized)
	}
	token := strings.TrimSpace(authorization[7:])
	if token == "" {
		return auth.Principal{}, fmt.Errorf("%w: token trống", auth.ErrUnauthorized)
	}
	return s.authService.AuthenticateAccessToken(token, requestContextFromRequest(r))
}

func (s *Server) authorizeEntityAction(
	principal *auth.Principal,
	entity string,
	action authz.EntityAction,
) error {
	if s.cfg.DisableAuthForData {
		return nil
	}
	if principal == nil {
		return fmt.Errorf("%w: thiếu thông tin phiên làm việc", auth.ErrUnauthorized)
	}
	// Entity routes honor the active role embedded in the current token so
	// context switching remains an explicit least-privilege operation.
	if authz.CanEntityAction(principal.User.Role, entity, action) {
		return nil
	}
	return fmt.Errorf(
		"%w: vai trò %s không có quyền %s trên %s",
		auth.ErrForbidden,
		principal.User.Role,
		action,
		entity,
	)
}

// ── CORS Middleware ───────────────────────────────────────────

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			if s.isAllowedOrigin(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token, X-Request-ID")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-RateLimit-Remaining")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) isAllowedOrigin(origin string) bool {
	// In production, wildcard origins are not allowed
	if _, ok := s.allowedOrigins["*"]; ok {
		env := strings.ToLower(s.cfg.Environment)
		if env == "production" || env == "staging" {
			s.logger.Warn("wildcard CORS origin rejected", slog.String("env", env))
			return false
		}
		return true
	}
	// Case-insensitive origin matching
	normalised := strings.ToLower(origin)
	for allowed := range s.allowedOrigins {
		if strings.ToLower(allowed) == normalised {
			return true
		}
	}
	return false
}

// ── CSRF Protection (Double Submit Cookie) ───────────────────

func (s *Server) withCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF for safe methods and WebSocket upgrade
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			next.ServeHTTP(w, r)
			return
		}

		// Skip CSRF for public auth endpoints — these are stateless (no session
		// cookies), already rate-limited, and must work from any origin.
		publicAuthPaths := []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/auth/send-otp",
			"/api/v1/auth/verify-otp",
			"/api/v1/auth/refresh",
		}
		for _, p := range publicAuthPaths {
			if r.URL.Path == p {
				next.ServeHTTP(w, r)
				return
			}
		}

		// For state-changing requests, validate Origin header against allowed origins
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin == "" {
			// Fall back to Referer if Origin is missing (some browsers)
			referer := strings.TrimSpace(r.Header.Get("Referer"))
			if referer != "" {
				// Extract origin from referer URL
				if idx := strings.Index(referer, "://"); idx >= 0 {
					rest := referer[idx+3:]
					if slashIdx := strings.Index(rest, "/"); slashIdx >= 0 {
						origin = referer[:idx+3+slashIdx]
					} else {
						origin = referer
					}
				}
			}
		}

		// Reject state-changing requests with no origin info (prevent CSRF)
		if origin == "" {
			// Allow requests with valid Authorization bearer tokens (API clients)
			if strings.HasPrefix(strings.ToLower(r.Header.Get("Authorization")), "bearer ") {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"missing origin header for state-changing request"}`))
			return
		}

		if !s.isAllowedOrigin(origin) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"origin not allowed"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ── Security Headers Middleware ──────────────────────────────

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

// ── Structured Logging Middleware ─────────────────────────────

type accessLogEntry struct {
	Timestamp string `json:"ts"`
	RequestID string `json:"rid,omitempty"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	IP        string `json:"ip"`
	UserAgent string `json:"ua,omitempty"`
}

func (s *Server) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		rw := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		latency := time.Since(started)

		entry := accessLogEntry{
			Timestamp: started.UTC().Format(time.RFC3339),
			RequestID: getRequestID(r),
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    rw.statusCode,
			LatencyMs: latency.Milliseconds(),
			IP:        extractClientIP(r),
			UserAgent: r.UserAgent(),
		}

		logJSON, err := json.Marshal(entry)
		if err != nil {
			s.logger.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.Duration("latency", latency),
			)
			return
		}
		s.logger.Info(string(logJSON))
	})
}

// ── Panic Recovery Middleware ─────────────────────────────────

func withRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())
				slog.Default().Error("panic recovered",
					slog.Any("error", rec),
					slog.String("stack", strings.ReplaceAll(stack, "\n", "\\n")),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ── Rate Limiting Middleware ──────────────────────────────────

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitorBucket
	rate     int           // tokens per interval
	interval time.Duration // refill interval
	burst    int           // max tokens
}

type visitorBucket struct {
	tokens   int
	lastSeen time.Time
}

func newRateLimiter(rate int, interval time.Duration, burst int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitorBucket),
		rate:     rate,
		interval: interval,
		burst:    burst,
	}
	// Periodic cleanup of stale visitors
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists {
		rl.visitors[ip] = &visitorBucket{
			tokens:   rl.burst - 1,
			lastSeen: now,
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(v.lastSeen)
	refill := int(elapsed/rl.interval) * rl.rate
	v.tokens += refill
	if v.tokens > rl.burst {
		v.tokens = rl.burst
	}
	v.lastSeen = now

	if v.tokens <= 0 {
		return false
	}
	v.tokens--
	return true
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	cutoff := time.Now().Add(-10 * time.Minute)
	for ip, v := range rl.visitors {
		if v.lastSeen.Before(cutoff) {
			delete(rl.visitors, ip)
		}
	}
}

func withRateLimit(limiter *rateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractClientIP(r)
			if !limiter.allow(ip) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"rate limit exceeded"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ── Request Body Size Limiter ────────────────────────────────

const maxRequestBodySize = 10 * 1024 * 1024 // 10MB

func withBodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		}
		next.ServeHTTP(w, r)
	})
}

// ── Request Context ──────────────────────────────────────────

func requestContextFromRequest(r *http.Request) auth.RequestContext {
	return auth.RequestContext{
		IP:        extractClientIP(r),
		UserAgent: strings.TrimSpace(r.UserAgent()),
	}
}

func extractClientIP(r *http.Request) string {
	remoteHost := strings.TrimSpace(r.RemoteAddr)
	if h, _, err := net.SplitHostPort(remoteHost); err == nil {
		remoteHost = h
	}

	// Only trust X-Forwarded-For from known private/loopback proxies
	xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xff != "" && isTrustedProxy(remoteHost) {
		// Use the leftmost (client) IP from XFF
		if strings.Contains(xff, ",") {
			return strings.TrimSpace(strings.Split(xff, ",")[0])
		}
		return xff
	}

	return remoteHost
}

// isTrustedProxy checks if the remote address is a private/loopback IP.
func isTrustedProxy(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	// Loopback
	if parsed.IsLoopback() {
		return true
	}
	// Private ranges: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	privateRanges := []struct{ start, end net.IP }{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
	}
	ip4 := parsed.To4()
	if ip4 == nil {
		return false
	}
	for _, pr := range privateRanges {
		s4 := pr.start.To4()
		e4 := pr.end.To4()
		if s4 != nil && e4 != nil &&
			bytesGTE(ip4, s4) && bytesGTE(e4, ip4) {
			return true
		}
	}
	return false
}

func bytesGTE(a, b net.IP) bool {
	for i := range a {
		if a[i] < b[i] {
			return false
		}
		if a[i] > b[i] {
			return true
		}
	}
	return true
}

// tokenFromRequest extracts the bearer token from the Authorization header.
// Tokens in query strings are no longer accepted for security reasons.
func tokenFromRequest(r *http.Request) string {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
		return strings.TrimSpace(authorization[7:])
	}
	return ""
}

// ── Response Recorder ────────────────────────────────────────

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// ── Role Groups (reusable across handlers) ───────────────────

// Federation management roles
var federationReadRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleFederationPresident, auth.RoleFederationSecretary,
	auth.RoleProvincialAdmin, auth.RoleBTC, auth.RoleTechnicalDirector,
}
var federationWriteRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleFederationPresident, auth.RoleFederationSecretary,
}

// BTC management roles
var btcReadRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleBTC, auth.RoleRefereeManager,
	auth.RoleDelegate, auth.RoleFederationPresident,
}
var btcWriteRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleBTC, auth.RoleRefereeManager,
}

// Club management roles
var clubReadRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleClubLeader, auth.RoleClubViceLeader,
	auth.RoleClubSecretary, auth.RoleClubAccountant,
	auth.RoleProvincialAdmin, auth.RoleFederationPresident, auth.RoleCoach,
}
var clubWriteRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleClubLeader, auth.RoleClubViceLeader, auth.RoleClubSecretary,
}

// Provincial management roles
var provincialReadRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleProvincialAdmin, auth.RoleProvincialPresident,
	auth.RoleProvincialVicePresident, auth.RoleProvincialSecretary,
	auth.RoleProvincialTechnicalHead, auth.RoleProvincialRefereeHead,
	auth.RoleProvincialCommitteeMember, auth.RoleProvincialAccountant,
	auth.RoleFederationPresident, auth.RoleFederationSecretary,
}
var provincialWriteRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleProvincialAdmin, auth.RoleProvincialPresident,
	auth.RoleProvincialSecretary,
}

// ── Route-Specific Body Limits ──────────────────────────────

func withBodyLimitSize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}
