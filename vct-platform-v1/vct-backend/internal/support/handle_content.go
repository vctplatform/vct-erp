package support

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/auth"
	"vct-platform/backend/internal/domain/support"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleCategoryRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/support/categories")
	path = strings.Trim(path, "/")

	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu xác thực")
		return
	}

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListCategories(w, r)
		case http.MethodPost:
			m.handleCreateCategory(w, r, p)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	id := strings.Split(path, "/")[0]
	switch r.Method {
	case http.MethodGet:
		m.handleGetCategory(w, r, id)
	case http.MethodPut, http.MethodPatch:
		m.handleUpdateCategory(w, r, id, p)
	case http.MethodDelete:
		m.handleDeleteCategory(w, r, id, p)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := m.service.ListCategories(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, cats)
}

func (m *Module) GetCategory(w http.ResponseWriter, r *http.Request, id string) {
	cat, err := m.service.GetCategory(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "SUPPORT_404", "Không tìm thấy danh mục")
		return
	}
	httputil.Success(w, http.StatusOK, cat)
}

func (m *Module) handleGetCategory(w http.ResponseWriter, r *http.Request, id string) {
	cat, err := m.service.GetCategory(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "SUPPORT_404", "Không tìm thấy danh mục")
		return
	}
	httputil.Success(w, http.StatusOK, cat)
}

func (m *Module) handleCreateCategory(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo danh mục")
		return
	}

	var c support.SupportCategory
	if err := httputil.DecodeJSON(r, &c); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}

	created, err := m.service.CreateCategory(r.Context(), c)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUpdateCategory(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật danh mục")
		return
	}

	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}

	updated, err := m.service.UpdateCategory(r.Context(), id, patch)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteCategory(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xóa danh mục")
		return
	}

	if err := m.service.DeleteCategory(r.Context(), id); err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// ── FAQ ──────────────────────────────────────────────────────

func (m *Module) handleFAQRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/support/faqs")
	path = strings.Trim(path, "/")

	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu xác thực")
		return
	}

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListFAQs(w, r)
		case http.MethodPost:
			m.handleCreateFAQ(w, r, p)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	id := strings.Split(path, "/")[0]
	switch r.Method {
	case http.MethodGet:
		m.handleGetFAQ(w, r, id)
	case http.MethodPut, http.MethodPatch:
		m.handleUpdateFAQ(w, r, id, p)
	case http.MethodDelete:
		m.handleDeleteFAQ(w, r, id, p)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListFAQs(w http.ResponseWriter, r *http.Request) {
	faqs, err := m.service.ListFAQs(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, faqs)
}

func (m *Module) handleGetFAQ(w http.ResponseWriter, r *http.Request, id string) {
	faq, err := m.service.GetFAQ(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "SUPPORT_404", "Không tìm thấy FAQ")
		return
	}
	httputil.Success(w, http.StatusOK, faq)
}

func (m *Module) handleCreateFAQ(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo FAQ")
		return
	}

	var f support.FAQ
	if err := httputil.DecodeJSON(r, &f); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}

	created, err := m.service.CreateFAQ(r.Context(), f)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleUpdateFAQ(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật FAQ")
		return
	}

	var patch map[string]any
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}

	updated, err := m.service.UpdateFAQ(r.Context(), id, patch)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "SUPPORT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, updated)
}

func (m *Module) handleDeleteFAQ(w http.ResponseWriter, r *http.Request, id string, p auth.Principal) {
	if !isAdmin(p.User.Role) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xóa FAQ")
		return
	}

	if err := m.service.DeleteFAQ(r.Context(), id); err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"message": "deleted"})
}
