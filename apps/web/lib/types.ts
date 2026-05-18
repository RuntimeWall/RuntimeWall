export type SandboxStatus =
  | "creating"
  | "running"
  | "stopped"
  | "exited"
  | "unknown";

export interface Sandbox {
  id: string;
  name: string;
  status: SandboxStatus;
  image: string;
  container_id?: string;
  created_at: string;
  started_at?: string | null;
}

export interface LaunchResult {
  id: string;
  container_id: string;
  image: string;
  status: SandboxStatus;
}

export type ThreatLevel =
  | "none"
  | "suspicious"
  | "destructive"
  | "exfiltration";

export type EventType =
  | "command"
  | "package_install"
  | "file_modify"
  | "process_launch"
  | "policy_violation";

export type CommandSource = "terminal";

export interface RuntimeEvent {
  id: string;
  sandbox_id: string;
  event: EventType;
  command?: string;
  threat: ThreatLevel;
  blocked: boolean;
  reason?: string;
  source: CommandSource;
  timestamp: string;
}

export interface SecurityPolicy {
  block_network_tools: boolean;
  block_package_installs: boolean;
  readonly_filesystem: boolean;
  block_destructive_commands: boolean;
  block_exfiltration: boolean;
}
