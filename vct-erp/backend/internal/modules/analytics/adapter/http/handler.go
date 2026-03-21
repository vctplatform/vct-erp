package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
)

// RevenueStreamService is consumed by the analytics HTTP adapter.
type RevenueStreamService interface {
	RevenueStream(ctx context.Context, companyCode string, from time.Time, to time.Time) ([]analyticsdomain.RevenueStreamPoint, error)
}

// CashRunwayService is consumed by the analytics HTTP adapter.
type CashRunwayService interface {
	CashRunway(ctx context.Context, companyCode string, asOf time.Time, months int) (analyticsdomain.CashRunway, error)
}

// Handler exposes dashboard-ready JSON endpoints.
type Handler struct {
	revenueUC RevenueStreamService
	runwayUC  CashRunwayService
}

// NewHandler constructs the analytics HTTP adapter.
func NewHandler(revenueUC RevenueStreamService, runwayUC CashRunwayService) *Handler {
	return &Handler{
		revenueUC: revenueUC,
		runwayUC:  runwayUC,
	}
}

// RevenueStream returns net revenue grouped by business cost center.
func (h *Handler) RevenueStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.revenueUC == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "service_not_wired"})
		return
	}

	query := r.URL.Query()
	points, err := h.revenueUC.RevenueStream(
		r.Context(),
		firstNonEmpty(query.Get("company_code"), "VCT_SIM"),
		parseDate(query.Get("date_from")),
		parseDate(query.Get("date_to")),
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "analytics_failed",
			"message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": points})
}

// CashRunway returns a 3-month contracted runway projection.
func (h *Handler) CashRunway(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.runwayUC == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "service_not_wired"})
		return
	}

	query := r.URL.Query()
	result, err := h.runwayUC.CashRunway(
		r.Context(),
		firstNonEmpty(query.Get("company_code"), "VCT_SIM"),
		parseDateTime(query.Get("as_of")),
		3,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "analytics_failed",
			"message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseDate(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func parseDateTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func firstNonEmpty(left string, right string) string {
	if left != "" {
		return left
	}
	return right
}
