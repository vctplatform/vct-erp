package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

const (
	defaultDashboardCacheTTL = 60 * time.Second
	defaultTrendPoints       = 7
	defaultCashflowPoints    = 6
	defaultRunwayMonths      = 6
)

var segmentColors = map[string]string{
	"saas":   "#0F766E",
	"dojo":   "#D97706",
	"retail": "#2563EB",
	"rental": "#BE123C",
}

// DashboardInput controls the live executive dashboard request.
type DashboardInput struct {
	Access AccessMetadata
	AsOf   time.Time
}

// Dashboard returns the plug-and-play live dashboard contract used by the Command Center.
func (s *Service) Dashboard(ctx context.Context, input DashboardInput) (analyticsdomain.CommandCenterDashboardData, error) {
	normalized := s.normalizeAccess(input.Access)
	asOf := input.AsOf
	if asOf.IsZero() {
		asOf = s.now().UTC()
	}

	if s.dashboardCache != nil {
		cached, err := s.dashboardCache.Get(ctx, normalized.CompanyCode)
		switch {
		case err == nil:
			if err := s.recordAccess(ctx, "finance.dashboard.live", normalized); err != nil {
				return analyticsdomain.CommandCenterDashboardData{}, err
			}
			return cached, nil
		case errors.Is(err, analyticsdomain.ErrDashboardCacheMiss):
		default:
			return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("read dashboard cache: %w", err)
		}
	}

	snapshot, err := s.buildDashboard(ctx, normalized.CompanyCode, asOf)
	if err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, err
	}
	if err := s.recordAccess(ctx, "finance.dashboard.live", normalized); err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, err
	}
	if s.dashboardCache != nil {
		if err := s.dashboardCache.Set(ctx, snapshot, s.dashboardCacheTTLOrDefault()); err != nil {
			return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("write dashboard cache: %w", err)
		}
	}
	return snapshot, nil
}

// DashboardCards returns the live KPI-card subset of the full dashboard contract.
func (s *Service) DashboardCards(ctx context.Context, input DashboardInput) (analyticsdomain.DashboardCardsResponse, error) {
	dashboard, err := s.Dashboard(ctx, input)
	if err != nil {
		return analyticsdomain.DashboardCardsResponse{}, err
	}
	return analyticsdomain.DashboardCardsResponse{
		CompanyCode:        dashboard.CompanyCode,
		GeneratedAt:        dashboard.GeneratedAt,
		DataMode:           dashboard.DataMode,
		Cards:              dashboard.Cards,
		RecommendedRefresh: dashboard.RecommendedRefresh,
	}, nil
}

// DashboardSegments returns the live revenue-mix subset of the full dashboard contract.
func (s *Service) DashboardSegments(ctx context.Context, input DashboardInput) (analyticsdomain.DashboardSegmentsResponse, error) {
	dashboard, err := s.Dashboard(ctx, input)
	if err != nil {
		return analyticsdomain.DashboardSegmentsResponse{}, err
	}
	return analyticsdomain.DashboardSegmentsResponse{
		CompanyCode:        dashboard.CompanyCode,
		GeneratedAt:        dashboard.GeneratedAt,
		DataMode:           dashboard.DataMode,
		RevenueMix:         dashboard.RevenueMix,
		RecommendedRefresh: dashboard.RecommendedRefresh,
	}, nil
}

// DashboardCashflow returns the live chart subset of the full dashboard contract.
func (s *Service) DashboardCashflow(ctx context.Context, input DashboardInput) (analyticsdomain.DashboardCashflowResponse, error) {
	dashboard, err := s.Dashboard(ctx, input)
	if err != nil {
		return analyticsdomain.DashboardCashflowResponse{}, err
	}
	return analyticsdomain.DashboardCashflowResponse{
		CompanyCode:        dashboard.CompanyCode,
		GeneratedAt:        dashboard.GeneratedAt,
		DataMode:           dashboard.DataMode,
		CashflowChart:      dashboard.CashflowChart,
		RunwayProjection:   dashboard.RunwayProjection,
		RecommendedRefresh: dashboard.RecommendedRefresh,
	}, nil
}

