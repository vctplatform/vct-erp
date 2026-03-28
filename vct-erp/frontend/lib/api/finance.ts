import { cache } from "react";

import {
  type FinanceDashboardSnapshot,
  type FinancePieSlice,
} from "@/lib/contracts/finance";
import { localizeFinanceDashboardSnapshot } from "@/lib/i18n/finance";
import type { AppLocale } from "@/lib/i18n/shared";

const BACKEND_BASE_URL =
  process.env.FINANCE_API_BASE_URL ?? "http://localhost:8080";
const BACKEND_ROLE = process.env.FINANCE_DASHBOARD_ROLE ?? "ceo";
const BACKEND_ACTOR_ID =
  process.env.FINANCE_DASHBOARD_ACTOR_ID ?? "command-center";
const FINANCE_COMPANY_CODE =
  process.env.FINANCE_COMPANY_CODE ?? "VCT_SIM";

async function requestSnapshot(path: string) {
  const url = new URL(path, BACKEND_BASE_URL);
  if (FINANCE_COMPANY_CODE) {
    url.searchParams.set("company_code", FINANCE_COMPANY_CODE);
  }

  const response = await fetch(url.toString(), {
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

export const getFinanceDashboardSnapshot = cache(async (
  locale: AppLocale = "vi",
) => {
  try {
    return localizeFinanceDashboardSnapshot(
      await requestSnapshot("/api/v1/finance/dashboard"),
      locale,
    );
  } catch {
    try {
      return localizeFinanceDashboardSnapshot(
        await requestSnapshot("/api/v1/finance/dashboard/mock"),
        locale,
      );
    } catch {
      return buildFallbackDashboardSnapshot(locale);
    }
  }
});

export function buildFallbackDashboardSnapshot(
  locale: AppLocale = "vi",
): FinanceDashboardSnapshot {
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
        title: locale === "vi" ? "Tổng tài sản hiện có" : "Current Cash Assets",
        value: 0,
        formatted_value: "0 VND",
        unit: "VND",
        description:
          locale === "vi"
            ? "Chờ kết nối đến backend tài chính"
            : "Waiting for the finance backend connection",
        trend: {
          direction: "flat",
          percentage: 0,
          delta: 0,
          period: locale === "vi" ? "vs thang truoc" : "vs previous month",
        },
        chart_data: [],
      },
      {
        key: "quarter_net_revenue",
        title: locale === "vi" ? "Doanh thu thuần quý" : "Quarter Net Revenue",
        value: 0,
        formatted_value: "0 VND",
        unit: "VND",
        description:
          locale === "vi"
            ? "Chờ kết nối đến backend tài chính"
            : "Waiting for the finance backend connection",
        trend: {
          direction: "flat",
          percentage: 0,
          delta: 0,
          period: locale === "vi" ? "vs quy truoc" : "vs previous quarter",
        },
        chart_data: [],
      },
      {
        key: "runway_index",
        title: locale === "vi" ? "Chỉ số runway" : "Runway Index",
        value: 0,
        formatted_value: locale === "vi" ? "0 tháng" : "0 months",
        unit: "months",
        status: "warning",
        description:
          locale === "vi"
            ? "Chờ kết nối đến backend tài chính"
            : "Waiting for the finance backend connection",
        trend: {
          direction: "flat",
          percentage: 0,
          delta: 0,
          period: locale === "vi" ? "vs thang truoc" : "vs previous month",
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
