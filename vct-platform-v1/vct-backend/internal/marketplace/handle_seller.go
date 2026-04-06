package marketplace

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/auth"
	"vct-platform/backend/internal/domain/marketplace"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) RegisterSellerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/marketplace/seller/dashboard", m.handleSellerDashboard)
	mux.HandleFunc("/api/v1/marketplace/seller/products", m.handleSellerProducts)
	mux.HandleFunc("/api/v1/marketplace/seller/products/", m.handleSellerProductDetail)
	mux.HandleFunc("/api/v1/marketplace/seller/orders", m.handleSellerOrders)
	mux.HandleFunc("/api/v1/marketplace/seller/orders/", m.handleSellerOrderDetail)
}

func (m *Module) handleSellerDashboard(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireMarketplaceManager(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền truy cập dashboard người bán")
		return
	}

	dashboard, err := m.service.SellerDashboard(r.Context(), marketplaceScopeSellerID(p))
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, dashboard)
}

func (m *Module) handleSellerProducts(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireMarketplaceManager(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý sản phẩm")
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := m.service.ListSellerProducts(r.Context(), marketplaceScopeSellerID(p))
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{
			"items": items,
			"total": len(items),
		})
	case http.MethodPost:
		var payload marketplace.Product
		if err := httputil.DecodeJSON(r, &payload); err != nil {
			httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
			return
		}

		if !isMarketplaceAdmin(p) {
			payload.SellerID = p.User.ID
			payload.SellerRole = string(p.User.Role)
			payload.SellerName = marketplacePrincipalName(p)
		} else {
			if strings.TrimSpace(payload.SellerID) == "" {
				payload.SellerID = p.User.ID
			}
			if strings.TrimSpace(payload.SellerRole) == "" {
				payload.SellerRole = string(p.User.Role)
			}
			if strings.TrimSpace(payload.SellerName) == "" {
				payload.SellerName = marketplacePrincipalName(p)
			}
		}

		created, err := m.service.CreateProduct(r.Context(), payload)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleSellerProductDetail(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireMarketplaceManager(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý sản phẩm")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/v1/marketplace/seller/products/")
	productID = strings.TrimSpace(strings.Trim(productID, "/"))

	switch r.Method {
	case http.MethodGet:
		items, err := m.service.ListSellerProducts(r.Context(), marketplaceScopeSellerID(p))
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		for _, item := range items {
			if item.ID == productID {
				httputil.Success(w, http.StatusOK, item)
				return
			}
		}
		httputil.Error(w, http.StatusNotFound, "MARKETPLACE_404", "Không tìm thấy sản phẩm hoặc không có quyền")
	case http.MethodPatch:
		patch := map[string]any{}
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
			return
		}
		updated, err := m.service.UpdateProduct(r.Context(), productID, patch)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, updated)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleSellerOrders(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireMarketplaceManager(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý đơn hàng")
		return
	}

	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	items, err := m.service.ListOrders(r.Context(), marketplaceScopeSellerID(p))
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{
		"items": items,
		"total": len(items),
	})
}

func (m *Module) handleSellerOrderDetail(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !requireMarketplaceManager(p) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền quản lý đơn hàng")
		return
	}

	if r.Method != http.MethodPatch {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	orderID := strings.TrimPrefix(r.URL.Path, "/api/v1/marketplace/seller/orders/")
	orderID = strings.TrimSpace(strings.Trim(orderID, "/"))

	patch := map[string]any{}
	if err := httputil.DecodeJSON(r, &patch); err != nil {
		httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
		return
	}

	updated, err := m.service.UpdateOrder(r.Context(), orderID, patch)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, updated)
}

// ── Helpers ──────────────────────────────────────────────────

func requireMarketplaceManager(p auth.Principal) bool {
	role := p.User.Role
	return role == auth.RoleAdmin ||
		role == auth.RoleFederationPresident ||
		role == auth.RoleFederationSecretary ||
		role == auth.RoleProvincialAdmin ||
		role == auth.RoleClubLeader ||
		role == auth.RoleClubViceLeader ||
		role == auth.RoleClubSecretary ||
		role == auth.RoleClubAccountant ||
		role == auth.RoleCoach
}

func isMarketplaceAdmin(p auth.Principal) bool {
	switch p.User.Role {
	case auth.RoleAdmin, auth.RoleFederationPresident, auth.RoleFederationSecretary, auth.RoleProvincialAdmin:
		return true
	default:
		return false
	}
}

func marketplaceScopeSellerID(p auth.Principal) string {
	if isMarketplaceAdmin(p) {
		return ""
	}
	return strings.TrimSpace(p.User.ID)
}

func marketplacePrincipalName(p auth.Principal) string {
	if strings.TrimSpace(p.User.DisplayName) != "" {
		return strings.TrimSpace(p.User.DisplayName)
	}
	if strings.TrimSpace(p.User.Username) != "" {
		return strings.TrimSpace(p.User.Username)
	}
	return "Seller"
}
