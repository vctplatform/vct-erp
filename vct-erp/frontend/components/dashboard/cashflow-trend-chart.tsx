"use client";

import { useDeferredValue } from "react";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { Card } from "@/components/ui/card";
import type { FinanceMultiLineChart } from "@/lib/contracts/finance";
import { formatCompactCurrency } from "@/lib/formatters";

type CashflowTrendChartProps = {
  chart: FinanceMultiLineChart;
};

function toChartRows(chart: FinanceMultiLineChart) {
  return chart.x_axis.map((label, index) => {
    const row: Record<string, string | number> = { label };
    for (const series of chart.series) {
      row[series.key] = series.values[index] ?? 0;
    }
    return row;
  });
}

export function CashflowTrendChart({ chart }: CashflowTrendChartProps) {
  const deferredChart = useDeferredValue(chart);
  const rows = toChartRows(deferredChart);

  return (
    <Card className="p-5">
      <div className="mb-5">
        <p className="text-xs font-medium uppercase tracking-[0.24em] text-[var(--color-ink-soft)]">
          Cashflow Trend
        </p>
        <h2 className="mt-2 text-xl font-semibold text-[var(--color-ink)]">
          Dòng tien 6 thang
        </h2>
      </div>

      <div className="h-80">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={rows}>
            <CartesianGrid stroke="rgba(95,114,151,0.16)" strokeDasharray="4 4" />
            <XAxis dataKey="label" tickLine={false} axisLine={false} />
            <YAxis
              tickFormatter={(value) => `${Math.round(value / 1_000_000)}M`}
              tickLine={false}
              axisLine={false}
              width={72}
            />
            <Tooltip
              formatter={(value: number) => formatCompactCurrency(value)}
            />
            <Legend />
            {deferredChart.series.map((series) => (
              <Line
                key={series.key}
                type="monotone"
                dataKey={series.key}
                stroke={series.color}
                strokeWidth={3}
                dot={false}
                activeDot={{ r: 5 }}
              />
            ))}
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
