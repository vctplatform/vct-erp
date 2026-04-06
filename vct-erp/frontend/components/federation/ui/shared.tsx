"use client";

import { cn } from "@/lib/utils";
import { STATUS_COLORS, type StatusKey } from "@/lib/federation/constants";
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, Search, X, SlidersHorizontal } from "lucide-react";

/* ═══════════════════════════════════════════════════════
 * StatusBadge
 * ═══════════════════════════════════════════════════════ */
interface StatusBadgeProps {
  status: string;
  label?: string;
  size?: "sm" | "md";
  className?: string;
}

export function StatusBadge({ status, label, size = "sm", className }: StatusBadgeProps) {
  const key = status as StatusKey;
  const colors = STATUS_COLORS[key] ?? STATUS_COLORS.active;

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full font-medium",
        colors.bg,
        colors.text,
        size === "sm" ? "px-2 py-0.5 text-[0.65rem]" : "px-3 py-1 text-xs",
        className
      )}
    >
      <span className={cn("size-1.5 rounded-full", colors.dot)} />
      {label ?? status.replace(/_/g, " ")}
    </span>
  );
}

/* ═══════════════════════════════════════════════════════
 * SectionPanel — Glass container for content sections
 * ═══════════════════════════════════════════════════════ */
interface SectionPanelProps {
  title: string;
  subtitle?: string;
  kicker?: string;
  children: React.ReactNode;
  actions?: React.ReactNode;
  className?: string;
  noPadding?: boolean;
}

export function SectionPanel({ title, subtitle, kicker, children, actions, className, noPadding }: SectionPanelProps) {
  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-2xl border border-[var(--color-border)] bg-[var(--color-panel)] transition-all duration-300",
        "dark:border-white/6 dark:bg-white/[0.02] dark:backdrop-blur-xl",
        "neon-hover",
        !noPadding && "p-5 md:p-6",
        className
      )}
    >
      {/* Top highlight line — dark mode only */}
      <div className="pointer-events-none absolute inset-x-0 top-0 hidden h-px bg-gradient-to-r from-transparent via-white/8 to-transparent dark:block" />

      <div className={cn("mb-5 flex items-start justify-between gap-4", noPadding && "px-5 pt-5 md:px-6 md:pt-6")}>
        <div>
          {kicker && (
            <p className="mb-1.5 text-[0.62rem] uppercase tracking-[0.3em] text-[var(--color-ink-soft)] dark:text-white/35">
              {kicker}
            </p>
          )}
          <h3 className="text-lg font-semibold text-[var(--color-ink)] dark:text-white">
            {title}
          </h3>
          {subtitle && (
            <p className="mt-1 text-sm text-[var(--color-ink-soft)] dark:text-white/50">
              {subtitle}
            </p>
          )}
        </div>
        {actions && <div className="flex items-center gap-2">{actions}</div>}
      </div>
      <div className={cn(noPadding && "px-5 pb-5 md:px-6 md:pb-6")}>
        {children}
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * FederationDataTable — Full-featured data table
 * ═══════════════════════════════════════════════════════ */
interface DataTableColumn {
  key: string;
  label: string;
  align?: "left" | "center" | "right";
  width?: string;
}

interface DataTableProps {
  columns: DataTableColumn[];
  rows: Record<string, React.ReactNode>[];
  compact?: boolean;
  onRowClick?: (row: Record<string, React.ReactNode>, index: number) => void;
  className?: string;
  zebraStripe?: boolean;
}

