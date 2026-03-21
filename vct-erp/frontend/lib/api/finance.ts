import {
  type FinanceDashboardSnapshot,
  type FinancePieSlice,
} from "@/lib/contracts/finance";

const BACKEND_BASE_URL =
  process.env.FINANCE_API_BASE_URL ?? "http://localhost:8080";
const BACKEND_ROLE = process.env.FINANCE_DASHBOARD_ROLE ?? "ceo";
const BACKEND_ACTOR_ID =
  process.env.FINANCE_DASHBOARD_ACTOR_ID ?? "command-center";

async function requestSnapshot(path: string) {
  const response = await fetch(`${BACKEND_BASE_URL}${path}`, {
    method: "GET",
    cache: "no-store",
    headers: {
      "X-App-Role": BACKEND_ROLE,
      "X-Actor-ID": BACKEND_ACTOR_ID,
    },
  });

  if (!response.ok) {
    throw new Error(`finance backend responded ${response.status}`);
  }

  return (await response.json()) as FinanceDashboardSnapshot;
}

export async function getFinanceDashboardSnapshot() {
  try {
    return await requestSnapshot("/api/v1/finance/dashboard");
  } catch {
    try {
      return await requestSnapshot("/api/v1/finance/dashboard/mock");
    } catch {
      return buildFallbackDashboardSnapshot();
    }
  }
}

export function buildFallbackDashboardSnapshot(): FinanceDashboardSnapshot {
  const revenueMix: FinancePieSlice[] = [
    { label: "SaaS", value: 0, color: "#0F766E" },
    { label: "Dojo", value: 0, color: "#D97706" },
    { label: "Retail", value: 0, color: "#2563EB" },
    { label: "Rental", value: 0, color: "#BE123C" },
  ];

  return {
    company_code: "VCT_GROUP",
    generated_at: new Date().toISOString(),
    data_mode: "fallback",
    recommended_refresh: "polling",
    cards: [
      {
        key: "cash_assets",
        title: "Tong tai san hien co",
        value: 0,
        formatted_value: "0 VND",
        unit: "VND",
        description: "Cho ket noi den backend tai chinh",
        trend: {
          direction: "flat",
          percentage: 0,
          delta: 0,
          period: "vs thang truoc",
        },
        chart_data: [],
      },
      {
        key: "quarter_net_revenue",
        title: "Doanh thu thuan quy",
        value: 0,
        formatted_value: "0 VND",
        unit: "VND",
        description: "Cho ket noi den backend tai chinh",
        trend: {
          direction: "flat",
          percentage: 0,
          delta: 0,
          period: "vs quy truoc",
        },
        chart_data: [],
      },
      {
        key: "runway_index",
        title: "Chi so runway",
        value: 0,
        formatted_value: "0 thang",
        unit: "months",
        status: "warning",
        description: "Cho ket noi den backend tai chinh",
        trend: {
          direction: "flat",
          percentage: 0,
          delta: 0,
          period: "vs thang truoc",
        },
        chart_data: [],
      },
    ],
    revenue_mix: revenueMix,
    cashflow_chart: {
      granularity: "month",
      x_axis: [],
      series: [
        { key: "revenue", label: "Revenue", color: "#0F766E", values: [] },
        { key: "expense", label: "Expense", color: "#D97706", values: [] },
        { key: "profit", label: "Profit", color: "#2563EB", values: [] },
      ],
    },
    runway_projection: [],
  };
}
