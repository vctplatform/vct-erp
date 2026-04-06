package scoring

// ── Feature: Submit Score ────────────────────────────────────
// Complexity: COMPLEX — validate input + domain rules + event + broadcast

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

// submitScoreRequest is the request payload for submitting a combat score.
type submitScoreRequest struct {
	Round  int     `json:"round"`
	Corner string  `json:"corner"` // "red" or "blue"
	Points float64 `json:"points"`
}

func (m *Module) handleSubmitScore(w http.ResponseWriter, r *http.Request, matchID, userID string) {
	var req submitScoreRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", err.Error())
		return
	}

	// Validation
	if req.Corner != "red" && req.Corner != "blue" {
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", "corner phải là 'red' hoặc 'blue'")
		return
	}
	if req.Round < 1 {
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", "round phải >= 1")
		return
	}

	// Domain operation
	if err := m.service.RecordCombatScore(r.Context(), matchID, userID, req.Round, req.Corner, req.Points); err != nil {
		httputil.WriteError(w, err)
		return
	}

	// Broadcast
	m.broadcast("combat_matches", "scored", matchID, map[string]any{
		"round":   req.Round,
		"corner":  req.Corner,
		"points":  req.Points,
		"user_id": userID,
	}, nil)

	httputil.JSON(w, http.StatusOK, map[string]any{
		"match_id": matchID,
		"round":    req.Round,
		"corner":   req.Corner,
		"points":   req.Points,
		"message":  "Đã ghi nhận điểm",
	})
}
