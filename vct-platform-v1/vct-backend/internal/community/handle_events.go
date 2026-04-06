package community

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain/community"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleCommunityEventRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/community-events")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListEvents(w, r)
		case http.MethodPost:
			m.handleCreateEvent(w, r)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	httputil.Error(w, http.StatusNotFound, "EVENT_404", "Không tìm thấy sự kiện")
}

func (m *Module) handleListEvents(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "community_events", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem sự kiện")
		return
	}

	list, err := m.service.ListEvents(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, list)
}

func (m *Module) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "community_events", authz.ActionCreate) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo sự kiện")
		return
	}

	var payload community.Event
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "EVENT_400", err.Error())
		return
	}

	created, err := m.service.CreateEvent(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "EVENT_400", err.Error())
		return
	}

	raw, _ := clubToMap(created)
	m.broadcast("community_events", "created", created.ID, raw, nil)
	httputil.Success(w, http.StatusCreated, created)
}