func (s *Service) buildDashboard(ctx context.Context, companyCode string, asOf time.Time) (analyticsdomain.CommandCenterDashboardData, error) {
	reconciliation, err := s.repo.ReconciliationBalances(ctx, companyCode)
	if err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("load finance reconciliation: %w", err)
	}
	segments, err := s.repo.Segments(ctx, companyCode)
	if err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("load gross profit segments: %w", err)
	}
	trendSeries, err := s.repo.FinanceSeries(ctx, companyCode, asOf.UTC(), defaultTrendPoints)
	if err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("load finance trend series: %w", err)
	}
	cashflowSeries, err := s.repo.FinanceSeries(ctx, companyCode, asOf.UTC(), defaultCashflowPoints)
	if err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("load cashflow chart series: %w", err)
	}
	runway, err := s.repo.CashRunway(ctx, companyCode, asOf.UTC(), defaultRunwayMonths)
	if err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("load cash runway: %w", err)
	}

	currentCash := sumCashBalances(reconciliation)
	currentQuarterRevenue, previousQuarterRevenue, err := s.quarterRevenueComparison(ctx, companyCode, asOf.UTC())
	if err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("load quarter revenue comparison: %w", err)
	}
	runwaySeries := runwayHistory(trendSeries)
	currentRunwayMonths := calculateRunwayMonths(runway.CurrentCash, runway.AverageMonthlyBurn)
	previousRunwayMonths := previousRunway(runwaySeries, currentRunwayMonths)

	cards := []analyticsdomain.DashboardCard{
		buildDashboardCard(
			"cash_assets",
			"Tong tai san hien co",
			currentCash,
			"VND",
			"",
			"Tong 111 va 112 tai thoi diem hien tai",
			trendFromMoney(currentCash, previousCashBalance(trendSeries), "vs thang truoc"),
			miniChartFromMoneySeries(trendSeries, func(point analyticsdomain.FinanceSeriesPoint) ledgerdomain.Money { return point.CashEnding }),
		),
		buildDashboardCard(
			"quarter_net_revenue",
			"Doanh thu thuan quy",
			currentQuarterRevenue,
			"VND",
			"",
			"SaaS + Dojo + Retail - giam tru doanh thu trong quy hien tai",
			trendFromMoney(currentQuarterRevenue, previousQuarterRevenue, "vs quy truoc"),
			miniChartFromMoneySeries(trendSeries, func(point analyticsdomain.FinanceSeriesPoint) ledgerdomain.Money { return point.Revenue }),
		),
		buildRunwayCard(currentRunwayMonths, previousRunwayMonths, runwaySeries),
	}

	return analyticsdomain.CommandCenterDashboardData{
		CompanyCode:        companyCode,
		GeneratedAt:        s.now().UTC(),
		DataMode:           "live",
		Cards:              cards,
		RevenueMix:         buildRevenueMix(segments),
		CashflowChart:      buildCashflowChart(cashflowSeries),
		RunwayProjection:   buildRunwayProjection(runway),
		RecommendedRefresh: "websocket",
	}, nil
}

func (s *Service) quarterRevenueComparison(ctx context.Context, companyCode string, asOf time.Time) (ledgerdomain.Money, ledgerdomain.Money, error) {
	currentQuarterStart := quarterStart(asOf)
	currentQuarterEnd := currentQuarterStart.AddDate(0, 3, -1)
	previousQuarterStart := currentQuarterStart.AddDate(0, -3, 0)
	previousQuarterEnd := currentQuarterStart.AddDate(0, 0, -1)

	currentPoints, err := s.repo.RevenueStream(ctx, companyCode, currentQuarterStart, currentQuarterEnd)
	if err != nil {
		return ledgerdomain.Money{}, ledgerdomain.Money{}, err
	}
	previousPoints, err := s.repo.RevenueStream(ctx, companyCode, previousQuarterStart, previousQuarterEnd)
	if err != nil {
		return ledgerdomain.Money{}, ledgerdomain.Money{}, err
	}

	return sumRevenuePoints(currentPoints), sumRevenuePoints(previousPoints), nil
}

