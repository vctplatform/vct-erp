package federation

import (
	"context"
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/approval"
	"vct-platform/backend/internal/shared/httputil"
)

// handleApprovalCRUD handles basic CRUD for approval requests.
func (m *Module) handleApprovalCRUD(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/approvals/")
	parts := strings.Split(path, "/")

	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	// POST /api/v1/approvals/ → submit new request
	if len(parts) == 1 && parts[0] == "" && r.Method == http.MethodPost {
		m.handleApprovalSubmit(w, r, userID)
		return
	}

	if len(parts) < 1 || parts[0] == "" {
		httputil.Error(w, http.StatusNotFound, "APPROVAL_404", "Không tìm thấy tài nguyên")
		return
	}
	requestID := parts[0]

	// GET /api/v1/approvals/:id
	if len(parts) == 1 && r.Method == http.MethodGet {
		m.handleApprovalGet(w, r, requestID)
		return
	}

	// GET /api/v1/approvals/:id/steps
	if len(parts) == 2 && parts[1] == "steps" && r.Method == http.MethodGet {
		m.handleApprovalSteps(w, r, requestID)
		return
	}

	// GET /api/v1/approvals/:id/history
	if len(parts) == 2 && parts[1] == "history" && r.Method == http.MethodGet {
		m.handleApprovalHistory(w, r, requestID)
		return
	}

	// POST /api/v1/approvals/:id/approve|reject|return|cancel
	if len(parts) == 2 && r.Method == http.MethodPost {
		action := parts[1]
		m.handleApprovalAction(w, r, userID, requestID, action)
		return
	}

	httputil.Error(w, http.StatusNotFound, "APPROVAL_404", "Không tìm thấy tài nguyên")
}

func (m *Module) handleApprovalSubmit(w http.ResponseWriter, r *http.Request, userID string) {
	var input approval.SubmitInput
	if err := httputil.DecodeJSON(r, &input); err != nil {
		httputil.Error(w, http.StatusBadRequest, "APPROVAL_400", err.Error())
		return
	}
	if input.WorkflowCode == "" || input.Title == "" {
		httputil.Error(w, http.StatusBadRequest, "APPROVAL_400", "workflow_code và title là bắt buộc")
		return
	}
	input.RequestedBy = userID

	created, err := m.approval.Submit(r.Context(), input)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "APPROVAL_400", err.Error())
		return
	}

	// Note: event broadcasting is handled within the domain service or here
	m.broadcaster.BroadcastEntityChange("approval", "submitted", created.ID, map[string]any{"requested_by": userID}, nil)

	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleApprovalGet(w http.ResponseWriter, r *http.Request, requestID string) {
	req, err := m.approval.GetRequest(r.Context(), requestID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "APPROVAL_404", "approval request not found")
		return
	}
	httputil.Success(w, http.StatusOK, req)
}

func (m *Module) handleApprovalSteps(w http.ResponseWriter, r *http.Request, requestID string) {
	steps, err := m.approval.GetSteps(r.Context(), requestID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"request_id": requestID, "steps": steps})
}

func (m *Module) handleApprovalHistory(w http.ResponseWriter, r *http.Request, requestID string) {
	history, err := m.approval.GetHistory(r.Context(), requestID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"request_id": requestID, "history": history})
}

func (m *Module) handleApprovalAction(w http.ResponseWriter, r *http.Request, userID, requestID, action string) {
	var body struct {
		Comment string `json:"comment"`
		Reason  string `json:"reason"`
	}
	_ = httputil.DecodeJSON(r, &body)

	var err error
	switch action {
	case "approve":
		err = m.approval.Approve(r.Context(), requestID, userID, body.Comment)
	case "reject":
		err = m.approval.Reject(r.Context(), requestID, userID, body.Reason)
	case "return":
		err = m.approval.Return(r.Context(), requestID, userID, body.Reason)
	case "cancel":
		err = m.approval.Cancel(r.Context(), requestID, userID)
	default:
		httputil.Error(w, http.StatusBadRequest, "APPROVAL_400", "unknown action: "+action)
		return
	}

	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "APPROVAL_400", err.Error())
		return
	}

	m.broadcaster.BroadcastEntityChange("approval", action, requestID, map[string]any{"actor_id": userID}, nil)

	httputil.Success(w, http.StatusOK, map[string]string{"status": action + "d", "request_id": requestID})
}

