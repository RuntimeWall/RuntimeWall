import type { ThreatLevel } from "@/lib/types";

const colors: Record<ThreatLevel, string> = {
  none: "bg-slate-900 text-slate-400 border-slate-700",
  suspicious: "bg-yellow-950 text-yellow-300 border-yellow-800",
  destructive: "bg-red-950 text-red-300 border-red-800",
  exfiltration: "bg-purple-950 text-purple-300 border-purple-800",
};

export function ThreatBadge({ threat }: { threat: ThreatLevel }) {
  return (
    <span
      className={`inline-flex rounded border px-1.5 py-0.5 text-[10px] uppercase tracking-wide ${colors[threat]}`}
    >
      {threat}
    </span>
  );
}
