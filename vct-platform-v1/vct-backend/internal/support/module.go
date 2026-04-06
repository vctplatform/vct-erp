package support

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/support"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Support module.
type Module struct {
	service *support.Service
	logger  *slog.Logger
}

// Deps holds the dependencies for the Support module.
type Deps struct {
	Service *support.Service
	Logger  *slog.Logger
}

// New creates a new Support module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service: deps.Service,
		logger:  deps.Logger.With(slog.String("module", "support")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers support routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/support/tickets", m.handleTicketRoutes)
	mux.HandleFunc("/api/v1/support/tickets/", m.handleTicketRoutes)
	mux.HandleFunc("/api/v1/support/categories", m.handleCategoryRoutes)
	mux.HandleFunc("/api/v1/support/categories/", m.handleCategoryRoutes)
	mux.HandleFunc("/api/v1/support/faqs", m.handleFAQRoutes)
	mux.HandleFunc("/api/v1/support/faqs/", m.handleFAQRoutes)
	mux.HandleFunc("/api/v1/support/stats", m.handleStats)
	m.logger.Info("support module routes registered")
}
