package httpapi

import (
	"net/http"

	"vct-platform/backend/internal/apiversioning"
	"vct-platform/backend/internal/metrics"
)

// Handler returns the fully wired HTTP handler with CORS, logging, and middleware chain.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/readyz", s.handleReadiness)
	mux.Handle("/metrics", s.metricsRegistry.ExposeHandler())

	// API Documentation (OpenAPI / Scalar UI)
	mux.HandleFunc("/api/docs", s.handleAPIDocs)
	mux.HandleFunc("/api/openapi.yaml", s.handleAPISpec)

	mux.HandleFunc("/api/v1/ws", s.handleWebSocket)
	// Auth routes — stricter rate limiting + smaller body limit for login/register
	loginRL := withRateLimit(s.loginRateLimiter)
	loginBody := withBodyLimitSize(2 * 1024) // 2KB body limit
	mux.Handle("/api/v1/auth/login", loginRL(loginBody(http.HandlerFunc(s.handleAuthLogin))))
	mux.Handle("/api/v1/auth/register", loginRL(loginBody(http.HandlerFunc(s.handleAuthRegister))))
	mux.HandleFunc("/api/v1/auth/refresh", s.handleAuthRefresh)
	mux.HandleFunc("/api/v1/auth/me", s.withAuth(s.handleAuthMe))
	mux.HandleFunc("/api/v1/auth/logout", s.withAuth(s.handleAuthLogout))
	mux.HandleFunc("/api/v1/auth/revoke", s.withAuth(s.handleAuthRevoke))
	mux.HandleFunc("/api/v1/auth/audit", s.withAuth(s.handleAuthAudit))
	mux.HandleFunc("/api/v1/auth/switch-context", s.withAuth(s.handleAuthSwitchContext))
	mux.HandleFunc("/api/v1/auth/my-roles", s.withAuth(s.handleAuthMyRoles))
	// OTP routes — same rate limit as login
	mux.Handle("/api/v1/auth/send-otp", loginRL(loginBody(http.HandlerFunc(s.handleAuthSendOTP))))
	mux.Handle("/api/v1/auth/verify-otp", loginRL(loginBody(http.HandlerFunc(s.handleAuthVerifyOTP))))
	// Modules are now registered via s.modules loop at the bottom.
	// Public API (no auth)
	mux.HandleFunc("/api/v1/public/", s.handlePublicRoutes)
	// Specific domain entities — Use subtree patterns only to avoid 307 redirects
	// in Go 1.22+ ServeMux (exact + subtree double registration causes 307).
	// Heritage
	// Community
	// VCT Marketplace
	// ── Core Routes ──────────────────────────────────────────
	// (handleInvoiceGet and handlePortalActivities migrated to domain modules)
	// ── Domain Events ────────────────────────────────────────
	mux.HandleFunc("/api/v1/events/recent", s.withAuth(s.handleRecentEvents))
	// Generic entity CRUD (catch-all for unmigrated entities)
	mux.HandleFunc("/api/v1/", s.handleEntityRoutes)

	// ── Decentralized Modules (Hybrid Architecture) ──────────
	for _, m := range s.modules {
		m.RegisterRoutes(mux)
	}

	return apiversioning.Middleware(s.versionRegistry)(
		metrics.HTTPMiddleware(s.metricsRegistry)(
			withRecover(
				withRequestID(
					withSecurityHeaders(
						withRateLimit(s.rateLimiter)(
							withBodyLimit(
								s.withCSRF(
									s.withCORS(
										s.withLogging(mux),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}
