package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	analyticsusecase "vct-platform/backend/internal/modules/analytics/usecase"
	sharedmiddleware "vct-platform/backend/internal/shared/middleware"
)

// RevenueStreamService is consumed by the analytics HTTP adapter.
type RevenueStreamService interface {
	RevenueStream(ctx context.Context, companyCode string, from time.Time, to time.Time) ([]analyticsdomain.RevenueStreamPoint, error)
}

// CashRunwayService is consumed by the analytics HTTP adapter.
type CashRunwayService interface {
	CashRunway(ctx context.Context, companyCode string, asOf time.Time, months int) (analyticsdomain.CashRunway, error)
}

// FinanceSummaryService is consumed by the executive dashboard HTTP adapter.
type FinanceSummaryService interface {
	FinanceSummary(ctx context.Context, access analyticsusecase.AccessMetadata) (analyticsdomain.FinanceSummary, error)
}

// SegmentService is consumed by the executive dashboard HTTP adapter.
type SegmentService interface {
	SegmentProfit(ctx context.Context, access analyticsusecase.AccessMetadata) ([]analyticsdomain.SegmentGrossProfit, error)
}

// DashboardCashRunwayService is consumed by the executive dashboard HTTP adapter.
type DashboardCashRunwayService interface {
	DashboardCashRunway(ctx context.Context, input analyticsusecase.CashRunwayInput) (analyticsdomain.CashRunway, error)
}

// DashboardService is consumed by the live command-center endpoint.
type DashboardService interface {
	Dashboard(ctx context.Context, input analyticsusecase.DashboardInput) (analyticsdomain.CommandCenterDashboardData, error)
}

// DashboardCardsService is consumed by the KPI-card endpoint.
type DashboardCardsService interface {
	DashboardCards(ctx context.Context, input analyticsusecase.DashboardInput) (analyticsdomain.DashboardCardsResponse, error)
}

// DashboardSegmentsService is consumed by the revenue-mix endpoint.
type DashboardSegmentsService interface {
	DashboardSegments(ctx context.Context, input analyticsusecase.DashboardInput) (analyticsdomain.DashboardSegmentsResponse, error)
}

// DashboardCashflowService is consumed by the live chart endpoint.
type DashboardCashflowService interface {
	DashboardCashflow(ctx context.Context, input analyticsusecase.DashboardInput) (analyticsdomain.DashboardCashflowResponse, error)
}

// Handler exposes dashboard-ready JSON endpoints.
type Handler struct {
	revenueUC       RevenueStreamService
	runwayUC        CashRunwayService
	summaryUC       FinanceSummaryService
	segmentsUC      SegmentService
	dashboardRunway DashboardCashRunwayService
	dashboardUC     DashboardService
	dashboardCards  DashboardCardsService
	dashboardMix    DashboardSegmentsService
	dashboardChart  DashboardCashflowService
	actorHeader     string
}

