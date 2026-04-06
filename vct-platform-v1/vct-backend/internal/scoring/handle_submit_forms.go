package scoring

// ── Feature: Submit Forms Score ──────────────────────────────
// Complexity: MEDIUM — validate + event + broadcast

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

type submitFormsScoreRequest struct {
	RefereeID string  `json:"referee_id"`
	AthleteID string  `json:"athlete_id"`
	Score     float64 `json:"score"`
}

func (m *Module) handleSubmitFormsScore(w http.ResponseWriter, r *http.Request, perfID, userID string) {
	var req submitFormsScoreRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", err.Error())
		return
	}

	// Validation
	if req.Score < 0 || req.Score > 10 {
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", "điểm phải từ 0 đến 10")
		return
	}

	if err := m.service.SubmitFormsScore(r.Context(), perfID, req.RefereeID, req.AthleteID, req.Score); err != nil {
		httputil.WriteError(w, err)
		return
	}

	m.broadcast("form_performances", "judge_scored", perfID, map[string]any{
		"referee_id": req.RefereeID,
		"athlete_id": req.AthleteID,
		"score":      req.Score,
		"user_id":    userID,
	}, nil)

	httputil.JSON(w, http.StatusOK, map[string]any{
		"performance_id": perfID,
		"referee_id":     req.RefereeID,
		"score":          req.Score,
		"message":        "Đã ghi nhận điểm giám khảo",
	})
}
