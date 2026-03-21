import { DashboardClient } from "@/components/dashboard/dashboard-client";
import { getFinanceDashboardSnapshot } from "@/lib/api/finance";

export default async function DashboardPage() {
  const snapshot = await getFinanceDashboardSnapshot();

  return (
    <section className="space-y-6">
      <div className="flex flex-col gap-2 md:flex-row md:items-end md:justify-between">
        <div>
          <p className="text-xs font-medium uppercase tracking-[0.28em] text-[var(--color-ink-soft)]">
            VCT Group Financial Intelligence
          </p>
          <h1 className="text-3xl font-semibold tracking-tight text-[var(--color-ink)] md:text-4xl">
            Command Center
          </h1>
        </div>
        <div className="rounded-full border border-[var(--color-border)] bg-[var(--color-panel)] px-4 py-2 text-sm text-[var(--color-ink-soft)]">
          Snapshot mode:{" "}
          <span className="font-semibold uppercase tracking-[0.22em] text-[var(--color-navy-700)] dark:text-white">
            {snapshot.data_mode}
          </span>
        </div>
      </div>

      <DashboardClient initialData={snapshot} />
    </section>
  );
}
