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
      const [gcpRes, awsRes] = await Promise.all([
        fetch('/api/gcp/connections').then(r => r.ok ? r.json() : []),
        fetch('/api/aws/connections').then(r => r.ok ? r.json() : []),
      ])
      const gcpList = (gcpRes || []).map(c => ({ ...c, provider: 'gcp' }))
      const awsList = (awsRes || []).map(c => ({ ...c, provider: 'aws' }))
      connections.value = [...gcpList, ...awsList]
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
      const endpoint = provider === 'gcp' ? '/api/gcp/test' : '/api/aws/test'
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
      const endpoint = provider === 'gcp' ? '/api/gcp/connection' : '/api/aws/connection'
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

  async function removeConnection(provider, id) {
    clearMessages()
    try {
      const endpoint = provider === 'gcp' ? `/api/gcp/connection/${id}` : `/api/aws/connection/${id}`
      const res = await fetch(endpoint, { method: 'DELETE' })
      if (res.ok) await fetchConnections()
    } catch (err) {
      error.value = 'Delete failed: ' + err.message
    }
  }

  // ── bucket operations ─────────────────────────────────────────

  async function browseObjects(provider, bucket, credentials, prefix = '') {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/browse' : '/api/aws/bucket/browse'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, prefix }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { prefix, entries: [...] }
  }

  async function getDownloadURL(provider, bucket, credentials, object) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/download' : '/api/aws/bucket/download'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
    return (await res.json()).url
  }

  async function deleteObject(provider, bucket, credentials, object) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/delete' : '/api/aws/bucket/delete'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials, object }),
    })
    if (!res.ok) throw new Error(await res.text())
  }

  async function uploadObjects(provider, bucket, credentials, prefix, files) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/upload' : '/api/aws/bucket/upload'
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
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/stats' : '/api/aws/bucket/stats'
    const res = await fetch(endpoint, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ bucket, credentials }),
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json() // { object_count, total_size, truncated }
  }

  // kept for compat
  async function listObjects(provider, bucket, credentials) {
    const endpoint = provider === 'gcp' ? '/api/gcp/bucket/objects' : '/api/aws/bucket/objects'
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
    fetchConnections, testConnection, saveConnection, removeConnection, clearMessages,
    browseObjects, getDownloadURL, deleteObject, uploadObjects, getBucketStats, listObjects,
  }
}
