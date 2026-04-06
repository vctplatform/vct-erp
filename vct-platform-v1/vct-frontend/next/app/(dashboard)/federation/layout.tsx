"use client";

import { FederationSidebar } from "@/components/federation/federation-sidebar";
import { FederationTopbar } from "@/components/federation/federation-topbar";
import "@/app/federation.css";

export default function FederationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="federation-layout-wrapper min-h-screen bg-[#020617] text-slate-50 antialiased">
      <div className="flex p-4 gap-6 h-screen overflow-hidden">
        {/* Sidebar */}
        <FederationSidebar />

        {/* Main Content */}
        <div className="flex-1 flex flex-col min-w-0 h-full overflow-hidden">
          <FederationTopbar />
          
          <main className="flex-1 overflow-y-auto mt-4 pr-2 vct-custom-scrollbar">
            <div className="max-w-[1600px] mx-auto pb-10">
              {children}
            </div>
          </main>
        </div>
      </div>
    </div>
  );
}
