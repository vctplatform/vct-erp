package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	financeusecase "vct-platform/backend/internal/modules/finance/usecase"
)

// CaptureService is the finance application boundary used by the HTTP adapter.
type CaptureService interface {
	Capture(ctx context.Context, req financedomain.CaptureRequest) (*financedomain.CaptureResult, error)
}

// VoidService is the privileged finance boundary used by the void endpoint.
type VoidService interface {
	VoidEntry(ctx context.Context, entryID string) error
}

// Handler exposes HTTP endpoints for the finance module.
type Handler struct {
	captureUC         CaptureService
	voidUC            VoidService
	idempotencyHeader string
}

// NewHandler constructs the finance HTTP adapter.
func NewHandler(captureUC CaptureService, voidUC VoidService, idempotencyHeader string) *Handler {
	if strings.TrimSpace(idempotencyHeader) == "" {
		idempotencyHeader = "Idempotency-Key"
	}

	return &Handler{
		captureUC:         captureUC,
		voidUC:            voidUC,
		idempotencyHeader: idempotencyHeader,
	}
}

// Capture receives financial events from subsidiary units and dispatches them through the idempotent finance capture use case.
func (h *Handler) Capture(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.captureUC == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error":   "service_not_wired",
			"message": "finance capture use case has not been composed yet",
		})
		return
	}

	var req financedomain.CaptureRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	req.IdempotencyKey = strings.TrimSpace(r.Header.Get(h.idempotencyHeader))
	result, err := h.captureUC.Capture(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case financeusecase.IsCaptureConflict(err):
			status = http.StatusConflict
		case financeusecase.IsCaptureValidationError(err):
			status = http.StatusUnprocessableEntity
		}

		writeJSON(w, status, map[string]string{
			"error":   "finance_capture_failed",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// VoidJournalEntry is reserved for chief accountants and delegates to a privileged void use case.
func (h *Handler) VoidJournalEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
		return
	}
	if h.voidUC == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error":   "service_not_wired",
			"message": "finance void use case has not been composed yet",
		})
		return
	}

	entryID := pathEntryID(r.URL.Path)
	if entryID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_request",
			"message": "journal entry id is required",
		})
		return
	}

	if err := h.voidUC.VoidEntry(r.Context(), entryID); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, financedomain.ErrUnsupportedOperation) {
			status = http.StatusNotImplemented
		}
		writeJSON(w, status, map[string]string{
			"error":   "finance_void_failed",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":   "accepted",
		"entry_id": entryID,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func pathEntryID(path string) string {
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 4 {
		return ""
	}
	return parts[len(parts)-2]
}
