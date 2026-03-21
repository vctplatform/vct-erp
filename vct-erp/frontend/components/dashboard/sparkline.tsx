"use client";

import type { FinanceMiniChartPoint } from "@/lib/contracts/finance";

type SparklineProps = {
  points: FinanceMiniChartPoint[];
  stroke: string;
};

export function Sparkline({ points, stroke }: SparklineProps) {
  if (points.length === 0) {
    return <div className="h-12 rounded-2xl bg-[var(--color-canvas-soft)]" />;
  }

  const width = 240;
  const height = 56;
  const values = points.map((point) => point.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const step = width / Math.max(points.length - 1, 1);

  const path = points
    .map((point, index) => {
      const range = max - min || 1;
      const x = index * step;
      const y = height - ((point.value - min) / range) * (height - 10) - 5;
      return `${index === 0 ? "M" : "L"} ${x.toFixed(1)} ${y.toFixed(1)}`;
    })
    .join(" ");

  return (
    <svg
      viewBox={`0 0 ${width} ${height}`}
      className="h-12 w-full overflow-visible"
      preserveAspectRatio="none"
      aria-hidden="true"
    >
      <path
        d={path}
        fill="none"
        stroke={stroke}
        strokeWidth="3"
        strokeLinecap="round"
      />
    </svg>
  );
}
