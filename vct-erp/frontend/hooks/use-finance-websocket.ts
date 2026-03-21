"use client";

import {
  startTransition,
  useEffect,
  useEffectEvent,
  useRef,
  useState,
} from "react";

import type { FinanceRealtimeEvent } from "@/lib/contracts/finance";

type UseFinanceWebSocketOptions = {
  enabled?: boolean;
  onTransaction?: (event: FinanceRealtimeEvent) => void;
};

export function useFinanceWebSocket(
  url: string | null,
  options: UseFinanceWebSocketOptions = {},
) {
  const [isConnected, setIsConnected] = useState(false);
  const reconnectRef = useRef<number | null>(null);
  const onTransaction = useEffectEvent((event: FinanceRealtimeEvent) => {
    options.onTransaction?.(event);
  });

  useEffect(() => {
    if (!url || options.enabled === false) {
      return;
    }

    let socket: WebSocket | null = null;
    let cancelled = false;

    function connect() {
      if (cancelled) {
        return;
      }

      socket = new WebSocket(url);

      socket.addEventListener("open", () => {
        startTransition(() => setIsConnected(true));
      });

      socket.addEventListener("close", () => {
        startTransition(() => setIsConnected(false));
        if (!cancelled) {
          reconnectRef.current = window.setTimeout(connect, 2000);
        }
      });

      socket.addEventListener("message", (message) => {
        try {
          const payload = JSON.parse(message.data) as FinanceRealtimeEvent;
          if (payload.event === "NEW_TRANSACTION") {
            startTransition(() => onTransaction(payload));
          }
        } catch {
          // Ignore malformed realtime payloads and wait for the next frame.
        }
      });
    }

    connect();

    return () => {
      cancelled = true;
      if (reconnectRef.current !== null) {
        window.clearTimeout(reconnectRef.current);
      }
      socket?.close();
    };
  }, [options.enabled, url]);

  return { isConnected };
}
