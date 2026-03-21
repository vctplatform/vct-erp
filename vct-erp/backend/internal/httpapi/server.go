package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	analyticshhttp "vct-platform/backend/internal/modules/analytics/adapter/http"
	financehttp "vct-platform/backend/internal/modules/finance/adapter/http"
	"vct-platform/backend/internal/modules/ledger/domain"
	"vct-platform/backend/internal/modules/ledger/usecase"
	sharedmiddleware "vct-platform/backend/internal/shared/middleware"
)

// PostEntryService is the application boundary consumed by the transport layer.
type PostEntryService interface {
	PostEntry(ctx context.Context, req usecase.PostEntryRequest) (*usecase.PostEntryResult, error)
}

// Dependencies holds the HTTP-facing application services.
type Dependencies struct {
	PostEntryUC       PostEntryService
	FinanceCaptureUC  financehttp.CaptureService
	FinanceVoidUC     financehttp.VoidService
	AnalyticsRevenue  analyticshhttp.RevenueStreamService
	AnalyticsRunway   analyticshhttp.CashRunwayService
	AppRoleHeader     string
	AppActorHeader    string
	IdempotencyHeader string
}

// Server wires the HTTP routes for the core ledger module.
type Server struct {
	mux           *http.ServeMux
	handler       http.Handler
	postEntryUC   PostEntryService
	financeHandle *financehttp.Handler
	analytics     *analyticshhttp.Handler
}

// New constructs the HTTP server with health and posting endpoints.
func New(deps Dependencies) *Server {
	server := &Server{
		mux:         http.NewServeMux(),
		postEntryUC: deps.PostEntryUC,
		financeHandle: financehttp.NewHandler(
			deps.FinanceCaptureUC,
			deps.FinanceVoidUC,
			deps.IdempotencyHeader,
			deps.AppActorHeader,
		),
		analytics: analyticshhttp.NewHandler(deps.AnalyticsRevenue, deps.AnalyticsRunway),
	}

	server.mux.HandleFunc("/healthz", server.handleHealth)
	server.mux.HandleFunc("/api/v1/ledger/journal-entries", server.handleJournalEntries)
	server.mux.HandleFunc("/v1/finance/capture", server.financeHandle.Capture)
	server.mux.HandleFunc("/v1/analytics/revenue-stream", server.analytics.RevenueStream)
	server.mux.HandleFunc("/v1/analytics/cash-runway", server.analytics.CashRunway)
	server.mux.Handle(
		"/v1/finance/journal-entries/",
		sharedmiddleware.RequireRoles("chief_accountant")(
			http.HandlerFunc(server.financeHandle.VoidJournalEntry),
		),
	)
	server.handler = sharedmiddleware.SecurityHeaders(
		sharedmiddleware.WithRoleFromHeader(roleHeaderOrDefault(deps.AppRoleHeader))(server.mux),
	)
	return server
}

// ServeHTTP delegates to the configured routes.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func (s *Server) handleJournalEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method_not_allowed",
		})
		return
	}

	if s.postEntryUC == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error":   "service_not_wired",
			"message": "ledger posting use case has not been composed in cmd/api yet",
		})
		return
	}

	var req usecase.PostEntryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	result, err := s.postEntryUC.PostEntry(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case usecase.IsValidationError(err):
			status = http.StatusUnprocessableEntity
		case errors.Is(err, domain.ErrAccountNotFound):
			status = http.StatusNotFound
		case errors.Is(err, domain.ErrAccountNotPostable), errors.Is(err, domain.ErrAccountInactive):
			status = http.StatusConflict
		}

		writeJSON(w, status, map[string]string{
			"error":   "post_entry_failed",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func roleHeaderOrDefault(header string) string {
	if header == "" {
		return "X-App-Role"
	}
	return header
}
