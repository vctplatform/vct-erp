import { MobileNav } from "@/components/dashboard/mobile-nav";
import { Sidebar } from "@/components/dashboard/sidebar";
import { Topbar } from "@/components/dashboard/topbar";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen">
      <div className="mx-auto flex min-h-screen max-w-[1720px] gap-5 px-4 py-4 md:px-6 md:py-5">
        <Sidebar />
        <div className="flex min-w-0 flex-1 flex-col gap-5">
          <Topbar />
          <MobileNav />
          <main className="min-h-[calc(100vh-7rem)] rounded-[2rem] border border-[var(--color-border)] bg-[linear-gradient(180deg,rgba(255,255,255,0.78),rgba(255,255,255,0.58))] p-4 shadow-[0_24px_70px_rgba(11,22,42,0.08)] backdrop-blur md:p-6 dark:bg-[linear-gradient(180deg,rgba(13,26,44,0.76),rgba(13,26,44,0.62))]">
            {children}
          </main>
        </div>
      </div>
    </div>
  );
}
