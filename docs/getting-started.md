# Getting Started

This guide walks you through running Anvesa Vestra locally for the first time.

---

## Option A — Docker (Quickest)

No toolchain needed. Pull the pre-built image and run it:

```bash
docker run -d \
  --name anveesa-vestra \
  -p 80:80 \
  -v anveesa-data:/data \
  pandhuwibowo/anveesa-vestra:latest
```

Open [http://localhost](http://localhost) in your browser, then jump straight to [step 4](#4-add-your-first-connection).

The `-v anveesa-data:/data` flag persists your SQLite database across container restarts. See [Deployment](./deployment.md) for more Docker options.

---

## Option B — From Source

### Prerequisites

| Requirement | Version | Notes |
|---|---|---|
| Go | 1.21+ | [golang.org](https://golang.org/dl/) |
| Bun | 1.0+ | [bun.sh](https://bun.sh) — or use npm/Node 18+ |
| Make | any | Available on macOS/Linux by default |

### 1. Clone the Repository

```bash
git clone https://github.com/PandhuWibowo/anveesa-vestra.git
cd anveesa-vestra
```

### 2. Install Frontend Dependencies

```bash
cd web
bun install
cd ..
```

> Using npm? Run `npm install` instead of `bun install`.

### 3. Start the Development Server

```bash
make dev
```

This command starts both the Go backend (port **8080**) and the Vite dev server (port **5173**) in parallel. It also waits for the backend to be ready before launching the frontend.

Open [http://localhost:5173](http://localhost:5173) in your browser.

---

## Add Your First Connection

When the app loads you will see the welcome screen. Click **New Connection** (or press `N`) to open the connection form.

### Google Cloud Storage

1. Select the **GCS** tab.
2. Enter a **Connection name** (e.g. `my-production-bucket`).
3. Paste your **GCS bucket name** (without `gs://`).
4. Paste your **Service account JSON** key into the Credentials field.
5. Click **Test Connection** — a green notice confirms access.
6. Click **Save**.

> For public buckets you can leave Credentials empty.

### Amazon S3 / S3-Compatible

1. Select the **AWS** tab.
2. Enter a **Connection name**.
3. Enter your **S3 bucket name**.
4. Paste a JSON credentials object:

```json
{
  "access_key_id": "AKIAIOSFODNN7EXAMPLE",
  "secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
  "region": "us-east-1"
}
```

5. Click **Test Connection**, then **Save**.

See [Managing Connections](./connections.md) for Cloudflare R2 and MinIO credential examples.

---

## Browse Your Bucket

Click any saved connection in the sidebar to open the file browser. From there you can:

- Navigate folders
- Upload files (drag-and-drop or click)
- Download, delete, rename, or view file metadata

See [File Browser](./browser.md) for a full feature walkthrough.

---

## Project Layout (Source)

```
anveesa-vestra/
├── server/              Go backend — API server and database
│   ├── main.go          Route definitions and server startup
│   ├── db/              SQLite schema and initialization
│   ├── handlers/        Request handlers (gcp.go, aws.go, huawei.go)
│   └── middleware/      CORS middleware
├── web/                 Vue 3 frontend
│   ├── src/
│   │   ├── App.vue              Root component and navigation logic
│   │   ├── components/          UI and feature components
│   │   └── composables/         Shared state and API logic
│   └── vite.config.js           Dev server and proxy config
├── docs/                This documentation
├── deploy/              Container runtime configuration
│   ├── nginx.conf       nginx: serve frontend + proxy /api to Go server
│   └── supervisord.conf Process manager: runs nginx and Go server together
├── .github/
│   └── workflows/
│       └── docker.yml   CI/CD: build and push Docker image to DockerHub
├── Dockerfile           Multi-stage image build (bun → go → nginx)
└── Makefile             Dev and build commands
```

---

## Keyboard Shortcuts

| Key | Action |
|---|---|
| `N` | Open the new connection form (when not typing in an input) |
| `Escape` | Close modals and dialogs |
