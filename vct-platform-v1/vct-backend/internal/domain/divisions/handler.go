package divisions

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handler exposes REST endpoints for administrative divisions.
type Handler struct {
	store *Store
}

// NewHandler creates a new divisions handler.
func NewHandler() *Handler {
	return &Handler{store: Default()}
}

// RegisterRoutes mounts the divisions endpoints onto the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/divisions/provinces", h.handleProvinces)
	mux.HandleFunc("/api/v1/divisions/provinces/", h.handleProvinceWards)
}

// GET /api/v1/divisions/provinces?q=...
func (h *Handler) handleProvinces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")
	var provinces []ProvinceInfo
	if q != "" {
		provinces = h.store.SearchProvinces(q)
	} else {
		provinces = h.store.Provinces()
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"count": len(provinces),
		"data":  provinces,
	})
}

// GET /api/v1/divisions/provinces/{code}/wards?q=...
func (h *Handler) handleProvinceWards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parse province code from URL: /api/v1/divisions/provinces/{code}/wards
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/divisions/provinces/")
	parts := strings.Split(path, "/")
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, `{"error":"missing province code"}`, http.StatusBadRequest)
		return
	}

	code, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, `{"error":"invalid province code"}`, http.StatusBadRequest)
		return
	}

	// If path is just /provinces/{code} — return province info
	if len(parts) == 1 || (len(parts) == 2 && parts[1] == "") {
		p := h.store.Province(code)
		if p == nil {
			http.Error(w, `{"error":"province not found"}`, http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, ProvinceInfo{
			Name:         p.Name,
			Code:         p.Code,
			DivisionType: p.DivisionType,
			Codename:     p.Codename,
			PhoneCode:    p.PhoneCode,
			WardCount:    len(p.Wards),
		})
		return
	}

	// /provinces/{code}/wards
	if parts[1] != "wards" {
		http.Error(w, `{"error":"unknown sub-resource"}`, http.StatusNotFound)
		return
	}

	p := h.store.Province(code)
	if p == nil {
		http.Error(w, `{"error":"province not found"}`, http.StatusNotFound)
		return
	}

	q := r.URL.Query().Get("q")
	var wards []Ward
	if q != "" {
		wards = h.store.SearchWards(code, q)
	} else {
		wards = h.store.Wards(code)
	}
	if wards == nil {
		wards = []Ward{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"province": p.Name,
		"count":    len(wards),
		"data":     wards,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
