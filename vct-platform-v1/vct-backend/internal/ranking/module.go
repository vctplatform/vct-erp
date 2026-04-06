// Package ranking implements the VCT Ranking Module using Hybrid Architecture.
// Manages athlete and team rankings for Võ Cổ Truyền tournaments.
//
// All features are SIMPLE — read-only query handlers.
package ranking

import (
	"log/slog"
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/ranking"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Ranking module.
type Module struct {
	service *ranking.Service
	authFn  func(r *http.Request) (string, error)
	logger  *slog.Logger
}

// Deps holds the dependencies for the Ranking module.
type Deps struct {
	Service *ranking.Service
	Logger  *slog.Logger
	AuthFn  func(r *http.Request) (string, error)
}

// New creates a new Ranking module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service: deps.Service,
		authFn:  deps.AuthFn,
		logger:  deps.Logger.With(slog.String("module", "ranking")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers ranking routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/rankings/", m.handleRankingRoutes)
	m.logger.Info("ranking module routes registered")
}

// ── Routing ──────────────────────────────────────────────────

func (m *Module) handleRankingRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/rankings")
	path = strings.Trim(path, "/")

	switch {
	case path == "" || path == "athletes":
		m.handleListAthleteRankings(w, r)

	case path == "teams":
		m.handleListTeamRankings(w, r)

	default:
		// /rankings/{id}
		segments := strings.Split(path, "/")
		m.handleGetAthleteRanking(w, r, segments[0])
	}
}

// ── Feature: List Athlete Rankings ───────────────────────────
// Complexity: SIMPLE — query only

func (m *Module) handleListAthleteRankings(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category != "" {
		list, err := m.service.ListAthleteRankingsByCategory(r.Context(), category)
		if err != nil {
			httputil.WriteError(w, err)
			return
		}
		httputil.JSON(w, http.StatusOK, list)
		return
	}

	list, err := m.service.ListAthleteRankings(r.Context())
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, list)
}

// ── Feature: List Team Rankings ──────────────────────────────
// Complexity: SIMPLE — query only

func (m *Module) handleListTeamRankings(w http.ResponseWriter, r *http.Request) {
	list, err := m.service.ListTeamRankings(r.Context())
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, list)
}

// ── Feature: Get Athlete Ranking ─────────────────────────────
// Complexity: SIMPLE — query by ID

func (m *Module) handleGetAthleteRanking(w http.ResponseWriter, r *http.Request, id string) {
	item, err := m.service.GetAthleteRanking(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "RANKING_404", "Không tìm thấy xếp hạng")
		return
	}
	httputil.JSON(w, http.StatusOK, item)
}
