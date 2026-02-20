# Anvesa Vestra

A self-hosted, open-source cloud storage manager. Browse, upload, download, rename, and manage files across **Google Cloud Storage** and **Amazon S3** (including Cloudflare R2 and MinIO) from a single clean web interface.

---

## Quick Start

```bash
# Install frontend dependencies
cd web && bun install && cd ..

# Start backend + frontend together
make dev
```

Open [http://localhost:5173](http://localhost:5173). Press **New Connection** to add your first bucket.

---

## Documentation

Full product documentation lives in the [`docs/`](./docs/) directory:

| Document | Description |
|---|---|
| [Overview](./docs/index.md) | What Anvesa Vestra is and what it does |
| [Getting Started](./docs/getting-started.md) | Installation, first connection, keyboard shortcuts |
| [Managing Connections](./docs/connections.md) | GCS, S3, Cloudflare R2, and MinIO credential setup |
| [File Browser](./docs/browser.md) | Upload, download, rename, bulk ops, metadata editor |
| [API Reference](./docs/api-reference.md) | Complete REST API for the backend |
| [Deployment](./docs/deployment.md) | Build a production binary and run as a service |
| [Contributing](./docs/contributing.md) | Local dev setup, conventions, adding providers |

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.23, `net/http` |
| Storage | `cloud.google.com/go/storage` · `aws-sdk-go-v2` |
| Database | SQLite (`modernc.org/sqlite`) |
| Frontend | Vue 3 · Vite 5 · Composition API |
| Package manager | Bun (or npm) |

---

## Development Commands

```bash
make dev      # Start backend (port 8080) + frontend (port 5173)
make build    # Compile Go binary to bin/server and build web/dist/
```
