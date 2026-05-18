import type {
  LaunchResult,
  RuntimeEvent,
  Sandbox,
  SecurityPolicy,
} from "./types";

const API_BASE =
  process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") ?? "http://localhost:8080";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      ...(init?.body ? { "Content-Type": "application/json" } : {}),
      ...(init?.headers ?? {}),
    },
    cache: "no-store",
  });

  if (!res.ok) {
    let detail = res.statusText;
    try {
      const body = (await res.json()) as { error?: string };
      if (body?.error) detail = body.error;
    } catch {
      // ignore
    }
    throw new Error(`${res.status} ${detail}`);
  }

  if (res.status === 204) return undefined as unknown as T;
  return (await res.json()) as T;
}

export async function listSandboxes(): Promise<Sandbox[]> {
  const data = await request<{ sandboxes: Sandbox[] }>("/api/v1/sandboxes");
  return data.sandboxes ?? [];
}

export function createSandbox(): Promise<LaunchResult> {
  return request<LaunchResult>("/sandbox/create", { method: "POST" });
}

export function stopSandbox(id: string): Promise<void> {
  return request<void>(`/api/v1/sandboxes/${id}/stop`, { method: "POST" });
}

export function deleteSandbox(id: string): Promise<void> {
  return request<void>(`/api/v1/sandboxes/${id}`, { method: "DELETE" });
}

export function cleanupStoppedSandboxes(): Promise<{
  removed: string[];
  count: number;
}> {
  return request<{ removed: string[]; count: number }>(
    "/api/v1/sandboxes/cleanup",
    { method: "POST" },
  );
}

export async function getPolicy(id: string): Promise<SecurityPolicy> {
  const data = await request<{ policy: SecurityPolicy }>(
    `/api/v1/sandboxes/${id}/policy`,
  );
  return data.policy;
}

export function setPolicy(id: string, policy: SecurityPolicy): Promise<void> {
  return request<void>(`/api/v1/sandboxes/${id}/policy`, {
    method: "PUT",
    body: JSON.stringify(policy),
  });
}

export async function listEvents(id: string): Promise<RuntimeEvent[]> {
  const data = await request<{ events: RuntimeEvent[] }>(
    `/api/v1/sandboxes/${id}/events`,
  );
  return data.events ?? [];
}

export function eventsStreamURL(id: string): string {
  return `${API_BASE}/api/v1/sandboxes/${id}/events/stream`;
}

export function attachWebSocketURL(id: string): string {
  const base = API_BASE.replace(/^http/, "ws");
  return `${base}/api/v1/sandboxes/${id}/attach`;
}
