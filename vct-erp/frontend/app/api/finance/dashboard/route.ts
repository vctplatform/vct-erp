import { NextResponse } from "next/server";

import { getFinanceDashboardSnapshot } from "@/lib/api/finance";

export const dynamic = "force-dynamic";

export async function GET() {
  const snapshot = await getFinanceDashboardSnapshot();
  return NextResponse.json(snapshot, {
    headers: {
      "Cache-Control": "no-store",
    },
  });
}
