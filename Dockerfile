# ── Stage 1: Build Vue frontend ───────────────────────────────────────────────
FROM oven/bun:1-alpine AS frontend-builder

# Use /app so that the relative import ../../../../docs/ in DocsViewer.vue
# resolves correctly to /app/docs/ at build time.
WORKDIR /app

COPY web/package.json web/bun.lock ./web/
RUN cd web && bun install --frozen-lockfile

# Copy both web source and docs (bundled as static imports at build time)
COPY web/ ./web/
COPY docs/ ./docs/

RUN cd web && bun run build

# ── Stage 2: Build Go server ──────────────────────────────────────────────────
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app/server

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

# modernc.org/sqlite is pure-Go, so CGO is not needed
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/server .

# ── Stage 3: Final image ──────────────────────────────────────────────────────
FROM nginx:1.27-alpine

# Install supervisord to manage both nginx and the Go server
RUN apk add --no-cache supervisor

# nginx: serve frontend + proxy /api to Go backend
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf

# supervisord: manage nginx and the Go server processes
COPY deploy/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Static frontend assets
COPY --from=frontend-builder /app/web/dist /usr/share/nginx/html

# Go server binary
COPY --from=backend-builder /app/bin/server /usr/local/bin/server

# SQLite data directory — mount a volume here for persistence
RUN mkdir -p /data
WORKDIR /data

EXPOSE 80

CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
