import { NextResponse } from "next/server";

import { getFinanceDashboardSnapshot } from "@/lib/api/finance";
import { getServerLocale } from "@/lib/i18n/server";

export const dynamic = "force-dynamic";

export async function GET() {
  const locale = await getServerLocale();
  const snapshot = await getFinanceDashboardSnapshot(locale);
  return NextResponse.json(snapshot, {
    headers: {
      "Cache-Control": "no-store",
    },
  });
}
