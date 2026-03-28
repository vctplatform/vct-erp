import { cn } from "@/lib/utils";

type StatusPillTone = "navy" | "emerald" | "amber" | "rose";

const toneClasses: Record<StatusPillTone, string> = {
  navy: "border-[rgba(23,59,112,0.16)] bg-[rgba(23,59,112,0.08)] text-[var(--color-navy-700)] dark:text-white",
  emerald:
    "border-emerald-500/20 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300",
  amber:
    "border-amber-500/25 bg-amber-500/12 text-amber-700 dark:text-amber-200",
  rose: "border-rose-500/20 bg-rose-500/12 text-rose-700 dark:text-rose-200",
};

export function StatusPill({
  children,
  tone = "navy",
  className,
}: {
  children: React.ReactNode;
  tone?: StatusPillTone;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em]",
        toneClasses[tone],
        className,
      )}
    >
      {children}
    </span>
  );
}
