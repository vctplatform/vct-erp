package scoring

// ── Feature: Start Match ─────────────────────────────────────
// Complexity: MEDIUM — validate + write event + broadcast

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleStartMatch(w http.ResponseWriter, r *http.Request, matchID, userID string) {
	if err := m.service.StartCombatMatch(r.Context(), matchID, userID); err != nil {
		httputil.WriteError(w, err)
		return
	}

	m.broadcast("combat_matches", "match_started", matchID, map[string]any{
		"id":      matchID,
		"status":  "dang_dau",
		"user_id": userID,
	}, nil)

	httputil.JSON(w, http.StatusOK, map[string]any{
		"match_id": matchID,
		"status":   "dang_dau",
		"message":  "Trận đấu đã bắt đầu",
	})
}
