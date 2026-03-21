export default function DashboardLoading() {
  return (
    <div className="space-y-6">
      <div className="h-20 animate-pulse rounded-[1.5rem] bg-[var(--color-panel)]" />
      <div className="grid gap-4 xl:grid-cols-3">
        <div className="h-56 animate-pulse rounded-[1.5rem] bg-[var(--color-panel)]" />
        <div className="h-56 animate-pulse rounded-[1.5rem] bg-[var(--color-panel)]" />
        <div className="h-56 animate-pulse rounded-[1.5rem] bg-[var(--color-panel)]" />
      </div>
      <div className="grid gap-6 xl:grid-cols-[1.05fr_1.3fr]">
        <div className="h-96 animate-pulse rounded-[1.5rem] bg-[var(--color-panel)]" />
        <div className="h-96 animate-pulse rounded-[1.5rem] bg-[var(--color-panel)]" />
      </div>
    </div>
  );
}
