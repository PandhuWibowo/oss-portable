# Managing Connections

A **connection** is a named reference to a bucket along with the credentials needed to access it. Connections are stored locally in a SQLite database (`server/data.db`) and never sent to any external service.

---

## Connection Fields

| Field | Required | Description |
|---|---|---|
| Name | Yes | A human-readable label shown in the sidebar |
| Provider | Yes | `GCS` or `AWS` (set on creation, cannot be changed) |
| Bucket | Yes | The bucket name (without protocol prefix) |
| Credentials | Depends | JSON key for GCS; JSON object for AWS |

---

## Google Cloud Storage (GCS)

### Obtaining a Service Account Key

1. Open [Google Cloud Console → IAM & Admin → Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts).
2. Select or create a service account.
3. Under **Keys**, click **Add Key → Create new key → JSON**.
4. Download the `.json` file.

### Required IAM Roles

| Role | When needed |
|---|---|
| `Storage Object Viewer` | Browse and download only |
| `Storage Object Creator` | Upload files |
| `Storage Object Admin` | Upload, delete, rename, metadata edit |
| `Storage Admin` | Full bucket management |

Assign the minimum role needed for your use case.

### Credential Format

Paste the full contents of the downloaded JSON key file into the Credentials field:

```json
{
  "type": "service_account",
  "project_id": "my-project",
  "private_key_id": "key-id",
  "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----\n",
  "client_email": "my-sa@my-project.iam.gserviceaccount.com",
  "client_id": "123456789",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token"
}
```

### Public Buckets

If your bucket is publicly readable, leave the Credentials field empty. Anvesa Vestra will access it without authentication. Note that upload, delete, and metadata-edit operations require credentials regardless of bucket visibility.

---

## Amazon S3

### Credential Format

Provide a JSON object with your AWS credentials:

```json
{
  "access_key_id": "AKIAIOSFODNN7EXAMPLE",
  "secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
  "region": "us-east-1"
}
```

| Field | Required | Default | Description |
|---|---|---|---|
| `access_key_id` | Yes | — | AWS IAM access key |
| `secret_access_key` | Yes | — | AWS IAM secret |
| `region` | No | `us-east-1` | AWS region where the bucket lives |
| `session_token` | No | — | Temporary session token (STS/assumed role) |
| `endpoint` | No | — | Custom endpoint for S3-compatible services |

### Required IAM Permissions

```json
{
  "Effect": "Allow",
  "Action": [
    "s3:ListBucket",
    "s3:GetObject",
    "s3:PutObject",
    "s3:DeleteObject",
    "s3:CopyObject",
    "s3:HeadObject",
    "s3:HeadBucket"
  ],
  "Resource": [
    "arn:aws:s3:::your-bucket-name",
    "arn:aws:s3:::your-bucket-name/*"
  ]
}
```

---

## Cloudflare R2

R2 uses the S3-compatible API. Set a custom endpoint to point to your R2 account.

```json
{
  "access_key_id": "your-r2-access-key-id",
  "secret_access_key": "your-r2-secret-access-key",
  "endpoint": "https://<ACCOUNT_ID>.r2.cloudflarestorage.com",
  "region": "auto"
}
```

Replace `<ACCOUNT_ID>` with your Cloudflare account ID, found in the R2 dashboard.

---

## MinIO

MinIO also uses the S3-compatible API. Use your MinIO server's address as the endpoint.

```json
{
  "access_key_id": "minioadmin",
  "secret_access_key": "minioadmin",
  "endpoint": "http://localhost:9000",
  "region": "us-east-1"
}
```

> For self-signed TLS, ensure the MinIO server certificate is trusted on the machine running Anvesa Vestra's backend.

---

## Testing a Connection

Before saving, click **Test Connection**. The backend will:

- **GCS**: Attempt to list the first object in the bucket using the provided credentials.
- **AWS/S3**: Run `HeadBucket` to verify bucket access.

A green notice confirms success. A red notice shows the error returned by the cloud provider.

---

## Editing a Connection

Click the **pencil icon** next to any connection in the sidebar to open the edit form. The provider tab is locked — you cannot switch from GCS to AWS on an existing connection. All other fields (name, bucket, credentials) can be updated.

After editing, Anvesa Vestra re-tests the connection before saving.

---

## Deleting a Connection

Click the **× icon** next to a connection in the sidebar. A confirmation dialog will appear. Deleting a connection removes it from the local database only — it does not delete any data from your cloud bucket.

---

## Data Storage

Connections are stored in `server/data.db` (SQLite). This file is created automatically on first run and is excluded from version control via `.gitignore`.

**Back up this file** if you want to preserve your connections across reinstalls.
