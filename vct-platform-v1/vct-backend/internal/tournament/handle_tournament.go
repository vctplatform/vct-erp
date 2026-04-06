package tournament

import (
	"encoding/json"
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain"
	"vct-platform/backend/internal/shared/auth"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleTournamentRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tournaments")
	path = strings.Trim(path, "/")

	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu xác thực")
		return
	}

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListTournaments(w, r, p)
		case http.MethodPost:
			m.handleCreateTournament(w, r, p)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	segments := strings.Split(path, "/")
	id := segments[0]

	if len(segments) > 1 {
		// Handle sub-resources if any
		httputil.Error(w, http.StatusNotFound, "TOURNAMENT_404", "Không tìm thấy tài nguyên")
		return
	}

	switch r.Method {
	case http.MethodGet:
		m.handleGetTournament(w, r, id, p)
	case http.MethodPatch:
		m.handleUpdateTournament(w, r, id, p)
	case http.MethodDelete:
		m.handleDeleteTournament(w, r, id, p)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListTournaments(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !authz.CanEntityAction(p.User.Role, "tournaments", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem giải đấu")
		return
	}

	list, err := m.service.List(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, list)
}

func (m *Module) handleGetTournament(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !authz.CanEntityAction(p.User.Role, "tournaments", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem giải đấu")
		return
	}

	t, err := m.service.GetByID(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "TOURNAMENT_404", "Không tìm thấy giải đấu")
		return
	}
	httputil.Success(w, http.StatusOK, t)
}

func (m *Module) handleCreateTournament(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !authz.CanEntityAction(p.User.Role, "tournaments", authz.ActionCreate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo giải đấu")
		return
	}

	var payload domain.Tournament
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
		return
	}

	created, err := m.service.Create(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
		return
	}

	raw, _ := tMap(created)
	m.broadcast("created", created.ID, raw, nil)
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUpdateTournament(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !authz.CanEntityAction(p.User.Role, "tournaments", authz.ActionUpdate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật giải đấu")
		return
	}

	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
		return
	}

	updated, err := m.service.Update(r.Context(), id, patch)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
		return
	}

	raw, _ := tMap(updated)
	m.broadcast("updated", id, raw, nil)
	httputil.Success(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteTournament(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !authz.CanEntityAction(p.User.Role, "tournaments", authz.ActionDelete) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xóa giải đấu")
		return
	}

	if err := m.service.Delete(r.Context(), id); err != nil {
		httputil.InternalError(w, err)
		return
	}

	m.broadcast("deleted", id, nil, nil)
	w.WriteHeader(http.StatusNoContent)
}

func tMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}
