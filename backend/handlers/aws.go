package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awsauth "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appdb "github.com/PandhuWibowo/oss-portable/db"
)

// ── helpers ──────────────────────────────────────────────────────

func awsCredsFromJSON(raw string) (map[string]string, error) {
	var creds map[string]string
	if err := json.Unmarshal([]byte(raw), &creds); err != nil {
		return nil, err
	}
	return creds, nil
}

// awsS3Client builds an S3 client from parsed credentials.
// The credentials map may include "endpoint" for R2 / MinIO / custom S3-compatible storage.
func awsS3Client(ctx context.Context, creds map[string]string) (*s3.Client, error) {
	accessKey := creds["access_key_id"]
	secretKey := creds["secret_access_key"]
	token := creds["session_token"]
	region := creds["region"]
	if region == "" {
		region = "us-east-1"
	}
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("missing access_key_id or secret_access_key")
	}

	provider := awsauth.NewStaticCredentialsProvider(accessKey, secretKey, token)
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(provider),
	)
	if err != nil {
		return nil, err
	}

	opts := []func(*s3.Options){}
	if ep := creds["endpoint"]; ep != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(ep)
			o.UsePathStyle = true // required for MinIO / R2
		})
	}
	return s3.NewFromConfig(cfg, opts...), nil
}

// testS3 verifies AWS / S3-compatible bucket access.
func testS3(bucket, credentialsJSON string) error {
	creds, err := awsCredsFromJSON(credentialsJSON)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := awsS3Client(ctx, creds)
	if err != nil {
		return err
	}
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	return err
}

// ── connection CRUD ───────────────────────────────────────────────

func ListAWS(w http.ResponseWriter, r *http.Request) {
	rows, err := appdb.DB.Query(
		"SELECT id, name, bucket, credentials, created_at FROM aws_connections ORDER BY created_at DESC",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type AWSConnection struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Bucket      string    `json:"bucket"`
		Credentials string    `json:"credentials"`
		CreatedAt   time.Time `json:"created_at"`
	}

	conns := []AWSConnection{}
	for rows.Next() {
		var c AWSConnection
		var created string
		if err := rows.Scan(&c.ID, &c.Name, &c.Bucket, &c.Credentials, &created); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c.CreatedAt, _ = time.Parse(time.RFC3339, created)
		conns = append(conns, c)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(conns)
}

func CreateAWS(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := testS3(req.Bucket, req.Credentials); err != nil {
		http.Error(w, fmt.Sprintf("test failed: %v", err), http.StatusBadRequest)
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := appdb.DB.Exec(
		"INSERT INTO aws_connections (name, bucket, credentials, created_at) VALUES (?, ?, ?, ?)",
		req.Name, req.Bucket, req.Credentials, now,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

func DeleteAWS(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if _, err = appdb.DB.Exec("DELETE FROM aws_connections WHERE id = ?", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func TestAWS(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := testS3(req.Bucket, req.Credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ── bucket operations ─────────────────────────────────────────────

type awsEntry struct {
	Type    string    `json:"type"` // "dir" | "file"
	Name    string    `json:"name"`
	Display string    `json:"display"`
	Size    int64     `json:"size,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

// BrowseAWSBucket lists entries (files + virtual folders) at a given prefix.
func BrowseAWSBucket(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
		Prefix      string `json:"prefix"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creds, err := awsCredsFromJSON(req.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := awsS3Client(ctx, creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(req.Bucket),
		Prefix:    aws.String(req.Prefix),
		Delimiter: aws.String("/"),
		MaxKeys:   aws.Int32(1000),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var entries []awsEntry

	for _, p := range result.CommonPrefixes {
		if p.Prefix == nil {
			continue
		}
		display := strings.TrimSuffix(strings.TrimPrefix(*p.Prefix, req.Prefix), "/")
		entries = append(entries, awsEntry{Type: "dir", Name: *p.Prefix, Display: display})
	}
	for _, obj := range result.Contents {
		if obj.Key == nil || *obj.Key == req.Prefix {
			continue
		}
		display := strings.TrimPrefix(*obj.Key, req.Prefix)
		var size int64
		if obj.Size != nil {
			size = *obj.Size
		}
		var updated time.Time
		if obj.LastModified != nil {
			updated = *obj.LastModified
		}
		entries = append(entries, awsEntry{Type: "file", Name: *obj.Key, Display: display, Size: size, Updated: updated})
	}
	if entries == nil {
		entries = []awsEntry{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"prefix": req.Prefix, "entries": entries})
}

// ListAWSObjects is kept for backward compat (flat listing).
func ListAWSObjects(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creds, err := awsCredsFromJSON(req.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := awsS3Client(ctx, creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type awsObject struct {
		Name    string    `json:"name"`
		Size    int64     `json:"size"`
		Updated time.Time `json:"updated"`
	}

	const maxResults = 1000
	var objects []awsObject
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket:  aws.String(req.Bucket),
		MaxKeys: aws.Int32(maxResults),
	})
	for paginator.HasMorePages() && len(objects) < maxResults {
		page, pageErr := paginator.NextPage(ctx)
		if pageErr != nil {
			break
		}
		for _, obj := range page.Contents {
			var updated time.Time
			if obj.LastModified != nil {
				updated = *obj.LastModified
			}
			var size int64
			if obj.Size != nil {
				size = *obj.Size
			}
			var name string
			if obj.Key != nil {
				name = *obj.Key
			}
			objects = append(objects, awsObject{Name: name, Size: size, Updated: updated})
		}
	}
	if objects == nil {
		objects = []awsObject{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"objects":   objects,
		"truncated": len(objects) == maxResults,
	})
}

// AWSDownloadURL generates a presigned GET URL (15 min expiry).
func AWSDownloadURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
		Object      string `json:"object"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creds, err := awsCredsFromJSON(req.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := awsS3Client(ctx, creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	psClient := s3.NewPresignClient(client)
	presigned, err := psClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(req.Object),
	}, func(o *s3.PresignOptions) { o.Expires = 15 * time.Minute })
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"url": presigned.URL})
}

// DeleteAWSObject deletes a single S3 object.
func DeleteAWSObject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
		Object      string `json:"object"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creds, err := awsCredsFromJSON(req.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := awsS3Client(ctx, creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(req.Object),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UploadAWSObject uploads a file to S3 via multipart form.
func UploadAWSObject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bucket := r.FormValue("bucket")
	rawCreds := r.FormValue("credentials")
	prefix := r.FormValue("prefix")

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	creds, err := awsCredsFromJSON(rawCreds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := awsS3Client(ctx, creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	objectName := prefix + header.Filename
	if _, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectName),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String(contentType),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"name": objectName})
}

// AWSBucketStats returns sampled object count and total size.
func AWSBucketStats(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creds, err := awsCredsFromJSON(req.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := awsS3Client(ctx, creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const maxSample = 10000
	var count, totalSize int64
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket:  aws.String(req.Bucket),
		MaxKeys: aws.Int32(1000),
	})
	for paginator.HasMorePages() && count < maxSample {
		page, pageErr := paginator.NextPage(ctx)
		if pageErr != nil {
			break
		}
		for _, obj := range page.Contents {
			count++
			if obj.Size != nil {
				totalSize += *obj.Size
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"object_count": count,
		"total_size":   totalSize,
		"truncated":    count == maxSample,
	})
}
