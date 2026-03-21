"use client";

import { startTransition } from "react";
import { Activity, RadioTower } from "lucide-react";
import { toast } from "sonner";

import { CashflowTrendChart } from "@/components/dashboard/cashflow-trend-chart";
import { KpiCardGrid } from "@/components/dashboard/kpi-card-grid";
import { RevenueMixChart } from "@/components/dashboard/revenue-mix-chart";
import { RunwayForecastChart } from "@/components/dashboard/runway-forecast-chart";
import { useFinanceDashboard } from "@/hooks/use-finance-dashboard";
import { useFinanceWebSocket } from "@/hooks/use-finance-websocket";
import type {
  FinanceDashboardSnapshot,
  FinanceRealtimeEvent,
} from "@/lib/contracts/finance";
import { formatCurrency } from "@/lib/formatters";
import { useFinanceDashboardStore } from "@/lib/store/finance-dashboard-store";

type DashboardClientProps = {
  initialData: FinanceDashboardSnapshot;
};

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
  const { snapshot, error, mutate, isValidating } =
    useFinanceDashboard(initialData);
  const applyRealtimeEvent = useFinanceDashboardStore(
    (state) => state.applyRealtimeEvent,
  );
  const wsUrl = resolveFinanceWsUrl();

  const realtime = useFinanceWebSocket(wsUrl, {
    enabled: snapshot.recommended_refresh === "websocket",
    onTransaction: (event: FinanceRealtimeEvent) => {
      toast.success(
        `Doanh thu moi: +${formatCurrency(event.amount)} (${event.segment})`,
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
              Financial signal feed
            </p>
            <p className="text-sm text-[var(--color-ink-soft)]">
              Last snapshot {new Date(snapshot.generated_at).toLocaleString("vi-VN")}
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
            {realtime.isConnected ? "Realtime live" : "Polling fallback"}
          </span>
          <span className="rounded-full border border-[var(--color-border)] px-3 py-2 text-[var(--color-ink-soft)]">
            {isValidating ? "Syncing..." : "Data ready"}
          </span>
        </div>
      </div>

      {error ? (
        <div className="rounded-[1.35rem] border border-rose-500/25 bg-rose-500/10 px-4 py-4 text-sm text-rose-700 dark:text-rose-200">
          Khong the dong bo live dashboard. Frontend dang hien thi snapshot gan nhat.
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
