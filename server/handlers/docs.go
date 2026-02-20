package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// docsDir resolves the docs directory regardless of working directory.
// In dev the server runs from server/, so docs are at ../docs.
// In production the binary runs from the project root, so docs are at ./docs.
func docsDir() string {
	for _, candidate := range []string{"../docs", "./docs", "docs"} {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return "../docs"
}

// ServeDocs handles GET /api/docs/{page} and returns the raw markdown file.
func ServeDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract page name — strip /api/docs/ prefix
	page := strings.TrimPrefix(r.URL.Path, "/api/docs/")
	page = strings.TrimSpace(page)
	if page == "" {
		page = "index"
	}

	// filepath.Base prevents path traversal (e.g. ../../etc/passwd)
	page = filepath.Base(page)
	if !strings.HasSuffix(page, ".md") {
		page += ".md"
	}

	fullPath := filepath.Join(docsDir(), page)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, `{"error":"page not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}
