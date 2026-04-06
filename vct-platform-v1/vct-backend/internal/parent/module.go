package parent

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/parent"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Parent module.
type Module struct {
	service     *parent.Service
	authFn      func(r *http.Request) (string, error)
	logger      *slog.Logger
}

// Deps holds the dependencies for the Parent module.
type Deps struct {
	Service     *parent.Service
	AuthFn      func(r *http.Request) (string, error)
	Logger      *slog.Logger
}

// New creates a new Parent module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		authFn:      deps.AuthFn,
		logger:      deps.Logger.With(slog.String("module", "parent")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers parent routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/parent/dashboard", m.handleParentDashboard)
	mux.HandleFunc("/api/v1/parent/children", m.handleParentChildren)
	mux.HandleFunc("/api/v1/parent/children/link", m.handleParentLinkChild)
	mux.HandleFunc("/api/v1/parent/children/", m.handleParentChildDetail)
	mux.HandleFunc("/api/v1/parent/consents", m.handleParentConsents)
	mux.HandleFunc("/api/v1/parent/consents/", m.handleParentConsentAction)
	m.logger.Info("parent module routes registered")
}

// authenticate extracts userID from request using the authFn.
func (m *Module) authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	if m.authFn == nil {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "authentication required")
		return "", false
	}
	userID, err := m.authFn(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", err.Error())
		return "", false
	}
	return userID, true
}
