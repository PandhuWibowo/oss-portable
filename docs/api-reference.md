# API Reference

The Anveesa Vestra backend exposes a REST API on port **8080**. All endpoints accept and return `application/json` unless noted otherwise. CORS is enabled for all origins in the default configuration.

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
| `bucket` | string | Bucket (or container) name |
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

---

### Bucket Operations (GCS)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/gcp/bucket/browse` | Browse objects (paginated) |
| `POST` | `/api/gcp/bucket/upload` | Upload file (multipart form) |
| `POST` | `/api/gcp/bucket/download` | Get signed download URL |
| `POST` | `/api/gcp/bucket/delete` | Delete object |
| `POST` | `/api/gcp/bucket/copy` | Copy or rename object |
| `POST` | `/api/gcp/bucket/stats` | Bucket statistics |
| `POST` | `/api/gcp/bucket/metadata` | Get object metadata |
| `POST` | `/api/gcp/bucket/metadata/update` | Update object metadata |

**Browse request body**
```json
{
  "bucket": "my-bucket",
  "credentials": "...",
  "prefix": "images/2024/",
  "page_token": ""
}
```

Pass `next_page_token` back in subsequent requests to page through results. An empty token means the listing is complete.

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

## Huawei OBS Endpoints

All Huawei endpoints follow the same pattern under `/api/huawei/`. See [Managing Connections](./connections.md#huawei-obs) for the credential format.

### Connections

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/huawei/connections` | List saved OBS connections |
| `POST` | `/api/huawei/connection` | Create OBS connection |
| `PUT` | `/api/huawei/connection/{id}` | Update OBS connection |
| `DELETE` | `/api/huawei/connection/{id}` | Delete OBS connection |
| `POST` | `/api/huawei/test` | Test credentials |

### Bucket Operations

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/huawei/bucket/browse` | Browse objects (paginated) |
| `POST` | `/api/huawei/bucket/upload` | Upload file (multipart form) |
| `POST` | `/api/huawei/bucket/download` | Get presigned download URL |
| `POST` | `/api/huawei/bucket/delete` | Delete object |
| `POST` | `/api/huawei/bucket/copy` | Copy or rename object |
| `POST` | `/api/huawei/bucket/stats` | Bucket statistics |
| `POST` | `/api/huawei/bucket/metadata` | Get object metadata |
| `POST` | `/api/huawei/bucket/metadata/update` | Update object metadata |

---

## Alibaba Cloud OSS Endpoints

All Alibaba endpoints follow the same pattern under `/api/alibaba/`. See [Managing Connections](./connections.md#alibaba-cloud-oss) for the credential format.

### Connections

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/alibaba/connections` | List saved OSS connections |
| `POST` | `/api/alibaba/connection` | Create OSS connection |
| `PUT` | `/api/alibaba/connection/{id}` | Update OSS connection |
| `DELETE` | `/api/alibaba/connection/{id}` | Delete OSS connection |
| `POST` | `/api/alibaba/test` | Test credentials |

### Bucket Operations

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/alibaba/bucket/browse` | Browse objects (paginated) |
| `POST` | `/api/alibaba/bucket/upload` | Upload file (multipart form) |
| `POST` | `/api/alibaba/bucket/download` | Get presigned download URL |
| `POST` | `/api/alibaba/bucket/delete` | Delete object |
| `POST` | `/api/alibaba/bucket/copy` | Copy or rename object |
| `POST` | `/api/alibaba/bucket/stats` | Bucket statistics |
| `POST` | `/api/alibaba/bucket/metadata` | Get object metadata |
| `POST` | `/api/alibaba/bucket/metadata/update` | Update object metadata |

---

## Azure Blob Storage Endpoints

All Azure endpoints follow the same pattern under `/api/azure/`. See [Managing Connections](./connections.md#azure-blob-storage) for the credential format.

### Connections

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/azure/connections` | List saved Azure connections |
| `POST` | `/api/azure/connection` | Create Azure connection |
| `PUT` | `/api/azure/connection/{id}` | Update Azure connection |
| `DELETE` | `/api/azure/connection/{id}` | Delete Azure connection |
| `POST` | `/api/azure/test` | Test credentials |

### Container Operations

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/azure/bucket/browse` | Browse blobs (paginated) |
| `POST` | `/api/azure/bucket/upload` | Upload blob (multipart form) |
| `POST` | `/api/azure/bucket/download` | Get SAS download URL |
| `POST` | `/api/azure/bucket/delete` | Delete blob |
| `POST` | `/api/azure/bucket/copy` | Copy or rename blob |
| `POST` | `/api/azure/bucket/stats` | Container statistics |
| `POST` | `/api/azure/bucket/metadata` | Get blob metadata |
| `POST` | `/api/azure/bucket/metadata/update` | Update blob metadata |

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