func quarterStart(asOf time.Time) time.Time {
	monthIndex := ((int(asOf.Month()) - 1) / 3) * 3
	return time.Date(asOf.Year(), time.Month(monthIndex+1), 1, 0, 0, 0, 0, time.UTC)
}

func buildDashboardCard(key string, title string, value ledgerdomain.Money, unit string, status string, description string, trend analyticsdomain.CardTrend, chartData []analyticsdomain.MiniChartPoint) analyticsdomain.DashboardCard {
	return analyticsdomain.DashboardCard{
		Key:            key,
		Title:          title,
		Value:          moneyToFloat64(value),
		FormattedValue: formatCompactVND(value, unit),
		Unit:           unit,
		Status:         status,
		Description:    description,
		Trend:          trend,
		ChartData:      chartData,
	}
}

func buildRunwayCard(current float64, previous float64, history []analyticsdomain.MiniChartPoint) analyticsdomain.DashboardCard {
	status := "critical"
	if current >= 6 {
		status = "healthy"
	} else if current >= 3 {
		status = "warning"
	}

	return analyticsdomain.DashboardCard{
		Key:            "runway_index",
		Title:          "Chi so runway",
		Value:          current,
		FormattedValue: fmt.Sprintf("%.1f thang", current),
		Unit:           "months",
		Status:         status,
		Description:    "Do du tien van hanh neu burn rate giu nguyen",
		Trend:          trendFromFloat(current, previous, "vs thang truoc"),
		ChartData:      history,
	}
}

func buildRevenueMix(segments []analyticsdomain.SegmentGrossProfit) []analyticsdomain.PieChartSlice {
	revenueMix := make([]analyticsdomain.PieChartSlice, 0, len(segments))
	for _, segment := range segments {
		revenueMix = append(revenueMix, analyticsdomain.PieChartSlice{
			Label: segment.SegmentLabel,
			Value: moneyToFloat64(segment.NetRevenue),
			Color: segmentColors[segment.SegmentKey],
		})
	}
	return revenueMix
}

func buildCashflowChart(points []analyticsdomain.FinanceSeriesPoint) analyticsdomain.MultiLineChart {
	xAxis := make([]string, 0, len(points))
	revenue := make([]float64, 0, len(points))
	expense := make([]float64, 0, len(points))
	profit := make([]float64, 0, len(points))

	for _, point := range points {
		xAxis = append(xAxis, point.PeriodLabel)
		revenue = append(revenue, moneyToFloat64(point.Revenue))
		expense = append(expense, moneyToFloat64(point.Expense))
		profit = append(profit, moneyToFloat64(point.Profit))
	}

	return analyticsdomain.MultiLineChart{
		Granularity: "month",
		XAxis:       xAxis,
		Series: []analyticsdomain.LineChartSeries{
			{Key: "revenue", Label: "Revenue", Color: "#0F766E", Values: revenue},
			{Key: "expense", Label: "Expense", Color: "#D97706", Values: expense},
			{Key: "profit", Label: "Profit", Color: "#2563EB", Values: profit},
		},
	}
}

func buildRunwayProjection(runway analyticsdomain.CashRunway) []analyticsdomain.RunwayProjectionPoint {
	points := make([]analyticsdomain.RunwayProjectionPoint, 0, len(runway.Months))
	for _, month := range runway.Months {
		points = append(points, analyticsdomain.RunwayProjectionPoint{
			Label:            month.MonthLabel,
			OpeningCash:      moneyToFloat64(month.OpeningCash),
			ContractedInflow: moneyToFloat64(month.ContractedInflow),
			ProjectedBurn:    moneyToFloat64(month.ProjectedBurn),
			ProjectedEnding:  moneyToFloat64(month.ProjectedEnding),
		})
	}
	return points
}

func sumCashBalances(balances []analyticsdomain.ReconciliationBalance) ledgerdomain.Money {
	total := ledgerdomain.MustParseMoney("0.0000")
	for _, balance := range balances {
		if strings.HasPrefix(balance.AccountCode, "111") || strings.HasPrefix(balance.AccountCode, "112") {
			total = total.Add(balance.EffectiveNetBalance)
		}
	}
	return total
}

