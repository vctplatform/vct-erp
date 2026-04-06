package scoring

// ── Feature: Finalize Forms Performance ──────────────────────
// Complexity: COMPLEX — calculation + result + event + broadcast

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleFinalizeForms(w http.ResponseWriter, r *http.Request, perfID, userID string) {
	result, err := m.service.FinalizeFormsPerformance(r.Context(), perfID)
	if err != nil {
		httputil.WriteError(w, err)
		return
	}

	m.broadcast("form_performances", "finalized", perfID, map[string]any{
		"id":          perfID,
		"status":      "da_cham",
		"final_score": result.FinalScore,
		"judge_count": result.JudgeCount,
		"user_id":     userID,
	}, nil)

	httputil.JSON(w, http.StatusOK, map[string]any{
		"performance_id": perfID,
		"status":         "da_cham",
		"final_score":    result.FinalScore,
		"judge_count":    result.JudgeCount,
		"message":        "Đã hoàn tất chấm điểm quyền",
	})
}
