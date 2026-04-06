package provincial

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/federation"
	"vct-platform/backend/internal/shared/httputil"
)

// ── Provincial Reports ───────────────────────────────────────

func (m *Module) handleProvReportRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/reports")
	id := strings.TrimPrefix(path, "/")

	switch {
	case r.Method == "GET" && id == "":
		provinceID := r.URL.Query().Get("province_id")
		reports, err := m.fedService.ListProvincialReports(r.Context(), provinceID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"reports": reports, "total": len(reports)})

	case r.Method == "POST" && id == "":
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		var report federation.ProvincialReport
		if err := httputil.DecodeJSON(r, &report); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", "invalid JSON: "+err.Error())
			return
		}
		created, err := m.fedService.CreateProvincialReport(r.Context(), report)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Provincial Statistics ────────────────────────────────────

func (m *Module) handleProvStats(w http.ResponseWriter, r *http.Request) {
	provinceID := r.URL.Query().Get("province_id")
	if provinceID == "" {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", "province_id is required")
		return
	}
	stats, err := m.fedService.GetProvincialStatistics(r.Context(), provinceID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}

// ── Implementation of other Federation-level Provincial Managers ──
// These were in federation_provincial_handler.go and are used by Federation admins.

func (m *Module) handleFederationProvClub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/clubs")
	id := strings.TrimPrefix(path, "/")

	switch {
	case r.Method == "DELETE" && id != "":
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		if err := m.fedService.DeleteProvincialClub(r.Context(), id); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"status": "deleted"})
	default:
		// Fallback to main handler for GET/POST
		m.handleProvincialClubs(w, r)
	}
}

// handleProvincialReferees, handleProvincialCommittee, handleProvincialTransfers
// ported/stubbed to satisfy RegisterRoutes and resolve lint errors.

func (m *Module) handleProvincialReferees(w http.ResponseWriter, r *http.Request) {
	httputil.Error(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Referees handler to be fully migrated")
}

func (m *Module) handleProvincialCommittee(w http.ResponseWriter, r *http.Request) {
	httputil.Error(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Committee handler to be fully migrated")
}

func (m *Module) handleProvincialTransfers(w http.ResponseWriter, r *http.Request) {
	httputil.Error(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Transfers handler to be fully migrated")
}
