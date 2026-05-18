"use client";

import { useState } from "react";
import { createSandbox } from "@/lib/api";

export function CreateSandboxButton({ onCreated }: { onCreated: () => void }) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  return (
    <div className="flex flex-col items-end gap-1">
      <button
        onClick={async () => {
          setBusy(true);
          setError(null);
          try {
            await createSandbox();
            onCreated();
          } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
          } finally {
            setBusy(false);
          }
        }}
        disabled={busy}
        className="rounded bg-accent px-3 py-2 text-sm font-medium text-slate-900 hover:bg-sky-400 disabled:opacity-50"
      >
        {busy ? "Creating..." : "+ New Sandbox"}
      </button>
      {error && <span className="text-xs text-red-300">{error}</span>}
    </div>
  );
}
