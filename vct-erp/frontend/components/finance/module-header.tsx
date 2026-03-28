import { StatusPill } from "@/components/finance/status-pill";
import { cn } from "@/lib/utils";

export function ModuleHeader({
  kicker,
  title,
  description,
  mode,
  className,
  actions,
}: {
  kicker: string;
  title: string;
  description: string;
  mode?: {
    label: string;
    tone?: "navy" | "emerald" | "amber" | "rose";
  };
  className?: string;
  actions?: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-[1.85rem] border border-[var(--color-border)] bg-[radial-gradient(circle_at_top_left,rgba(23,59,112,0.18),transparent_38%),linear-gradient(180deg,rgba(255,255,255,0.88),rgba(255,255,255,0.7))] p-6 shadow-[0_18px_56px_rgba(10,24,50,0.08)] dark:bg-[radial-gradient(circle_at_top_left,rgba(14,165,233,0.15),transparent_35%),linear-gradient(180deg,rgba(11,22,42,0.94),rgba(11,22,42,0.82))] md:p-7",
        className,
      )}
    >
      <div className="absolute -right-12 top-0 h-32 w-32 rounded-full bg-[rgba(16,185,129,0.14)] blur-3xl" />
      <div className="absolute left-0 top-12 h-24 w-24 rounded-full bg-[rgba(23,59,112,0.14)] blur-3xl" />

      <div className="relative flex flex-col gap-5 md:flex-row md:items-end md:justify-between">
        <div className="max-w-3xl">
          <p className="text-xs font-medium uppercase tracking-[0.32em] text-[var(--color-ink-soft)]">
            {kicker}
          </p>
          <h1 className="mt-3 text-3xl font-semibold tracking-tight text-[var(--color-ink)] md:text-5xl">
            {title}
          </h1>
          <p className="mt-3 max-w-2xl text-sm leading-6 text-[var(--color-ink-soft)] md:text-base">
            {description}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-3">
          {mode ? <StatusPill tone={mode.tone}>{mode.label}</StatusPill> : null}
          {actions}
        </div>
      </div>
    </div>
  );
}
