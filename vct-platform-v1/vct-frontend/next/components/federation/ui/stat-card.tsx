"use client";

import { cn } from "@/lib/utils";
import type { LucideIcon } from "lucide-react";
import { TrendingUp, TrendingDown, Minus } from "lucide-react";
import { useEffect, useRef, useState } from "react";

interface StatCardProps {
  icon: LucideIcon;
  label: string;
  value: number | string;
  suffix?: string;
  trend?: { direction: "up" | "down" | "stable"; value: string };
  color?: "cyan" | "emerald" | "amber" | "rose" | "indigo" | "sky";
  sparkline?: number[];
  className?: string;
}

const COLORS = {
  cyan: {
    iconBg: "bg-cyan-500/10 dark:bg-cyan-500/15",
    iconText: "text-cyan-600 dark:text-cyan-400",
    glow: "dark:group-hover:shadow-[0_0_30px_rgba(6,182,212,0.18)]",
    borderGlow: "dark:group-hover:border-cyan-500/20",
    sparkStroke: "#06B6D4",
    sparkFill: "rgba(6, 182, 212, 0.15)",
  },
  emerald: {
    iconBg: "bg-emerald-500/10 dark:bg-emerald-500/15",
    iconText: "text-emerald-600 dark:text-emerald-400",
    glow: "dark:group-hover:shadow-[0_0_30px_rgba(16,185,129,0.18)]",
    borderGlow: "dark:group-hover:border-emerald-500/20",
    sparkStroke: "#10B981",
    sparkFill: "rgba(16, 185, 129, 0.15)",
  },
  amber: {
    iconBg: "bg-amber-500/10 dark:bg-amber-500/15",
    iconText: "text-amber-600 dark:text-amber-400",
    glow: "dark:group-hover:shadow-[0_0_30px_rgba(245,158,11,0.18)]",
    borderGlow: "dark:group-hover:border-amber-500/20",
    sparkStroke: "#F59E0B",
    sparkFill: "rgba(245, 158, 11, 0.15)",
  },
  rose: {
    iconBg: "bg-rose-500/10 dark:bg-rose-500/15",
    iconText: "text-rose-600 dark:text-rose-400",
    glow: "dark:group-hover:shadow-[0_0_30px_rgba(225,29,72,0.18)]",
    borderGlow: "dark:group-hover:border-rose-500/20",
    sparkStroke: "#E11D48",
    sparkFill: "rgba(225, 29, 72, 0.15)",
  },
  indigo: {
    iconBg: "bg-indigo-500/10 dark:bg-indigo-500/15",
    iconText: "text-indigo-600 dark:text-indigo-400",
    glow: "dark:group-hover:shadow-[0_0_30px_rgba(99,102,241,0.18)]",
    borderGlow: "dark:group-hover:border-indigo-500/20",
    sparkStroke: "#6366F1",
    sparkFill: "rgba(99, 102, 241, 0.15)",
  },
  sky: {
    iconBg: "bg-sky-500/10 dark:bg-sky-500/15",
    iconText: "text-sky-600 dark:text-sky-400",
    glow: "dark:group-hover:shadow-[0_0_30px_rgba(14,165,233,0.18)]",
    borderGlow: "dark:group-hover:border-sky-500/20",
    sparkStroke: "#0EA5E9",
    sparkFill: "rgba(14, 165, 233, 0.15)",
  },
};

/* ── Animated Number Counter ── */
function AnimatedNumber({ value }: { value: number }) {
  const [display, setDisplay] = useState(0);

  useEffect(() => {
    const duration = 1200;
    const start = performance.now();

    function update(now: number) {
      const elapsed = now - start;
      const progress = Math.min(elapsed / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3);
      setDisplay(Math.round((value) * eased));
      if (progress < 1) requestAnimationFrame(update);
    }

    requestAnimationFrame(update);
  }, [value]);

  return <span>{display.toLocaleString("vi-VN")}</span>;
}

/* ── Mini Sparkline SVG ── */
function MiniSparkline({
  data,
  stroke,
  fill,
}: {
  data: number[];
  stroke: string;
  fill: string;
}) {
  if (data.length < 2) return null;

  const w = 100;
  const h = 32;
  const max = Math.max(...data);
  const min = Math.min(...data);
  const range = max - min || 1;

  const points = data.map((v, i) => ({
    x: (i / (data.length - 1)) * w,
    y: h - ((v - min) / range) * (h - 4) - 2,
  }));

  const linePath = points.map((p, i) => `${i === 0 ? "M" : "L"} ${p.x} ${p.y}`).join(" ");
  const areaPath = `${linePath} L ${w} ${h} L 0 ${h} Z`;

  return (
    <svg viewBox={`0 0 ${w} ${h}`} className="h-8 w-full" preserveAspectRatio="none">
      <defs>
        <linearGradient id={`spark-${stroke}`} x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor={fill} stopOpacity={1} />
          <stop offset="100%" stopColor={fill} stopOpacity={0} />
        </linearGradient>
      </defs>
      <path d={areaPath} fill={`url(#spark-${stroke})`} />
      <path d={linePath} fill="none" stroke={stroke} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

/* ── StatCard Component ── */
export function StatCard({
  icon: Icon,
  label,
  value,
  suffix,
  trend,
  color = "cyan",
  sparkline,
  className,
}: StatCardProps) {
  const c = COLORS[color];

  return (
    <div
      className={cn(
        "group relative overflow-hidden rounded-2xl border bg-white/5 backdrop-blur-xl transition-all duration-300 hover:-translate-y-0.5",
        "dark:border-white/10 dark:bg-white/[0.03]",
        "group-hover:shadow-[0_12px_30px_rgba(0,0,0,0.15)]",
        c.glow,
        c.borderGlow,
        className
      )}
    >
      <div className="p-5 relative">
        <div className="flex items-center justify-between">
          <div className={cn("rounded-xl p-2.5", c.iconBg)}>
            <Icon className={cn("size-5", c.iconText)} />
          </div>
          {trend && (
            <div
              className={cn(
                "flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium",
                trend.direction === "up" && "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
                trend.direction === "down" && "bg-rose-500/10 text-rose-600 dark:text-rose-400",
                trend.direction === "stable" && "bg-slate-500/10 text-slate-600 dark:text-slate-400"
              )}
            >
              {trend.direction === "up" && <TrendingUp className="size-3" />}
              {trend.direction === "down" && <TrendingDown className="size-3" />}
              {trend.direction === "stable" && <Minus className="size-3" />}
              {trend.value}
            </div>
          )}
        </div>
        <div className="mt-4">
          <p className="text-3xl font-bold tracking-tight text-white dark:text-white">
            {typeof value === "number" ? <AnimatedNumber value={value} /> : value}
            {suffix && (
              <span className="ml-1.5 text-base font-medium text-white/50">
                {suffix}
              </span>
            )}
          </p>
          <p className="mt-1.5 text-sm text-white/60">
            {label}
          </p>
        </div>

        {sparkline && sparkline.length > 1 && (
          <div className="mt-3 opacity-60 transition-opacity group-hover:opacity-100">
            <MiniSparkline data={sparkline} stroke={c.sparkStroke} fill={c.sparkFill} />
          </div>
        )}
      </div>
    </div>
  );
}
