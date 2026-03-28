"use client";

import { LocaleProvider } from "@/components/i18n/locale-provider";
import { ThemeProvider } from "next-themes";
import { SWRConfig } from "swr";
import { Toaster } from "sonner";

import type { AppLocale } from "@/lib/i18n/shared";

export function Providers({
  children,
  locale,
}: {
  children: React.ReactNode;
  locale: AppLocale;
}) {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <LocaleProvider initialLocale={locale}>
        <SWRConfig
          value={{
            revalidateOnFocus: false,
            shouldRetryOnError: false,
          }}
        >
          {children}
          <Toaster
            richColors
            position="top-right"
            toastOptions={{
              className: "font-sans",
            }}
          />
        </SWRConfig>
      </LocaleProvider>
    </ThemeProvider>
  );
}
