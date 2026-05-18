# RuntimeWall Dashboard

Next.js dashboard for RuntimeWall — sandbox table, live security events feed,
and embedded xterm.js terminal attached to the RuntimeWall WebSocket.

## Prerequisites

- Node.js 20+
- RuntimeWall API running on `http://localhost:8080` (see `apps/api`)

## Run

```bash
cd apps/web
npm install
cp .env.local.example .env.local
npm run dev
```

Open <http://localhost:3000>.

## Features

- Sandbox table with **Create / Open / Stop / Delete** actions
- Live security events feed (Server-Sent Events)
- Embedded xterm.js terminal connected to the RuntimeWall WebSocket attach
  endpoint (`/api/v1/sandboxes/{id}/attach`)
- Per-sandbox **policy editor** for runtime governance (toggle block rules)

## Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | RuntimeWall API base URL |
