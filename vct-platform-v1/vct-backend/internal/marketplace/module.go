package marketplace

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/marketplace"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Marketplace module.
type Module struct {
	service     *marketplace.Service
	broadcaster httputil.EventBroadcaster
	logger      *slog.Logger
}

// Deps holds the dependencies for the Marketplace module.
type Deps struct {
	Service     *marketplace.Service
	Broadcaster httputil.EventBroadcaster
	Logger      *slog.Logger
}

// New creates a new Marketplace module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		broadcaster: deps.Broadcaster,
		logger:      deps.Logger.With(slog.String("module", "marketplace")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers marketplace routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/marketplace/products", m.handleProductRoutes)
	mux.HandleFunc("/api/v1/marketplace/products/", m.handleProductRoutes)
	mux.HandleFunc("/api/v1/marketplace/orders", m.handleOrderRoutes)
	mux.HandleFunc("/api/v1/marketplace/orders/", m.handleOrderRoutes)

	// Seller sub-routes
	m.RegisterSellerRoutes(mux)

	m.logger.Info("marketplace module routes registered")
}

// broadcast sends a real-time entity change event.
func (m *Module) broadcast(entity, action, id string, data map[string]any, meta map[string]any) {
	if m.broadcaster != nil {
		m.broadcaster.BroadcastEntityChange(entity, action, id, data, meta)
	}
}
