import { cn } from "@/lib/utils";

export function Card({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "rounded-[var(--radius-card)] border border-[var(--color-border)] bg-[var(--color-panel)] shadow-[0_14px_48px_rgba(13,26,44,0.08)]",
        className,
      )}
      {...props}
    />
  );
}
