package main

import (
	"log"
	"net/http"

	appdb "github.com/PandhuWibowo/oss-portable/db"
	"github.com/PandhuWibowo/oss-portable/handlers"
	"github.com/PandhuWibowo/oss-portable/middleware"
)

func main() {
	if err := appdb.Init(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	mux := http.NewServeMux()

	// ── GCP connections ───────────────────────────────────────────
	mux.HandleFunc("/api/gcp/connections",        middleware.CORS(handlers.ListGCP))
	mux.HandleFunc("/api/gcp/connection",         middleware.CORS(handlers.CreateGCP))
	mux.HandleFunc("/api/gcp/connection/",        middleware.CORS(handlers.DeleteGCPConn))
	mux.HandleFunc("/api/gcp/test",               middleware.CORS(handlers.TestGCP))

	// ── GCP bucket operations ─────────────────────────────────────
	mux.HandleFunc("/api/gcp/bucket/browse",      middleware.CORS(handlers.BrowseGCPBucket))
	mux.HandleFunc("/api/gcp/bucket/objects",     middleware.CORS(handlers.ListGCPObjects))
	mux.HandleFunc("/api/gcp/bucket/download",    middleware.CORS(handlers.GCPDownloadURL))
	mux.HandleFunc("/api/gcp/bucket/delete",      middleware.CORS(handlers.DeleteGCPObject))
	mux.HandleFunc("/api/gcp/bucket/upload",      middleware.CORS(handlers.UploadGCPObject))
	mux.HandleFunc("/api/gcp/bucket/stats",       middleware.CORS(handlers.GCPBucketStats))

	// ── AWS connections ───────────────────────────────────────────
	mux.HandleFunc("/api/aws/connections",        middleware.CORS(handlers.ListAWS))
	mux.HandleFunc("/api/aws/connection",         middleware.CORS(handlers.CreateAWS))
	mux.HandleFunc("/api/aws/connection/",        middleware.CORS(handlers.DeleteAWS))
	mux.HandleFunc("/api/aws/test",               middleware.CORS(handlers.TestAWS))

	// ── AWS bucket operations ─────────────────────────────────────
	mux.HandleFunc("/api/aws/bucket/browse",      middleware.CORS(handlers.BrowseAWSBucket))
	mux.HandleFunc("/api/aws/bucket/objects",     middleware.CORS(handlers.ListAWSObjects))
	mux.HandleFunc("/api/aws/bucket/download",    middleware.CORS(handlers.AWSDownloadURL))
	mux.HandleFunc("/api/aws/bucket/delete",      middleware.CORS(handlers.DeleteAWSObject))
	mux.HandleFunc("/api/aws/bucket/upload",      middleware.CORS(handlers.UploadAWSObject))
	mux.HandleFunc("/api/aws/bucket/stats",       middleware.CORS(handlers.AWSBucketStats))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	log.Printf("starting backend on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
