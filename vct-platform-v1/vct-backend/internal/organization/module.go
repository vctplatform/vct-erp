package organization

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/organization"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Organization module.
type Module struct {
	service     *organization.Service
	logger      *slog.Logger
	broadcaster httputil.EventBroadcaster
}

// Deps holds the dependencies for the Organization module.
type Deps struct {
	Service     *organization.Service
	Logger      *slog.Logger
	Broadcaster httputil.EventBroadcaster
}

// New creates a new Organization module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		logger:      deps.Logger.With(slog.String("module", "organization")),
		broadcaster: deps.Broadcaster,
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers all organization-related routes.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/teams", m.handleTeamRoutes)
	mux.HandleFunc("/api/v1/teams/", m.handleTeamRoutes)
	mux.HandleFunc("/api/v1/referees", m.handleRefereeRoutes)
	mux.HandleFunc("/api/v1/referees/", m.handleRefereeRoutes)
	mux.HandleFunc("/api/v1/arenas", m.handleArenaRoutes)
	mux.HandleFunc("/api/v1/arenas/", m.handleArenaRoutes)

	m.logger.Info("organization module routes registered")
}
