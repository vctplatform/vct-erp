package community

import (
	"encoding/json"
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain/community"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleClubRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/clubs")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListClubs(w, r)
		case http.MethodPost:
			m.handleCreateClub(w, r)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	id := strings.Split(path, "/")[0]
	switch r.Method {
	case http.MethodGet:
		m.handleGetClub(w, r, id)
	case http.MethodPatch:
		m.handleUpdateClub(w, r, id)
	case http.MethodDelete:
		m.handleDeleteClub(w, r, id)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListClubs(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "clubs", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem CLB")
		return
	}

	list, err := m.service.ListClubs(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, list)
}

func (m *Module) handleGetClub(w http.ResponseWriter, r *http.Request, id string) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "clubs", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem CLB")
		return
	}

	club, err := m.service.GetClub(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "CLUB_404", "Không tìm thấy CLB")
		return
	}
	httputil.Success(w, http.StatusOK, club)
}

func (m *Module) handleCreateClub(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "clubs", authz.ActionCreate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo CLB")
		return
	}

	var payload community.Club
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
		return
	}

	created, err := m.service.CreateClub(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
		return
	}

	raw, _ := clubToMap(created)
	m.broadcast("clubs", "created", created.ID, raw, nil)
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUpdateClub(w http.ResponseWriter, r *http.Request, id string) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "clubs", authz.ActionUpdate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật CLB")
		return
	}

	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
		return
	}

	updated, err := m.service.UpdateClub(r.Context(), id, patch)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
		return
	}

	raw, _ := clubToMap(updated)
	m.broadcast("clubs", "updated", id, raw, nil)
	httputil.Success(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteClub(w http.ResponseWriter, r *http.Request, id string) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "clubs", authz.ActionDelete) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xóa CLB")
		return
	}

	if err := m.service.DeleteClub(r.Context(), id); err != nil {
		httputil.InternalError(w, err)
		return
	}

	m.broadcast("clubs", "deleted", id, nil, nil)
	httputil.Success(w, http.StatusOK, map[string]any{"message": "deleted"})
}

func clubToMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}
