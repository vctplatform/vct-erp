package usecase

import (
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
)

// GetMockDashboardData returns a high-fidelity contract for plug-and-play frontend work.
func GetMockDashboardData() analyticsdomain.CommandCenterDashboardData {
	return analyticsdomain.CommandCenterDashboardData{
		CompanyCode:        "VCT_GROUP",
		GeneratedAt:        time.Now().UTC(),
		DataMode:           "mock",
		RecommendedRefresh: "websocket",
		Cards: []analyticsdomain.DashboardCard{
			{
				Key:            "cash_assets",
				Title:          "Tong tai san hien co",
				Value:          8420000000,
				FormattedValue: "8.42 ty VND",
				Unit:           "VND",
				Description:    "Tong 1111 va 1121 tai thoi diem hien tai",
				Trend: analyticsdomain.CardTrend{
					Direction:  "up",
					Percentage: 6.8,
					Delta:      536000000,
					Period:     "vs thang truoc",
				},
				ChartData: []analyticsdomain.MiniChartPoint{
					{Label: "2026-03-16", Value: 7860000000},
					{Label: "2026-03-17", Value: 7920000000},
					{Label: "2026-03-18", Value: 8010000000},
					{Label: "2026-03-19", Value: 8090000000},
					{Label: "2026-03-20", Value: 8210000000},
					{Label: "2026-03-21", Value: 8330000000},
					{Label: "2026-03-22", Value: 8420000000},
				},
			},
			{
				Key:            "quarter_net_revenue",
				Title:          "Doanh thu thuan quy",
				Value:          2180000000,
				FormattedValue: "2.18 ty VND",
				Unit:           "VND",
				Description:    "SaaS + Dojo + Retail - giam tru doanh thu",
				Trend: analyticsdomain.CardTrend{
					Direction:  "up",
					Percentage: 12.4,
					Delta:      240000000,
					Period:     "vs quy truoc",
				},
				ChartData: []analyticsdomain.MiniChartPoint{
					{Label: "W-6", Value: 268000000},
					{Label: "W-5", Value: 284000000},
					{Label: "W-4", Value: 296000000},
					{Label: "W-3", Value: 314000000},
					{Label: "W-2", Value: 327000000},
					{Label: "W-1", Value: 342000000},
					{Label: "Now", Value: 349000000},
				},
			},
			{
				Key:            "runway_index",
				Title:          "Chi so runway",
				Value:          8.7,
				FormattedValue: "8.7 thang",
				Unit:           "months",
				Status:         "healthy",
				Description:    "Do du tien van hanh neu burn rate giu nguyen",
				Trend: analyticsdomain.CardTrend{
					Direction:  "up",
					Percentage: 9.1,
					Delta:      0.7,
					Period:     "vs thang truoc",
				},
				ChartData: []analyticsdomain.MiniChartPoint{
					{Label: "M-6", Value: 7.1},
					{Label: "M-5", Value: 7.4},
					{Label: "M-4", Value: 7.7},
					{Label: "M-3", Value: 8.0},
					{Label: "M-2", Value: 8.2},
					{Label: "M-1", Value: 8.4},
					{Label: "Now", Value: 8.7},
				},
			},
		},
		RevenueMix: []analyticsdomain.PieChartSlice{
			{Label: "SaaS", Value: 1320000000, Color: "#0F766E"},
			{Label: "Dojo", Value: 510000000, Color: "#D97706"},
			{Label: "Retail", Value: 290000000, Color: "#2563EB"},
			{Label: "Rental", Value: 60000000, Color: "#BE123C"},
		},
		CashflowChart: analyticsdomain.MultiLineChart{
			Granularity: "month",
			XAxis:       []string{"2025-10", "2025-11", "2025-12", "2026-01", "2026-02", "2026-03"},
			Series: []analyticsdomain.LineChartSeries{
				{
					Key:    "revenue",
					Label:  "Revenue",
					Color:  "#0F766E",
					Values: []float64{1480000000, 1560000000, 1620000000, 1710000000, 1790000000, 1880000000},
				},
				{
					Key:    "expense",
					Label:  "Expense",
					Color:  "#D97706",
					Values: []float64{890000000, 910000000, 930000000, 955000000, 980000000, 1010000000},
				},
				{
					Key:    "profit",
					Label:  "Profit",
					Color:  "#2563EB",
					Values: []float64{590000000, 650000000, 690000000, 755000000, 810000000, 870000000},
				},
			},
		},
		RunwayProjection: []analyticsdomain.RunwayProjectionPoint{
			{Label: "2026-04", OpeningCash: 8420000000, ContractedInflow: 640000000, ProjectedBurn: 950000000, ProjectedEnding: 8110000000},
			{Label: "2026-05", OpeningCash: 8110000000, ContractedInflow: 620000000, ProjectedBurn: 940000000, ProjectedEnding: 7790000000},
			{Label: "2026-06", OpeningCash: 7790000000, ContractedInflow: 610000000, ProjectedBurn: 935000000, ProjectedEnding: 7465000000},
			{Label: "2026-07", OpeningCash: 7465000000, ContractedInflow: 590000000, ProjectedBurn: 930000000, ProjectedEnding: 7125000000},
			{Label: "2026-08", OpeningCash: 7125000000, ContractedInflow: 575000000, ProjectedBurn: 925000000, ProjectedEnding: 6775000000},
			{Label: "2026-09", OpeningCash: 6775000000, ContractedInflow: 560000000, ProjectedBurn: 920000000, ProjectedEnding: 6415000000},
		},
	}
}
