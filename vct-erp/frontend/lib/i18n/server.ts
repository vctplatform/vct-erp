import { cookies } from "next/headers";

import { localeCookieName, normalizeLocale, type AppLocale } from "@/lib/i18n/shared";

export async function getServerLocale(): Promise<AppLocale> {
  const cookieStore = await cookies();
  return normalizeLocale(cookieStore.get(localeCookieName)?.value);
}
