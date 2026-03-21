package domain

import (
	"context"
	"errors"
	"time"
)

// ErrDashboardCacheMiss reports a missing dashboard snapshot in the cache tier.
var ErrDashboardCacheMiss = errors.New("dashboard cache miss")

// DashboardCache stores live dashboard snapshots for short-lived executive reads.
type DashboardCache interface {
	Get(ctx context.Context, companyCode string) (CommandCenterDashboardData, error)
	Set(ctx context.Context, snapshot CommandCenterDashboardData, ttl time.Duration) error
	Delete(ctx context.Context, companyCode string) error
}
