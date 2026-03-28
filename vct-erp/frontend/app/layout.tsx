import type { Metadata } from "next";
import { IBM_Plex_Mono, IBM_Plex_Sans } from "next/font/google";

import { Providers } from "@/app/providers";
import "@/app/globals.css";
import { getServerLocale } from "@/lib/i18n/server";
import { cn } from "@/lib/utils";

const sans = IBM_Plex_Sans({
  subsets: ["latin", "vietnamese"],
  variable: "--font-sans",
});

const mono = IBM_Plex_Mono({
  subsets: ["latin", "vietnamese"],
  weight: ["400", "500"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "VCT Command Center",
  description: "Real-time financial dashboard for VCT Group executives.",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const locale = await getServerLocale();

  return (
    <html lang={locale} suppressHydrationWarning>
      <body
        className={cn(
          sans.variable,
          mono.variable,
          "min-h-screen bg-[radial-gradient(circle_at_top,_rgba(24,59,112,0.16),_transparent_48%),linear-gradient(180deg,_var(--color-canvas),_var(--color-canvas-soft))] font-sans text-[var(--color-ink)] antialiased",
        )}
      >
        <Providers locale={locale}>{children}</Providers>
      </body>
    </html>
  );
}
