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

New features

- Temporary SQLite DB: the backend uses a local SQLite database (`backend/data.db`) to store saved GCP connections.
- GCP Bucket Dashboard: the frontend includes a simple dashboard to add, test, save, list, and delete GCP bucket connections.

- AWS S3 support: the backend stores AWS S3 connections in `aws_connections` and the dashboard can test/save/list/delete AWS connections. For AWS the `credentials` field expects a JSON object with `access_key_id`, `secret_access_key`, and optional `region` and `session_token`.

API endpoints

- `GET /api/gcp/connections` — list saved connections
- `POST /api/gcp/connection` — save a connection (body: `{name,bucket,credentials}`) — the backend tests credentials before saving
- `DELETE /api/gcp/connection/{id}` — delete connection
- `POST /api/gcp/test` — test credentials without saving (body: `{bucket,credentials}`)

- `GET /api/aws/connections` — list saved AWS connections
- `POST /api/aws/connection` — save an AWS connection (body: `{name,bucket,credentials}`)
- `DELETE /api/aws/connection/{id}` — delete AWS connection
- `POST /api/aws/test` — test AWS credentials without saving (body: `{bucket,credentials}`)

Notes

- Requests expect a service-account JSON for `credentials` (paste the JSON into the form).
- The backend uses the Google Cloud Storage client to validate credentials by fetching bucket attributes.
- To run the frontend without Bun, use `npm install && npm run dev` in the `frontend` folder.
