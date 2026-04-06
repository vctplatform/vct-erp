"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard, Users, Building2, UserCog, GraduationCap,
  Trophy, Award, Wallet, FileText, Megaphone, BarChart3, Settings,
  ChevronRight, Shield, Zap, Globe, Activity,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { FEDERATION_NAV } from "@/lib/federation/constants";

const ICON_MAP: Record<string, React.ElementType> = {
  LayoutDashboard, Users, Building2, UserCog, GraduationCap,
  Trophy, Award, Wallet, FileText, Megaphone, BarChart3, Settings,
};

export function FederationSidebar() {
  const pathname = usePathname();

  return (
    <aside className="sticky top-4 hidden h-[calc(100vh-2rem)] w-[18rem] shrink-0 lg:block">
      <div className="relative flex h-full flex-col overflow-hidden rounded-[2rem] border border-[var(--color-border)] bg-[var(--color-panel)] shadow-[0_24px_70px_rgba(11,22,42,0.08)] dark:border-white/6 dark:bg-[radial-gradient(circle_at_top_left,rgba(6,182,212,0.1),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(20,184,166,0.06),transparent_28%),linear-gradient(180deg,rgba(8,16,32,0.98),rgba(5,10,22,0.98))] dark:shadow-[0_30px_90px_rgba(0,0,0,0.5)]">
        {/* Ambient glow orbs — dark mode only */}
        <div className="pointer-events-none absolute -left-10 top-6 hidden h-32 w-32 rounded-full bg-cyan-500/12 blur-[60px] dark:block" />
        <div className="pointer-events-none absolute -right-8 bottom-16 hidden h-28 w-28 rounded-full bg-teal-500/8 blur-[60px] dark:block" />
        <div className="pointer-events-none absolute left-1/2 top-1/3 hidden h-20 w-20 -translate-x-1/2 rounded-full bg-indigo-500/6 blur-[50px] dark:block" />

        {/* Top highlight line */}
        <div className="pointer-events-none absolute inset-x-4 top-0 hidden h-px bg-gradient-to-r from-transparent via-cyan-500/20 to-transparent dark:block" />

        <div className="relative flex h-full flex-col overflow-y-auto px-4 py-5 scrollbar-hide">
          {/* ── Header Branding ──────────────────────── */}
          <div className="mb-5 animate-fade-in rounded-2xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4 dark:border-white/6 dark:bg-white/[0.03]">
            <div className="flex items-center gap-3">
              <span className="inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border border-[var(--color-border)] bg-gradient-to-br from-cyan-500/20 to-teal-500/10 shadow-[0_8px_20px_rgba(6,182,212,0.12)] dark:border-white/10">
                <Shield className="size-5 text-cyan-600 dark:text-cyan-400" />
              </span>
              <div className="min-w-0">
                <p className="text-[0.62rem] uppercase tracking-[0.3em] text-[var(--color-ink-soft)] dark:text-white/40">
                  Liên đoàn Quốc gia
                </p>
                <h2 className="mt-0.5 truncate text-sm font-semibold text-[var(--color-ink)] dark:text-white">
                  VCT Việt Nam
                </h2>
              </div>
            </div>
            <div className="mt-3 flex items-center gap-2">
              <div className="flex items-center gap-1.5 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2.5 py-1 text-[0.65rem] font-medium text-emerald-700 dark:border-emerald-400/15 dark:text-emerald-300">
                <Zap className="size-3" />
                Online
              </div>
              <div className="flex items-center gap-1.5 rounded-full border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-2.5 py-1 text-[0.65rem] font-medium text-[var(--color-ink-soft)] dark:border-white/6 dark:bg-white/[0.03] dark:text-white/50">
                <Globe className="size-3" />
                63 Tỉnh/TP
              </div>
            </div>
          </div>

          {/* ── Navigation Groups ────────────────────── */}
          <nav className="flex-1 space-y-4">
            {FEDERATION_NAV.map((group, gi) => (
              <div key={group.group} style={{ animationDelay: `${(gi + 1) * 80}ms` }} className="animate-fade-in">
                <p className="mb-2 px-2 text-[0.62rem] uppercase tracking-[0.28em] text-[var(--color-ink-soft)] dark:text-white/30">
                  {group.group}
                </p>
                <div className="space-y-0.5">
                  {group.items.map((item) => {
                    const Icon = ICON_MAP[item.icon];
                    const isActive =
                      item.href === "/federation"
                        ? pathname === "/federation"
                        : pathname.startsWith(item.href);
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        className={cn(
                          "group/link relative flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-200",
                          isActive
                            ? "border border-cyan-500/20 bg-cyan-500/10 text-cyan-700 shadow-[0_0_20px_rgba(6,182,212,0.06)] dark:border-cyan-400/15 dark:bg-cyan-500/10 dark:text-cyan-300 dark:shadow-[0_0_25px_rgba(6,182,212,0.12)]"
                            : "border border-transparent text-[var(--color-ink-soft)] hover:border-[var(--color-border)] hover:bg-[var(--color-canvas-soft)] dark:text-white/50 dark:hover:border-white/6 dark:hover:bg-white/[0.03]"
                        )}
                      >
                        {/* Active indicator bar */}
                        {isActive && (
                          <span className="absolute -left-0.5 top-1/2 h-5 w-[3px] -translate-y-1/2 rounded-full bg-gradient-to-b from-cyan-400 to-teal-500 shadow-[0_0_8px_rgba(6,182,212,0.5)]" />
                        )}
                        {Icon && (
                          <Icon
                            className={cn(
                              "size-4 shrink-0 transition-colors duration-200",
                              isActive
                                ? "text-cyan-600 dark:text-cyan-400"
                                : "text-[var(--color-ink-soft)] group-hover/link:text-[var(--color-ink)] dark:text-white/35 dark:group-hover/link:text-white/65"
                            )}
                          />
                        )}
                        <span className="flex-1 truncate">{item.label}</span>
                        {isActive && (
                          <ChevronRight className="size-3.5 text-cyan-500/50 dark:text-cyan-400/40" />
                        )}
                      </Link>
                    );
                  })}
                </div>
              </div>
            ))}
          </nav>

          {/* ── Footer Status ────────────────────────── */}
          <div className="mt-4 space-y-3">
            {/* Live cluster status */}
            <div className="rounded-xl border border-emerald-500/15 bg-emerald-500/5 p-3 dark:border-emerald-400/10 dark:bg-emerald-500/[0.06]">
              <div className="flex items-center gap-2">
                <Activity className="size-3.5 text-emerald-600 dark:text-emerald-400" />
                <span className="text-[0.62rem] uppercase tracking-[0.24em] text-emerald-600 dark:text-emerald-400">
                  Hệ thống
                </span>
                <span className="ml-auto rounded-full border border-emerald-400/15 bg-emerald-500/10 px-2 py-0.5 text-[0.6rem] font-medium text-emerald-600 dark:text-emerald-300">
                  Hoạt động
                </span>
              </div>
            </div>

            {/* Build info */}
            <div className="rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-3 text-xs dark:border-white/5 dark:bg-white/[0.02]">
              <p className="font-medium text-[var(--color-ink)] dark:text-white/75">
                VCT Platform v2.0
              </p>
              <p className="mt-1 text-[var(--color-ink-soft)] dark:text-white/35">
                Federation Module — Build 2026.04
              </p>
            </div>
          </div>
        </div>
      </div>
    </aside>
  );
}
