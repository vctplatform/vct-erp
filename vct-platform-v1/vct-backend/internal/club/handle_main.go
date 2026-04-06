package club

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleClubDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}
	d, err := m.service.GetDashboard(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, d)
}
