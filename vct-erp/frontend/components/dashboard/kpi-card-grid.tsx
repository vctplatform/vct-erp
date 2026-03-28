"use client";

import { ArrowDownRight, ArrowUpRight, Minus } from "lucide-react";

import { useLocale } from "@/components/i18n/locale-provider";
import { AnimatedNumber } from "@/components/dashboard/animated-number";
import { Sparkline } from "@/components/dashboard/sparkline";
import { Card } from "@/components/ui/card";
import type { FinanceDashboardCard } from "@/lib/contracts/finance";
import { formatCompactCurrency, formatPercent, formatRunway } from "@/lib/formatters";
import {
  translateFinanceStatus,
  translateFinanceTrendPeriod,
} from "@/lib/i18n/finance";

type KpiCardGridProps = {
  cards: FinanceDashboardCard[];
};

function formatCardValue(
  card: FinanceDashboardCard,
  value: number,
  locale: "vi" | "en",
) {
  if (card.unit === "months") {
    return formatRunway(value, locale);
  }

  return formatCompactCurrency(value, locale);
}

function iconForTrend(direction: FinanceDashboardCard["trend"]["direction"]) {
  switch (direction) {
    case "up":
      return ArrowUpRight;
    case "down":
      return ArrowDownRight;
    default:
      return Minus;
  }
}

export function KpiCardGrid({ cards }: KpiCardGridProps) {
  const { locale } = useLocale();

  return (
    <div className="grid gap-4 xl:grid-cols-3">
      {cards.map((card) => {
        const TrendIcon = iconForTrend(card.trend.direction);
        const trendColor =
          card.trend.direction === "up"
            ? "text-emerald-500"
            : card.trend.direction === "down"
              ? "text-rose-500"
              : "text-[var(--color-ink-soft)]";
        const sparkColor =
          card.trend.direction === "down" ? "#e11d48" : "#10b981";

        return (
          <Card key={card.key} className="overflow-hidden p-5">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-sm font-medium text-[var(--color-ink-soft)]">
                  {card.title}
                </p>
                <div className="mt-3 text-3xl font-semibold tracking-tight text-[var(--color-ink)]">
                  <AnimatedNumber
                    value={card.value}
                    formatter={(value) => formatCardValue(card, value, locale)}
                  />
                </div>
              </div>
              {card.status ? (
                <span className="rounded-full border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-[var(--color-ink-soft)]">
                  {translateFinanceStatus(card.status, locale)}
                </span>
              ) : null}
            </div>

            <p className="mt-3 text-sm text-[var(--color-ink-soft)]">
              {card.description}
            </p>

            <div className="mt-4 flex items-center justify-between gap-4">
              <div className={`inline-flex items-center gap-1.5 text-sm font-medium ${trendColor}`}>
                <TrendIcon className="size-4" />
                <span>{formatPercent(card.trend.percentage)}</span>
                <span className="text-[var(--color-ink-soft)]">
                  {translateFinanceTrendPeriod(card.trend.period, locale)}
                </span>
              </div>
              <span className="font-mono text-xs text-[var(--color-ink-soft)]">
                {card.trend.delta >= 0 ? "+" : ""}
                {card.unit === "months"
                  ? `${card.trend.delta.toFixed(1)}`
                  : formatCompactCurrency(card.trend.delta, locale)}
              </span>
            </div>

            <div className="mt-4">
              <Sparkline points={card.chart_data} stroke={sparkColor} />
            </div>
          </Card>
        );
      })}
    </div>
  );
}
