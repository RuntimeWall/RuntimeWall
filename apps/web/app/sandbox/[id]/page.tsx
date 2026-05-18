import Link from "next/link";
import { Terminal } from "@/components/terminal";
import { EventsFeed } from "@/components/events-feed";
import { PolicyPanel } from "@/components/policy-panel";

export default function SandboxPage({ params }: { params: { id: string } }) {
  const { id } = params;

  return (
    <div className="flex h-screen flex-col">
      <header className="flex items-center justify-between border-b border-border bg-panel px-6 py-3">
        <div>
          <Link href="/" className="text-xs text-slate-400 hover:text-accent">
            ← Dashboard
          </Link>
          <h1 className="mt-1 text-sm font-semibold text-accent">
            Sandbox <span className="font-mono text-slate-200">{id}</span>
          </h1>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        <main className="flex flex-1 flex-col">
          <div className="border-b border-border bg-panel px-4 py-2 text-[10px] uppercase tracking-wide text-slate-500">
            Terminal
          </div>
          <div className="flex-1 overflow-hidden">
            <Terminal sandboxId={id} />
          </div>
        </main>

        <aside className="flex w-96 flex-col border-l border-border bg-bg">
          <div className="border-b border-border bg-panel px-4 py-2 text-[10px] uppercase tracking-wide text-slate-500">
            Security Events
          </div>
          <EventsFeed sandboxId={id} />
          <PolicyPanel sandboxId={id} />
        </aside>
      </div>
    </div>
  );
}
