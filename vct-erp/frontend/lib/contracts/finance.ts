export type FinanceMiniChartPoint = {
  label: string;
  value: number;
};

export type FinanceCardTrend = {
  direction: "up" | "down" | "flat";
  percentage: number;
  delta: number;
  period: string;
};

export type FinanceDashboardCard = {
  key: string;
  title: string;
  value: number;
  formatted_value: string;
  unit: string;
  status?: string;
  description?: string;
  trend: FinanceCardTrend;
  chart_data: FinanceMiniChartPoint[];
};

export type FinancePieSlice = {
  label: string;
  value: number;
  color: string;
};

export type FinanceLineSeries = {
  key: string;
  label: string;
  color: string;
  values: number[];
};

export type FinanceMultiLineChart = {
  granularity: string;
  x_axis: string[];
  series: FinanceLineSeries[];
};

export type FinanceRunwayPoint = {
  label: string;
  opening_cash: number;
  contracted_inflow: number;
  projected_burn: number;
  projected_ending: number;
};

export type FinanceDashboardSnapshot = {
  company_code: string;
  generated_at: string;
  data_mode: "live" | "mock" | "fallback";
  cards: FinanceDashboardCard[];
  revenue_mix: FinancePieSlice[];
  cashflow_chart: FinanceMultiLineChart;
  runway_projection: FinanceRunwayPoint[];
  recommended_refresh: string;
};

export type FinanceRealtimeEvent = {
  event: "NEW_TRANSACTION";
  company_code: string;
  entry_id: string;
  reference_no?: string;
  amount: number;
  segment: string;
  source_module?: string;
  timestamp: string;
};
