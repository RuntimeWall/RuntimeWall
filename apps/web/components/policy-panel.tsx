"use client";

import { useCallback, useEffect, useState } from "react";
import { getPolicy, setPolicy } from "@/lib/api";
import type { SecurityPolicy } from "@/lib/types";

const fields: { key: keyof SecurityPolicy; label: string }[] = [
  { key: "block_destructive_commands", label: "Block destructive commands" },
  { key: "block_exfiltration", label: "Block exfiltration patterns" },
  { key: "block_network_tools", label: "Block network tools" },
  { key: "block_package_installs", label: "Block package installs" },
  { key: "readonly_filesystem", label: "Readonly filesystem" },
];

export function PolicyPanel({ sandboxId }: { sandboxId: string }) {
  const [policy, setLocal] = useState<SecurityPolicy | null>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    getPolicy(sandboxId)
      .then((p) => {
        if (!cancelled) setLocal(p);
      })
      .catch((err) => {
        if (!cancelled)
          setError(err instanceof Error ? err.message : String(err));
      });
    return () => {
      cancelled = true;
    };
  }, [sandboxId]);

  const toggle = useCallback(
    async (key: keyof SecurityPolicy) => {
      if (!policy) return;
      const next: SecurityPolicy = { ...policy, [key]: !policy[key] };
      setLocal(next);
      setSaving(true);
      setError(null);
      try {
        await setPolicy(sandboxId, next);
      } catch (err) {
        setError(err instanceof Error ? err.message : String(err));
      } finally {
        setSaving(false);
      }
    },
    [policy, sandboxId],
  );

  return (
    <section className="border-t border-border bg-panel p-3">
      <div className="mb-2 flex items-center justify-between">
        <h3 className="text-xs uppercase tracking-wide text-slate-500">
          Security Policy
        </h3>
        {saving && (
          <span className="text-[10px] text-slate-500">saving…</span>
        )}
      </div>
      {error && (
        <div className="mb-2 rounded border border-red-800 bg-red-950/40 p-2 text-[11px] text-red-200">
          {error}
        </div>
      )}
      {!policy ? (
        <div className="text-xs text-slate-500">Loading policy…</div>
      ) : (
        <ul className="space-y-2 text-xs">
          {fields.map((f) => (
            <li key={f.key} className="flex items-center justify-between gap-3">
              <span className="text-slate-300">{f.label}</span>
              <button
                onClick={() => toggle(f.key)}
                aria-pressed={policy[f.key]}
                className={`relative h-5 w-9 rounded-full border transition ${
                  policy[f.key]
                    ? "border-emerald-600 bg-emerald-700/50"
                    : "border-slate-700 bg-slate-800"
                }`}
              >
                <span
                  className={`absolute top-0.5 block h-4 w-4 rounded-full bg-white transition ${
                    policy[f.key] ? "left-4" : "left-0.5"
                  }`}
                />
              </button>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
