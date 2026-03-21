import { Sidebar } from "@/components/dashboard/sidebar";
import { Topbar } from "@/components/dashboard/topbar";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen">
      <div className="mx-auto flex min-h-screen max-w-[1680px] gap-4 px-4 py-4 md:px-6">
        <Sidebar />
        <div className="flex min-w-0 flex-1 flex-col gap-4">
          <Topbar />
          <main className="min-h-[calc(100vh-7rem)] rounded-[1.75rem] border border-[var(--color-border)] bg-[var(--color-panel-soft)] p-4 shadow-[0_18px_60px_rgba(11,22,42,0.08)] backdrop-blur md:p-6">
            {children}
          </main>
        </div>
      </div>
    </div>
  );
}
