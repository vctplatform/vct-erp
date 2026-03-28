import type { AppLocale } from "@/lib/i18n/shared";

function localeCode(locale: AppLocale) {
  return locale === "vi" ? "vi-VN" : "en-US";
}

export function formatCurrency(value: number, locale: AppLocale = "vi") {
  return new Intl.NumberFormat(localeCode(locale), {
    style: "currency",
    currency: "VND",
    maximumFractionDigits: 0,
  }).format(value);
}

export function formatCompactCurrency(
  value: number,
  locale: AppLocale = "vi",
) {
  const absolute = Math.abs(value);

  if (absolute >= 1_000_000_000) {
    return locale === "vi"
      ? `${(value / 1_000_000_000).toFixed(2)} tỷ VND`
      : `${(value / 1_000_000_000).toFixed(2)}B VND`;
  }

  if (absolute >= 1_000_000) {
    return locale === "vi"
      ? `${(value / 1_000_000).toFixed(2)} tr VND`
      : `${(value / 1_000_000).toFixed(2)}M VND`;
  }

  return formatCurrency(value, locale);
}

export function formatPercent(value: number) {
  return `${value >= 0 ? "+" : ""}${value.toFixed(1)}%`;
}

export function formatRunway(value: number, locale: AppLocale = "vi") {
  return `${value.toFixed(1)} ${locale === "vi" ? "tháng" : "months"}`;
}
