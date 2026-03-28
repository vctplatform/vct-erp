export const localeCookieName = "vct-locale";

export const appLocales = ["vi", "en"] as const;

export type AppLocale = (typeof appLocales)[number];

export function isAppLocale(value: string): value is AppLocale {
  return appLocales.includes(value as AppLocale);
}

export function normalizeLocale(value?: string | null): AppLocale {
  if (value && isAppLocale(value)) {
    return value;
  }
  return "vi";
}
