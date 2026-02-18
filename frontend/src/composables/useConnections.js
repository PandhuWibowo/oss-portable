import { ref, computed } from 'vue'

export function useConnections() {
  const connections = ref([])
  const loading = ref(false)
  const testing = ref(false)
  const saving = ref(false)
  const error = ref('')
  const notice = ref('')

  function clearMessages() {
    error.value = ''
    notice.value = ''
  }

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
      console.error(err)
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
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ bucket, credentials }),
      })
      if (!res.ok) {
        const txt = await res.text()
        error.value = 'Test failed: ' + txt
      } else {
        notice.value = 'Connection test succeeded ✓'
      }
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
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      })
      if (!res.ok) {
        const txt = await res.text()
        error.value = 'Save failed: ' + txt
        return false
      }
      notice.value = 'Connection saved successfully ✓'
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
      const endpoint = provider === 'gcp'
        ? `/api/gcp/connection/${id}`
        : `/api/aws/connection/${id}`
      const res = await fetch(endpoint, { method: 'DELETE' })
      if (res.ok) {
        await fetchConnections()
      }
    } catch (err) {
      error.value = 'Delete failed: ' + err.message
      console.error(err)
    }
  }

  return {
    connections,
    loading,
    testing,
    saving,
    error,
    notice,
    fetchConnections,
    testConnection,
    saveConnection,
    removeConnection,
    clearMessages,
  }
}
