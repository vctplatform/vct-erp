"use client";

import { useDeferredValue } from "react";
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts";

import { useLocale } from "@/components/i18n/locale-provider";
import { Card } from "@/components/ui/card";
import type { FinancePieSlice } from "@/lib/contracts/finance";
import { formatCompactCurrency } from "@/lib/formatters";
import { translateFinanceSegment } from "@/lib/i18n/finance";

type RevenueMixChartProps = {
  data: FinancePieSlice[];
};

export function RevenueMixChart({ data }: RevenueMixChartProps) {
  const { locale } = useLocale();
  const deferredData = useDeferredValue(data);
  const copy =
    locale === "vi"
      ? {
          kicker: "Cơ cấu doanh thu",
          title: "Phân bổ doanh thu theo mảng",
        }
      : {
          kicker: "Revenue Mix",
          title: "Revenue Composition by Segment",
        };

  return (
    <Card className="p-5">
      <div className="mb-5 flex items-center justify-between">
        <div>
          <p className="text-xs font-medium uppercase tracking-[0.24em] text-[var(--color-ink-soft)]">
            {copy.kicker}
          </p>
          <h2 className="mt-2 text-xl font-semibold text-[var(--color-ink)]">
            {copy.title}
          </h2>
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
        <div className="h-72">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={deferredData}
                dataKey="value"
                nameKey="label"
                innerRadius={78}
                outerRadius={112}
                paddingAngle={3}
              >
                {deferredData.map((entry) => (
                  <Cell key={entry.label} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip
                formatter={(value) =>
                  formatCompactCurrency(
                    Number(Array.isArray(value) ? value[0] : value ?? 0),
                    locale,
                  )
                }
                labelFormatter={(value) =>
                  translateFinanceSegment(String(value), locale) ?? String(value)
                }
              />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="space-y-3">
          {deferredData.map((slice) => (
            <div
              key={slice.label}
              className="rounded-2xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4"
            >
              <div className="flex items-center justify-between gap-3">
                <div className="flex items-center gap-3">
                  <span
                    className="size-3 rounded-full"
                    style={{ backgroundColor: slice.color }}
                  />
                  <span className="font-medium text-[var(--color-ink)]">
                    {translateFinanceSegment(slice.label, locale)}
                  </span>
                </div>
                <span className="font-mono text-sm text-[var(--color-ink-soft)]">
                  {formatCompactCurrency(slice.value, locale)}
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </Card>
  );
}
