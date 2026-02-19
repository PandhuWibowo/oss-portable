package main

import (
	"log"
	"net/http"

	appdb "github.com/PandhuWibowo/oss-portable/db"
	"github.com/PandhuWibowo/oss-portable/handlers"
	"github.com/PandhuWibowo/oss-portable/middleware"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html><head><meta charset="utf-8"><title>OSS Portable API</title></head>
<body style="font-family:system-ui,Arial;margin:2rem;max-width:600px">
<h1>OSS Portable</h1>
<p>Backend API is running. Available endpoints:</p>
<ul>
  <li><a href="/api/gcp/connections">GET /api/gcp/connections</a></li>
  <li>POST /api/gcp/connection</li>
  <li>POST /api/gcp/test</li>
  <li>DELETE /api/gcp/connection/{id}</li>
  <li><a href="/api/aws/connections">GET /api/aws/connections</a></li>
  <li>POST /api/aws/connection</li>
  <li>POST /api/aws/test</li>
  <li>DELETE /api/aws/connection/{id}</li>
</ul>
<p>For the frontend dashboard, run the Vite dev server and open <a href="http://localhost:5173">http://localhost:5173</a>.</p>
</body></html>`))
}

func main() {
	if err := appdb.Init(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	// GCP routes
	mux.HandleFunc("/api/gcp/connections", middleware.CORS(handlers.ListGCP))
	mux.HandleFunc("/api/gcp/connection", middleware.CORS(handlers.CreateGCP))
	mux.HandleFunc("/api/gcp/connection/", middleware.CORS(handlers.DeleteGCP))
	mux.HandleFunc("/api/gcp/test", middleware.CORS(handlers.TestGCP))
	mux.HandleFunc("/api/gcp/bucket/objects", middleware.CORS(handlers.ListGCPObjects))

	// AWS routes
	mux.HandleFunc("/api/aws/connections", middleware.CORS(handlers.ListAWS))
	mux.HandleFunc("/api/aws/connection", middleware.CORS(handlers.CreateAWS))
	mux.HandleFunc("/api/aws/connection/", middleware.CORS(handlers.DeleteAWS))
	mux.HandleFunc("/api/aws/test", middleware.CORS(handlers.TestAWS))
	mux.HandleFunc("/api/aws/bucket/objects", middleware.CORS(handlers.ListAWSObjects))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	log.Printf("starting backend on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
