import { BellDot, Search, ShieldCheck } from "lucide-react";

import { ThemeToggle } from "@/components/ui/theme-toggle";

export function Topbar() {
  return (
    <header className="flex flex-col gap-3 rounded-[1.5rem] border border-[var(--color-border)] bg-[var(--color-panel)] px-4 py-4 shadow-[0_12px_48px_rgba(13,26,44,0.06)] md:flex-row md:items-center md:justify-between md:px-5">
      <form className="flex items-center gap-3 rounded-full border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-4 py-2 md:min-w-[22rem]">
        <Search className="size-4 text-[var(--color-ink-soft)]" />
        <input
          type="search"
          placeholder="Search reports, cost center, voucher no..."
          className="h-10 w-full bg-transparent text-sm text-[var(--color-ink)] outline-none placeholder:text-[var(--color-ink-soft)]"
        />
      </form>

      <div className="flex items-center gap-3">
        <div className="inline-flex items-center gap-2 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-4 py-2 text-sm text-emerald-700 dark:text-emerald-300">
          <ShieldCheck className="size-4" />
          <span>Board access</span>
        </div>
        <button
          type="button"
          className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--color-border)] bg-[var(--color-panel)] text-[var(--color-ink)]"
          aria-label="Notifications"
        >
          <BellDot className="size-4" />
        </button>
        <ThemeToggle />
        <div className="rounded-full border border-[var(--color-border)] bg-[linear-gradient(135deg,rgba(23,59,112,0.96),rgba(14,165,233,0.72))] px-4 py-2 text-sm text-white">
          <p className="font-medium">VCT Executive</p>
          <p className="text-xs text-white/72">ceo@vct.group</p>
        </div>
      </div>
    </header>
  );
}
