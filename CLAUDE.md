# CLAUDE.md — Anvesa Vestra

## Project Overview

**Anvesa Vestra** is a self-hosted, open-source cloud storage manager. It provides a single web interface to browse, upload, download, rename, and manage files across multiple cloud providers: Google Cloud Storage, Amazon S3 (including Cloudflare R2 and MinIO), Azure Blob Storage, Alibaba Cloud OSS, and Huawei OBS.

Repository: `github.com/PandhuWibowo/oss-portable`

---

## Tech Stack

| Layer           | Technology                                      |
| --------------- | ----------------------------------------------- |
| Backend         | Go 1.24, `net/http` (no framework)              |
| Storage SDKs    | `cloud.google.com/go/storage`, `aws-sdk-go-v2`, Azure SDK, Alibaba OSS, Huawei OBS |
| Database        | SQLite via `modernc.org/sqlite` (pure Go, no CGO) |
| Frontend        | Vue 3 (Composition API), Naive UI component library |
| Build Tool      | Vite 5                                          |
| Package Manager | Bun (preferred) or npm                          |
| Deployment      | Docker (multi-stage), Nginx + supervisord       |

---

## Project Structure

```
anveesa-oss/
├── server/                  # Go backend
│   ├── main.go              # Entry point, route registration
│   ├── handlers/            # HTTP handlers per cloud provider
│   │   ├── gcp.go           # Google Cloud Storage
│   │   ├── aws.go           # Amazon S3 / R2 / MinIO
│   │   ├── azure.go         # Azure Blob Storage
│   │   ├── alibaba.go       # Alibaba Cloud OSS
│   │   ├── huawei.go        # Huawei OBS
│   │   └── docs.go          # Documentation endpoint
│   ├── db/
│   │   └── db.go            # SQLite database setup
│   ├── middleware/
│   │   └── cors.go          # CORS middleware
│   ├── go.mod / go.sum
│   └── data.db              # SQLite database file (gitignored)
├── web/                     # Vue 3 frontend
│   ├── src/
│   │   ├── App.vue          # Root component
│   │   ├── main.js          # Vue app entry
│   │   ├── styles.css       # Global styles
│   │   ├── components/      # UI components
│   │   │   ├── connections/ # Connection management views
│   │   │   ├── docs/        # Documentation viewer
│   │   │   ├── layout/      # App layout components
│   │   │   └── ui/          # Reusable UI components
│   │   └── composables/     # Vue composables
│   │       ├── useConnections.js
│   │       ├── useConfirm.js
│   │       ├── useTheme.js
│   │       └── useToast.js
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
├── docs/                    # Product documentation (markdown)
├── deploy/                  # Deployment configs
│   ├── nginx.conf
│   └── supervisord.conf
├── Dockerfile               # Multi-stage production build
├── Makefile                 # Dev and build commands
└── README.md
```

---

## Development Commands

```bash
# Install frontend dependencies
cd web && bun install && cd ..

# Start backend (port 8080) + frontend (port 5173) together
make dev

# Build production binary + frontend dist
make build
```

The Vite dev server proxies `/api` requests to `http://localhost:8080` (the Go backend).

---

## Architecture Notes

- **Backend**: All routes registered in `server/main.go`. Each cloud provider has its own handler file in `server/handlers/`. The backend exposes a REST API under `/api`.
- **Frontend**: Single-page Vue 3 app using Naive UI for components. State management via composables (no Vuex/Pinia). The app fetches from `/api` endpoints.
- **Database**: SQLite stores connection configurations. The `data.db` file lives in the server directory during development and in `/data` in Docker.
- **No CGO**: The SQLite driver (`modernc.org/sqlite`) is pure Go — no C compiler needed.
- **Docker**: Multi-stage build: Bun builds the frontend, Go compiles the backend, final image runs both behind Nginx with supervisord.

---

## Code Conventions

- **Go**: Standard library HTTP patterns, no framework. Handlers return JSON. Use `log.Printf` for logging.
- **Vue**: Composition API with `<script setup>`. Naive UI components for UI. Composables for shared logic.
- **CSS**: Single global `styles.css` file, no CSS framework.
- **Docs**: Markdown files in `docs/` directory, bundled into the frontend at build time.

---

## Key Ports

| Service        | Port |
| -------------- | ---- |
| Go backend     | 8080 |
| Vite dev server| 5173 |
| Docker (prod)  | 80   |
