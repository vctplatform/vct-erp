package scoring

// ── Feature: Get Match State ─────────────────────────────────
// Complexity: SIMPLE — query only, no mutation

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleGetState(w http.ResponseWriter, r *http.Request, matchID string) {
	state, err := m.service.BuildCombatState(r.Context(), matchID)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, state)
}
