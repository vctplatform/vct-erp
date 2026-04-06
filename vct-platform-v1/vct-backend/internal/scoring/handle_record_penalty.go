package scoring

// ── Feature: Record Penalty ──────────────────────────────────
// Complexity: MEDIUM — validate + event + broadcast

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

type recordPenaltyRequest struct {
	Round     int     `json:"round"`
	Corner    string  `json:"corner"`
	Deduction float64 `json:"deduction"`
	Reason    string  `json:"reason"`
}

func (m *Module) handleRecordPenalty(w http.ResponseWriter, r *http.Request, matchID, userID string) {
	var req recordPenaltyRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", err.Error())
		return
	}

	if err := m.service.RecordPenalty(r.Context(), matchID, userID, req.Round, req.Corner, req.Deduction, req.Reason); err != nil {
		httputil.WriteError(w, err)
		return
	}

	m.broadcast("combat_matches", "penalty", matchID, map[string]any{
		"round":     req.Round,
		"corner":    req.Corner,
		"deduction": req.Deduction,
		"reason":    req.Reason,
		"user_id":   userID,
	}, nil)

	httputil.JSON(w, http.StatusOK, map[string]any{
		"match_id":  matchID,
		"corner":    req.Corner,
		"deduction": req.Deduction,
		"message":   "Đã ghi nhận lỗi/phạt",
	})
}
