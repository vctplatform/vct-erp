package support

import (
	"net/http"
	"strconv"
	"strings"

	"vct-platform/backend/internal/auth"
	"vct-platform/backend/internal/domain/support"
	"vct-platform/backend/internal/shared/httputil"
)

// Roles that can manage ALL tickets
var adminRoles = []auth.UserRole{
	auth.RoleAdmin, auth.RoleBTC, auth.RoleFederationPresident,
	auth.RoleFederationSecretary, auth.RoleTechnicalDirector,
}

func isAdmin(role auth.UserRole) bool {
	for _, r := range adminRoles {
		if r == role {
			return true
		}
	}
	return false
}

func (m *Module) handleTicketRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/support/tickets")
	path = strings.Trim(path, "/")

	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu xác thực")
		return
	}

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListTickets(w, r, p)
		case http.MethodPost:
			m.handleCreateTicket(w, r, p)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	parts := strings.Split(path, "/")
	id := parts[0]

	if len(parts) > 1 {
		switch parts[1] {
		case "reply":
			m.handleCreateReply(w, r, id, p)
			return
		case "replies":
			m.handleListReplies(w, r, id, p)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		m.handleGetTicket(w, r, id, p)
	case http.MethodPut, http.MethodPatch:
		m.handleUpdateTicket(w, r, id, p)
	case http.MethodDelete:
		m.handleDeleteTicket(w, r, id, p)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListTickets(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	filter := support.ListFilter{
		Page:     page,
		Limit:    limit,
		Status:   q.Get("status"),
		Priority: q.Get("priority"),
		Type:     q.Get("type"),
		Search:   q.Get("search"),
	}

	if !isAdmin(p.User.Role) {
		filter.UserID = p.User.ID
	}

	result, err := m.service.ListTickets(r.Context(), filter)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, result)
}

func (m *Module) handleCreateTicket(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	var t support.SupportTicket
	if err := httputil.DecodeJSON(r, &t); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}

	t.NguoiTaoID = p.User.ID
	if t.NguoiTaoTen == "" {
		t.NguoiTaoTen = p.User.Username
	}

	created, err := m.service.CreateTicket(r.Context(), t)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleGetTicket(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	ticket, err := m.service.GetTicket(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "SUPPORT_404", "Không tìm thấy ticket")
		return
	}

	if !isAdmin(p.User.Role) && ticket.NguoiTaoID != p.User.ID {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem ticket này")
		return
	}

	httputil.Success(w, http.StatusOK, ticket)
}

func (m *Module) handleUpdateTicket(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Chỉ admin mới có quyền cập nhật ticket")
		return
	}

	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}

	updated, err := m.service.UpdateTicket(r.Context(), id, patch)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteTicket(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Chỉ admin mới có quyền xóa ticket")
		return
	}

	if err := m.service.DeleteTicket(r.Context(), id); err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (m *Module) handleCreateReply(w http.ResponseWriter, r *http.Request, ticketID string, p auth.Principal) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	var reply support.TicketReply
	if err := httputil.DecodeJSON(r, &reply); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}

	reply.TicketID = ticketID
	reply.NguoiTraID = p.User.ID
	if reply.NguoiTra == "" {
		reply.NguoiTra = p.User.Username
	}
	reply.IsStaff = isAdmin(p.User.Role)

	created, err := m.service.CreateReply(r.Context(), reply)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleListReplies(w http.ResponseWriter, r *http.Request, ticketID string, p auth.Principal) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	ticket, err := m.service.GetTicket(r.Context(), ticketID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "SUPPORT_404", "Không tìm thấy ticket")
		return
	}

	if !isAdmin(p.User.Role) && ticket.NguoiTaoID != p.User.ID {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem phản hồi")
		return
	}

	replies, err := m.service.ListReplies(r.Context(), ticketID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, replies)
}

func (m *Module) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	p, ok := httputil.GetPrincipal(r)
	if !ok || !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem thống kê")
		return
	}

	stats, err := m.service.GetStats(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}
