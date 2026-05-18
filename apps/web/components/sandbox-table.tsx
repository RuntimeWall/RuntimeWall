"use client";

import Link from "next/link";
import { deleteSandbox, stopSandbox } from "@/lib/api";
import type { Sandbox } from "@/lib/types";
import { StatusBadge } from "@/components/status-badge";

export function SandboxTable({
  sandboxes,
  onChange,
}: {
  sandboxes: Sandbox[];
  onChange: () => void;
}) {
  if (sandboxes.length === 0) {
    return (
      <div className="p-6 text-sm text-slate-400">
        No sandboxes yet. Click <strong>+ New Sandbox</strong> to create one.
      </div>
    );
  }

  return (
    <table className="w-full text-sm">
      <thead className="text-left text-xs uppercase text-slate-500">
        <tr>
          <th className="px-4 py-2">ID</th>
          <th className="px-4 py-2">Image</th>
          <th className="px-4 py-2">Status</th>
          <th className="px-4 py-2">Container</th>
          <th className="px-4 py-2 text-right">Actions</th>
        </tr>
      </thead>
      <tbody>
        {sandboxes.map((sb) => (
          <tr
            key={sb.id}
            className="border-t border-border hover:bg-slate-900/40"
          >
            <td className="px-4 py-3 font-mono text-xs">
              <Link
                href={`/sandbox/${sb.id}`}
                className="text-accent hover:underline"
              >
                {sb.id.slice(0, 8)}…
              </Link>
            </td>
            <td className="px-4 py-3 text-slate-300">{sb.image}</td>
            <td className="px-4 py-3">
              <StatusBadge status={sb.status} />
            </td>
            <td className="px-4 py-3 font-mono text-xs text-slate-500">
              {sb.container_id?.slice(0, 12)}
            </td>
            <td className="px-4 py-3 text-right">
              <Link
                href={`/sandbox/${sb.id}`}
                className="mr-2 rounded border border-border px-2 py-1 text-xs hover:border-accent hover:text-accent"
              >
                Open
              </Link>
              <button
                onClick={() => stopSandbox(sb.id).then(onChange).catch(() => {})}
                className="mr-2 rounded border border-border px-2 py-1 text-xs hover:border-yellow-500 hover:text-yellow-300"
              >
                Stop
              </button>
              <button
                onClick={() =>
                  deleteSandbox(sb.id).then(onChange).catch(() => {})
                }
                className="rounded border border-border px-2 py-1 text-xs hover:border-red-500 hover:text-red-300"
              >
                Delete
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
