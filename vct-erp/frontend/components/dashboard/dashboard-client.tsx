"use client";

import { startTransition } from "react";
import dynamic from "next/dynamic";
import { Activity, RadioTower } from "lucide-react";
import { toast } from "sonner";

import { useLocale } from "@/components/i18n/locale-provider";
import { KpiCardGrid } from "@/components/dashboard/kpi-card-grid";
import { useFinanceDashboard } from "@/hooks/use-finance-dashboard";
import { useFinanceWebSocket } from "@/hooks/use-finance-websocket";
import type {
  FinanceDashboardSnapshot,
  FinanceRealtimeEvent,
} from "@/lib/contracts/finance";
import { formatCurrency } from "@/lib/formatters";
import { getFinanceLocaleCode } from "@/lib/i18n/finance";
import { useFinanceDashboardStore } from "@/lib/store/finance-dashboard-store";

type DashboardClientProps = {
  initialData: FinanceDashboardSnapshot;
};

function chartLoadingCard(heightClass = "h-80") {
  return (
    <div className="rounded-[1.5rem] border border-[var(--color-border)] bg-[var(--color-panel)] p-5">
      <div
        className={`${heightClass} animate-pulse rounded-[1.25rem] bg-[var(--color-canvas-soft)]`}
      />
    </div>
  );
}

const RevenueMixChart = dynamic(
  () =>
    import("@/components/dashboard/revenue-mix-chart").then((module) => ({
      default: module.RevenueMixChart,
    })),
  {
    ssr: false,
    loading: () => chartLoadingCard("h-72"),
  },
);

const CashflowTrendChart = dynamic(
  () =>
    import("@/components/dashboard/cashflow-trend-chart").then((module) => ({
      default: module.CashflowTrendChart,
    })),
  {
    ssr: false,
    loading: () => chartLoadingCard(),
  },
);

const RunwayForecastChart = dynamic(
  () =>
    import("@/components/dashboard/runway-forecast-chart").then((module) => ({
      default: module.RunwayForecastChart,
    })),
  {
    ssr: false,
    loading: () => chartLoadingCard(),
  },
);

function resolveFinanceWsUrl() {
  if (typeof window === "undefined") {
    return null;
  }

  const base = process.env.NEXT_PUBLIC_FINANCE_WS_URL;
  if (!base) {
    return null;
  }

  const url = new URL(base);
  url.searchParams.set(
    "role",
    process.env.NEXT_PUBLIC_FINANCE_WS_ROLE ?? "ceo",
  );
  url.searchParams.set(
    "actor_id",
    process.env.NEXT_PUBLIC_FINANCE_WS_ACTOR_ID ?? "command-center",
  );
  return url.toString();
}

export function DashboardClient({ initialData }: DashboardClientProps) {
  const { locale } = useLocale();
  const { snapshot, error, mutate, isValidating } =
    useFinanceDashboard(initialData);
  const applyRealtimeEvent = useFinanceDashboardStore(
    (state) => state.applyRealtimeEvent,
  );
  const wsUrl = resolveFinanceWsUrl();
  const copy =
    locale === "vi"
      ? {
          newRevenue: "Doanh thu mới",
          signalFeed: "Dòng tín hiệu tài chính",
          lastSnapshot: "Ảnh chụp gần nhất",
          realtimeLive: "Realtime trực tiếp",
          pollingFallback: "Rơi về polling",
          syncing: "Đang đồng bộ...",
          dataReady: "Dữ liệu sẵn sàng",
          syncError:
            "Không thể đồng bộ dashboard live. Frontend đang hiển thị snapshot gần nhất.",
        }
      : {
          newRevenue: "New revenue",
          signalFeed: "Financial Signal Feed",
          lastSnapshot: "Last snapshot",
          realtimeLive: "Realtime live",
          pollingFallback: "Polling fallback",
          syncing: "Syncing...",
          dataReady: "Data ready",
          syncError:
            "Unable to sync the live dashboard. The frontend is showing the latest snapshot.",
        };

  const realtime = useFinanceWebSocket(wsUrl, {
    enabled: snapshot.recommended_refresh === "websocket",
    onTransaction: (event: FinanceRealtimeEvent) => {
      toast.success(
        `${copy.newRevenue}: +${formatCurrency(event.amount, locale)} (${event.segment})`,
      );

      startTransition(() => {
        applyRealtimeEvent(event);
      });

      void mutate();
    },
  });

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 rounded-[1.4rem] border border-[var(--color-border)] bg-[var(--color-panel)] px-4 py-4 md:flex-row md:items-center md:justify-between">
        <div className="flex items-center gap-3">
          <span className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-[var(--color-canvas-soft)] text-[var(--color-navy-700)] dark:text-white">
            <Activity className="size-4" />
          </span>
          <div>
            <p className="text-sm font-medium text-[var(--color-ink)]">
              {copy.signalFeed}
            </p>
            <p className="text-sm text-[var(--color-ink-soft)]">
              {copy.lastSnapshot}{" "}
              {new Date(snapshot.generated_at).toLocaleString(
                getFinanceLocaleCode(locale),
              )}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3 text-sm">
          <span className="inline-flex items-center gap-2 rounded-full border border-[var(--color-border)] px-3 py-2 text-[var(--color-ink-soft)]">
            <RadioTower
              className={`size-4 ${
                realtime.isConnected ? "text-emerald-500" : "text-rose-500"
              }`}
            />
            {realtime.isConnected ? copy.realtimeLive : copy.pollingFallback}
          </span>
          <span className="rounded-full border border-[var(--color-border)] px-3 py-2 text-[var(--color-ink-soft)]">
            {isValidating ? copy.syncing : copy.dataReady}
          </span>
        </div>
      </div>

      {error ? (
        <div className="rounded-[1.35rem] border border-rose-500/25 bg-rose-500/10 px-4 py-4 text-sm text-rose-700 dark:text-rose-200">
          {copy.syncError}
        </div>
      ) : null}

      <KpiCardGrid cards={snapshot.cards} />

      <div className="grid gap-6 xl:grid-cols-[1.05fr_1.3fr]">
        <RevenueMixChart data={snapshot.revenue_mix} />
        <CashflowTrendChart chart={snapshot.cashflow_chart} />
      </div>

      <RunwayForecastChart data={snapshot.runway_projection} />
    </div>
  );
}