func (m *Module) handleApprovalMyPending(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	_, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	role := r.URL.Query().Get("role")
	// If no role specified, we might need a way to get user role.
	// In Vertical Slice, we can fetch it via main service or assume it's in the principal (if we passed principal).
	// For now, I'll assume role is passed or handled by domain.

	requests, err := m.approval.ListPendingForRole(r.Context(), role)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{
		"role":     role,
		"requests": requests,
		"total":    len(requests),
	})
}

func (m *Module) handleApprovalMyRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	requests, err := m.approval.ListMyRequests(r.Context(), userID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{
		"user_id":  userID,
		"requests": requests,
		"total":    len(requests),
	})
}

func (m *Module) handleWorkflowDefinitions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	workflows := m.getDefaultWorkflowDefinitions()
	httputil.Success(w, http.StatusOK, workflows)
}

// m.getDefaultWorkflowDefinitions ported from legacy approval_handler.go
func (m *Module) getDefaultWorkflowDefinitions() []map[string]any {
	return []map[string]any{
		{
			"code": "club_registration", "entity_type": "club",
			"display_name": "Đăng ký thành lập CLB / Võ đường",
			"steps": []map[string]any{
				{"step": 1, "name": "LĐ Tỉnh xem xét", "role": "provincial_admin"},
				{"step": 2, "name": "LĐ Quốc gia xác nhận", "role": "federation_secretary"},
			},
		},
		{
			"code": "member_registration", "entity_type": "member",
			"display_name": "Đăng ký hội viên vào CLB",
			"steps": []map[string]any{
				{"step": 1, "name": "HLV xác nhận", "role": "coach"},
			},
		},
		{
			"code": "referee_card", "entity_type": "referee",
			"display_name": "Cấp/Gia hạn Thẻ Trọng tài",
			"steps": []map[string]any{
				{"step": 1, "name": "GĐ Kỹ thuật kiểm tra", "role": "technical_director"},
				{"step": 2, "name": "LĐ Quốc gia ký", "role": "president"},
			},
		},
		{
			"code": "tournament_hosting", "entity_type": "tournament",
			"display_name": "Đăng ký Tổ chức Giải đấu",
			"steps": []map[string]any{
				{"step": 1, "name": "Thư ký LĐ thẩm định", "role": "federation_secretary"},
				{"step": 2, "name": "GĐ Kỹ thuật kiểm tra", "role": "technical_director"},
				{"step": 3, "name": "Chủ tịch LĐ phê duyệt", "role": "president"},
			},
		},
		{
			"code": "team_registration", "entity_type": "team",
			"display_name": "Đoàn Đăng ký Tham gia Giải",
			"steps": []map[string]any{
				{"step": 1, "name": "BTC giải duyệt", "role": "tournament_director"},
			},
		},
		{
			"code": "athlete_registration", "entity_type": "registration",
			"display_name": "VĐV Đăng ký Nội dung Thi đấu",
			"steps": []map[string]any{
				{"step": 1, "name": "Auto-validate + BTC duyệt", "role": "tournament_director"},
			},
		},
		{
			"code": "result_approval", "entity_type": "match_result",
			"display_name": "Phê duyệt Kết quả Thi đấu",
			"steps": []map[string]any{
				{"step": 1, "name": "Trọng tài chính xác nhận", "role": "chief_referee"},
				{"step": 2, "name": "BTC phê duyệt", "role": "tournament_director"},
			},
		},
		{
			"code": "expense_small", "entity_type": "transaction",
			"display_name": "Phê duyệt Chi tiêu (≤ 5 triệu)",
			"steps": []map[string]any{
				{"step": 1, "name": "Thư ký LĐ duyệt", "role": "federation_secretary"},
			},
		},
		{
			"code": "expense_medium", "entity_type": "transaction",
			"display_name": "Phê duyệt Chi tiêu (5-50 triệu)",
			"steps": []map[string]any{
				{"step": 1, "name": "Thư ký LĐ duyệt", "role": "federation_secretary"},
				{"step": 2, "name": "Chủ tịch LĐ phê duyệt", "role": "president"},
			},
		},
		{
			"code": "expense_large", "entity_type": "transaction",
			"display_name": "Phê duyệt Chi tiêu (> 50 triệu)",
			"steps": []map[string]any{
				{"step": 1, "name": "Thư ký LĐ duyệt", "role": "federation_secretary"},
				{"step": 2, "name": "Chủ tịch LĐ phê duyệt", "role": "president"},
				{"step": 3, "name": "Ban thường vụ (2/3 đồng ý)", "role": "executive_board", "requires_all": false, "min_approvals": 3},
			},
		},
		{
			"code": "fee_confirmation", "entity_type": "payment",
			"display_name": "Xác nhận Đóng Lệ phí Đoàn",
			"steps": []map[string]any{
				{"step": 1, "name": "Kế toán xác nhận", "role": "accountant"},
			},
		},
		{
			"code": "belt_promotion", "entity_type": "belt_exam",
			"display_name": "Thi Thăng Đai",
			"steps": []map[string]any{
				{"step": 1, "name": "Ban chuyên môn xét điều kiện", "role": "technical_director"},
				{"step": 2, "name": "LĐ cấp bằng", "role": "president"},
			},
		},
		{
			"code": "training_class", "entity_type": "training_class",
			"display_name": "Mở Lớp Đào tạo / Tập huấn",
			"steps": []map[string]any{
				{"step": 1, "name": "GĐ Kỹ thuật thẩm định", "role": "technical_director"},
				{"step": 2, "name": "LĐ phê duyệt", "role": "president"},
			},
		},
		{
			"code": "news_publish", "entity_type": "news",
			"display_name": "Phê duyệt Tin tức / Thông báo",
			"steps": []map[string]any{
				{"step": 1, "name": "Thư ký duyệt", "role": "federation_secretary"},
			},
		},
		{
			"code": "complaint", "entity_type": "complaint",
			"display_name": "Khiếu nại & Kháng nghị",
			"steps": []map[string]any{
				{"step": 1, "name": "Ban giải quyết KN xem xét", "role": "discipline_board"},
				{"step": 2, "name": "BTC ra quyết định", "role": "tournament_director"},
			},
		},
	}
}

// SeedDefaultWorkflows seeds the workflow definitions into the approval database.
func (m *Module) SeedDefaultWorkflows(ctx context.Context, repo approval.WorkflowRepository) {
	defs := m.getDefaultWorkflowDefinitions()
	for _, d := range defs {
		code, _ := d["code"].(string)
		entityType, _ := d["entity_type"].(string)
		displayName, _ := d["display_name"].(string)
		rawSteps, _ := d["steps"].([]map[string]any)

		var steps []approval.StepTemplate
		for _, rs := range rawSteps {
			stepNum, _ := rs["step"].(int)
			name, _ := rs["name"].(string)
			role, _ := rs["role"].(string)
			requiresAll, _ := rs["requires_all"].(bool)
			minApprovals, _ := rs["min_approvals"].(int)
			steps = append(steps, approval.StepTemplate{
				StepNumber:   stepNum,
				StepName:     name,
				ApproverRole: role,
				RequiresAll:  requiresAll,
				MinApprovals: minApprovals,
			})
		}

		wf := approval.WorkflowDefinition{
			WorkflowCode: code,
			EntityType:   entityType,
			DisplayName:  displayName,
			Steps:        steps,
			IsActive:     true,
		}
		_ = repo.Create(ctx, wf)
	}
}
