package scoring

import (
	"database/sql"
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/scoring"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Scoring module implementing httputil.Module.
// It owns all scoring-related routes, handlers, and domain logic.
type Module struct {
	service         *scoring.Service
	registration    *scoring.RegistrationService
	broadcaster     httputil.EventBroadcaster
	authFn          func(r *http.Request) (string, error) // extracts userID from request
	logger          *slog.Logger
}

// Deps holds the dependencies needed to create a Scoring module.
type Deps struct {
	DB           *sql.DB        // nil = use in-memory
	Logger       *slog.Logger
	Broadcaster  httputil.EventBroadcaster
	AuthFn       func(r *http.Request) (string, error)
	Config       scoring.ScoringConfig
	Registration *scoring.RegistrationService
}

// New creates a new Scoring module with the given dependencies.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.Config == (scoring.ScoringConfig{}) {
		deps.Config = scoring.DefaultScoringConfig()
	}

	var repo scoring.ScoringRepository
	if deps.DB != nil {
		repo = newPgRepository(deps.DB)
	} else {
		repo = newMemRepository()
	}

	return &Module{
		service:      scoring.NewService(repo, deps.Config),
		registration: deps.Registration,
		broadcaster:  deps.Broadcaster,
		authFn:       deps.AuthFn,
		logger:       deps.Logger.With(slog.String("module", "scoring")),
	}
}

// Service returns the underlying scoring service for backward compatibility
// with the existing wire.go during migration.
func (m *Module) Service() *scoring.Service {
	return m.service
}

// ── httputil.Module Interface ────────────────────────────────

// Compile-time check: Module implements httputil.Module.
var _ httputil.Module = (*Module)(nil)

// authenticate extracts userID from request, returns error response if fails.
func (m *Module) authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	if m.authFn == nil {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "yêu cầu xác thực")
		return "", false
	}
	userID, err := m.authFn(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", err.Error())
		return "", false
	}
	return userID, true
}

// broadcast sends a real-time entity change event if broadcaster is available.
func (m *Module) broadcast(entityType, action, entityID string, data map[string]any, meta map[string]any) {
	if m.broadcaster != nil {
		m.broadcaster.BroadcastEntityChange(entityType, action, entityID, data, meta)
	}
}
