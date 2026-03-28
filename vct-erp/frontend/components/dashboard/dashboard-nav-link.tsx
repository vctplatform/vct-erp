"use client";

import {
  BookText,
  ChartPie,
  ChevronRight,
  LayoutDashboard,
  Radar,
  RefreshCcw,
  Scale,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

import type { FinanceNavIcon } from "@/lib/finance/navigation";
import { cn } from "@/lib/utils";

const iconMap: Record<FinanceNavIcon, React.ComponentType<{ className?: string }>> = {
  command: LayoutDashboard,
  segments: ChartPie,
  ledger: BookText,
  reconciliation: RefreshCcw,
  reports: Scale,
  control: Radar,
};

function isActivePath(pathname: string, href: string) {
  if (href === "/") {
    return pathname === "/";
  }

  return pathname === href || pathname.startsWith(`${href}/`);
}

export function DashboardNavLink({
  href,
  label,
  caption,
  icon,
  variant = "sidebar",
}: {
  href: string;
  label: string;
  caption: string;
  icon: FinanceNavIcon;
  variant?: "sidebar" | "mobile";
}) {
  const pathname = usePathname();
  const active = isActivePath(pathname, href);
  const Icon = iconMap[icon];

  if (variant === "mobile") {
    return (
      <Link
        href={href}
        aria-current={active ? "page" : undefined}
        className={cn(
          "group flex min-w-[11rem] items-center gap-3 rounded-[1.3rem] border px-3 py-3.5 transition duration-300",
          active
            ? "border-[rgba(23,59,112,0.18)] bg-[linear-gradient(135deg,rgba(23,59,112,0.14),rgba(14,165,233,0.10))] text-[var(--color-ink)] shadow-[0_12px_32px_rgba(13,26,44,0.10)]"
            : "border-transparent bg-[var(--color-canvas-soft)] text-[var(--color-ink-soft)] hover:border-[rgba(23,59,112,0.14)] hover:bg-[var(--color-panel)] hover:text-[var(--color-ink)]",
        )}
      >
        <span
          className={cn(
            "inline-flex h-11 w-11 items-center justify-center rounded-[1rem] border shadow-[0_8px_22px_rgba(13,26,44,0.08)] transition",
            active
              ? "border-[rgba(23,59,112,0.16)] bg-[var(--color-panel)] text-[var(--color-navy-700)] dark:text-white"
              : "border-[var(--color-border)] bg-[var(--color-panel)] text-[var(--color-navy-700)] dark:text-white",
          )}
        >
          <Icon className="size-4" />
        </span>
        <span className="min-w-0">
          <span className="block truncate text-sm font-semibold text-[var(--color-ink)]">
            {label}
          </span>
          <span className="mt-0.5 block truncate text-xs text-[var(--color-ink-soft)]">
            {caption}
          </span>
        </span>
      </Link>
    );
  }

  return (
    <Link
      href={href}
      aria-current={active ? "page" : undefined}
      className={cn(
        "group relative block overflow-hidden rounded-[1.35rem] border px-3 py-3.5 transition duration-300",
        active
          ? "border-white/14 bg-[linear-gradient(135deg,rgba(255,255,255,0.14),rgba(255,255,255,0.06))] shadow-[0_14px_34px_rgba(5,12,24,0.28)]"
          : "border-transparent bg-transparent hover:border-white/10 hover:bg-white/6",
      )}
    >
      <div
        className={cn(
          "absolute inset-y-3 left-0 w-1 rounded-full transition",
          active ? "bg-emerald-300" : "bg-transparent group-hover:bg-white/18",
        )}
      />
      <div className="flex items-center gap-3">
        <span
          className={cn(
            "inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-[1rem] border text-white transition",
            active
              ? "border-white/16 bg-white/14 shadow-[0_10px_28px_rgba(10,24,50,0.28)]"
              : "border-white/8 bg-white/8 group-hover:border-white/14 group-hover:bg-white/10",
          )}
        >
          <Icon className="size-4" />
        </span>
        <span className="min-w-0 flex-1">
          <span className="block text-sm font-semibold text-white">
            {label}
          </span>
          <span
            className={cn(
              "mt-1 block text-xs leading-5 transition",
              active ? "text-white/78" : "text-white/52 group-hover:text-white/68",
            )}
          >
            {caption}
          </span>
        </span>
        <ChevronRight
          className={cn(
            "size-4 shrink-0 transition",
            active
              ? "translate-x-0 text-white/70"
              : "-translate-x-1 text-white/28 group-hover:translate-x-0 group-hover:text-white/52",
          )}
        />
      </div>
    </Link>
  );
}
