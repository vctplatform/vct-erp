import { NextResponse } from "next/server";

import {
  isAppLocale,
  localeCookieName,
  normalizeLocale,
} from "@/lib/i18n/shared";

export async function POST(request: Request) {
  const payload = (await request.json().catch(() => null)) as
    | { locale?: string }
    | null;

  if (!payload?.locale || !isAppLocale(payload.locale)) {
    return NextResponse.json(
      { error: "invalid_locale" },
      { status: 400 },
    );
  }

  const locale = normalizeLocale(payload.locale);
  const response = NextResponse.json({ ok: true, locale });
  response.cookies.set({
    name: localeCookieName,
    value: locale,
    path: "/",
    maxAge: 60 * 60 * 24 * 365,
    sameSite: "lax",
  });
  return response;
}
