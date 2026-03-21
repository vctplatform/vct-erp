package domain

import "time"

// MiniChartPoint is the smallest chart contract used by KPI cards.
type MiniChartPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// CardTrend captures the delta versus the previous comparison window.
type CardTrend struct {
	Direction  string  `json:"direction"`
	Percentage float64 `json:"percentage"`
	Delta      float64 `json:"delta"`
	Period     string  `json:"period"`
}

// DashboardCard represents one KPI card on the Command Center.
type DashboardCard struct {
	Key            string           `json:"key"`
	Title          string           `json:"title"`
	Value          float64          `json:"value"`
	FormattedValue string           `json:"formatted_value"`
	Unit           string           `json:"unit"`
	Status         string           `json:"status,omitempty"`
	Description    string           `json:"description,omitempty"`
	Trend          CardTrend        `json:"trend"`
	ChartData      []MiniChartPoint `json:"chart_data"`
}

// PieChartSlice is the plug-and-play contract expected by the frontend pie chart.
type PieChartSlice struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

// LineChartSeries describes one Y-axis dataset in the multi-line chart.
type LineChartSeries struct {
	Key    string    `json:"key"`
	Label  string    `json:"label"`
	Color  string    `json:"color"`
	Values []float64 `json:"values"`
}

// MultiLineChart is the standard line chart contract for the Command Center.
type MultiLineChart struct {
	Granularity string            `json:"granularity"`
	XAxis       []string          `json:"x_axis"`
	Series      []LineChartSeries `json:"series"`
}

// RunwayProjectionPoint provides the 6-month runway table/line data.
type RunwayProjectionPoint struct {
	Label            string  `json:"label"`
	OpeningCash      float64 `json:"opening_cash"`
	ContractedInflow float64 `json:"contracted_inflow"`
	ProjectedBurn    float64 `json:"projected_burn"`
	ProjectedEnding  float64 `json:"projected_ending"`
}

// CommandCenterDashboardData is the top-level plug-and-play contract for the alpha dashboard.
type CommandCenterDashboardData struct {
	CompanyCode        string                  `json:"company_code"`
	GeneratedAt        time.Time               `json:"generated_at"`
	DataMode           string                  `json:"data_mode"`
	Cards              []DashboardCard         `json:"cards"`
	RevenueMix         []PieChartSlice         `json:"revenue_mix"`
	CashflowChart      MultiLineChart          `json:"cashflow_chart"`
	RunwayProjection   []RunwayProjectionPoint `json:"runway_projection"`
	RecommendedRefresh string                  `json:"recommended_refresh"`
}

// DashboardCashflowResponse is the live contract for widgets that only need chart data.
type DashboardCashflowResponse struct {
	CompanyCode        string                  `json:"company_code"`
	GeneratedAt        time.Time               `json:"generated_at"`
	DataMode           string                  `json:"data_mode"`
	CashflowChart      MultiLineChart          `json:"cashflow_chart"`
	RunwayProjection   []RunwayProjectionPoint `json:"runway_projection"`
	RecommendedRefresh string                  `json:"recommended_refresh"`
}

// DashboardCardsResponse is the live contract for the executive KPI strip.
type DashboardCardsResponse struct {
	CompanyCode        string          `json:"company_code"`
	GeneratedAt        time.Time       `json:"generated_at"`
	DataMode           string          `json:"data_mode"`
	Cards              []DashboardCard `json:"cards"`
	RecommendedRefresh string          `json:"recommended_refresh"`
}

// DashboardSegmentsResponse is the live contract for revenue mix widgets.
type DashboardSegmentsResponse struct {
	CompanyCode        string          `json:"company_code"`
	GeneratedAt        time.Time       `json:"generated_at"`
	DataMode           string          `json:"data_mode"`
	RevenueMix         []PieChartSlice `json:"revenue_mix"`
	RecommendedRefresh string          `json:"recommended_refresh"`
}
