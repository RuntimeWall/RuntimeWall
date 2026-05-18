"use client";

import { useCallback, useEffect, useState } from "react";
import { cleanupStoppedSandboxes, listSandboxes } from "@/lib/api";
import type { Sandbox } from "@/lib/types";
import { CreateSandboxButton } from "@/components/create-sandbox-button";
import { SandboxTable } from "@/components/sandbox-table";

export default function HomePage() {
  const [sandboxes, setSandboxes] = useState<Sandbox[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    try {
      const data = await listSandboxes();
      setSandboxes(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
    const interval = setInterval(refresh, 5000);
    return () => clearInterval(interval);
  }, [refresh]);

  return (
    <div className="min-h-screen p-6">
      <header className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-accent">RuntimeWall</h1>
          <p className="text-sm text-slate-400">
            Security-first runtime for autonomous AI agents
          </p>
        </div>
        <CreateSandboxButton onCreated={refresh} />
      </header>

      {error && (
        <div className="mb-4 rounded border border-red-700 bg-red-950/40 p-3 text-sm text-red-200">
          API error: {error}
        </div>
      )}

      <div className="overflow-hidden rounded-lg border border-border bg-panel">
        <div className="flex items-center justify-between border-b border-border px-4 py-3">
          <h2 className="text-sm font-semibold text-slate-300">Sandboxes</h2>
          <div className="flex items-center gap-3">
            <button
              onClick={async () => {
                try {
                  await cleanupStoppedSandboxes();
                  refresh();
                } catch (err) {
                  setError(err instanceof Error ? err.message : String(err));
                }
              }}
              className="text-xs text-slate-400 hover:text-red-300"
              title="Remove all stopped/exited sandboxes"
            >
              Cleanup stopped
            </button>
            <button
              onClick={refresh}
              className="text-xs text-slate-400 hover:text-accent"
            >
              Refresh
            </button>
          </div>
        </div>

        {loading ? (
          <div className="p-6 text-sm text-slate-400">Loading…</div>
        ) : (
          <SandboxTable sandboxes={sandboxes} onChange={refresh} />
        )}
      </div>
    </div>
  );
}
