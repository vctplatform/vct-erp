package usecase

import (
	"context"
	"fmt"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
)

// Service exposes dashboard-ready analytics responses.
type Service struct {
	repo           analyticsdomain.Repository
	dashboardCache analyticsdomain.DashboardCache
	dashboardTTL   time.Duration
	now            func() time.Time
}

// NewService constructs the analytics application service.
func NewService(repo analyticsdomain.Repository, options ...ServiceOption) *Service {
	service := &Service{
		repo: repo,
		now:  time.Now,
	}

	for _, option := range options {
		if option != nil {
			option(service)
		}
	}

	return service
}

// RevenueStream returns net revenue by business cost center.
func (s *Service) RevenueStream(ctx context.Context, companyCode string, from time.Time, to time.Time) ([]analyticsdomain.RevenueStreamPoint, error) {
	if from.IsZero() {
		from = time.Date(s.now().UTC().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if to.IsZero() {
		to = s.now().UTC()
	}
	return s.repo.RevenueStream(ctx, companyCode, from.UTC(), to.UTC())
}

// FinanceSummary loads the executive finance card metrics and writes a mandatory audit trail.
func (s *Service) FinanceSummary(ctx context.Context, access AccessMetadata) (analyticsdomain.FinanceSummary, error) {
	normalized := s.normalizeAccess(access)
	summary, err := s.repo.FinanceSummary(ctx, normalized.CompanyCode)
	if err != nil {
		return analyticsdomain.FinanceSummary{}, err
	}
	if err := s.recordAccess(ctx, "finance.summary", normalized); err != nil {
		return analyticsdomain.FinanceSummary{}, err
	}
	return summary, nil
}

// SegmentProfit loads gross profit structure by operating segment and writes a mandatory audit trail.
func (s *Service) SegmentProfit(ctx context.Context, access AccessMetadata) ([]analyticsdomain.SegmentGrossProfit, error) {
	normalized := s.normalizeAccess(access)
	segments, err := s.repo.Segments(ctx, normalized.CompanyCode)
	if err != nil {
		return nil, err
	}
	if err := s.recordAccess(ctx, "finance.segments", normalized); err != nil {
		return nil, err
	}
	return segments, nil
}

// CashRunway returns a 3-month forward-looking contracted runway view.
func (s *Service) CashRunway(ctx context.Context, companyCode string, asOf time.Time, months int) (analyticsdomain.CashRunway, error) {
	if asOf.IsZero() {
		asOf = s.now().UTC()
	}
	if months <= 0 {
		months = 3
	}
	return s.repo.CashRunway(ctx, companyCode, asOf.UTC(), months)
}

// DashboardCashRunway returns a dashboard-ready 6-month projection and writes a mandatory audit trail.
func (s *Service) DashboardCashRunway(ctx context.Context, input CashRunwayInput) (analyticsdomain.CashRunway, error) {
	normalized := s.normalizeAccess(input.Access)
	asOf := input.AsOf
	if asOf.IsZero() {
		asOf = s.now().UTC()
	}
	months := input.Months
	if months <= 0 {
		months = 6
	}

	runway, err := s.repo.CashRunway(ctx, normalized.CompanyCode, asOf.UTC(), months)
	if err != nil {
		return analyticsdomain.CashRunway{}, err
	}
	if err := s.recordAccess(ctx, "finance.cash_runway", normalized); err != nil {
		return analyticsdomain.CashRunway{}, err
	}
	return runway, nil
}

func (s *Service) normalizeAccess(access AccessMetadata) AccessMetadata {
	if access.CompanyCode == "" {
		access.CompanyCode = "VCT_GROUP"
	}
	if access.Filters == nil {
		access.Filters = map[string]string{}
	}
	return access
}

func (s *Service) recordAccess(ctx context.Context, reportCode string, access AccessMetadata) error {
	if s.repo == nil {
		return fmt.Errorf("analytics repository is not configured")
	}
	if access.ActorID == "" {
		return fmt.Errorf("actor id is required for finance audit logging")
	}
	if access.ActorRole == "" {
		return fmt.Errorf("actor role is required for finance audit logging")
	}
	return s.repo.RecordReportAccess(ctx, analyticsdomain.ReportAccessLog{
		CompanyCode: access.CompanyCode,
		ReportCode:  reportCode,
		ActorID:     access.ActorID,
		ActorRole:   access.ActorRole,
		IPAddress:   access.IPAddress,
		UserAgent:   access.UserAgent,
		Filters:     access.Filters,
		AccessedAt:  s.now().UTC(),
	})
}

// ServiceOption customizes the analytics application service.
type ServiceOption func(*Service)

// WithDashboardCache enables short-lived Redis caching for executive dashboard snapshots.
func WithDashboardCache(cache analyticsdomain.DashboardCache, ttl time.Duration) ServiceOption {
	return func(service *Service) {
		service.dashboardCache = cache
		service.dashboardTTL = ttl
	}
}

func (s *Service) dashboardCacheTTLOrDefault() time.Duration {
	if s.dashboardTTL <= 0 {
		return defaultDashboardCacheTTL
	}
	return s.dashboardTTL
}
