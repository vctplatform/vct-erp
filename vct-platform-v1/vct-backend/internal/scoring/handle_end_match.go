package scoring

// ── Feature: End Match ───────────────────────────────────────
// Complexity: COMPLEX — calculate result + event + broadcast

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleEndMatch(w http.ResponseWriter, r *http.Request, matchID, userID string) {
	result, err := m.service.EndCombatMatch(r.Context(), matchID, userID)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	m.broadcast("combat_matches", "match_ended", matchID, map[string]any{
		"id":      matchID,
		"status":  "ket_thuc",
		"winner":  result.Winner,
		"method":  result.Method,
		"user_id": userID,
	}, nil)

	httputil.JSON(w, http.StatusOK, map[string]any{
		"match_id": matchID,
		"status":   "ket_thuc",
		"winner":   result.Winner,
		"method":   result.Method,
		"message":  "Trận đấu đã kết thúc",
	})
}