export function FederationDataTable({ columns, rows, compact, onRowClick, className, zebraStripe = true }: DataTableProps) {
  return (
    <div className={cn("overflow-x-auto", className)}>
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-[var(--color-border)] dark:border-white/6">
            {columns.map((col) => (
              <th
                key={col.key}
                className={cn(
                  "whitespace-nowrap px-4 py-3 text-[0.68rem] font-semibold uppercase tracking-[0.2em] text-[var(--color-ink-soft)] dark:text-white/40",
                  col.align === "right" && "text-right",
                  col.align === "center" && "text-center"
                )}
                style={col.width ? { width: col.width } : undefined}
              >
                {col.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr
              key={i}
              onClick={() => onRowClick?.(row, i)}
              className={cn(
                "border-b border-[var(--color-border)] transition-colors duration-200 last:border-b-0 dark:border-white/4",
                onRowClick && "cursor-pointer",
                // Zebra stripe
                zebraStripe && i % 2 === 1 && "bg-[var(--color-canvas-soft)]/50 dark:bg-white/[0.01]",
                // Hover glow
                "hover:bg-[var(--color-canvas-soft)] dark:hover:bg-white/[0.03]",
                compact ? "text-sm" : ""
              )}
            >
              {columns.map((col) => (
                <td
                  key={col.key}
                  className={cn(
                    "whitespace-nowrap px-4 text-[var(--color-ink)] dark:text-white/80",
                    compact ? "py-2.5" : "py-3.5",
                    col.align === "right" && "text-right",
                    col.align === "center" && "text-center"
                  )}
                >
                  {row[col.key]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
      {rows.length === 0 && (
        <div className="py-12 text-center text-sm text-[var(--color-ink-soft)] dark:text-white/40">
          Không có dữ liệu
        </div>
      )}
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * EmptyState
 * ═══════════════════════════════════════════════════════ */
interface EmptyStateProps {
  icon: React.ElementType;
  title: string;
  description?: string;
  action?: React.ReactNode;
}

export function EmptyState({ icon: Icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="rounded-2xl bg-[var(--color-canvas-soft)] p-4 dark:bg-white/5">
        <Icon className="size-8 text-[var(--color-ink-soft)] dark:text-white/30" />
      </div>
      <h3 className="mt-4 text-lg font-semibold text-[var(--color-ink)] dark:text-white">
        {title}
      </h3>
      {description && (
        <p className="mt-2 max-w-sm text-sm text-[var(--color-ink-soft)] dark:text-white/50">
          {description}
        </p>
      )}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * PageHeader — Hero section for each page
 * ═══════════════════════════════════════════════════════ */
export function PageHeader({
  kicker,
  title,
  description,
  actions,
}: {
  kicker?: string;
  title: string;
  description?: string;
  actions?: React.ReactNode;
}) {
  return (
    <div className="animate-fade-in mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        {kicker && (
          <p className="mb-1.5 text-[0.62rem] uppercase tracking-[0.3em] text-[var(--color-ink-soft)] dark:text-cyan-400/60">
            {kicker}
          </p>
        )}
        <h1 className="text-2xl font-bold tracking-tight text-[var(--color-ink)] dark:text-white md:text-3xl">
          {title}
        </h1>
        {description && (
          <p className="mt-2 max-w-2xl text-sm leading-relaxed text-[var(--color-ink-soft)] dark:text-white/50">
            {description}
          </p>
        )}
      </div>
      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * GlassButton — Premium button with variants
 * ═══════════════════════════════════════════════════════ */
export function GlassButton({
  children,
  variant = "primary",
  size = "md",
  className,
  ...props
}: {
  children: React.ReactNode;
  variant?: "primary" | "secondary" | "ghost" | "danger";
  size?: "sm" | "md" | "lg";
  className?: string;
} & React.ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <button
      className={cn(
        "inline-flex items-center justify-center gap-2 rounded-xl font-medium transition-all duration-200 active:scale-[0.97]",
        // Size
        size === "sm" && "px-3 py-1.5 text-xs",
        size === "md" && "px-4 py-2.5 text-sm",
        size === "lg" && "px-6 py-3 text-sm",
        // Variant
        variant === "primary" &&
          "bg-gradient-to-r from-cyan-500 to-teal-500 text-white shadow-[0_8px_25px_rgba(6,182,212,0.25)] hover:-translate-y-0.5 hover:shadow-[0_12px_35px_rgba(6,182,212,0.35)]",
        variant === "secondary" &&
          "border border-[var(--color-border)] bg-[var(--color-panel)] text-[var(--color-ink)] hover:bg-[var(--color-canvas-soft)] dark:border-white/10 dark:bg-white/5 dark:text-white dark:hover:bg-white/8",
        variant === "ghost" &&
          "text-[var(--color-ink-soft)] hover:bg-[var(--color-canvas-soft)] hover:text-[var(--color-ink)] dark:text-white/50 dark:hover:bg-white/5 dark:hover:text-white",
        variant === "danger" &&
          "bg-rose-500/10 text-rose-600 hover:bg-rose-500/20 dark:text-rose-400",
        className
      )}
      {...props}
    >
      {children}
    </button>
  );
}

/* ═══════════════════════════════════════════════════════
 * Tabs — Segmented control
 * ═══════════════════════════════════════════════════════ */
export function Tabs({
  tabs,
  activeTab,
  onChange,
}: {
  tabs: { id: string; label: string; count?: number }[];
  activeTab: string;
  onChange: (id: string) => void;
}) {
  return (
    <div className="flex gap-1 rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-1 dark:border-white/6 dark:bg-white/[0.03]">
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => onChange(tab.id)}
          className={cn(
            "flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-all duration-200",
            activeTab === tab.id
              ? "bg-[var(--color-panel)] text-[var(--color-ink)] shadow-sm dark:bg-white/8 dark:text-white"
              : "text-[var(--color-ink-soft)] hover:text-[var(--color-ink)] dark:text-white/40 dark:hover:text-white/70"
          )}
        >
          {tab.label}
          {tab.count !== undefined && (
            <span
              className={cn(
                "rounded-full px-1.5 py-0.5 text-[0.6rem] font-semibold",
                activeTab === tab.id
                  ? "bg-cyan-500/15 text-cyan-600 dark:text-cyan-400"
                  : "bg-[var(--color-border)] text-[var(--color-ink-soft)] dark:bg-white/6 dark:text-white/40"
              )}
            >
              {tab.count}
            </span>
          )}
        </button>
      ))}
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * SearchInput — Reusable search field
 * ═══════════════════════════════════════════════════════ */
export function SearchInput({
  value,
  onChange,
  placeholder = "Tìm kiếm...",
  className,
}: {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex items-center gap-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-3 py-2 transition-colors focus-within:border-cyan-500/40 dark:border-white/8 dark:bg-white/5 dark:focus-within:border-cyan-400/30",
        className
      )}
    >
      <Search className="size-4 shrink-0 text-[var(--color-ink-soft)] dark:text-white/40" />
      <input
        type="text"
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-full min-w-0 bg-transparent text-sm text-[var(--color-ink)] outline-none placeholder:text-[var(--color-ink-soft)] dark:text-white dark:placeholder:text-white/30"
      />
      {value && (
        <button
          onClick={() => onChange("")}
          className="rounded p-0.5 hover:bg-[var(--color-border)] dark:hover:bg-white/8"
        >
          <X className="size-3 text-[var(--color-ink-soft)] dark:text-white/40" />
        </button>
      )}
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * FilterChip — Small filter tag
 * ═══════════════════════════════════════════════════════ */
export function FilterChip({
  label,
  selected,
  onClick,
}: {
  label: string;
  selected: boolean;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        "inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium transition-all duration-200",
        selected
          ? "border border-cyan-500/25 bg-cyan-500/10 text-cyan-700 dark:border-cyan-400/20 dark:text-cyan-300"
          : "border border-[var(--color-border)] bg-[var(--color-canvas-soft)] text-[var(--color-ink-soft)] hover:bg-[var(--color-border)] dark:border-white/6 dark:bg-white/[0.03] dark:text-white/50 dark:hover:bg-white/6"
      )}
    >
      {label}
      {selected && <X className="size-3" />}
    </button>
  );
}

/* ═══════════════════════════════════════════════════════
 * FilterBar — Composed filter bar
 * ═══════════════════════════════════════════════════════ */
export function FilterBar({
  children,
  onReset,
  activeCount = 0,
}: {
  children: React.ReactNode;
  onReset?: () => void;
  activeCount?: number;
}) {
  return (
    <div className="flex flex-wrap items-center gap-2">
      <div className="flex items-center gap-1.5 text-xs font-medium text-[var(--color-ink-soft)] dark:text-white/40">
        <SlidersHorizontal className="size-3.5" />
        Bộ lọc
        {activeCount > 0 && (
          <span className="rounded-full bg-cyan-500/15 px-1.5 py-0.5 text-[0.6rem] font-bold text-cyan-600 dark:text-cyan-400">
            {activeCount}
          </span>
        )}
      </div>
      {children}
      {onReset && activeCount > 0 && (
        <button
          onClick={onReset}
          className="text-xs text-[var(--color-ink-soft)] underline decoration-dotted hover:text-[var(--color-ink)] dark:text-white/40 dark:hover:text-white/70"
        >
          Xóa tất cả
        </button>
      )}
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * Pagination — Reusable pagination
 * ═══════════════════════════════════════════════════════ */
export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  totalItems,
  pageSize,
}: {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  totalItems?: number;
  pageSize?: number;
}) {
  if (totalPages <= 1) return null;

  const getVisiblePages = () => {
    const pages: (number | "...")[] = [];
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
    } else {
      pages.push(1);
      if (currentPage > 3) pages.push("...");
      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);
      for (let i = start; i <= end; i++) pages.push(i);
      if (currentPage < totalPages - 2) pages.push("...");
      pages.push(totalPages);
    }
    return pages;
  };

  return (
    <div className="flex flex-col items-center gap-3 pt-4 sm:flex-row sm:justify-between">
      {totalItems !== undefined && pageSize !== undefined && (
        <p className="text-xs text-[var(--color-ink-soft)] dark:text-white/40">
          Hiển thị {Math.min((currentPage - 1) * pageSize + 1, totalItems)}–{Math.min(currentPage * pageSize, totalItems)} / {totalItems.toLocaleString("vi-VN")} kết quả
        </p>
      )}
      <div className="flex items-center gap-1">
        <button
          onClick={() => onPageChange(1)}
          disabled={currentPage === 1}
          className="rounded-lg p-1.5 text-[var(--color-ink-soft)] transition-colors hover:bg-[var(--color-canvas-soft)] disabled:opacity-30 dark:text-white/40 dark:hover:bg-white/5"
        >
          <ChevronsLeft className="size-4" />
        </button>
        <button
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className="rounded-lg p-1.5 text-[var(--color-ink-soft)] transition-colors hover:bg-[var(--color-canvas-soft)] disabled:opacity-30 dark:text-white/40 dark:hover:bg-white/5"
        >
          <ChevronLeft className="size-4" />
        </button>
        {getVisiblePages().map((page, i) =>
          page === "..." ? (
            <span key={`dots-${i}`} className="px-1 text-xs text-[var(--color-ink-soft)] dark:text-white/30">
              •••
            </span>
          ) : (
            <button
              key={page}
              onClick={() => onPageChange(page)}
              className={cn(
                "flex h-8 w-8 items-center justify-center rounded-lg text-xs font-medium transition-all duration-200",
                currentPage === page
                  ? "bg-gradient-to-r from-cyan-500 to-teal-500 text-white shadow-[0_4px_12px_rgba(6,182,212,0.3)]"
                  : "text-[var(--color-ink-soft)] hover:bg-[var(--color-canvas-soft)] dark:text-white/50 dark:hover:bg-white/5"
              )}
            >
              {page}
            </button>
          )
        )}
        <button
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="rounded-lg p-1.5 text-[var(--color-ink-soft)] transition-colors hover:bg-[var(--color-canvas-soft)] disabled:opacity-30 dark:text-white/40 dark:hover:bg-white/5"
        >
          <ChevronRight className="size-4" />
        </button>
        <button
          onClick={() => onPageChange(totalPages)}
          disabled={currentPage === totalPages}
          className="rounded-lg p-1.5 text-[var(--color-ink-soft)] transition-colors hover:bg-[var(--color-canvas-soft)] disabled:opacity-30 dark:text-white/40 dark:hover:bg-white/5"
        >
          <ChevronsRight className="size-4" />
        </button>
      </div>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * ViewToggle — Table/Grid view selector
 * ═══════════════════════════════════════════════════════ */
export function ViewToggle({
  view,
  onChange,
}: {
  view: "table" | "grid";
  onChange: (view: "table" | "grid") => void;
}) {
  return (
    <div className="flex rounded-lg border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-0.5 dark:border-white/6 dark:bg-white/[0.03]">
      <button
        onClick={() => onChange("table")}
        className={cn(
          "rounded-md px-2.5 py-1.5 text-xs font-medium transition-all",
          view === "table"
            ? "bg-[var(--color-panel)] text-[var(--color-ink)] shadow-sm dark:bg-white/8 dark:text-white"
            : "text-[var(--color-ink-soft)] dark:text-white/40"
        )}
      >
        <svg className="size-4" viewBox="0 0 16 16" fill="currentColor"><path d="M0 1h16v2H0V1zm0 4h16v2H0V5zm0 4h16v2H0V9zm0 4h16v2H0v-2z"/></svg>
      </button>
      <button
        onClick={() => onChange("grid")}
        className={cn(
          "rounded-md px-2.5 py-1.5 text-xs font-medium transition-all",
          view === "grid"
            ? "bg-[var(--color-panel)] text-[var(--color-ink)] shadow-sm dark:bg-white/8 dark:text-white"
            : "text-[var(--color-ink-soft)] dark:text-white/40"
        )}
      >
        <svg className="size-4" viewBox="0 0 16 16" fill="currentColor"><path d="M0 0h7v7H0V0zm9 0h7v7H9V0zM0 9h7v7H0V9zm9 0h7v7H9V9z"/></svg>
      </button>
    </div>
  );
}

/* ═══════════════════════════════════════════════════════
 * ProgressBar — Mini progress indicator
 * ═══════════════════════════════════════════════════════ */
export function ProgressBar({
  value,
  max = 100,
  color = "cyan",
  className,
}: {
  value: number;
  max?: number;
  color?: "cyan" | "emerald" | "amber" | "rose";
  className?: string;
}) {
  const percentage = Math.min((value / max) * 100, 100);
  const gradients = {
    cyan: "from-cyan-500 to-teal-500",
    emerald: "from-emerald-500 to-green-500",
    amber: "from-amber-500 to-orange-500",
    rose: "from-rose-500 to-pink-500",
  };

  return (
    <div className={cn("h-2 overflow-hidden rounded-full bg-[var(--color-canvas-soft)] dark:bg-white/6", className)}>
      <div
        className={cn(
          "h-full rounded-full bg-gradient-to-r transition-all duration-700 ease-out",
          gradients[color]
        )}
        style={{ width: `${percentage}%` }}
      />
    </div>
  );
}
