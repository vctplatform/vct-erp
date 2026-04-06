import { FederationSidebar } from "@/components/federation/federation-sidebar";
import { FederationTopbar } from "@/components/federation/federation-topbar";

export default function FederationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen">
      <div className="mx-auto flex min-h-screen max-w-[1720px] gap-5 px-4 py-4 md:px-6 md:py-5">
        <FederationSidebar />
        <div className="flex min-w-0 flex-1 flex-col gap-4">
          <FederationTopbar />
          <main className="min-h-[calc(100vh-7rem)] rounded-[2rem] border border-[var(--color-border)] bg-[linear-gradient(180deg,rgba(255,255,255,0.78),rgba(255,255,255,0.58))] p-5 shadow-[0_24px_70px_rgba(11,22,42,0.08)] backdrop-blur md:p-7 dark:border-white/8 dark:bg-[linear-gradient(180deg,rgba(8,16,32,0.92),rgba(5,10,22,0.88))] dark:shadow-[0_30px_90px_rgba(0,0,0,0.3)]">
            {children}
          </main>
        </div>
      </div>
    </div>
  );
}
