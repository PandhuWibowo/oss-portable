package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awsauth "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var db *sql.DB

type GCPConnection struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Bucket      string    `json:"bucket"`
	Credentials string    `json:"credentials"`
	CreatedAt   time.Time `json:"created_at"`
}

type AWSConnection struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Bucket      string    `json:"bucket"`
	Credentials string    `json:"credentials"`
	CreatedAt   time.Time `json:"created_at"`
}

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "file:data.db?_foreign_keys=1")
	if err != nil {
		return err
	}
	create := `
	CREATE TABLE IF NOT EXISTS gcp_connections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		bucket TEXT NOT NULL,
		credentials TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`
	_, err = db.Exec(create)
	if err != nil {
		return err
	}
	createAws := `
	CREATE TABLE IF NOT EXISTS aws_connections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		bucket TEXT NOT NULL,
		credentials TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`
	_, err = db.Exec(createAws)
	return err
}

func allowCORS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func listConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		allowCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	allowCORS(w)
	rows, err := db.Query("SELECT id, name, bucket, credentials, created_at FROM gcp_connections ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	conns := []GCPConnection{}
	for rows.Next() {
		var c GCPConnection
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

func createConnectionHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	var req struct {
		Name        string `json:"name"`
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Test connection first
	if err := testGCP(req.Bucket, req.Credentials); err != nil {
		http.Error(w, fmt.Sprintf("test failed: %v", err), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	res, err := db.Exec("INSERT INTO gcp_connections (name, bucket, credentials, created_at) VALUES (?, ?, ?, ?)", req.Name, req.Bucket, req.Credentials, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	out := map[string]any{"id": id}
	_ = json.NewEncoder(w).Encode(out)
}

func deleteConnectionHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	// expecting /api/gcp/connection/{id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	idStr := parts[4]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM gcp_connections WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func testGCP(bucket string, credentialsJSON string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(credentialsJSON)))
	if err != nil {
		return err
	}
	defer client.Close()
	// Try to fetch bucket attrs to verify access
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = client.Bucket(bucket).Attrs(ctx)
	return err
}

func testS3(bucket string, credentialsJSON string) error {
	// credentialsJSON expected to be JSON with keys: access_key_id, secret_access_key, region (optional), session_token (optional)
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

	ctx := context.Background()
	provider := awsauth.NewStaticCredentialsProvider(accessKey, secretKey, token)
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region), awsconfig.WithCredentialsProvider(provider))
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	return err
}

func listAWSConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		allowCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	allowCORS(w)
	rows, err := db.Query("SELECT id, name, bucket, credentials, created_at FROM aws_connections ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
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

func createAWSConnectionHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
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
	res, err := db.Exec("INSERT INTO aws_connections (name, bucket, credentials, created_at) VALUES (?, ?, ?, ?)", req.Name, req.Bucket, req.Credentials, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	out := map[string]any{"id": id}
	_ = json.NewEncoder(w).Encode(out)
}

func deleteAWSConnectionHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	idStr := parts[4]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM aws_connections WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	allowCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := testGCP(req.Bucket, req.Credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!doctype html>
<html><head><meta charset="utf-8"><title>OSS Portable</title></head><body style="font-family:system-ui,Arial;margin:2rem;">
<h1>OSS Portable</h1>
<p>Backend API running. Available endpoints:</p>
<ul>
	<li><a href="/api/gcp/connections">/api/gcp/connections</a></li>
	<li><a href="/api/gcp/test">/api/gcp/test</a></li>
	<li><a href="/api/aws/connections">/api/aws/connections</a></li>
	<li><a href="/api/aws/test">/api/aws/test</a></li>
</ul>
<p>If you want the frontend dashboard, run the Vite dev server and open <a href="http://localhost:5173">http://localhost:5173</a>.</p>
</body></html>`
	_, _ = w.Write([]byte(html))
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	// AWS endpoints
	mux.HandleFunc("/api/aws/connections", listAWSConnectionsHandler)
	mux.HandleFunc("/api/aws/connection", createAWSConnectionHandler)
	mux.HandleFunc("/api/aws/connection/", deleteAWSConnectionHandler)
	mux.HandleFunc("/api/aws/test", func(w http.ResponseWriter, r *http.Request) {
		allowCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
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
	})

	mux.HandleFunc("/api/gcp/connections", listConnectionsHandler)
	mux.HandleFunc("/api/gcp/connection", createConnectionHandler)
	mux.HandleFunc("/api/gcp/connection/", deleteConnectionHandler) // trailing slash to catch id
	mux.HandleFunc("/api/gcp/test", testHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Printf("starting backend on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