// NewHandler constructs the analytics HTTP adapter.
func NewHandler(
	revenueUC RevenueStreamService,
	runwayUC CashRunwayService,
	summaryUC FinanceSummaryService,
	segmentsUC SegmentService,
	dashboardRunway DashboardCashRunwayService,
	dashboardUC DashboardService,
	dashboardCards DashboardCardsService,
	dashboardMix DashboardSegmentsService,
	dashboardChart DashboardCashflowService,
	actorHeader string,
) *Handler {
	if strings.TrimSpace(actorHeader) == "" {
		actorHeader = "X-Actor-ID"
	}

	return &Handler{
		revenueUC:       revenueUC,
		runwayUC:        runwayUC,
		summaryUC:       summaryUC,
		segmentsUC:      segmentsUC,
		dashboardRunway: dashboardRunway,
		dashboardUC:     dashboardUC,
		dashboardCards:  dashboardCards,
		dashboardMix:    dashboardMix,
		dashboardChart:  dashboardChart,
		actorHeader:     actorHeader,
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

// Summary returns the top-level finance dashboard cards for CFO, CEO, and system administrators.
func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.summaryUC == nil {
		if h.dashboardCards == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "service_not_wired"})
			return
		}
	}

	access, err := h.accessMetadata(r, "VCT_GROUP")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	if h.dashboardCards != nil {
		result, err := h.dashboardCards.DashboardCards(r.Context(), analyticsusecase.DashboardInput{
			Access: access,
			AsOf:   parseDateTime(r.URL.Query().Get("as_of")),
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":   "analytics_failed",
				"message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}

	result, err := h.summaryUC.FinanceSummary(r.Context(), access)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "analytics_failed",
			"message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// Segments returns the gross profit structure used by pie and stacked bar charts.
func (h *Handler) Segments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.segmentsUC == nil {
		if h.dashboardMix == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "service_not_wired"})
			return
		}
	}

	access, err := h.accessMetadata(r, "VCT_GROUP")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	if h.dashboardMix != nil {
		result, err := h.dashboardMix.DashboardSegments(r.Context(), analyticsusecase.DashboardInput{
			Access: access,
			AsOf:   parseDateTime(r.URL.Query().Get("as_of")),
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":   "analytics_failed",
				"message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}

	items, err := h.segmentsUC.SegmentProfit(r.Context(), access)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "analytics_failed",
			"message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
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

// FinanceCashRunway returns a 6-month dashboard projection and records audit access.
func (h *Handler) FinanceCashRunway(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.dashboardRunway == nil {
		if h.dashboardChart == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "service_not_wired"})
			return
		}
	}

	access, err := h.accessMetadata(r, "VCT_GROUP")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	if h.dashboardChart != nil {
		result, err := h.dashboardChart.DashboardCashflow(r.Context(), analyticsusecase.DashboardInput{
			Access: access,
			AsOf:   parseDateTime(r.URL.Query().Get("as_of")),
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":   "analytics_failed",
				"message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, result)
		return
	}

	result, err := h.dashboardRunway.DashboardCashRunway(r.Context(), analyticsusecase.CashRunwayInput{
		Access: access,
		AsOf:   parseDateTime(r.URL.Query().Get("as_of")),
		Months: parseMonths(r.URL.Query().Get("months"), 6),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "analytics_failed",
			"message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// Dashboard returns the live plug-and-play dashboard payload for the executive Command Center.
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.dashboardUC == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "service_not_wired"})
		return
	}

	access, err := h.accessMetadata(r, "VCT_GROUP")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	result, err := h.dashboardUC.Dashboard(r.Context(), analyticsusecase.DashboardInput{
		Access: access,
		AsOf:   parseDateTime(r.URL.Query().Get("as_of")),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "analytics_failed",
			"message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// DashboardMock returns a plug-and-play dashboard payload for frontend development.
func (h *Handler) DashboardMock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	writeJSON(w, http.StatusOK, analyticsusecase.GetMockDashboardData())
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

func parseMonths(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	var months int
	if _, err := fmt.Sscanf(value, "%d", &months); err != nil || months <= 0 {
		return fallback
	}
	return months
}

func (h *Handler) accessMetadata(r *http.Request, fallbackCompany string) (analyticsusecase.AccessMetadata, error) {
	actorID := strings.TrimSpace(r.Header.Get(h.actorHeader))
	if actorID == "" {
		return analyticsusecase.AccessMetadata{}, fmt.Errorf("%s header is required", h.actorHeader)
	}

	filters := make(map[string]string, len(r.URL.Query()))
	for key, values := range r.URL.Query() {
		if len(values) == 0 {
			continue
		}
		filters[key] = values[0]
	}

	return analyticsusecase.AccessMetadata{
		CompanyCode: firstNonEmpty(strings.TrimSpace(r.URL.Query().Get("company_code")), fallbackCompany),
		ActorID:     actorID,
		ActorRole:   sharedmiddleware.RoleFromContext(r.Context()),
		IPAddress:   clientIP(r),
		UserAgent:   strings.TrimSpace(r.UserAgent()),
		Filters:     filters,
	}, nil
}

func clientIP(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
