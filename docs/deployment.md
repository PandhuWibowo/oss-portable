# Deployment

Anvesa Vestra can be built into a single Go binary that serves both the API and the pre-built frontend as static files. No Node.js or Bun runtime is needed in production.

---

## Build

Run the following from the project root:

```bash
make build
```

This command:
1. Compiles the Go backend to `bin/server`
2. Runs `bun run build` in `web/` to produce `web/dist/`

The resulting `bin/server` binary is self-contained for the backend. Copy it alongside the `web/dist/` folder (or embed the dist files in the binary — see below).

---

## Running the Production Binary

```bash
./bin/server
```

The server listens on port **8080** by default. The SQLite database is created as `server/data.db` relative to the working directory on first run.

> Set the `PORT` environment variable to change the listening port if your environment requires it (you may need to add this support to `main.go` — see [Contributing](./contributing.md)).

---

## Serving the Frontend

In development the Vite dev server proxies API requests. In production you have two options:

### Option A — Reverse Proxy (Recommended)

Run Nginx (or Caddy) in front of the Go server:

```nginx
server {
    listen 80;
    server_name storage.example.com;

    # Serve built frontend
    root /var/www/anvesa-vestra/web/dist;
    index index.html;

    # API — proxy to Go backend
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # SPA fallback
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

Copy `web/dist/` to `/var/www/anvesa-vestra/web/dist/` after each build.

### Option B — Serve from Go (Embedded)

Embed the `web/dist/` directory into the binary using Go's `embed` package. Add the following to `server/main.go`:

```go
import "embed"

//go:embed ../web/dist
var staticFiles embed.FS

// Then serve with http.FileServer(http.FS(staticFiles))
```

This packages everything into one binary — no separate static file deployment needed.

---

## Running as a System Service

### systemd (Linux)

Create `/etc/systemd/system/anvesa-vestra.service`:

```ini
[Unit]
Description=Anvesa Vestra Cloud Storage Manager
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/anvesa-vestra
ExecStart=/opt/anvesa-vestra/bin/server
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable anvesa-vestra
sudo systemctl start anvesa-vestra
```

### macOS launchd

Create `~/Library/LaunchAgents/com.anvesa.vestra.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.anvesa.vestra</string>
  <key>ProgramArguments</key>
  <array>
    <string>/opt/anvesa-vestra/bin/server</string>
  </array>
  <key>WorkingDirectory</key>
  <string>/opt/anvesa-vestra</string>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
</dict>
</plist>
```

```bash
launchctl load ~/Library/LaunchAgents/com.anvesa.vestra.plist
```

---

## Data Persistence

The SQLite database file (`server/data.db`) holds all saved connections. In production, ensure this path is:

- **Writable** by the user running the server process
- **Backed up** regularly — losing it means losing all saved connections (not bucket data)
- **Not publicly accessible** — it contains credentials in plaintext

---

## Environment Reference

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listening port (requires code support) |

---

## Cross-Platform Builds

To build for a different OS/architecture from macOS:

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 ./server

# Linux arm64 (Raspberry Pi, AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o bin/server-linux-arm64 ./server

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/server.exe ./server
```
