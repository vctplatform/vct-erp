"use client";

import { create } from "zustand";

import type {
  FinanceDashboardSnapshot,
  FinanceRealtimeEvent,
} from "@/lib/contracts/finance";
import { formatCompactCurrency, formatRunway } from "@/lib/formatters";

type FinanceDashboardState = {
  snapshot: FinanceDashboardSnapshot | null;
  setSnapshot: (snapshot: FinanceDashboardSnapshot) => void;
  applyRealtimeEvent: (event: FinanceRealtimeEvent) => void;
};

export const useFinanceDashboardStore = create<FinanceDashboardState>(
  (set) => ({
    snapshot: null,
    setSnapshot: (snapshot) => set({ snapshot }),
    applyRealtimeEvent: (event) =>
      set((state) => {
        if (!state.snapshot || event.event !== "NEW_TRANSACTION") {
          return state;
        }

        const nextCards = state.snapshot.cards.map((card) => {
          if (card.key === "quarter_net_revenue") {
            const nextValue = card.value + event.amount;
            return {
              ...card,
              value: nextValue,
              formatted_value: formatCompactCurrency(nextValue),
            };
          }

          if (card.key === "cash_assets") {
            const nextValue = card.value + event.amount;
            return {
              ...card,
              value: nextValue,
              formatted_value: formatCompactCurrency(nextValue),
            };
          }

          if (card.key === "runway_index" && card.unit === "months") {
            const nextValue = card.value + 0.1;
            return {
              ...card,
              value: nextValue,
              formatted_value: formatRunway(nextValue),
            };
          }

          return card;
        });

        const nextRevenueMix = state.snapshot.revenue_mix.map((slice) =>
          slice.label.toUpperCase() === event.segment.toUpperCase()
            ? { ...slice, value: slice.value + event.amount }
            : slice,
        );

        return {
          snapshot: {
            ...state.snapshot,
            generated_at: event.timestamp,
            cards: nextCards,
            revenue_mix: nextRevenueMix,
          },
        };
      }),
  }),
);
