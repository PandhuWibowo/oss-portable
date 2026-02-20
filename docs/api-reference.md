# API Reference

The Anvesa Vestra backend exposes a REST API on port **8080**. All endpoints accept and return `application/json` unless noted otherwise. CORS is enabled for all origins in the default configuration.

---

## Base URL

```
http://localhost:8080
```

---

## Authentication

The API has no authentication layer of its own — it is designed to run on localhost or a private network. Protect the server with a firewall, reverse proxy authentication, or VPN if you expose it beyond localhost.

---

## Common Request Fields

Most `POST` endpoints that operate on bucket objects share these fields:

| Field | Type | Description |
|---|---|---|
| `bucket` | string | Bucket name |
| `credentials` | string | JSON-encoded credentials (see [Connections](./connections.md)) |

---

## GCS Endpoints

### Connections

#### List GCS Connections
```
GET /api/gcp/connections
```
Returns all saved GCS connections.

**Response**
```json
[
  {
    "id": 1,
    "name": "my-prod-bucket",
    "bucket": "my-bucket",
    "provider": "gcp",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

---

#### Create GCS Connection
```
POST /api/gcp/connection
```
Tests credentials first. Saves only on success.

**Request Body**
```json
{
  "name": "my-prod-bucket",
  "bucket": "my-bucket",
  "credentials": "{\"type\":\"service_account\",...}"
}
```

**Response** `200 OK`
```json
{ "id": 1 }
```

---

#### Update GCS Connection
```
PUT /api/gcp/connection/{id}
```
Re-tests credentials before saving the update.

**Request Body** — same shape as create.

**Response** `200 OK`
```json
{ "ok": true }
```

---

#### Delete GCS Connection
```
DELETE /api/gcp/connection/{id}
```

**Response** `200 OK`
```json
{ "ok": true }
```

---

#### Test GCS Credentials
```
POST /api/gcp/test
```
Tests credentials without saving.

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "{\"type\":\"service_account\",...}"
}
```

**Response** `200 OK`
```json
{ "ok": true }
```

**Error Response** `400`
```json
{ "error": "storage: bucket doesn't exist" }
```

---

### Bucket Operations

#### Browse Bucket (GCS)
```
POST /api/gcp/bucket/browse
```
Returns a single page of objects and folder prefixes under the given prefix.

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "...",
  "prefix": "images/2024/",
  "page_token": ""
}
```

**Response**
```json
{
  "prefix": "images/2024/",
  "entries": [
    { "name": "images/2024/",     "type": "prefix", "size": 0,      "updated": "" },
    { "name": "images/2024/photo.jpg", "type": "object", "size": 204800, "updated": "2024-01-15T10:00:00Z" }
  ],
  "next_page_token": "CjEKL2ltYWdlcy8yMDI0L..."
}
```

Pass `next_page_token` back in the next request to load the following page. An empty `next_page_token` means the listing is complete.

---

#### Upload File (GCS)
```
POST /api/gcp/bucket/upload
```
Multipart form upload.

**Form Fields**
| Field | Description |
|---|---|
| `bucket` | Bucket name |
| `credentials` | JSON credentials string |
| `prefix` | Destination folder prefix (e.g. `images/2024/`) |
| `file` | File binary (one or more) |

**Response** `200 OK`
```json
{ "ok": true }
```

---

#### Download URL (GCS)
```
POST /api/gcp/bucket/download
```
Returns a signed URL that grants read access for 15 minutes.

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "...",
  "object": "images/2024/photo.jpg"
}
```

**Response**
```json
{
  "url": "https://storage.googleapis.com/my-bucket/images/2024/photo.jpg?X-Goog-Signature=..."
}
```

> For public objects without credentials, returns a plain public URL.

---

#### Delete Object (GCS)
```
POST /api/gcp/bucket/delete
```

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "...",
  "object": "images/2024/photo.jpg"
}
```

**Response** `200 OK`
```json
{ "ok": true }
```

---

#### Copy / Rename Object (GCS)
```
POST /api/gcp/bucket/copy
```

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "...",
  "source": "images/2024/old-name.jpg",
  "destination": "images/2024/new-name.jpg",
  "delete_source": true
}
```

Set `delete_source: true` to rename (copy then delete). Set `false` to copy only.

**Response** `200 OK`
```json
{ "ok": true }
```

---

#### Bucket Statistics (GCS)
```
POST /api/gcp/bucket/stats
```
Samples up to 1,000 objects to estimate counts and size.

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "..."
}
```

**Response**
```json
{
  "object_count": 842,
  "total_size": 10485760,
  "truncated": false
}
```

---

#### Get Object Metadata (GCS)
```
POST /api/gcp/bucket/metadata
```

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "...",
  "object": "images/2024/photo.jpg"
}
```

**Response**
```json
{
  "content_type": "image/jpeg",
  "cache_control": "public, max-age=86400",
  "metadata": { "author": "alice" },
  "size": 204800,
  "updated": "2024-01-15T10:00:00Z",
  "etag": "\"abc123\"",
  "md5": "rL0Y20zC+Fzt72VPzMSk2A=="
}
```

---

#### Update Object Metadata (GCS)
```
POST /api/gcp/bucket/metadata/update
```

**Request Body**
```json
{
  "bucket": "my-bucket",
  "credentials": "...",
  "object": "images/2024/photo.jpg",
  "content_type": "image/jpeg",
  "cache_control": "public, max-age=3600",
  "metadata": { "author": "bob" }
}
```

**Response** `200 OK`
```json
{ "ok": true }
```

---

## AWS / S3-Compatible Endpoints

All AWS endpoints mirror the GCS endpoints under the `/api/aws/` prefix. The credential format differs — see [Managing Connections](./connections.md#amazon-s3).

### Connections

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/aws/connections` | List saved AWS connections |
| `POST` | `/api/aws/connection` | Create AWS connection |
| `PUT` | `/api/aws/connection/{id}` | Update AWS connection |
| `DELETE` | `/api/aws/connection/{id}` | Delete AWS connection |
| `POST` | `/api/aws/test` | Test credentials |

### Bucket Operations

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/aws/bucket/browse` | Browse objects (paginated) |
| `POST` | `/api/aws/bucket/upload` | Upload file (multipart form) |
| `POST` | `/api/aws/bucket/download` | Get presigned download URL |
| `POST` | `/api/aws/bucket/delete` | Delete object |
| `POST` | `/api/aws/bucket/copy` | Copy or rename object |
| `POST` | `/api/aws/bucket/stats` | Bucket statistics |
| `POST` | `/api/aws/bucket/metadata` | Get object metadata |
| `POST` | `/api/aws/bucket/metadata/update` | Update object metadata |

> AWS metadata updates are implemented as a copy-to-self with `MetadataDirective: REPLACE` because S3 does not allow in-place metadata edits.

---

## Error Responses

All endpoints return errors in this format:

```json
{ "error": "descriptive error message" }
```

| HTTP Status | Meaning |
|---|---|
| `400` | Bad request — missing fields or credential/connection error |
| `404` | Connection ID not found |
| `405` | Method not allowed |
| `500` | Unexpected server error |
