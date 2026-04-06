// Package heritage implements the VCT Heritage Module using Hybrid Architecture.
// Manages belt ranks (đai) and techniques (kỹ thuật) for Võ Cổ Truyền.
//
// All features are SIMPLE CRUD — thin handlers delegating to domain service.
package heritage

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/heritage"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Heritage module.
type Module struct {
	service     *heritage.Service
	broadcaster httputil.EventBroadcaster
	authFn      func(r *http.Request) (string, error)
	logger      *slog.Logger
}

// Deps holds the dependencies for the Heritage module.
type Deps struct {
	Service     *heritage.Service
	Logger      *slog.Logger
	Broadcaster httputil.EventBroadcaster
	AuthFn      func(r *http.Request) (string, error)
}

// New creates a new Heritage module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:     deps.Service,
		broadcaster: deps.Broadcaster,
		authFn:      deps.AuthFn,
		logger:      deps.Logger.With(slog.String("module", "heritage")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers heritage routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/belts/", m.handleBeltRoutes)
	mux.HandleFunc("/api/v1/techniques/", m.handleTechniqueRoutes)
	m.logger.Info("heritage module routes registered")
}

// ── Belt Routes ──────────────────────────────────────────────

func (m *Module) handleBeltRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/belts")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListBelts(w, r)
		case http.MethodPost:
			m.handleCreateBelt(w, r)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	// GET /api/v1/belts/{id}
	id := strings.Split(path, "/")[0]
	m.handleGetBelt(w, r, id)
}

func (m *Module) handleListBelts(w http.ResponseWriter, r *http.Request) {
	list, err := m.service.ListBelts(r.Context())
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, list)
}

func (m *Module) handleGetBelt(w http.ResponseWriter, r *http.Request, id string) {
	belt, err := m.service.GetBelt(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "HERITAGE_404", "Không tìm thấy đai")
		return
	}
	httputil.JSON(w, http.StatusOK, belt)
}

func (m *Module) handleCreateBelt(w http.ResponseWriter, r *http.Request) {
	var payload heritage.BeltRank
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "HERITAGE_400", err.Error())
		return
	}

	created, err := m.service.CreateBelt(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "HERITAGE_400", err.Error())
		return
	}

	if m.broadcaster != nil {
		raw, _ := toMap(created)
		m.broadcaster.BroadcastEntityChange("belts", "created", created.ID, raw, nil)
	}

	httputil.JSON(w, http.StatusCreated, created)
}

// ── Technique Routes ─────────────────────────────────────────

func (m *Module) handleTechniqueRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/techniques")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			category := r.URL.Query().Get("loai")
			if category != "" {
				m.handleListTechniquesByCategory(w, r, category)
			} else {
				m.handleListTechniques(w, r)
			}
		case http.MethodPost:
			m.handleCreateTechnique(w, r)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	id := strings.Split(path, "/")[0]
	m.handleGetTechnique(w, r, id)
}

func (m *Module) handleListTechniques(w http.ResponseWriter, r *http.Request) {
	list, err := m.service.ListTechniques(r.Context())
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, list)
}

func (m *Module) handleListTechniquesByCategory(w http.ResponseWriter, r *http.Request, category string) {
	list, err := m.service.ListTechniquesByCategory(r.Context(), category)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, list)
}

func (m *Module) handleGetTechnique(w http.ResponseWriter, r *http.Request, id string) {
	tech, err := m.service.GetTechnique(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "HERITAGE_404", "Không tìm thấy kỹ thuật")
		return
	}
	httputil.JSON(w, http.StatusOK, tech)
}

func (m *Module) handleCreateTechnique(w http.ResponseWriter, r *http.Request) {
	var payload heritage.Technique
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "HERITAGE_400", err.Error())
		return
	}

	created, err := m.service.CreateTechnique(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "HERITAGE_400", err.Error())
		return
	}

	if m.broadcaster != nil {
		raw, _ := toMap(created)
		m.broadcaster.BroadcastEntityChange("techniques", "created", created.ID, raw, nil)
	}

	httputil.JSON(w, http.StatusCreated, created)
}

// ── Helpers ──────────────────────────────────────────────────

func toMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}
