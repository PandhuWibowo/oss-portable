# oss-portable

This repository contains a minimal scaffold for a Go backend and a Vue 3 frontend. The frontend is intended to run with Bun as the runtime/package manager (Bun is optional but recommended for fast frontend installs and runs).

Quick start

Backend (Go):

1. Change into the backend directory:

```bash
cd backend
```

2. Run the server:

```bash
go run .
```

The backend listens on port `8080` and provides a sample endpoint at `/api/hello`.

Frontend (Vue + Vite, run with Bun):

1. Change into the frontend directory:

```bash
cd frontend
```

2. Install dependencies with Bun (or use `npm install` / `pnpm install`):

```bash
bun install
```

3. Start the dev server with Bun:

```bash
bun run dev
```

Vite dev server proxies `/api` to the Go backend (http://localhost:8080) for local development.

Run both backend and frontend in separate terminals.
