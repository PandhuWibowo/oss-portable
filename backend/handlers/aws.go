package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

// testS3 verifies AWS S3 bucket access using credentials provided as JSON.
func testS3(bucket, credentialsJSON string) error {
	var creds map[string]string
	if err := json.Unmarshal([]byte(credentialsJSON), &creds); err != nil {
		return err
	}
	accessKey := creds["access_key_id"]
	secretKey := creds["secret_access_key"]
	token := creds["session_token"]
	region := creds["region"]
	if region == "" {
		region = "us-east-1"
	}
	if accessKey == "" || secretKey == "" {
		return fmt.Errorf("missing access_key_id or secret_access_key")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provider := awsauth.NewStaticCredentialsProvider(accessKey, secretKey, token)
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(provider),
	)
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	return err
}

// awsObject represents a single S3 object returned to the client.
type awsObject struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Updated time.Time `json:"updated"`
}

// ListAWSObjects lists all objects inside an S3 bucket.
func ListAWSObjects(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var creds map[string]string
	if err := json.Unmarshal([]byte(req.Credentials), &creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accessKey := creds["access_key_id"]
	secretKey := creds["secret_access_key"]
	token := creds["session_token"]
	region := creds["region"]
	if region == "" {
		region = "us-east-1"
	}
	if accessKey == "" || secretKey == "" {
		http.Error(w, "missing access_key_id or secret_access_key", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	provider := awsauth.NewStaticCredentialsProvider(accessKey, secretKey, token)
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(provider),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := s3.NewFromConfig(cfg)

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

// ListAWS returns all saved AWS connections.
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
	_ = json.NewEncoder(w).Encode(conns)
}

// CreateAWS tests and saves a new AWS connection.
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
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

// DeleteAWS removes an AWS connection by ID.
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

// TestAWS tests an AWS connection without saving it.
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
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
