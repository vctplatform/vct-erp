package provincial

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/auth"
	"vct-platform/backend/internal/domain/provincial"
	"vct-platform/backend/internal/shared/httputil"
)

// RBAC Helpers
func requireProvincialRead(p auth.Principal) bool {
	role := p.User.Role
	return role == auth.RoleAdmin ||
		role == auth.RoleProvincialAdmin ||
		role == auth.RoleProvincialPresident ||
		role == auth.RoleProvincialSecretary ||
		role == auth.RoleClubLeader ||
		role == auth.RoleClubViceLeader
}

func requireProvincialWrite(p auth.Principal) bool {
	role := p.User.Role
	return role == auth.RoleAdmin ||
		role == auth.RoleProvincialAdmin ||
		role == auth.RoleProvincialPresident ||
		role == auth.RoleProvincialSecretary
}

func resolveProvinceID(r *http.Request) string {
	if prov := r.URL.Query().Get("province_id"); prov != "" {
		return prov
	}
	return "PROV-HCM" // Default for testing/demo
}

// ── Dashboard Handlers ───────────────────────────────────────

func (m *Module) handleProvincialDashboard(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireProvincialRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem dashboard tỉnh thành")
		return
	}
	provID := resolveProvinceID(r)
	stats, err := m.service.GetDashboard(r.Context(), provID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}

// ── Club Handlers ───────────────────────────────────────────

func (m *Module) handleProvincialClubs(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireProvincialRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý CLB")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/clubs")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch {
	case id == "":
		if r.Method == http.MethodGet {
			m.handleListClubs(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateClub(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "":
		if r.Method == http.MethodGet {
			m.handleGetClub(w, r, id)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "approve":
		m.handleClubAction(w, r, id, "approve", p)
	case action == "suspend":
		m.handleClubAction(w, r, id, "suspend", p)
	default:
		httputil.Error(w, http.StatusNotFound, "PROVINCIAL_404", "Không tìm thấy tài nguyên")
	}
}

func (m *Module) handleListClubs(w http.ResponseWriter, r *http.Request) {
	provID := resolveProvinceID(r)
	clubs, err := m.service.ListClubs(r.Context(), provID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"clubs": clubs, "total": len(clubs)})
}

func (m *Module) handleGetClub(w http.ResponseWriter, r *http.Request, id string) {
	club, err := m.service.GetClub(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "PROVINCIAL_404", "Không tìm thấy CLB")
		return
	}
	httputil.Success(w, http.StatusOK, club)
}

func (m *Module) handleCreateClub(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo CLB")
		return
	}
	var club provincial.ProvincialClub
	if err := httputil.DecodeJSON(r, &club); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	if club.ProvinceID == "" {
		club.ProvinceID = resolveProvinceID(r)
	}
	created, err := m.service.CreateClub(r.Context(), club)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleClubAction(w http.ResponseWriter, r *http.Request, id, action string, p auth.Principal) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền thực hiện thao tác này")
		return
	}

	var err error
	switch action {
	case "approve":
		err = m.service.ApproveClub(r.Context(), id)
	case "suspend":
		err = m.service.SuspendClub(r.Context(), id)
	}

	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": action + "_done"})
}

// ── Athlete Handlers ─────────────────────────────────────────

func (m *Module) handleProvincialAthletes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireProvincialRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý VĐV")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/athletes")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch {
	case id == "":
		if r.Method == http.MethodGet {
			m.handleListAthletes(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateAthlete(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "":
		if r.Method == http.MethodGet {
			m.handleGetAthlete(w, r, id)
		} else if r.Method == http.MethodPatch || r.Method == http.MethodPut {
			m.handleUpdateAthlete(w, r, id, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
	case action == "approve":
		m.handleAthleteAction(w, r, id, "approve", p)
	case action == "deactivate":
		m.handleAthleteAction(w, r, id, "deactivate", p)
	case action == "reactivate":
		m.handleAthleteAction(w, r, id, "reactivate", p)
	default:
		httputil.Error(w, http.StatusNotFound, "PROVINCIAL_404", "Không tìm thấy tài nguyên")
	}
}

func (m *Module) handleListAthletes(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	var athletes []provincial.ProvincialAthlete
	var err error

	if clubID != "" {
		athletes, err = m.service.ListAthletesByClub(r.Context(), clubID)
	} else {
		provID := resolveProvinceID(r)
		athletes, err = m.service.ListAthletes(r.Context(), provID)
	}

	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"athletes": athletes, "total": len(athletes)})
}

func (m *Module) handleGetAthlete(w http.ResponseWriter, r *http.Request, id string) {
	athlete, err := m.service.GetAthlete(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "PROVINCIAL_404", "Không tìm thấy VĐV")
		return
	}
	httputil.Success(w, http.StatusOK, athlete)
}

func (m *Module) handleCreateAthlete(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo VĐV")
		return
	}
	var athlete provincial.ProvincialAthlete
	if err := httputil.DecodeJSON(r, &athlete); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	if athlete.ProvinceID == "" {
		athlete.ProvinceID = resolveProvinceID(r)
	}
	created, err := m.service.CreateAthlete(r.Context(), athlete)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUpdateAthlete(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật VĐV")
		return
	}
	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	if err := m.service.UpdateAthlete(r.Context(), id, patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (m *Module) handleAthleteAction(w http.ResponseWriter, r *http.Request, id, action string, p auth.Principal) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền thực hiện thao tác này")
		return
	}

	var err error
	switch action {
	case "approve":
		err = m.service.ApproveAthlete(r.Context(), id)
	case "deactivate":
		err = m.service.DeactivateAthlete(r.Context(), id)
	case "reactivate":
		err = m.service.ReactivateAthlete(r.Context(), id)
	}

	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": action + "_done"})
}