func previousCashBalance(points []analyticsdomain.FinanceSeriesPoint) ledgerdomain.Money {
	if len(points) < 2 {
		return ledgerdomain.MustParseMoney("0.0000")
	}
	return points[len(points)-2].CashEnding
}

func previousRunway(points []analyticsdomain.MiniChartPoint, fallback float64) float64 {
	if len(points) < 2 {
		return fallback
	}
	return points[len(points)-2].Value
}

func runwayHistory(points []analyticsdomain.FinanceSeriesPoint) []analyticsdomain.MiniChartPoint {
	history := make([]analyticsdomain.MiniChartPoint, 0, len(points))
	for index, point := range points {
		start := 0
		if index >= 2 {
			start = index - 2
		}
		burnWindow := ledgerdomain.MustParseMoney("0.0000")
		windowCount := 0
		for cursor := start; cursor <= index; cursor++ {
			if points[cursor].Expense.Sign() <= 0 {
				continue
			}
			burnWindow = burnWindow.Add(points[cursor].Expense)
			windowCount++
		}
		averageBurn := burnWindow
		if windowCount > 0 {
			averageBurn = burnWindow.DivInt(windowCount)
		}
		history = append(history, analyticsdomain.MiniChartPoint{
			Label: point.PeriodLabel,
			Value: calculateRunwayMonths(point.CashEnding, averageBurn),
		})
	}
	return history
}

func calculateRunwayMonths(currentCash ledgerdomain.Money, averageBurn ledgerdomain.Money) float64 {
	if averageBurn.Sign() <= 0 {
		return 0
	}

	numerator, err := strconv.ParseFloat(currentCash.String(), 64)
	if err != nil {
		return 0
	}
	denominator, err := strconv.ParseFloat(averageBurn.String(), 64)
	if err != nil || denominator == 0 {
		return 0
	}

	return numerator / denominator
}

func miniChartFromMoneySeries(points []analyticsdomain.FinanceSeriesPoint, picker func(analyticsdomain.FinanceSeriesPoint) ledgerdomain.Money) []analyticsdomain.MiniChartPoint {
	chart := make([]analyticsdomain.MiniChartPoint, 0, len(points))
	for _, point := range points {
		chart = append(chart, analyticsdomain.MiniChartPoint{
			Label: point.PeriodLabel,
			Value: moneyToFloat64(picker(point)),
		})
	}
	return chart
}

func trendFromMoney(current ledgerdomain.Money, previous ledgerdomain.Money, period string) analyticsdomain.CardTrend {
	return trendFromFloat(moneyToFloat64(current), moneyToFloat64(previous), period)
}

func trendFromFloat(current float64, previous float64, period string) analyticsdomain.CardTrend {
	delta := current - previous
	direction := "flat"
	switch {
	case delta > 0:
		direction = "up"
	case delta < 0:
		direction = "down"
	}

	percentage := 0.0
	if previous != 0 {
		percentage = (delta / previous) * 100
	}

	return analyticsdomain.CardTrend{
		Direction:  direction,
		Percentage: percentage,
		Delta:      delta,
		Period:     period,
	}
}

func sumRevenuePoints(points []analyticsdomain.RevenueStreamPoint) ledgerdomain.Money {
	total := ledgerdomain.MustParseMoney("0.0000")
	for _, point := range points {
		total = total.Add(point.Amount)
	}
	return total
}

func moneyToFloat64(amount ledgerdomain.Money) float64 {
	parsed, err := strconv.ParseFloat(amount.String(), 64)
	if err != nil {
		return 0
	}
	return parsed
}

func formatCompactVND(amount ledgerdomain.Money, unit string) string {
	value := moneyToFloat64(amount)
	absolute := value
	if absolute < 0 {
		absolute = -absolute
	}

	switch {
	case absolute >= 1_000_000_000:
		return fmt.Sprintf("%.2f ty %s", value/1_000_000_000, unit)
	case absolute >= 1_000_000:
		return fmt.Sprintf("%.2f tr %s", value/1_000_000, unit)
	default:
		return fmt.Sprintf("%.0f %s", value, unit)
	}
}
