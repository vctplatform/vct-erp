package organization

import (
	"fmt"
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain"
	"vct-platform/backend/internal/shared/httputil"
)

// ── Team Handlers ───────────────────────────────────────────

func (m *Module) handleTeamRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu đăng nhập")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/teams")
	id := strings.Trim(path, "/")

	if id == "" {
		switch r.Method {
		case http.MethodGet:
			if !authz.CanEntityAction(p.User.Role, "teams", authz.ActionView) {
				httputil.Error(w, http.StatusForbidden, "AUTH_403", fmt.Sprintf("Vai trò %s không có quyền xem đội", p.User.Role))
				return
			}
			list, err := m.service.ListTeams(r.Context())
			if err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, list)
		case http.MethodPost:
			if !authz.CanEntityAction(p.User.Role, "teams", authz.ActionCreate) {
				httputil.Error(w, http.StatusForbidden, "AUTH_403", fmt.Sprintf("Vai trò %s không có quyền tạo đội", p.User.Role))
				return
			}
			var payload domain.Team
			if err := httputil.DecodeJSON(r, &payload); err != nil {
				httputil.Error(w, http.StatusBadRequest, "ORG_400", err.Error())
				return
			}
			created, err := m.service.CreateTeam(r.Context(), payload)
			if err != nil {
				httputil.Error(w, http.StatusBadRequest, "ORG_400", err.Error())
				return
			}
			m.broadcaster.BroadcastEntityChange("teams", "created", created.ID, httputil.ToMap(created), nil)
			httputil.Success(w, http.StatusCreated, created)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	// Detail routes (Get, Patch, Delete)
	switch r.Method {
	case http.MethodGet:
		if !authz.CanEntityAction(p.User.Role, "teams", authz.ActionView) {
			httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem chi tiết đội")
			return
		}
		team, err := m.service.GetTeam(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "ORG_404", "Không tìm thấy đội")
			return
		}
		httputil.Success(w, http.StatusOK, team)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Referee Handlers ────────────────────────────────────────

func (m *Module) handleRefereeRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu đăng nhập")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/referees")
	id := strings.Trim(path, "/")

	if id == "" {
		if r.Method == http.MethodGet {
			if !authz.CanEntityAction(p.User.Role, "referees", authz.ActionView) {
				httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem trọng tài")
				return
			}
			list, err := m.service.ListReferees(r.Context())
			if err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, list)
			return
		}
	}

	httputil.Error(w, http.StatusNotFound, "ORG_404", "Không tìm thấy tài nguyên trọng tài")
}

// ── Arena Handlers ──────────────────────────────────────────

func (m *Module) handleArenaRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		list, err := m.service.ListArenas(r.Context())
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, list)
		return
	}
	httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
}
