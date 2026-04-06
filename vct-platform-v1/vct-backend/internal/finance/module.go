package finance

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/finance"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Finance module.
type Module struct {
	service      *finance.Service
	subscription *finance.SubscriptionService
	logger       *slog.Logger
}

// Deps holds the dependencies for the Finance module.
type Deps struct {
	Service      *finance.Service
	Subscription *finance.SubscriptionService
	Logger       *slog.Logger
}

// New creates a new Finance module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:      deps.Service,
		subscription: deps.Subscription,
		logger:       deps.Logger.With(slog.String("module", "finance")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers finance routes on the mux.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/finance/transactions", m.handleTransactionRoutes)
	mux.HandleFunc("/api/v1/finance/transactions/", m.handleTransactionRoutes)
	mux.HandleFunc("/api/v1/finance/subscriptions", m.handleSubscriptionRoutes)
	mux.HandleFunc("/api/v1/finance/subscriptions/", m.handleSubscriptionRoutes)
	mux.HandleFunc("/api/v1/finance/subscriptions/expiring", m.handleExpiringSubscriptions)

	mux.HandleFunc("/api/v1/finance/plans", m.handlePlanRoutes)
	mux.HandleFunc("/api/v1/finance/plans/", m.handlePlanRoutes)

	mux.HandleFunc("/api/v1/finance/billing-cycles", m.handleBillingCycleRoutes)
	mux.HandleFunc("/api/v1/finance/billing-cycles/", m.handleBillingCycleRoutes)

	mux.HandleFunc("/api/v1/finance/invoices", m.handleInvoiceRoutes)
	mux.HandleFunc("/api/v1/finance/invoices/", m.handleInvoiceRoutes)

	mux.HandleFunc("/api/v1/finance/payments", m.handlePaymentRoutes)
	mux.HandleFunc("/api/v1/finance/payments/", m.handlePaymentRoutes)

	mux.HandleFunc("/api/v1/finance/fee-schedules", m.handleFeeScheduleRoutes)
	mux.HandleFunc("/api/v1/finance/budgets", m.handleBudgetRoutes)
	mux.HandleFunc("/api/v1/finance/sponsorships", m.handleSponsorshipRoutes)
	mux.HandleFunc("/api/v1/finance/sponsorships/", m.handleSponsorshipRoutes)

	m.logger.Info("finance module routes registered")
}
