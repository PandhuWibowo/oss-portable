# Contributing

Thank you for your interest in contributing to Anveesa Vestra. This guide explains how to set up a local development environment, understand the project structure, and follow the conventions used throughout the codebase.

---

## Development Setup

### 1. Fork and Clone

```bash
git clone https://github.com/PandhuWibowo/anveesa-vestra.git
cd anveesa-vestra
```

### 2. Install Dependencies

**Backend** — Go modules are fetched automatically on first build/run.

**Frontend**
```bash
cd web && bun install
```

### 3. Start in Dev Mode

```bash
make dev
```

Both processes run in parallel. `make dev` waits for the backend to be ready on port 8080 before starting the frontend. Press `Ctrl+C` to stop both.

---

## Project Structure

```
anveesa-vestra/
├── server/
│   ├── main.go              Route registration, server startup
│   ├── go.mod               Go module definition
│   ├── db/
│   │   └── db.go            SQLite init and schema migrations
│   ├── handlers/
│   │   ├── gcp.go           All GCS request handlers
│   │   ├── aws.go           All S3/R2/MinIO request handlers
│   │   ├── huawei.go        All Huawei OBS request handlers
│   │   ├── alibaba.go       All Alibaba Cloud OSS request handlers
│   │   └── azure.go         All Azure Blob Storage request handlers
│   └── middleware/
│       └── cors.go          CORS headers middleware
├── web/
│   ├── package.json
│   ├── vite.config.js       Proxy config (dev: /api → :8080)
│   └── src/
│       ├── App.vue           Root component, navigation state machine
│       ├── main.js           App entry point
│       ├── styles.css        Global CSS with custom property tokens
│       ├── components/
│       │   ├── layout/
│       │   │   └── AppHeader.vue        Sidebar with connection list and provider filter chips
│       │   ├── connections/
│       │   │   ├── AddConnectionForm.vue  Create/edit connection form (2-column provider grid)
│       │   │   └── BucketBrowser.vue      Main file browser
│       │   └── ui/
│       │       ├── BaseButton.vue
│       │       ├── BaseInput.vue
│       │       ├── BaseModal.vue
│       │       ├── BaseBadge.vue          Provider label badge (GCS / S3 / OBS / OSS / Azure)
│       │       ├── ProviderIcon.vue       SVG icon for each cloud provider
│       │       ├── SkeletonLoader.vue
│       │       ├── StatusNotice.vue
│       │       ├── ToastContainer.vue
│       │       └── ConfirmModal.vue
│       └── composables/
│           ├── useConnections.js   API calls + shared state
│           ├── useToast.js         Module-level singleton toast queue
│           ├── useConfirm.js       Module-level singleton confirm dialog
│           └── useTheme.js         Theme toggle + localStorage persist
└── docs/                    This documentation
```

---

## Backend Conventions

### Adding a New Endpoint

1. Write the handler function in the relevant file under `handlers/`.
2. Register the route in `server/main.go`.
3. Wrap the handler with `middleware.CORS(...)`.

Handler signature:
```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
    var req struct { ... }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
        return
    }

    // Do work...

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{"ok": true})
}
```

### Database

`db.Init()` returns an `*sql.DB`. Pass it as a package-level variable or inject it through a struct if you add new tables. Do not use an ORM — keep queries as plain SQL strings.

### Error Handling

Return errors as JSON:
```go
http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
```

Never log credentials or user data.

---

## Frontend Conventions

### State Management

There is no Vuex or Pinia. Shared state lives in composables:

- `useConnections.js` — connection list and all bucket API calls
- `useToast.js` — notification queue (module-level singleton)
- `useConfirm.js` — confirmation dialog state (module-level singleton)
- `useTheme.js` — dark/light preference

Module-level `ref()` makes composables behave as singletons — any component that calls `useToast()` gets the same reactive state.

### Component Guidelines

- Use `<script setup>` (Composition API) for all components.
- Props are declared with `defineProps`. Emits with `defineEmits`.
- Avoid inline styles — use CSS custom properties via `var(--token)`.
- UI-only components go in `components/ui/`. Feature components go in `components/connections/`.

### CSS Tokens

All colors, spacing, and shadows are defined as CSS custom properties in `styles.css`. Use them instead of hardcoded values:

```css
/* Good */
color: var(--text);
background: var(--surface);
border: 1px solid var(--border);

/* Avoid */
color: #1c1917;
background: white;
```

Provider brand colors are also available as tokens:

```css
--gcp:         #4285f4;   --gcp-bg:      rgba(66, 133, 244, .1);
--aws:         #d97706;   --aws-bg:      rgba(217, 119, 6, .1);
--huawei:      #cf0a2c;   --huawei-bg:   rgba(207, 10, 44, .1);
--alibaba:     #ff6a00;   --alibaba-bg:  rgba(255, 106, 0, .1);
--azure:       #0078d4;   --azure-bg:    rgba(0, 120, 212, .1);
```

Light and dark mode are switched by toggling `data-theme="light"` on `:root`.

### Adding a New Provider

1. Create `server/handlers/myprovider.go` with the full set of handlers (test, CRUD connections, browse, upload, download, delete, copy, stats, metadata).
2. Add a new table in `server/db/db.go`.
3. Register routes in `server/main.go`.
4. Add the provider card to the `PROVIDERS` array in `AddConnectionForm.vue`.
5. Update `useConnections.js` to call the new endpoints.
6. Add the provider SVG icon to `ProviderIcon.vue`.
7. Add provider color tokens to `styles.css` (`:root` block and dark-mode overrides).

---

## Code Style

**Go**
- `gofmt` before committing.
- Keep handlers focused — one responsibility per function.
- Prefer explicit error returns over panics.

**JavaScript / Vue**
- No TypeScript — plain JS with JSDoc comments where helpful.
- `const` over `let` where possible.
- Avoid deep nesting — extract helper functions.

---

## Pull Request Checklist

- [ ] `make build` succeeds without errors
- [ ] New API endpoints are documented in [api-reference.md](./api-reference.md)
- [ ] New UI features are documented in [browser.md](./browser.md) or [connections.md](./connections.md)
- [ ] No credentials, keys, or `data.db` files are committed
- [ ] CSS uses `var(--token)` — no hardcoded color values

---

## Reporting Issues

Open an issue on GitHub with:
- Steps to reproduce
- Expected behaviour
- Actual behaviour
- Go version (`go version`) and browser + OS

Do not include cloud credentials or bucket names in issue reports.
