package usecase

import (
	"context"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
)

// Service exposes dashboard-ready analytics responses.
type Service struct {
	repo analyticsdomain.Repository
	now  func() time.Time
}

// NewService constructs the analytics application service.
func NewService(repo analyticsdomain.Repository) *Service {
	return &Service{
		repo: repo,
		now:  time.Now,
	}
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
