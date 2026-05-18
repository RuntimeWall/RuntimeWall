"use client";

import { useEffect, useRef, useState } from "react";
import "@xterm/xterm/css/xterm.css";
import { attachWebSocketURL } from "@/lib/api";

type Status = "connecting" | "connected" | "disconnected";

export function Terminal({ sandboxId }: { sandboxId: string }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [status, setStatus] = useState<Status>("connecting");

  useEffect(() => {
    if (!containerRef.current) return;

    let disposed = false;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    let term: any = null;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    let fit: any = null;
    let ws: WebSocket | null = null;
    let onWindowResize: (() => void) | null = null;

    (async () => {
      const { Terminal: XTerm } = await import("@xterm/xterm");
      const { FitAddon } = await import("@xterm/addon-fit");
      if (disposed || !containerRef.current) return;

      term = new XTerm({
        cursorBlink: true,
        convertEol: true,
        fontSize: 13,
        fontFamily:
          "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
        theme: {
          background: "#0b1020",
          foreground: "#e2e8f0",
          cursor: "#38bdf8",
          black: "#0b1020",
          brightBlack: "#475569",
        },
      });
      fit = new FitAddon();
      term.loadAddon(fit);
      term.open(containerRef.current);
      fit.fit();

      ws = new WebSocket(attachWebSocketURL(sandboxId));
      ws.binaryType = "arraybuffer";

      const sendResize = () => {
        if (!ws || ws.readyState !== WebSocket.OPEN || !term) return;
        ws.send(
          JSON.stringify({ type: "resize", cols: term.cols, rows: term.rows }),
        );
      };

      ws.onopen = () => {
        setStatus("connected");
        sendResize();
      };
      ws.onclose = () => setStatus("disconnected");
      ws.onerror = () => setStatus("disconnected");
      ws.onmessage = (ev) => {
        if (!term) return;
        if (ev.data instanceof ArrayBuffer) {
          term.write(new Uint8Array(ev.data));
        } else if (typeof ev.data === "string") {
          term.write(ev.data);
        }
      };

      term.onData((data: string) => {
        if (ws && ws.readyState === WebSocket.OPEN) {
          ws.send(new TextEncoder().encode(data));
        }
      });

      onWindowResize = () => {
        try {
          fit?.fit();
          sendResize();
        } catch {
          /* ignore */
        }
      };
      window.addEventListener("resize", onWindowResize);
    })();

    return () => {
      disposed = true;
      if (onWindowResize) window.removeEventListener("resize", onWindowResize);
      try {
        ws?.close();
      } catch {
        /* ignore */
      }
      try {
        term?.dispose();
      } catch {
        /* ignore */
      }
    };
  }, [sandboxId]);

  const statusColor =
    status === "connected"
      ? "text-emerald-300"
      : status === "connecting"
        ? "text-yellow-300"
        : "text-red-300";

  return (
    <div className="relative h-full w-full bg-bg">
      <div ref={containerRef} className="absolute inset-0 p-2" />
      <div className="pointer-events-none absolute right-3 top-2 rounded bg-panel/80 px-2 py-1 text-[10px] uppercase">
        <span className={statusColor}>● {status}</span>
      </div>
    </div>
  );
}
