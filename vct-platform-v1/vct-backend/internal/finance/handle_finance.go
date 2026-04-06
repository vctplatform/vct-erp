package finance

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain/finance"
	"vct-platform/backend/internal/shared/auth"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleTransactionRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok || !authz.CanEntityAction(p.User.Role, "finance", authz.ActionView) {
		httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem giao dịch tài chính")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/finance/transactions")
	path = strings.Trim(path, "/")

	if path == "" {
		if r.Method == http.MethodGet {
			m.handleListTransactions(w, r)
		} else if r.Method == http.MethodPost {
			m.handleCreateTransaction(w, r)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	httputil.Error(w, http.StatusNotFound, "FINANCE_404", "Không tìm thấy giao dịch")
}

func (m *Module) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	list, err := m.service.ListTransactions(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, list)
}

func (m *Module) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var payload finance.Transaction
	if err := httputil.DecodeJSON(r, &payload); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FINANCE_400", err.Error())
		return
	}

	created, err := m.service.CreateTransaction(r.Context(), payload)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "FINANCE_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleSubscriptionRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu xác thực")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/finance/subscriptions")
	path = strings.Trim(path, "/")

	if path == "" {
		if r.Method == http.MethodGet {
			m.handleListSubscriptions(w, r, p)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	httputil.Error(w, http.StatusNotFound, "FINANCE_404", "Không tìm thấy gói đăng ký")
}

func (m *Module) handleListSubscriptions(w http.ResponseWriter, r *http.Request, p auth.Principal) {
	filter := finance.SubscriptionFilter{
		EntityType: r.URL.Query().Get("entity_type"),
	}
	list, err := m.subscription.ListSubscriptions(r.Context(), filter)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, list)
}

func (m *Module) handleExpiringSubscriptions(w http.ResponseWriter, r *http.Request) {
	list, err := m.subscription.ListExpiringSubscriptions(r.Context(), 30) // 30 days
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, list)
}

// ── Plan Handlers ────────────────────────────────────────────

func (m *Module) handlePlanRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/finance/plans")
	id := strings.Trim(path, "/")

	if id == "" {
		if r.Method == http.MethodGet {
			list, err := m.subscription.ListPlans(r.Context(), r.URL.Query().Get("entity_type"))
			if err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, list)
		} else {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	// Handle by ID if needed (GetPlan, UpdatePlan, DeletePlan)
	httputil.Error(w, http.StatusNotFound, "FINANCE_404", "Không tìm thấy gói")
}

// ── Billing Cycle Handlers ───────────────────────────────────

func (m *Module) handleBillingCycleRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/finance/billing-cycles")
	id := strings.Trim(path, "/")

	if id == "" {
		if r.Method == http.MethodGet {
			filter := finance.BillingCycleFilter{
				SubscriptionID: r.URL.Query().Get("subscription_id"),
			}
			list, err := m.subscription.ListAllBillingCycles(r.Context(), filter)
			if err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, list)
		}
		return
	}

	// Mark paid if /pay suffix
	if strings.HasSuffix(path, "/pay") {
		p, _ := httputil.GetPrincipal(r)
		bcID := strings.TrimSuffix(id, "/pay")
		updated, err := m.subscription.MarkBillingCyclePaid(r.Context(), bcID, p.User.ID)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "FINANCE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, updated)
		return
	}

	httputil.Error(w, http.StatusNotFound, "FINANCE_404", "Không tìm thấy billing cycle")
}

// ── Placeholder routes for invoices/payments/etc ─────────────

func (m *Module) handleInvoiceRoutes(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{"module": "finance", "feature": "invoices"})
}

func (m *Module) handlePaymentRoutes(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{"module": "finance", "feature": "payments"})
}

func (m *Module) handleFeeScheduleRoutes(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{"module": "finance", "feature": "fee-schedules"})
}

func (m *Module) handleBudgetRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		list, err := m.service.ListBudgets(r.Context())
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, list)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"module": "finance", "feature": "budgets"})
}

func (m *Module) handleSponsorshipRoutes(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{"module": "finance", "feature": "sponsorships"})
}
