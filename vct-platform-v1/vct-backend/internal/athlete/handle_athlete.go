package athlete

import (
	"encoding/json"
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleAthleteRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/athletes")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListAthletes(w, r)
		case http.MethodPost:
			m.handleCreateAthlete(w, r)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	segments := strings.Split(path, "/")
	id := segments[0]

	switch r.Method {
	case http.MethodGet:
		m.handleGetAthlete(w, r, id)
	case http.MethodPatch:
		m.handleUpdateAthlete(w, r, id)
	case http.MethodDelete:
		m.handleDeleteAthlete(w, r, id)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListAthletes(w http.ResponseWriter, r *http.Request) {
	// Standardized RBAC middleware will be used in RegisterRoutes if we move to chained handlers
	// For now, let's keep it self-contained within handlers for legacy compatibility
	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu xác thực")
		return
	}

	if !authz.CanEntityAction(p.User.Role, "athletes", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem vận động viên")
		return
	}

	teamID := r.URL.Query().Get("teamId")
	tournamentID := r.URL.Query().Get("tournamentId")

	var list []domain.Athlete
	var err error

	if teamID != "" {
		list, err = m.service.ListByTeam(r.Context(), teamID)
	} else if tournamentID != "" {
		list, err = m.service.ListByTournament(r.Context(), tournamentID)
	} else {
		list, err = m.service.ListAthletes(r.Context())
	}

	if err != nil {
		httputil.InternalError(w, err)
		return
	}

	httputil.Success(w, http.StatusOK, list)
}

func (m *Module) handleGetAthlete(w http.ResponseWriter, r *http.Request, id string) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "athletes", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem vận động viên")
		return
	}

	athlete, err := m.service.GetAthlete(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "ATHLETE_404", "Không tìm thấy vận động viên")
		return
	}

	httputil.Success(w, http.StatusOK, athlete)
}

func (m *Module) handleCreateAthlete(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "athletes", authz.ActionCreate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo vận động viên")
		return
	}

	var payload domain.Athlete
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
		return
	}

	created, err := m.service.CreateAthlete(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
		return
	}

	raw, _ := toMap(created)
	m.broadcast("created", created.ID, raw, nil)

	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUpdateAthlete(w http.ResponseWriter, r *http.Request, id string) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "athletes", authz.ActionUpdate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật vận động viên")
		return
	}

	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
		return
	}

	// Logic from legacy athlete_handler.go: handle specific status update or generic patch
	var updated *domain.Athlete
	var err error

	if status, ok := patch["trang_thai"].(string); ok && len(patch) == 1 {
		updated, err = m.service.UpdateStatus(r.Context(), id, domain.TrangThaiVDV(status))
	} else {
		// Generic update currently doesn't exist in service, but we should add it if we want to remove store dependence.
		// For now, assume service will be extended or module will handle simple repo update.
		// TODO: Implement generic update in service.
		updated, err = m.service.UpdateAthlete(r.Context(), id, patch)
	}

	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
		return
	}

	raw, _ := toMap(updated)
	m.broadcast("updated", id, raw, nil)

	httputil.Success(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteAthlete(w http.ResponseWriter, r *http.Request, id string) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "athletes", authz.ActionDelete) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xóa vận động viên")
		return
	}

	// Assuming service.DeleteAthlete will exist
	if err := m.service.DeleteAthlete(r.Context(), id); err != nil {
		httputil.InternalError(w, err)
		return
	}

	m.broadcast("deleted", id, nil, nil)
	w.WriteHeader(http.StatusNoContent)
}

func toMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}
