import { cn } from "@/lib/utils";

type FinanceTableColumn = {
  key: string;
  label: string;
  align?: "left" | "right";
};

type FinanceTableRow = Record<string, React.ReactNode>;

export function DataTable({
  columns,
  rows,
  compact = false,
}: {
  columns: FinanceTableColumn[];
  rows: FinanceTableRow[];
  compact?: boolean;
}) {
  return (
    <div className="overflow-x-auto">
      <table className="min-w-full border-separate border-spacing-0">
        <thead>
          <tr>
            {columns.map((column) => (
              <th
                key={column.key}
                className={cn(
                  "border-b border-[var(--color-border)] px-4 py-3 text-xs font-semibold uppercase tracking-[0.2em] text-[var(--color-ink-soft)]",
                  column.align === "right" ? "text-right" : "text-left",
                  compact ? "first:pl-0 last:pr-0" : "",
                )}
              >
                {column.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={index}>
              {columns.map((column) => (
                <td
                  key={column.key}
                  className={cn(
                    "border-b border-[var(--color-border)] px-4 py-4 text-sm text-[var(--color-ink)]",
                    column.align === "right" ? "text-right" : "text-left",
                    compact ? "first:pl-0 last:pr-0" : "",
                  )}
                >
                  {row[column.key]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
