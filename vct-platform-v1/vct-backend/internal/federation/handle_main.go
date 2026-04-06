package federation

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/auth"
	"vct-platform/backend/internal/domain/federation"
	"vct-platform/backend/internal/shared/httputil"
)

// RBAC Helpers
func requireFederationRead(p auth.Principal) bool {
	role := p.User.Role
	return role == auth.RoleAdmin ||
		role == auth.RoleBTC ||
		role == auth.RoleFederationPresident ||
		role == auth.RoleFederationSecretary ||
		role == auth.RoleTechnicalDirector
}

func requireFederationWrite(p auth.Principal) bool {
	role := p.User.Role
	return role == auth.RoleAdmin ||
		role == auth.RoleFederationPresident ||
		role == auth.RoleFederationSecretary
}

// ── Infrastructure Handlers ──────────────────────────────────

func (m *Module) handleProvinceRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem thông tin tỉnh thành")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/federation/provinces")
	id := strings.TrimSpace(strings.Trim(path, "/"))

	if id == "" {
		if r.Method == http.MethodGet {
			m.handleListProvinces(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateProvince(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	if r.Method == http.MethodGet {
		m.handleGetProvince(w, r, id)
	} else {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListProvinces(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	var provinces []federation.Province
	var err error

	if region != "" {
		provinces, err = m.main.ListProvincesByRegion(r.Context(), federation.RegionCode(region))
	} else {
		provinces, err = m.main.ListProvinces(r.Context())
	}

	if err != nil {
		httputil.InternalError(w, err)
		return
	}

	search := r.URL.Query().Get("search")
	if search != "" {
		q := strings.ToLower(search)
		var filtered []federation.Province
		for _, pv := range provinces {
			if strings.Contains(strings.ToLower(pv.Name), q) || strings.Contains(strings.ToLower(pv.Code), q) {
				filtered = append(filtered, pv)
			}
		}
		provinces = filtered
	}

	httputil.Success(w, http.StatusOK, provinces)
}

func (m *Module) handleGetProvince(w http.ResponseWriter, r *http.Request, id string) {
	prov, err := m.main.GetProvince(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "Không tìm thấy tỉnh thành")
		return
	}
	httputil.Success(w, http.StatusOK, prov)
}

func (m *Module) handleCreateProvince(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireFederationWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo tỉnh thành")
		return
	}

	var prov federation.Province
	if err := httputil.DecodeJSON(r, &prov); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}

	created, err := m.main.CreateProvince(r.Context(), prov)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUnitRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem chi hội")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/federation/units")
	id := strings.TrimSpace(strings.Trim(path, "/"))

	if id == "" {
		if r.Method == http.MethodGet {
			m.handleListUnits(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateUnit(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	if r.Method == http.MethodGet {
		m.handleGetUnit(w, r, id)
	} else {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListUnits(w http.ResponseWriter, r *http.Request) {
	units, err := m.main.ListUnits(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, units)
}

func (m *Module) handleGetUnit(w http.ResponseWriter, r *http.Request, id string) {
	unit, err := m.main.GetUnit(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "Không tìm thấy chi hội")
		return
	}
	httputil.Success(w, http.StatusOK, unit)
}

func (m *Module) handleCreateUnit(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !requireFederationWrite(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo chi hội")
		return
	}

	var unit federation.FederationUnit
	if err := httputil.DecodeJSON(r, &unit); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}

	created, err := m.main.CreateUnit(r.Context(), unit)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleOrgChart(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem sơ đồ tổ chức")
		return
	}

	chart, err := m.main.BuildOrgChart(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"root": chart})
}

func (m *Module) handleStatistics(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem thống kê")
		return
	}

	stats, err := m.main.GetNationalStatistics(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}

func (m *Module) handlePersonnelRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý nhân sự")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/federation/personnel")
	unitID := strings.Trim(path, "/")

	if r.Method == http.MethodGet {
		if unitID == "" {
			unitID = r.URL.Query().Get("unit_id")
		}
		list, err := m.main.ListPersonnel(r.Context(), unitID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, list)
	} else if r.Method == http.MethodPost {
		if !requireFederationWrite(p) {
			httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền bổ nhiệm nhân sự")
			return
		}
		var assign federation.PersonnelAssignment
		if err := httputil.DecodeJSON(r, &assign); err != nil {
			httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
			return
		}
		if err := m.main.AssignPersonnel(r.Context(), assign); err != nil {
			httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, map[string]string{"status": "personnel_assigned"})
	} else {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Master Data Handlers ─────────────────────────────────────

func (m *Module) handleMasterBelts(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem quy định đai")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/federation/master/belts")
	id := strings.Trim(path, "/")

	if id == "" {
		if r.Method == http.MethodGet {
			list, err := m.main.ListMasterBelts(r.Context())
			if err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, list)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	if r.Method == http.MethodGet {
		belt, err := m.main.GetMasterBelt(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "Không tìm thấy quy định đai")
			return
		}
		httputil.Success(w, http.StatusOK, belt)
	} else {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleMasterWeights(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireFederationRead(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem quy định hạng cân")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/federation/master/weights")
	id := strings.Trim(path, "/")

	if id == "" {
		if r.Method == http.MethodGet {
			list, err := m.main.ListMasterWeights(r.Context())
			if err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, list)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	if r.Method == http.MethodGet {
		weight, err := m.main.GetMasterWeight(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "Không tìm thấy quy định hạng cân")
			return
		}
		httputil.Success(w, http.StatusOK, weight)
	} else {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}
