"use client";

import { Wallet, TrendingUp, TrendingDown, ArrowUpRight, ArrowDownRight, Calendar, FileCheck, Banknote, Receipt } from "lucide-react";
import { AreaChart, Area, BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from "recharts";

import { PageHeader, SectionPanel, StatusBadge, FederationDataTable, GlassButton, Tabs } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_FINANCE, MOCK_DASHBOARD } from "@/lib/federation/mock-data";
import { cn } from "@/lib/utils";
import { useState } from "react";

const formatVND = (v: number) => {
  if (v >= 1000000000) return `${(v / 1000000000).toFixed(1)} tỷ`;
  if (v >= 1000000) return `${(v / 1000000).toFixed(0)} tr`;
  return v.toLocaleString("vi-VN");
};

export default function FinancePage() {
  const [activeTab, setActiveTab] = useState("all");
  const d = MOCK_DASHBOARD.financeSummary;

  const filtered = MOCK_FINANCE.filter((f) => {
    if (activeTab === "income" && f.type !== "income") return false;
    if (activeTab === "expense" && f.type !== "expense") return false;
    return true;
  });

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Tài chính liên đoàn"
        title="Tài chính"
        description="Quản lý thu chi, ngân sách, quyết toán của Liên đoàn — Theo dõi dòng tiền và cơ cấu nguồn thu."
      />

      <div className="stagger-fed grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={TrendingUp} label="Tổng thu" value={formatVND(d.totalIncome)} color="emerald" trend={{ direction: "up", value: "+8.2%" }} />
        <StatCard icon={TrendingDown} label="Tổng chi" value={formatVND(d.totalExpense)} color="rose" />
        <StatCard icon={Wallet} label="Số dư" value={formatVND(d.balance)} color="cyan" />
        <StatCard icon={Receipt} label="Giao dịch" value={MOCK_FINANCE.length} color="amber" />
      </div>

      {/* Revenue Trend Chart */}
      <SectionPanel
        title="Xu hướng Thu nhập theo Tháng"
        subtitle="12 tháng gần nhất"
        kicker="Biểu đồ doanh thu"
      >
        <div className="h-[280px]">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={d.monthlyIncome} margin={{ top: 5, right: 5, left: -10, bottom: 0 }}>
              <defs>
                <linearGradient id="incomeGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor="#10B981" stopOpacity={0.25} />
                  <stop offset="100%" stopColor="#10B981" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#ffffff10" vertical={false} />
              <XAxis dataKey="month" tick={{ fontSize: 12, fill: "#94a3b8" }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fontSize: 11, fill: "#94a3b8" }} axisLine={false} tickLine={false} tickFormatter={(v: number) => `${(v / 1000000000).toFixed(1)}`} />
              <Tooltip
                contentStyle={{
                  background: "#081020",
                  border: "1px solid rgba(255,255,255,0.1)",
                  borderRadius: "12px",
                  fontSize: "13px",
                }}
                itemStyle={{ color: "#fff" }}
                formatter={(value) => [formatVND(Number(value)) + " VND", ""]}
              />
              <Area type="monotone" dataKey="value" stroke="#10B981" strokeWidth={2.5} fill="url(#incomeGrad)" dot={false} name="Thu nhập" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </SectionPanel>

      {/* Transaction Table */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          tabs={[
            { id: "all", label: "Tất cả", count: MOCK_FINANCE.length },
            { id: "income", label: "Thu", count: MOCK_FINANCE.filter((f) => f.type === "income").length },
            { id: "expense", label: "Chi", count: MOCK_FINANCE.filter((f) => f.type === "expense").length },
          ]}
          activeTab={activeTab}
          onChange={setActiveTab}
        />
      </div>

      <SectionPanel title="Danh sách Giao dịch" kicker="Sổ cái">
        <FederationDataTable
          columns={[
            { key: "date", label: "Ngày" },
            { key: "type", label: "Loại", width: "80px" },
            { key: "category", label: "Danh mục" },
            { key: "description", label: "Diễn giải" },
            { key: "amount", label: "Số tiền", align: "right" },
            { key: "status", label: "Trạng thái", align: "center" },
          ]}
          rows={filtered.map((f) => ({
            date: (
              <span className="text-xs text-white/50">
                {new Date(f.date).toLocaleDateString("vi-VN")}
              </span>
            ),
            type: (
              <span className={cn(
                "inline-flex items-center gap-1 rounded-lg px-2 py-0.5 text-[0.65rem] font-medium",
                f.type === "income" ? "bg-emerald-500/10 text-emerald-400" : "bg-rose-500/10 text-rose-400"
              )}>
                {f.type === "income" ? <ArrowUpRight className="size-3" /> : <ArrowDownRight className="size-3" />}
                {f.type === "income" ? "Thu" : "Chi"}
              </span>
            ),
            category: <span className="text-sm text-white/80">{f.category}</span>,
            description: <span className="text-sm text-white/60">{f.description}</span>,
            amount: (
              <span className={cn(
                "font-mono text-sm font-semibold",
                f.type === "income" ? "text-emerald-400" : "text-rose-400"
              )}>
                {f.type === "income" ? "+" : "-"}{formatVND(f.amount)}
              </span>
            ),
            status: <StatusBadge status={f.status} />,
          }))}
        />
      </SectionPanel>
    </section>
  );
}
