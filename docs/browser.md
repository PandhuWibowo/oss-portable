# File Browser

The file browser is the main workspace in Anvesa Vestra. Select any connection from the sidebar to open it.

---

## Navigation

### Folders

Buckets are shown as a hierarchical folder tree. Click any folder (displayed with a `/` suffix) to enter it. The current path is shown in the **breadcrumb bar** at the top of the file list.

Click any segment in the breadcrumb to jump back to that level. Click the bucket name at the far left to return to the root.

### Pagination

Anvesa Vestra loads **200 objects per page**. As you scroll down, the next page is fetched automatically when the sentinel row at the bottom comes into view (infinite scroll). A spinner appears during loading; "All files loaded" appears when the list is complete.

---

## Toolbar Actions

| Button | Description |
|---|---|
| Upload | Open the file picker (also accepts drag-and-drop onto the table) |
| New Folder | Create an empty folder (uploads a hidden `.keep` placeholder) |
| Stats | Fetch and display object count and total bucket size |

---

## File Operations

### Upload

- Click **Upload** or drag files onto the file table.
- Multiple files can be selected at once.
- Files are uploaded to the **current folder prefix** — navigate into a folder before uploading to place files there.
- A toast notification confirms each successful upload.

### Download

Click the **download icon** in a file's action column. Anvesa Vestra generates a **signed URL** (GCS) or **presigned URL** (S3) that expires after 15 minutes and opens the download in a new tab.

> Signed URLs bypass public-access restrictions — the file does not need to be publicly readable.

### Delete

Click the **trash icon** next to a file. A confirmation dialog appears before the delete is executed. Folders cannot be deleted directly — delete all files inside them first.

### Rename / Move

Click the **rename icon** (double-arrow) next to a file to open the rename dialog. Enter the new name and click **Rename**. The operation is implemented as a **copy + delete**:

1. The object is copied to `current-prefix/new-name`.
2. The original object is deleted.

Moving a file to a different folder is not yet supported from this dialog — use copy path + re-upload for cross-folder moves.

---

## Bulk Operations

### Selecting Files

Click the **checkbox** on any row to select it. Click the **header checkbox** to select all files currently loaded on the page.

The **selection bar** appears at the bottom of the screen when one or more files are selected. It shows the count of selected files and offers:

| Action | Description |
|---|---|
| Delete selected | Delete all selected files with a single confirmation |
| Download selected | Download each selected file sequentially (300 ms apart to avoid rate limits) |

---

## Search

Type in the **search box** above the file list to filter the currently loaded objects by name. The filter is applied client-side against the visible page — it does not search the full bucket. Scroll down to load more objects, then search again for broader results.

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

Click **Save** to write changes back to the bucket. For S3-compatible providers, metadata is updated via a copy-to-self operation with `MetadataDirective: REPLACE`.

---

## Bucket Statistics

Click the **Stats** button in the toolbar to fetch bucket statistics. The backend samples up to **1,000 objects** to compute:

- **Object count** — total number of objects in the bucket
- **Total size** — sum of all object sizes (formatted as KB / MB / GB)

If the bucket contains more than 1,000 objects the result is marked as **estimated**.

---

## Copying a File Path

Click the **copy path icon** (clipboard) on any file row to copy its full object key (e.g. `images/2024/photo.jpg`) to the clipboard. A toast confirms the copy.

---

## Drag and Drop Upload

Drag one or more files from your desktop and drop them anywhere on the file table. The upload begins immediately and places files in the current folder prefix. Each file's progress is confirmed with a toast notification.

---

## Provider Badge

A colored badge in the browser header indicates which provider the active connection uses:

- **GCS** — blue badge
- **AWS** — orange badge

The current bucket name and connection name are also shown in the header.
