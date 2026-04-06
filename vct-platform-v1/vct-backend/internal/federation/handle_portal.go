package federation

import (
	"net/http"
	"time"

	"vct-platform/backend/internal/shared/httputil"
)

// ActivityItem represents a single item in the portal activity feed.
type ActivityItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"` // alert, update, match, etc.
}

// handlePortalActivities handles GET /api/v1/portal/activities
func (m *Module) handlePortalActivities(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.authenticate(w, r); !ok {
		return
	}

	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	// Mocked activity feed for now
	feed := []ActivityItem{
		{
			ID:          "ACT-1",
			Title:       "Cảnh báo hệ thống",
			Description: "Hãy thiết lập mã bảo vệ tài khoản (2FA) để đảm bảo an toàn.",
			Timestamp:   time.Now().Add(-1 * time.Hour),
			Type:        "alert",
		},
		{
			ID:          "ACT-2",
			Title:       "Truy cập gần đây",
			Description: "Bạn vừa đăng nhập cách đây 1 giờ trên trình duyệt mới.",
			Timestamp:   time.Now().Add(-1*time.Hour - 5*time.Minute),
			Type:        "update",
		},
		{
			ID:          "ACT-3",
			Title:       "Chào mừng trở lại!",
			Description: "Hệ thống VCT Platform v3 đã cập nhật thành công.",
			Timestamp:   time.Now().Add(-24 * time.Hour),
			Type:        "match",
		},
	}

	httputil.Success(w, http.StatusOK, map[string]any{
		"items": feed,
	})
}
