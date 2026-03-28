import type {
  FinanceDashboardSnapshot,
  FinancePieSlice,
} from "@/lib/contracts/finance";
import type { AppLocale } from "@/lib/i18n/shared";
import {
  translateFinanceRefreshMode,
} from "@/lib/i18n/finance";

export type FinanceSegmentShare = FinancePieSlice & {
  share: number;
};

export type FinanceInsightTone = "navy" | "emerald" | "amber" | "rose";

export function getDashboardCard(
  snapshot: FinanceDashboardSnapshot,
  key: string,
) {
  return snapshot.cards.find((card) => card.key === key);
}

export function requireDashboardCard(
  snapshot: FinanceDashboardSnapshot,
  key: string,
) {
  const card = getDashboardCard(snapshot, key);
  if (!card) {
    throw new Error(`missing finance dashboard card: ${key}`);
  }
  return card;
}

export function getSeries(
  snapshot: FinanceDashboardSnapshot,
  key: string,
) {
  return snapshot.cashflow_chart.series.find((series) => series.key === key);
}

export function totalRevenue(snapshot: FinanceDashboardSnapshot) {
  return snapshot.revenue_mix.reduce((sum, slice) => sum + slice.value, 0);
}

export function buildSegmentShares(snapshot: FinanceDashboardSnapshot) {
  const total = totalRevenue(snapshot);

  return snapshot.revenue_mix.map((slice) => ({
    ...slice,
    share: total > 0 ? (slice.value / total) * 100 : 0,
  }));
}

export function latestSeriesValue(
  snapshot: FinanceDashboardSnapshot,
  key: string,
) {
  const series = getSeries(snapshot, key);
  if (!series || series.values.length === 0) {
    return 0;
  }
  return series.values[series.values.length - 1] ?? 0;
}

export function averageSeriesValue(
  snapshot: FinanceDashboardSnapshot,
  key: string,
) {
  const series = getSeries(snapshot, key);
  if (!series || series.values.length === 0) {
    return 0;
  }
  return (
    series.values.reduce((sum, value) => sum + value, 0) / series.values.length
  );
}

export function latestRunwayEnding(snapshot: FinanceDashboardSnapshot) {
  if (snapshot.runway_projection.length === 0) {
    return 0;
  }
  return snapshot.runway_projection[snapshot.runway_projection.length - 1]
    ?.projected_ending ?? 0;
}

export function latestLabel(snapshot: FinanceDashboardSnapshot) {
  if (snapshot.cashflow_chart.x_axis.length === 0) {
    return "current cycle";
  }
  return snapshot.cashflow_chart.x_axis[snapshot.cashflow_chart.x_axis.length - 1];
}

export function previousLabel(snapshot: FinanceDashboardSnapshot) {
  if (snapshot.cashflow_chart.x_axis.length < 2) {
    return "previous cycle";
  }
  return snapshot.cashflow_chart.x_axis[snapshot.cashflow_chart.x_axis.length - 2];
}

export function statusTone(status?: string) {
  switch ((status ?? "").toLowerCase()) {
    case "healthy":
    case "live":
    case "ready":
      return "emerald" as const;
    case "warning":
    case "watch":
      return "amber" as const;
    case "critical":
    case "risk":
      return "rose" as const;
    default:
      return "navy" as const;
  }
}

export function strongestSegment(snapshot: FinanceDashboardSnapshot) {
  return [...buildSegmentShares(snapshot)].sort((left, right) => right.value - left.value)[0];
}

export function weakestSegment(snapshot: FinanceDashboardSnapshot) {
  return [...buildSegmentShares(snapshot)].sort((left, right) => left.value - right.value)[0];
}

export function buildLedgerLanes(
  snapshot: FinanceDashboardSnapshot,
  locale: AppLocale = "vi",
) {
  const segments = buildSegmentShares(snapshot);

  return segments.map((segment, index) => ({
    segment: segment.label,
    voucherType: index % 2 === 0 ? "PT" : "PC",
    amount: segment.value,
    checkpoint:
      segment.label === "SaaS"
        ? locale === "vi"
          ? "Phân bổ doanh thu chưa thực hiện theo nhịp hàng tháng"
          : "Recognize deferred revenue on a monthly cadence"
        : segment.label === "Dojo"
          ? locale === "vi"
            ? "Kiểm tra tuổi nợ học phí trước khi khóa sổ"
            : "Check receivable aging before closing"
          : segment.label === "Retail"
            ? locale === "vi"
              ? "Rà soát giảm trừ doanh thu và bút toán đảo POS"
              : "Review deductions and POS reversals"
            : locale === "vi"
              ? "Xác minh giải phóng cọc và bù trừ hư hại"
              : "Verify deposit release and damage offsets",
    confidence: 96 + index,
  }));
}

export function buildControlSignals(
  snapshot: FinanceDashboardSnapshot,
  locale: AppLocale = "vi",
) {
  const runwayCard = requireDashboardCard(snapshot, "runway_index");
  const refreshTone: FinanceInsightTone =
    snapshot.recommended_refresh === "websocket" ? "emerald" : "amber";

  return [
    {
      label: locale === "vi" ? "Tầng dữ liệu" : "Data Plane",
      value: snapshot.data_mode.toUpperCase(),
      caption:
        locale === "vi"
          ? "Nguồn runtime đang cấp dữ liệu cho dashboard"
          : "Runtime source for the dashboard payload",
      tone: statusTone(snapshot.data_mode) as FinanceInsightTone,
    },
    {
      label: locale === "vi" ? "Chế độ làm mới" : "Refresh Mode",
      value: translateFinanceRefreshMode(snapshot.recommended_refresh, locale) ?? snapshot.recommended_refresh,
      caption:
        locale === "vi"
          ? "Cách giao diện điều hành giữ đồng bộ dữ liệu"
          : "How the executive UI should stay synchronized",
      tone: refreshTone,
    },
    {
      label: locale === "vi" ? "Tư thế runway" : "Runway Posture",
      value:
        locale === "vi"
          ? `${runwayCard.value.toFixed(1)} tháng`
          : `${runwayCard.value.toFixed(1)} months`,
      caption:
        locale === "vi"
          ? "Số tháng vận hành theo mức đốt tiền hiện tại"
          : "Months of operating capacity under current burn",
      tone: statusTone(runwayCard.status) as FinanceInsightTone,
    },
  ];
}
