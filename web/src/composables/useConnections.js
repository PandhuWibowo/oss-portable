import { ref } from 'vue'

export function useConnections() {
  const connections = ref([])
  const loading     = ref(false)
  const testing     = ref(false)
  const saving      = ref(false)
  const error       = ref('')
  const notice      = ref('')

  function clearMessages() { error.value = ''; notice.value = '' }

  // ── connection list ──────────────────────────────────────────

  async function fetchConnections() {
    loading.value = true
    clearMessages()
    try {
      const [gcpRes, awsRes, huaweiRes] = await Promise.all([
        fetch('/api/gcp/connections').then(r => r.ok ? r.json() : []),
        fetch('/api/aws/connections').then(r => r.ok ? r.json() : []),
        fetch('/api/huawei/connections').then(r => r.ok ? r.json() : []),
      ])
      const gcpList    = (gcpRes    || []).map(c => ({ ...c, provider: 'gcp' }))
      const awsList    = (awsRes    || []).map(c => ({ ...c, provider: 'aws' }))
      const huaweiList = (huaweiRes || []).map(c => ({ ...c, provider: 'huawei' }))
      connections.value = [...gcpList, ...awsList, ...huaweiList]
    } catch (err) {
      error.value = 'Failed to load connections.'
    } finally {
      loading.value = false
    }
  }

  async function testConnection(provider, bucket, credentials) {
    testing.value = true
    clearMessages()
    try {
      const endpoint = provider === 'gcp' ? '/api/gcp/test' : provider === 'huawei' ? '/api/huawei/test' : '/api/aws/test'
      const res = await fetch(endpoint, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({ bucket, credentials }),
      })
      if (!res.ok) error.value = 'Test failed: ' + await res.text()
      else notice.value = 'Connection test succeeded ✓'
    } catch (err) {
      error.value = 'Error: ' + err.message
    } finally {
      testing.value = false
    }
  }

  async function saveConnection(provider, form) {
    saving.value = true
    clearMessages()
    try {
      const endpoint = provider === 'gcp' ? '/api/gcp/connection' : provider === 'huawei' ? '/api/huawei/connection' : '/api/aws/connection'
      const res = await fetch(endpoint, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify(form),
      })
      if (!res.ok) { error.value = 'Save failed: ' + await res.text(); return false }
      notice.value = 'Connection saved ✓'
      await fetchConnections()
      return true
    } catch (err) {
      error.value = 'Error: ' + err.message
      return false
    } finally {
      saving.value = false
    }
  }

  async function updateConnection(provider, id, form) {
    saving.value = true
    clearMessages()
    try {
      const endpoint = provider === 'gcp'
        ? `/api/gcp/connection/${id}`
        : provider === 'huawei'
        ? `/api/huawei/connection/${id}`
        : `/api/aws/connection/${id}`
      const res = await fetch(endpoint, {
        method:  'PUT',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify(form),
      })
      if (!res.ok) { error.value = 'Update failed: ' + await res.text(); return false }
      notice.value = 'Connection updated ✓'
      await fetchConnections()
      return true
    } catch (err) {
      error.value = 'Error: ' + err.message
      return false
    } finally {
      saving.value = false
    }
  }

  async function removeConnection(provider, id) {
    clearMessages()
    try {
      const endpoint = provider === 'gcp' ? `/api/gcp/connection/${id}` : provider === 'huawei' ? `/api/huawei/connection/${id}` : `/api/aws/connection/${id}`
      const res = await fetch(endpoint, { method: 'DELETE' })
      if (res.ok) await fetchConnections()
    } catch (err) {
      error.value = 'Delete failed: ' + err.message
    }
  }

  // ── bucket browsing ──────────────────────────────────────────

  async function browseObjects(provider, bucket, credentials, prefix = '', pageToken = '') {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/browse' : provider === 'huawei' ? '/api/huawei/bucket/browse' : '/api/aws/bucket/browse'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, prefix, page_token: pageToken }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { prefix, entries, next_page_token }
  }

  async function getDownloadURL(provider, bucket, credentials, object) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/download' : provider === 'huawei' ? '/api/huawei/bucket/download' : '/api/aws/bucket/download'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
    return (await res.json()).url
  }

  async function deleteObject(provider, bucket, credentials, object) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/delete' : provider === 'huawei' ? '/api/huawei/bucket/delete' : '/api/aws/bucket/delete'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
  }

  async function copyObject(provider, bucket, credentials, source, destination, deleteSource = true) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/copy' : provider === 'huawei' ? '/api/huawei/bucket/copy' : '/api/aws/bucket/copy'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, source, destination, delete_source: deleteSource }),
    })
    if (!res.ok) throw new Error(await res.text())
  }

  async function uploadObjects(provider, bucket, credentials, prefix, files) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/upload' : provider === 'huawei' ? '/api/huawei/bucket/upload' : '/api/aws/bucket/upload'
    await Promise.all(Array.from(files).map(file => {
      const form = new FormData()
      form.append('bucket',      bucket)
      form.append('credentials', credentials)
      form.append('prefix',      prefix)
      form.append('file',        file)
      return fetch(endpoint, { method: 'POST', body: form }).then(r => {
        if (!r.ok) return r.text().then(t => { throw new Error(t) })
      })
    }))
  }

  async function getBucketStats(provider, bucket, credentials) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/stats' : provider === 'huawei' ? '/api/huawei/bucket/stats' : '/api/aws/bucket/stats'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { object_count, total_size, truncated }
  }

  // ── metadata ─────────────────────────────────────────────────

  async function getObjectMetadata(provider, bucket, credentials, object) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/metadata' : provider === 'huawei' ? '/api/huawei/bucket/metadata' : '/api/aws/bucket/metadata'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { content_type, cache_control, metadata, size, updated, etag, md5? }
  }

  async function updateObjectMetadata(provider, bucket, credentials, object, patch) {
    const endpoint = provider === 'gcp'
      ? '/api/gcp/bucket/metadata/update'
      : provider === 'huawei'
      ? '/api/huawei/bucket/metadata/update'
      : '/api/aws/bucket/metadata/update'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object, ...patch }),
    })
    if (!res.ok) throw new Error(await res.text())
  }

  // ── compat (flat listing) ────────────────────────────────────

  async function listObjects(provider, bucket, credentials) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/objects' : provider === 'huawei' ? '/api/huawei/bucket/objects' : '/api/aws/bucket/objects'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials }),
    })
    if (!res.ok) throw new Error(await res.text())
    const data = await res.json()
    return { objects: data.objects ?? [], truncated: data.truncated ?? false }
  }

  return {
    connections, loading, testing, saving, error, notice,
    fetchConnections, testConnection, saveConnection, updateConnection,
    removeConnection, clearMessages,
    browseObjects, getDownloadURL, deleteObject, copyObject,
    uploadObjects, getBucketStats, listObjects,
    getObjectMetadata, updateObjectMetadata,
  }
}
