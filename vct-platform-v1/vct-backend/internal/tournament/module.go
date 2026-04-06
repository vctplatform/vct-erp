package tournament

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/adapter"
	"vct-platform/backend/internal/domain/btc"
	tournamentdom "vct-platform/backend/internal/domain/tournament"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Tournament module.
type Module struct {
	service     adapter.TournamentCRUD
	mgmt        *tournamentdom.MgmtService
	btc         *btc.Service
	broadcaster httputil.EventBroadcaster
	logger      *slog.Logger
}

// Deps holds the dependencies for the Tournament module.
type Deps struct {
	Service     adapter.TournamentCRUD
	Mgmt        *tournamentdom.MgmtService
	BTC         *btc.Service
	Broadcaster httputil.EventBroadcaster
	Logger      *slog.Logger
}

// New creates a new Tournament module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		mgmt:        deps.Mgmt,
		btc:         deps.BTC,
		broadcaster: deps.Broadcaster,
		logger:      deps.Logger.With(slog.String("module", "tournament")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers tournament routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/tournaments", m.handleTournamentRoutes)
	mux.HandleFunc("/api/v1/tournaments/", m.handleTournamentRoutes)
	
	// Advanced management
	mux.HandleFunc("/api/v1/tournament-mgmt/", m.handleTournamentMgmt)
	
	// BTC (Ban Tổ Chức) logic
	m.RegisterBTCRoutes(mux)
	
	// Brackets & orchestrated actions
	m.RegisterBracketRoutes(mux)

	m.logger.Info("tournament module routes registered")
}

// broadcast sends a real-time entity change event.
func (m *Module) broadcast(action, id string, data map[string]any, meta map[string]any) {
	if m.broadcaster != nil {
		m.broadcaster.BroadcastEntityChange("tournaments", action, id, data, meta)
	}
}
