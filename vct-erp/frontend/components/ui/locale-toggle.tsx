"use client";

import { Languages } from "lucide-react";

import { useLocale } from "@/components/i18n/locale-provider";
import { cn } from "@/lib/utils";

const labels = {
  vi: {
    button: "Đổi ngôn ngữ",
    vi: "VI",
    en: "EN",
  },
  en: {
    button: "Switch language",
    vi: "VI",
    en: "EN",
  },
} as const;

export function LocaleToggle() {
  const { locale, setLocale } = useLocale();
  const copy = labels[locale];

  return (
    <div className="inline-flex items-center gap-1 rounded-full border border-[var(--color-border)] bg-[var(--color-panel)] p-1">
      <span
        className="inline-flex h-8 w-8 items-center justify-center rounded-full text-[var(--color-ink-soft)]"
        aria-hidden="true"
      >
        <Languages className="size-4" />
      </span>
      {(["vi", "en"] as const).map((item) => (
        <button
          key={item}
          type="button"
          onClick={() => void setLocale(item)}
          aria-label={copy.button}
          className={cn(
            "rounded-full px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.18em] transition",
            locale === item
              ? "bg-[var(--color-navy-700)] text-white"
              : "text-[var(--color-ink-soft)] hover:bg-[var(--color-canvas-soft)]",
          )}
        >
          {copy[item]}
        </button>
      ))}
    </div>
  );
}
