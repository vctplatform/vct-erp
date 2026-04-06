package marketplace

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/marketplace"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleProductRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/marketplace/products")
	path = strings.Trim(path, "/")

	if path == "" {
		if r.Method == http.MethodGet {
			m.handleListCatalog(w, r)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	slug := strings.Split(path, "/")[0]
	if r.Method == http.MethodGet {
		m.handleGetProduct(w, r, slug)
	} else {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListCatalog(w http.ResponseWriter, r *http.Request) {
	result, err := m.service.ListCatalog(r.Context(), marketplace.CatalogFilter{
		Search:       strings.TrimSpace(r.URL.Query().Get("search")),
		Category:     strings.TrimSpace(r.URL.Query().Get("category")),
		Condition:    strings.TrimSpace(r.URL.Query().Get("condition")),
		Status:       strings.TrimSpace(r.URL.Query().Get("status")),
		FeaturedOnly: strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("featured")), "true"),
	})
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, result)
}

func (m *Module) handleGetProduct(w http.ResponseWriter, r *http.Request, slug string) {
	product, err := m.service.GetProductBySlug(r.Context(), slug)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "MARKETPLACE_404", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, product)
}

func (m *Module) handleOrderRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	var payload marketplace.CreateOrderInput
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
		return
	}

	order, err := m.service.CreateOrder(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "MARKETPLACE_400", err.Error())
		return
	}

	// Broadcast order event if needed (optional for marketplace)
	// raw, _ := toMap(order)
	// m.broadcast("orders", "created", order.ID, raw, nil)

	httputil.Success(w, http.StatusCreated, order)
}
