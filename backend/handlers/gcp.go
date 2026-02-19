package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	appdb "github.com/PandhuWibowo/oss-portable/db"
)

// testGCP verifies GCP bucket access.
// If credentialsJSON is empty, it connects anonymously (for public buckets).
func testGCP(bucket, credentialsJSON string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *storage.Client
	var err error
	if strings.TrimSpace(credentialsJSON) == "" {
		// Public bucket — no authentication needed
		client, err = storage.NewClient(ctx, option.WithoutAuthentication())
	} else {
		client, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(credentialsJSON)))
	}
	if err != nil {
		return err
	}
	defer client.Close()

	// First try Attrs() — works for authenticated connections and some public buckets.
	_, attrsErr := client.Bucket(bucket).Attrs(ctx)
	if attrsErr == nil {
		return nil
	}

	// If Attrs() failed (e.g. public bucket without storage.buckets.get permission),
	// try listing one object — allUsers with storage.objects.list can still read.
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{})
	_, listErr := it.Next()
	if listErr == nil || listErr == iterator.Done {
		// Successfully iterated or bucket is empty — access confirmed
		return nil
	}

	// Both checks failed — return the original Attrs error
	return attrsErr
}

// gcpObject represents a single GCS object returned to the client.
type gcpObject struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Updated     time.Time `json:"updated"`
	ContentType string    `json:"content_type"`
}

// ListGCPObjects lists all objects inside a GCS bucket.
func ListGCPObjects(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var client *storage.Client
	var err error
	if strings.TrimSpace(req.Credentials) == "" {
		client, err = storage.NewClient(ctx, option.WithoutAuthentication())
	} else {
		client, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(req.Credentials)))
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	const maxResults = 1000
	it := client.Bucket(req.Bucket).Objects(ctx, nil)
	var objects []gcpObject
	for len(objects) < maxResults {
		attrs, iterErr := it.Next()
		if iterErr == iterator.Done {
			break
		}
		if iterErr != nil {
			// Return whatever we collected before the error
			break
		}
		objects = append(objects, gcpObject{
			Name:        attrs.Name,
			Size:        attrs.Size,
			Updated:     attrs.Updated,
			ContentType: attrs.ContentType,
		})
	}
	if objects == nil {
		objects = []gcpObject{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"objects":   objects,
		"truncated": len(objects) == maxResults,
	})
}

// ListGCP returns all saved GCP connections.
func ListGCP(w http.ResponseWriter, r *http.Request) {
	rows, err := appdb.DB.Query(
		"SELECT id, name, bucket, credentials, created_at FROM gcp_connections ORDER BY created_at DESC",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type GCPConnection struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Bucket      string    `json:"bucket"`
		Credentials string    `json:"credentials"`
		CreatedAt   time.Time `json:"created_at"`
	}

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

// CreateGCP tests and saves a new GCP connection.
func CreateGCP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Bucket      string `json:"bucket"`
		Credentials string `json:"credentials"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := testGCP(req.Bucket, req.Credentials); err != nil {
		http.Error(w, fmt.Sprintf("test failed: %v", err), http.StatusBadRequest)
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := appdb.DB.Exec(
		"INSERT INTO gcp_connections (name, bucket, credentials, created_at) VALUES (?, ?, ?, ?)",
		req.Name, req.Bucket, req.Credentials, now,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
}

// DeleteGCP removes a GCP connection by ID.
func DeleteGCP(w http.ResponseWriter, r *http.Request) {
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
	if _, err = appdb.DB.Exec("DELETE FROM gcp_connections WHERE id = ?", id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// TestGCP tests a GCP connection without saving it.
func TestGCP(w http.ResponseWriter, r *http.Request) {
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
