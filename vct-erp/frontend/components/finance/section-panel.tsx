import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";

export function SectionPanel({
  kicker,
  title,
  description,
  className,
  children,
  aside,
}: {
  kicker?: string;
  title: string;
  description?: string;
  className?: string;
  children: React.ReactNode;
  aside?: React.ReactNode;
}) {
  return (
    <Card className={cn("overflow-hidden p-5 md:p-6", className)}>
      <div className="flex flex-col gap-4 border-b border-[var(--color-border)] pb-5 md:flex-row md:items-start md:justify-between">
        <div>
          {kicker ? (
            <p className="text-xs font-medium uppercase tracking-[0.28em] text-[var(--color-ink-soft)]">
              {kicker}
            </p>
          ) : null}
          <h2 className="mt-2 text-2xl font-semibold tracking-tight text-[var(--color-ink)]">
            {title}
          </h2>
          {description ? (
            <p className="mt-2 max-w-2xl text-sm leading-6 text-[var(--color-ink-soft)]">
              {description}
            </p>
          ) : null}
        </div>
        {aside ? <div className="shrink-0">{aside}</div> : null}
      </div>
      <div className="pt-5">{children}</div>
    </Card>
  );
}
