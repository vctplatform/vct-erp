import { BellDot, CalendarDays, Search, ShieldCheck } from "lucide-react";

import { LocaleToggle } from "@/components/ui/locale-toggle";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { formatFinanceDate } from "@/lib/i18n/finance";
import { getServerLocale } from "@/lib/i18n/server";

export async function Topbar() {
  const locale = await getServerLocale();
  const copy =
    locale === "vi"
      ? {
          search: "Tìm báo cáo, cost center, số chứng từ...",
          boardAccess: "Truy cập ban điều hành",
          notifications: "Thông báo",
          profileName: "Ban điều hành VCT",
        }
      : {
          search: "Search reports, cost center, voucher no...",
          boardAccess: "Board access",
          notifications: "Notifications",
          profileName: "VCT Executive",
        };

  const today = formatFinanceDate(new Date(), locale, {
    weekday: "long",
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });

  return (
    <header className="flex flex-col gap-3 rounded-[1.5rem] border border-[var(--color-border)] bg-[var(--color-panel)] px-4 py-4 shadow-[0_12px_48px_rgba(13,26,44,0.06)] md:flex-row md:items-center md:justify-between md:px-5">
      <form className="flex items-center gap-3 rounded-full border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-4 py-2 md:min-w-[22rem]">
        <Search className="size-4 text-[var(--color-ink-soft)]" />
        <input
          type="search"
          placeholder={copy.search}
          className="h-10 w-full bg-transparent text-sm text-[var(--color-ink)] outline-none placeholder:text-[var(--color-ink-soft)]"
        />
      </form>

      <div className="flex items-center gap-3">
        <div className="hidden items-center gap-2 rounded-full border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-4 py-2 text-sm text-[var(--color-ink-soft)] xl:inline-flex">
          <CalendarDays className="size-4" />
          <span>{today}</span>
        </div>
        <div className="inline-flex items-center gap-2 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-4 py-2 text-sm text-emerald-700 dark:text-emerald-300">
          <ShieldCheck className="size-4" />
          <span>{copy.boardAccess}</span>
        </div>
        <button
          type="button"
          className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--color-border)] bg-[var(--color-panel)] text-[var(--color-ink)]"
          aria-label={copy.notifications}
        >
          <BellDot className="size-4" />
        </button>
        <LocaleToggle />
        <ThemeToggle />
        <div className="rounded-full border border-[var(--color-border)] bg-[linear-gradient(135deg,rgba(23,59,112,0.96),rgba(14,165,233,0.72))] px-4 py-2 text-sm text-white">
          <p className="font-medium">{copy.profileName}</p>
          <p className="text-xs text-white/72">ceo@vct.group</p>
        </div>
      </div>
    </header>
  );
}
