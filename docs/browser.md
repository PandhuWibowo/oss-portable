# File Browser

The file browser is the main workspace in Anveesa Vestra. Select any connection from the sidebar to open it.

---

## Navigation

### Folders

Buckets are shown as a hierarchical folder tree. Click any folder (displayed with a `/` suffix) to enter it. The current path is shown in the **breadcrumb bar** at the top of the file list.

Click any segment in the breadcrumb to jump back to that level. Click the bucket name at the far left to return to the root.

### Pagination

Anveesa Vestra loads **200 objects per page**. As you scroll down, the next page is fetched automatically when the sentinel row at the bottom comes into view (infinite scroll). A spinner appears during loading.

---

## Toolbar Actions

| Button | Description |
|---|---|
| Upload | Open the file picker (also accepts drag-and-drop onto the table) |
| New Folder | Create an empty folder (uploads a hidden `.keep` placeholder) |
| Stats | Fetch and display object count and total bucket size |
| Refresh | Reload the current folder listing |

---

## File Operations

### Upload

- Click **Upload** or drag files onto the file table.
- Multiple files can be selected at once.
- Files are uploaded to the **current folder prefix** — navigate into a folder before uploading to place files there.
- A toast notification confirms each successful upload.

### Download

Click the **download icon** in a file's action column. Anveesa Vestra generates a signed URL (GCS), presigned URL (S3/OBS/OSS), or a direct SAS URL (Azure) that expires after 15 minutes and opens the download in a new tab.

> Signed URLs bypass public-access restrictions — the file does not need to be publicly readable.

### Delete

Click the **trash icon** next to a file. A confirmation dialog appears before the delete is executed. Folders cannot be deleted directly — delete all files inside them first.

### Rename / Move

Click the **rename icon** next to a file to open the rename dialog. Enter the new name and click **Move**. The operation is implemented as a **copy + delete**:

1. The object is copied to `current-prefix/new-name`.
2. The original object is deleted.

---

## Bulk Operations

### Selecting Files

Click the **checkbox** on any row to select it. Click the **header checkbox** to select all files currently loaded on the page.

The **selection bar** appears at the top of the file list when one or more files are selected. It shows the count of selected files and offers:

| Action | Description |
|---|---|
| Download all | Download each selected file sequentially |
| Delete all | Delete all selected files with a single confirmation |

---

## Search

Type in the **search box** above the file list to filter the currently loaded objects by name. The filter is applied client-side against the visible page — it does not search the full bucket. Scroll down to load more objects, then search again for broader results.

Press `/` to focus the search box from anywhere in the browser. Press `Escape` to clear it.

---

## Metadata Editor

Click the **info icon** (ℹ) next to any file to open its metadata panel. The panel loads and displays:

| Field | Editable | Description |
|---|---|---|
| Content-Type | Yes | MIME type used when the file is served |
| Cache-Control | Yes | HTTP caching directive |
| Custom metadata | Yes | Arbitrary key-value pairs stored on the object |
| Size | No | File size in bytes |
| Last modified | No | UTC timestamp of the last write |
| ETag | No | Entity tag for cache validation |
| MD5 | No (GCS only) | Base64-encoded MD5 hash |

Click **Save** to write changes back to the bucket. For S3-compatible providers and OBS/OSS, metadata is updated via a copy-to-self operation with `MetadataDirective: REPLACE`.

---

## Bucket Statistics

Click the **Stats** button in the toolbar to fetch bucket statistics. The backend samples up to **1,000 objects** to compute:

- **Object count** — total number of objects in the bucket
- **Total size** — sum of all object sizes (formatted as KB / MB / GB)

If the bucket contains more than 1,000 objects the result is marked as **estimated**.

---

## Copying a File Path

Click the **copy path icon** (clipboard) on any file row to copy its full object path to the clipboard. The format depends on the provider:

| Provider | Path format |
|---|---|
| GCS | `gs://bucket-name/path/to/file.txt` |
| AWS / R2 / MinIO / OBS / OSS | `s3://bucket-name/path/to/file.txt` |
| Azure | `az://container-name/path/to/file.txt` |

---

## Drag and Drop Upload

Drag one or more files from your desktop and drop them anywhere on the file table. The upload begins immediately and places files in the current folder prefix. Each file's progress is confirmed with a toast notification.

---

## Provider Badge & Icon

The browser header shows a colored icon badge indicating which provider the active connection uses:

| Provider | Badge label | Color |
|---|---|---|
| Google Cloud Storage | GCS | Blue |
| Amazon S3 / R2 / MinIO | S3 | Amber |
| Huawei OBS | OBS | Red |
| Alibaba Cloud OSS | OSS | Orange |
| Azure Blob Storage | Azure | Sky blue |

The current bucket name and connection name are also shown in the header.

---

## Sidebar Provider Filter

When you have connections from more than one provider, **filter chips** appear above the connection list in the sidebar. Click a chip to show only connections from that provider. Click it again to deselect. Multiple chips can be active simultaneously.
