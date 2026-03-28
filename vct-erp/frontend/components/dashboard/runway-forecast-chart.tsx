"use client";

import { useDeferredValue } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { useLocale } from "@/components/i18n/locale-provider";
import { Card } from "@/components/ui/card";
import type { FinanceRunwayPoint } from "@/lib/contracts/finance";
import { formatCompactCurrency } from "@/lib/formatters";

type RunwayForecastChartProps = {
  data: FinanceRunwayPoint[];
};

export function RunwayForecastChart({ data }: RunwayForecastChartProps) {
  const { locale } = useLocale();
  const deferredData = useDeferredValue(data);
  const copy =
    locale === "vi"
      ? {
          kicker: "Dự báo runway",
          title: "Dự báo dòng tiền",
        }
      : {
          kicker: "Runway Forecast",
          title: "Cash Projection",
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

      <div className="h-80">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={deferredData}>
            <CartesianGrid stroke="rgba(95,114,151,0.16)" strokeDasharray="4 4" />
            <XAxis dataKey="label" tickLine={false} axisLine={false} />
            <YAxis
              tickFormatter={(value) => `${Math.round(value / 1_000_000)}M`}
              tickLine={false}
              axisLine={false}
              width={72}
            />
            <Tooltip
              formatter={(value) =>
                formatCompactCurrency(
                  Number(Array.isArray(value) ? value[0] : value ?? 0),
                  locale,
                )
              }
            />
            <Bar
              dataKey="projected_ending"
              fill="#173b70"
              radius={[12, 12, 4, 4]}
            />
            <Bar
              dataKey="contracted_inflow"
              fill="#10b981"
              radius={[12, 12, 4, 4]}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
