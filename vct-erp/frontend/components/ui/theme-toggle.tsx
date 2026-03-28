"use client";

import { Moon, SunMedium } from "lucide-react";
import { useTheme } from "next-themes";

import { useLocale } from "@/components/i18n/locale-provider";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  const { locale } = useLocale();
  const isDark = resolvedTheme === "dark";
  const ariaLabel =
    locale === "vi" ? "Đổi chế độ sáng tối" : "Toggle color mode";

  return (
    <button
      type="button"
      onClick={() => setTheme(isDark ? "light" : "dark")}
      className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--color-border)] bg-[var(--color-panel)] text-[var(--color-ink)] transition hover:-translate-y-0.5"
      aria-label={ariaLabel}
    >
      {isDark ? <SunMedium className="size-4" /> : <Moon className="size-4" />}
    </button>
  );
}
