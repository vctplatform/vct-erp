import { ArrowDownRight, ArrowUpRight, Minus } from "lucide-react";

import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";

type InsightTone = "navy" | "emerald" | "amber" | "rose";

const toneClasses: Record<InsightTone, string> = {
  navy: "border-[rgba(23,59,112,0.14)] bg-[linear-gradient(180deg,rgba(23,59,112,0.08),transparent)]",
  emerald:
    "border-emerald-500/20 bg-[linear-gradient(180deg,rgba(16,185,129,0.12),transparent)]",
  amber:
    "border-amber-500/20 bg-[linear-gradient(180deg,rgba(245,158,11,0.12),transparent)]",
  rose: "border-rose-500/20 bg-[linear-gradient(180deg,rgba(225,29,72,0.12),transparent)]",
};

const directionIcons = {
  up: ArrowUpRight,
  down: ArrowDownRight,
  flat: Minus,
};

export function InsightCard({
  label,
  value,
  caption,
  tone = "navy",
  trend,
}: {
  label: string;
  value: string;
  caption: string;
  tone?: InsightTone;
  trend?: {
    direction: "up" | "down" | "flat";
    label: string;
  };
}) {
  const TrendIcon = trend ? directionIcons[trend.direction] : null;

  return (
    <Card className={cn("p-5", toneClasses[tone])}>
      <p className="text-xs font-medium uppercase tracking-[0.24em] text-[var(--color-ink-soft)]">
        {label}
      </p>
      <div className="mt-4 flex items-start justify-between gap-3">
        <p className="text-2xl font-semibold tracking-tight text-[var(--color-ink)] md:text-3xl">
          {value}
        </p>
        {TrendIcon && trend ? (
          <span className="inline-flex items-center gap-1 rounded-full border border-[var(--color-border)] bg-[var(--color-panel)] px-2.5 py-1 text-xs text-[var(--color-ink-soft)]">
            <TrendIcon className="size-3.5" />
            {trend.label}
          </span>
        ) : null}
      </div>
      <p className="mt-3 text-sm leading-6 text-[var(--color-ink-soft)]">
        {caption}
      </p>
    </Card>
  );
}
