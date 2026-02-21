# Anvesa Vestra

**Anvesa Vestra** is an open-source, self-hosted cloud storage manager. It gives you a clean, unified web interface to browse, upload, download, rename, and manage files across Google Cloud Storage and Amazon S3 — without leaving your browser.

---

## Why Anvesa Vestra?

Most cloud consoles are built for administrators, not daily users. They are slow, bloated, and require full cloud-provider accounts to access. Anvesa Vestra is different:

- **Self-hosted** — runs entirely on your machine or server. No data leaves your infrastructure.
- **Multi-provider** — connect GCS, S3, Cloudflare R2, and MinIO from one interface.
- **Credential-safe** — credentials are stored locally in SQLite and never sent to a third-party service.
- **Lightweight** — a single Go binary + static Vue files. Run natively or via Docker.

---

## Features

| Feature | Description |
|---|---|
| Connection management | Save, test, edit, and delete named bucket connections |
| File browser | Navigate folders, view files, breadcrumb paths |
| Upload | Drag-and-drop or click-to-upload, multi-file support |
| Download | One-click secure signed URLs (15-minute expiry) |
| Delete | Single-file or bulk delete with confirmation |
| Rename / Move | Copy-then-delete within the same bucket |
| Metadata editor | View and edit `content-type`, `cache-control`, and custom headers |
| Bucket stats | Object count and total storage size |
| Search | Filter files by name inline |
| Pagination | Infinite scroll — loads 200 objects at a time |
| Dark mode | System-aware toggle, persisted to local storage |

---

## Supported Providers

| Provider | Type | Notes |
|---|---|---|
| Google Cloud Storage | Native GCS | Service account JSON key |
| Amazon S3 | Native S3 | Access key + secret |
| Cloudflare R2 | S3-compatible | Set a custom endpoint |
| MinIO | S3-compatible | Set a custom endpoint |

---

## Documentation

- [Getting Started](./getting-started.md) — install, run, and add your first connection
- [Managing Connections](./connections.md) — GCS and S3 credential setup
- [File Browser](./browser.md) — browsing, uploading, and managing files
- [API Reference](./api-reference.md) — complete REST API for the backend
- [Deployment](./deployment.md) — build a production binary and serve it
- [Contributing](./contributing.md) — local development and project structure

---

## Quick Start

**Option A — Docker (recommended)**

```bash
docker run -d -p 80:80 -v anveesa-data:/data pandhuwibowo/anveesa-vestra:latest
```

Open [http://localhost](http://localhost) in your browser.

**Option B — From source**

```bash
git clone https://github.com/PandhuWibowo/anveesa-vestra.git
cd anveesa-vestra
make dev
```

Open [http://localhost:5173](http://localhost:5173) in your browser.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.23, `net/http` |
| Storage SDK | `cloud.google.com/go/storage`, `aws-sdk-go-v2` |
| Database | SQLite (`modernc.org/sqlite`) |
| Frontend | Vue 3, Vite 5, Composition API |
| Styling | Plain CSS with custom properties |
| Package manager | Bun (or npm) |
| Container | Docker + nginx + supervisord |
| CI/CD | GitHub Actions → DockerHub |
