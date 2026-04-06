package athlete

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/athlete"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Athlete module.
type Module struct {
	service     *athlete.Service
	profile     *athlete.ProfileService
	training    *athlete.TrainingService
	broadcaster httputil.EventBroadcaster
	authFn      func(r *http.Request) (string, error)
	logger      *slog.Logger
}

// Deps holds the dependencies for the Athlete module.
type Deps struct {
	Service     *athlete.Service
	Profile     *athlete.ProfileService
	Training    *athlete.TrainingService
	Broadcaster httputil.EventBroadcaster
	AuthFn      func(r *http.Request) (string, error)
	Logger      *slog.Logger
}

// New creates a new Athlete module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		profile:     deps.Profile,
		training:    deps.Training,
		broadcaster: deps.Broadcaster,
		authFn:      deps.AuthFn,
		logger:      deps.Logger.With(slog.String("module", "athlete")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers athlete routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/athletes/", m.handleAthleteRoutes)
	m.handleAthleteProfileRoutes(mux)
	m.handleTrainingSessionRoutes(mux)
	m.logger.Info("athlete module routes registered")
}

// broadcast sends a real-time entity change event.
func (m *Module) broadcast(action, id string, data map[string]any, meta map[string]any) {
	if m.broadcaster != nil {
		m.broadcaster.BroadcastEntityChange("athletes", action, id, data, meta)
	}
}
