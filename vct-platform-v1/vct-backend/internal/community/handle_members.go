package community

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain/community"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleMemberRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/members")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListMembers(w, r)
		case http.MethodPost:
			m.handleCreateMember(w, r)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	httputil.Error(w, http.StatusNotFound, "MEMBER_404", "Không tìm thấy thành viên")
}

func (m *Module) handleListMembers(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "members", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem thành viên")
		return
	}

	clubID := r.URL.Query().Get("clubId")
	var list []community.Member
	var err error

	if clubID != "" {
		list, err = m.service.ListMembersByClub(r.Context(), clubID)
	} else {
		list, err = m.service.ListMembers(r.Context())
	}

	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, list)
}

func (m *Module) handleCreateMember(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "members", authz.ActionCreate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo thành viên")
		return
	}

	var payload community.Member
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "MEMBER_400", err.Error())
		return
	}

	created, err := m.service.CreateMember(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "MEMBER_400", err.Error())
		return
	}

	raw, _ := clubToMap(created)
	m.broadcast("members", "created", created.ID, raw, nil)
	httputil.Success(w, http.StatusCreated, created)
}
