"use client";

import { Moon, SunMedium } from "lucide-react";
import { useTheme } from "next-themes";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  const isDark = resolvedTheme === "dark";

  return (
    <button
      type="button"
      onClick={() => setTheme(isDark ? "light" : "dark")}
      className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-[var(--color-border)] bg-[var(--color-panel)] text-[var(--color-ink)] transition hover:-translate-y-0.5"
      aria-label="Toggle color mode"
    >
      {isDark ? <SunMedium className="size-4" /> : <Moon className="size-4" />}
    </button>
  );
}
