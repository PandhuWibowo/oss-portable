import { ref } from 'vue'

export function useConnections() {
  const connections = ref([])
  const loading     = ref(false)
  const testing     = ref(false)
  const saving      = ref(false)
  const error       = ref('')
  const notice      = ref('')

  function clearMessages() { error.value = ''; notice.value = '' }

  const BASE = {
    gcp:     '/api/gcp',
    aws:     '/api/aws',
    huawei:  '/api/huawei',
    alibaba: '/api/alibaba',
    azure:   '/api/azure',
  }

  // ── connection list ──────────────────────────────────────────

  async function fetchConnections() {
    loading.value = true
    clearMessages()
    try {
      const [gcpRes, awsRes, huaweiRes, alibabaRes, azureRes] = await Promise.all([
        fetch('/api/gcp/connections').then(r => r.ok ? r.json() : []),
        fetch('/api/aws/connections').then(r => r.ok ? r.json() : []),
        fetch('/api/huawei/connections').then(r => r.ok ? r.json() : []),
        fetch('/api/alibaba/connections').then(r => r.ok ? r.json() : []),
        fetch('/api/azure/connections').then(r => r.ok ? r.json() : []),
      ])
      const gcpList     = (gcpRes     || []).map(c => ({ ...c, provider: 'gcp' }))
      const awsList     = (awsRes     || []).map(c => ({ ...c, provider: 'aws' }))
      const huaweiList  = (huaweiRes  || []).map(c => ({ ...c, provider: 'huawei' }))
      const alibabaList = (alibabaRes || []).map(c => ({ ...c, provider: 'alibaba' }))
      const azureList   = (azureRes   || []).map(c => ({ ...c, provider: 'azure' }))
      connections.value = [...gcpList, ...awsList, ...huaweiList, ...alibabaList, ...azureList]
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
      const res = await fetch(BASE[provider] + '/test', {
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
      const res = await fetch(BASE[provider] + '/connection', {
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
      const res = await fetch(`${BASE[provider]}/connection/${id}`, {
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
      const res = await fetch(`${BASE[provider]}/connection/${id}`, { method: 'DELETE' })
      if (res.ok) await fetchConnections()
    } catch (err) {
      error.value = 'Delete failed: ' + err.message
    }
  }

  // ── bucket browsing ──────────────────────────────────────────

  async function browseObjects(provider, bucket, credentials, prefix = '', pageToken = '') {
    const res = await fetch(BASE[provider] + '/bucket/browse', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, prefix, page_token: pageToken }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { prefix, entries, next_page_token }
  }

  async function getDownloadURL(provider, bucket, credentials, object) {
    const res = await fetch(BASE[provider] + '/bucket/download', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
    return (await res.json()).url
  }

  async function deleteObject(provider, bucket, credentials, object) {
    const res = await fetch(BASE[provider] + '/bucket/delete', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
  }

  async function copyObject(provider, bucket, credentials, source, destination, deleteSource = true) {
    const res = await fetch(BASE[provider] + '/bucket/copy', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, source, destination, delete_source: deleteSource }),
    })
    if (!res.ok) throw new Error(await res.text())
  }

  async function uploadObjects(provider, bucket, credentials, prefix, files) {
    await Promise.all(Array.from(files).map(file => {
      const form = new FormData()
      form.append('bucket',      bucket)
      form.append('credentials', credentials)
      form.append('prefix',      prefix)
      form.append('file',        file)
      return fetch(BASE[provider] + '/bucket/upload', { method: 'POST', body: form }).then(r => {
        if (!r.ok) return r.text().then(t => { throw new Error(t) })
      })
    }))
  }

  async function getBucketStats(provider, bucket, credentials) {
    const res = await fetch(BASE[provider] + '/bucket/stats', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { object_count, total_size, truncated }
  }

  // ── metadata ─────────────────────────────────────────────────

  async function getObjectMetadata(provider, bucket, credentials, object) {
    const res = await fetch(BASE[provider] + '/bucket/metadata', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { content_type, cache_control, metadata, size, updated, etag, md5? }
  }

  async function updateObjectMetadata(provider, bucket, credentials, object, patch) {
    const res = await fetch(BASE[provider] + '/bucket/metadata/update', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object, ...patch }),
    })
    if (!res.ok) throw new Error(await res.text())
  }

  // ── compat (flat listing) ────────────────────────────────────

  async function listObjects(provider, bucket, credentials) {
    const res = await fetch(BASE[provider] + '/bucket/objects', {
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
