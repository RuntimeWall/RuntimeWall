import type { SandboxStatus } from "@/lib/types";

const colors: Record<SandboxStatus, string> = {
  running: "bg-emerald-950 text-emerald-300 border-emerald-800",
  creating: "bg-yellow-950 text-yellow-300 border-yellow-800",
  stopped: "bg-slate-900 text-slate-300 border-slate-700",
  exited: "bg-red-950 text-red-300 border-red-800",
  unknown: "bg-slate-900 text-slate-400 border-slate-700",
};

export function StatusBadge({ status }: { status: SandboxStatus }) {
  return (
    <span
      className={`inline-flex items-center rounded border px-2 py-0.5 text-[11px] uppercase tracking-wide ${
        colors[status] ?? colors.unknown
      }`}
    >
      {status}
    </span>
  );
}
