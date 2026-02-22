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

	// ── Alibaba Cloud OSS connections ─────────────────────────────
	mux.HandleFunc("/api/alibaba/connections",  middleware.CORS(handlers.ListAlibaba))
	mux.HandleFunc("/api/alibaba/connection",   middleware.CORS(handlers.CreateAlibaba))
	mux.HandleFunc("/api/alibaba/connection/",  middleware.CORS(handlers.AlibabaConnByID))
	mux.HandleFunc("/api/alibaba/test",         middleware.CORS(handlers.TestAlibaba))

	// ── Alibaba OSS bucket operations ─────────────────────────────
	mux.HandleFunc("/api/alibaba/bucket/browse",          middleware.CORS(handlers.BrowseAlibabaBucket))
	mux.HandleFunc("/api/alibaba/bucket/objects",         middleware.CORS(handlers.ListAlibabaObjects))
	mux.HandleFunc("/api/alibaba/bucket/download",        middleware.CORS(handlers.AlibabaDownloadURL))
	mux.HandleFunc("/api/alibaba/bucket/delete",          middleware.CORS(handlers.DeleteAlibabaObject))
	mux.HandleFunc("/api/alibaba/bucket/copy",            middleware.CORS(handlers.CopyAlibabaObject))
	mux.HandleFunc("/api/alibaba/bucket/upload",          middleware.CORS(handlers.UploadAlibabaObject))
	mux.HandleFunc("/api/alibaba/bucket/stats",           middleware.CORS(handlers.AlibabaBucketStats))
	mux.HandleFunc("/api/alibaba/bucket/metadata",        middleware.CORS(handlers.GetAlibabaMetadata))
	mux.HandleFunc("/api/alibaba/bucket/metadata/update", middleware.CORS(handlers.UpdateAlibabaMetadata))

	// ── Azure Blob Storage connections ────────────────────────────
	mux.HandleFunc("/api/azure/connections",  middleware.CORS(handlers.ListAzure))
	mux.HandleFunc("/api/azure/connection",   middleware.CORS(handlers.CreateAzure))
	mux.HandleFunc("/api/azure/connection/",  middleware.CORS(handlers.AzureConnByID))
	mux.HandleFunc("/api/azure/test",         middleware.CORS(handlers.TestAzure))

	// ── Azure Blob Storage bucket operations ──────────────────────
	mux.HandleFunc("/api/azure/bucket/browse",          middleware.CORS(handlers.BrowseAzureBucket))
	mux.HandleFunc("/api/azure/bucket/objects",         middleware.CORS(handlers.ListAzureObjects))
	mux.HandleFunc("/api/azure/bucket/download",        middleware.CORS(handlers.AzureDownloadURL))
	mux.HandleFunc("/api/azure/bucket/delete",          middleware.CORS(handlers.DeleteAzureObject))
	mux.HandleFunc("/api/azure/bucket/copy",            middleware.CORS(handlers.CopyAzureObject))
	mux.HandleFunc("/api/azure/bucket/upload",          middleware.CORS(handlers.UploadAzureObject))
	mux.HandleFunc("/api/azure/bucket/stats",           middleware.CORS(handlers.AzureBucketStats))
	mux.HandleFunc("/api/azure/bucket/metadata",        middleware.CORS(handlers.GetAzureMetadata))
	mux.HandleFunc("/api/azure/bucket/metadata/update", middleware.CORS(handlers.UpdateAzureMetadata))

	// ── Google Drive connections ──────────────────────────────────
	mux.HandleFunc("/api/gdrive/connections",  middleware.CORS(handlers.ListGDrive))
	mux.HandleFunc("/api/gdrive/connection",   middleware.CORS(handlers.CreateGDrive))
	mux.HandleFunc("/api/gdrive/connection/",  middleware.CORS(handlers.GDriveConnByID))
	mux.HandleFunc("/api/gdrive/test",         middleware.CORS(handlers.TestGDrive))

	// ── Google Drive file operations ──────────────────────────────
	mux.HandleFunc("/api/gdrive/bucket/browse",          middleware.CORS(handlers.BrowseGDriveBucket))
	mux.HandleFunc("/api/gdrive/bucket/objects",         middleware.CORS(handlers.ListGDriveObjects))
	mux.HandleFunc("/api/gdrive/bucket/download",        middleware.CORS(handlers.GDriveDownloadURL))
	mux.HandleFunc("/api/gdrive/bucket/delete",          middleware.CORS(handlers.DeleteGDriveObject))
	mux.HandleFunc("/api/gdrive/bucket/copy",            middleware.CORS(handlers.CopyGDriveObject))
	mux.HandleFunc("/api/gdrive/bucket/upload",          middleware.CORS(handlers.UploadGDriveObject))
	mux.HandleFunc("/api/gdrive/bucket/stats",           middleware.CORS(handlers.GDriveBucketStats))
	mux.HandleFunc("/api/gdrive/bucket/metadata",        middleware.CORS(handlers.GetGDriveMetadata))
	mux.HandleFunc("/api/gdrive/bucket/metadata/update", middleware.CORS(handlers.UpdateGDriveMetadata))

	// ── Docs ──────────────────────────────────────────────────────
	mux.HandleFunc("/api/docs/", middleware.CORS(handlers.ServeDocs))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	log.Printf("starting backend on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
