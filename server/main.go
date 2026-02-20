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
	mux.HandleFunc("/api/gcp/connections",   middleware.CORS(handlers.ListGCP))
	mux.HandleFunc("/api/gcp/connection",    middleware.CORS(handlers.CreateGCP))
	mux.HandleFunc("/api/gcp/connection/",   middleware.CORS(handlers.GCPConnByID))
	mux.HandleFunc("/api/gcp/test",          middleware.CORS(handlers.TestGCP))

	// ── GCP bucket operations ─────────────────────────────────────
	mux.HandleFunc("/api/gcp/bucket/browse",           middleware.CORS(handlers.BrowseGCPBucket))
	mux.HandleFunc("/api/gcp/bucket/objects",          middleware.CORS(handlers.ListGCPObjects))
	mux.HandleFunc("/api/gcp/bucket/download",         middleware.CORS(handlers.GCPDownloadURL))
	mux.HandleFunc("/api/gcp/bucket/delete",           middleware.CORS(handlers.DeleteGCPObject))
	mux.HandleFunc("/api/gcp/bucket/copy",             middleware.CORS(handlers.CopyGCPObject))
	mux.HandleFunc("/api/gcp/bucket/upload",           middleware.CORS(handlers.UploadGCPObject))
	mux.HandleFunc("/api/gcp/bucket/stats",            middleware.CORS(handlers.GCPBucketStats))
	mux.HandleFunc("/api/gcp/bucket/metadata",         middleware.CORS(handlers.GetGCPMetadata))
	mux.HandleFunc("/api/gcp/bucket/metadata/update",  middleware.CORS(handlers.UpdateGCPMetadata))

	// ── AWS connections ───────────────────────────────────────────
	mux.HandleFunc("/api/aws/connections",   middleware.CORS(handlers.ListAWS))
	mux.HandleFunc("/api/aws/connection",    middleware.CORS(handlers.CreateAWS))
	mux.HandleFunc("/api/aws/connection/",   middleware.CORS(handlers.AWSConnByID))
	mux.HandleFunc("/api/aws/test",          middleware.CORS(handlers.TestAWS))

	// ── AWS bucket operations ─────────────────────────────────────
	mux.HandleFunc("/api/aws/bucket/browse",           middleware.CORS(handlers.BrowseAWSBucket))
	mux.HandleFunc("/api/aws/bucket/objects",          middleware.CORS(handlers.ListAWSObjects))
	mux.HandleFunc("/api/aws/bucket/download",         middleware.CORS(handlers.AWSDownloadURL))
	mux.HandleFunc("/api/aws/bucket/delete",           middleware.CORS(handlers.DeleteAWSObject))
	mux.HandleFunc("/api/aws/bucket/copy",             middleware.CORS(handlers.CopyAWSObject))
	mux.HandleFunc("/api/aws/bucket/upload",           middleware.CORS(handlers.UploadAWSObject))
	mux.HandleFunc("/api/aws/bucket/stats",            middleware.CORS(handlers.AWSBucketStats))
	mux.HandleFunc("/api/aws/bucket/metadata",         middleware.CORS(handlers.GetAWSMetadata))
	mux.HandleFunc("/api/aws/bucket/metadata/update",  middleware.CORS(handlers.UpdateAWSMetadata))

	// ── Huawei OBS connections ────────────────────────────────────
	mux.HandleFunc("/api/huawei/connections",  middleware.CORS(handlers.ListHuawei))
	mux.HandleFunc("/api/huawei/connection",   middleware.CORS(handlers.CreateHuawei))
	mux.HandleFunc("/api/huawei/connection/",  middleware.CORS(handlers.HuaweiConnByID))
	mux.HandleFunc("/api/huawei/test",         middleware.CORS(handlers.TestHuawei))

	// ── Huawei OBS bucket operations ──────────────────────────────
	mux.HandleFunc("/api/huawei/bucket/browse",          middleware.CORS(handlers.BrowseHuaweiBucket))
	mux.HandleFunc("/api/huawei/bucket/objects",         middleware.CORS(handlers.ListHuaweiObjects))
	mux.HandleFunc("/api/huawei/bucket/download",        middleware.CORS(handlers.HuaweiDownloadURL))
	mux.HandleFunc("/api/huawei/bucket/delete",          middleware.CORS(handlers.DeleteHuaweiObject))
	mux.HandleFunc("/api/huawei/bucket/copy",            middleware.CORS(handlers.CopyHuaweiObject))
	mux.HandleFunc("/api/huawei/bucket/upload",          middleware.CORS(handlers.UploadHuaweiObject))
	mux.HandleFunc("/api/huawei/bucket/stats",           middleware.CORS(handlers.HuaweiBucketStats))
	mux.HandleFunc("/api/huawei/bucket/metadata",        middleware.CORS(handlers.GetHuaweiMetadata))
	mux.HandleFunc("/api/huawei/bucket/metadata/update", middleware.CORS(handlers.UpdateHuaweiMetadata))

	// ── Docs ──────────────────────────────────────────────────────
	mux.HandleFunc("/api/docs/", middleware.CORS(handlers.ServeDocs))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	log.Printf("starting backend on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
