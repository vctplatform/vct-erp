"use client";

import { startTransition, useEffect } from "react";
import useSWR from "swr";

import type { FinanceDashboardSnapshot } from "@/lib/contracts/finance";
import { useFinanceDashboardStore } from "@/lib/store/finance-dashboard-store";

async function fetcher(url: string) {
  const response = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`dashboard request failed with ${response.status}`);
  }

  return (await response.json()) as FinanceDashboardSnapshot;
}

export function useFinanceDashboard(initialData: FinanceDashboardSnapshot) {
  const snapshot = useFinanceDashboardStore((state) => state.snapshot);
  const setSnapshot = useFinanceDashboardStore((state) => state.setSnapshot);

  const swr = useSWR("/api/finance/dashboard", fetcher, {
    fallbackData: initialData,
    refreshInterval: 60_000,
    revalidateOnFocus: false,
    keepPreviousData: true,
  });

  useEffect(() => {
    setSnapshot(initialData);
  }, [initialData, setSnapshot]);

  useEffect(() => {
    if (!swr.data) {
      return;
    }

    startTransition(() => {
      setSnapshot(swr.data);
    });
  }, [setSnapshot, swr.data]);

  return {
    snapshot: snapshot ?? swr.data ?? initialData,
    error: swr.error,
    isLoading: !snapshot && !swr.data,
    isValidating: swr.isValidating,
    mutate: swr.mutate,
  };
}
