import { ActivitySquare } from "lucide-react";

import { DashboardNavLink } from "@/components/dashboard/dashboard-nav-link";
import { getFinanceNavigation } from "@/lib/finance/navigation";
import { getServerLocale } from "@/lib/i18n/server";

export async function MobileNav() {
  const locale = await getServerLocale();
  const { all } = getFinanceNavigation(locale);
  const copy =
    locale === "vi"
      ? {
          label: "Command rail",
          title: "Điều hướng tài chính",
        }
      : {
          label: "Command rail",
          title: "Finance navigation",
        };

  return (
    <nav className="space-y-3 lg:hidden">
      <div className="flex items-center gap-3 rounded-[1.35rem] border border-[var(--color-border)] bg-[var(--color-panel)] px-4 py-3 shadow-[0_12px_36px_rgba(13,26,44,0.06)]">
        <span className="inline-flex h-11 w-11 items-center justify-center rounded-[1rem] border border-[rgba(23,59,112,0.10)] bg-[linear-gradient(135deg,rgba(23,59,112,0.10),rgba(14,165,233,0.08))] text-[var(--color-navy-700)] dark:text-white">
          <ActivitySquare className="size-4" />
        </span>
        <div>
          <p className="text-[0.64rem] uppercase tracking-[0.28em] text-[var(--color-ink-soft)]">
            {copy.label}
          </p>
          <p className="mt-1 text-sm font-semibold text-[var(--color-ink)]">
            {copy.title}
          </p>
        </div>
      </div>

      <div className="overflow-x-auto">
        <div className="flex min-w-max gap-2 rounded-[1.5rem] border border-[var(--color-border)] bg-[var(--color-panel)] p-2.5 shadow-[0_14px_42px_rgba(13,26,44,0.06)]">
          {all.map((item) => (
            <DashboardNavLink
              key={item.href}
              href={item.href}
              label={item.label}
              caption={item.caption}
              icon={item.icon}
              variant="mobile"
            />
          ))}
        </div>
      </div>
    </nav>
  );
}
