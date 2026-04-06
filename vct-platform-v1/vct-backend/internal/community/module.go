package community

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/community"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Community module.
type Module struct {
	service     *community.Service
	broadcaster httputil.EventBroadcaster
	logger      *slog.Logger
}

// Deps holds the dependencies for the Community module.
type Deps struct {
	Service     *community.Service
	Broadcaster httputil.EventBroadcaster
	Logger      *slog.Logger
}

// New creates a new Community module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		broadcaster: deps.Broadcaster,
		logger:      deps.Logger.With(slog.String("module", "community")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers community routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/clubs/", m.handleClubRoutes)
	mux.HandleFunc("/api/v1/members/", m.handleMemberRoutes)
	mux.HandleFunc("/api/v1/community-events/", m.handleCommunityEventRoutes)
	m.logger.Info("community module routes registered")
}

// broadcast sends a real-time entity change event.
func (m *Module) broadcast(entity, action, id string, data map[string]any, meta map[string]any) {
	if m.broadcaster != nil {
		m.broadcaster.BroadcastEntityChange(entity, action, id, data, meta)
	}
}
