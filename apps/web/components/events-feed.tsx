"use client";

import { useEffect, useRef, useState } from "react";
import { eventsStreamURL } from "@/lib/api";
import { ThreatBadge } from "@/components/threat-badge";
import type { RuntimeEvent } from "@/lib/types";

const MAX_EVENTS = 200;

export function EventsFeed({ sandboxId }: { sandboxId: string }) {
  const [events, setEvents] = useState<RuntimeEvent[]>([]);
  const [connected, setConnected] = useState(false);
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const es = new EventSource(eventsStreamURL(sandboxId));

    es.onopen = () => setConnected(true);
    es.onerror = () => setConnected(false);

    const handler = (ev: MessageEvent) => {
      try {
        const data = JSON.parse(ev.data) as RuntimeEvent;
        setEvents((prev) => {
          const next = [...prev, data];
          return next.length > MAX_EVENTS
            ? next.slice(next.length - MAX_EVENTS)
            : next;
        });
      } catch {
        /* ignore malformed event */
      }
    };

    es.addEventListener("runtime", handler);
    return () => {
      es.removeEventListener("runtime", handler);
      es.close();
    };
  }, [sandboxId]);

  useEffect(() => {
    const el = listRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [events]);

  return (
    <div className="flex flex-1 flex-col overflow-hidden">
      <div className="flex items-center justify-between border-b border-border px-3 py-1.5 text-[10px] uppercase">
        <span className="text-slate-500">{events.length} events</span>
        <span
          className={
            connected ? "text-emerald-300" : "text-slate-500"
          }
        >
          ● {connected ? "live" : "offline"}
        </span>
      </div>
      <div ref={listRef} className="flex-1 overflow-y-auto p-2 text-xs">
        {events.length === 0 ? (
          <div className="p-2 text-slate-500">
            Waiting for events… run a command in the terminal.
          </div>
        ) : (
          <ul className="space-y-2">
            {events.map((ev) => (
              <li
                key={ev.id}
                className={`rounded border p-2 ${
                  ev.blocked
                    ? "border-red-800 bg-red-950/30"
                    : ev.threat !== "none"
                      ? "border-yellow-800 bg-yellow-950/20"
                      : "border-border bg-panel"
                }`}
              >
                <div className="mb-1 flex items-center justify-between gap-2">
                  <span className="text-[10px] uppercase tracking-wide text-slate-500">
                    {ev.event}
                  </span>
                  <ThreatBadge threat={ev.threat} />
                </div>
                {ev.command && (
                  <div className="break-all font-mono text-slate-200">
                    {ev.command}
                  </div>
                )}
                {ev.reason && (
                  <div className="mt-1 text-[11px] text-slate-400">
                    {ev.reason}
                  </div>
                )}
                <div className="mt-1 flex items-center justify-between text-[10px] text-slate-500">
                  <span>{new Date(ev.timestamp).toLocaleTimeString()}</span>
                  {ev.blocked && (
                    <span className="font-semibold text-red-400">BLOCKED</span>
                  )}
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
