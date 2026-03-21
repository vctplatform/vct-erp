package middleware

import (
	"context"
	"net/http"
	"strings"
)

type roleContextKey struct{}

// WithRoleFromHeader loads an application role from the configured header into request context.
func WithRoleFromHeader(headerName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := strings.TrimSpace(r.Header.Get(headerName))
			if role == "" {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), roleContextKey{}, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleFromContext reads the application role injected by middleware.
func RoleFromContext(ctx context.Context) string {
	role, _ := ctx.Value(roleContextKey{}).(string)
	return role
}

// RequireRoles ensures the caller owns one of the allowed application roles.
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		if trimmed := strings.TrimSpace(role); trimmed != "" {
			allowed[trimmed] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := RoleFromContext(r.Context())
			if _, ok := allowed[role]; !ok {
				http.Error(w, `{"error":"forbidden","message":"insufficient permissions"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