// ── Võ Sinh Handlers ─────────────────────────────────────────

func (m *Module) handleProvincialVoSinh(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireProvincialRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý võ sinh")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/vo-sinh")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	id := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch {
	case id == "stats":
		m.handleVoSinhStats(w, r)
	case id == "":
		if r.Method == http.MethodGet {
			m.handleListVoSinh(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateVoSinh(w, r, p)
		}
	case action == "":
		if r.Method == http.MethodGet {
			m.handleGetVoSinh(w, r, id)
		} else if r.Method == http.MethodPatch || r.Method == http.MethodPut {
			m.handleUpdateVoSinh(w, r, id, p)
		}
	case action == "approve":
		m.handleVoSinhAction(w, r, id, "approve", p)
	case action == "deactivate":
		m.handleVoSinhAction(w, r, id, "deactivate", p)
	case action == "reactivate":
		m.handleVoSinhAction(w, r, id, "reactivate", p)
	case action == "belt-history":
		m.handleVoSinhBeltHistory(w, r, id)
	default:
		httputil.Error(w, http.StatusNotFound, "PROVINCIAL_404", "Không tìm thấy tài nguyên")
	}
}

func (m *Module) handleVoSinhStats(w http.ResponseWriter, r *http.Request) {
	provID := resolveProvinceID(r)
	stats, err := m.service.GetVoSinhStats(r.Context(), provID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}

func (m *Module) handleListVoSinh(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	var list []provincial.VoSinh
	var err error

	if clubID != "" {
		list, err = m.service.ListVoSinhByClub(r.Context(), clubID)
	} else {
		provID := resolveProvinceID(r)
		list, err = m.service.ListVoSinh(r.Context(), provID)
	}

	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"vo_sinh": list, "total": len(list)})
}

func (m *Module) handleGetVoSinh(w http.ResponseWriter, r *http.Request, id string) {
	vs, err := m.service.GetVoSinh(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "PROVINCIAL_404", "Không tìm thấy võ sinh")
		return
	}
	httputil.Success(w, http.StatusOK, vs)
}

func (m *Module) handleCreateVoSinh(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo võ sinh")
		return
	}
	var vs provincial.VoSinh
	if err := httputil.DecodeJSON(r, &vs); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	if vs.ProvinceID == "" {
		vs.ProvinceID = resolveProvinceID(r)
	}
	created, err := m.service.CreateVoSinh(r.Context(), vs)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUpdateVoSinh(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật võ sinh")
		return
	}
	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	if err := m.service.UpdateVoSinh(r.Context(), id, patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	updated, _ := m.service.GetVoSinh(r.Context(), id)
	httputil.Success(w, http.StatusOK, updated)
}

func (m *Module) handleVoSinhAction(w http.ResponseWriter, r *http.Request, id, action string, p auth.Principal) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	if !requireProvincialWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền thực hiện thao tác này")
		return
	}

	var err error
	switch action {
	case "approve":
		err = m.service.ApproveVoSinh(r.Context(), id)
	case "deactivate":
		err = m.service.DeactivateVoSinh(r.Context(), id)
	case "reactivate":
		err = m.service.ReactivateVoSinh(r.Context(), id)
	}

	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": action + "_done"})
}

func (m *Module) handleVoSinhBeltHistory(w http.ResponseWriter, r *http.Request, id string) {
	hist, err := m.service.ListBeltHistory(r.Context(), id)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"belt_history": hist, "total": len(hist)})
}

// ── Remaining Provincial Handlers (Coaches, Referees, Committee, Transfers) ──────────

func (m *Module) handleProvincialCoaches(w http.ResponseWriter, r *http.Request) {
	// Implementation similar to athletes
	httputil.Error(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Coaches handler to be migrated")
}
