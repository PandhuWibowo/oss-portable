<template>
  <div class="container">
    <Header />
    <div class="layout">
        <section class="card">
          <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:12px">
            <div>
              <h3 style="margin:0">Connections</h3>
              <div class="small">Manage saved cloud connections</div>
            </div>
            <div style="display:flex;gap:8px;align-items:center">
              <input v-model="q" placeholder="Search name or bucket" style="padding:8px;border-radius:8px;min-width:220px" />
            </div>
          </div>

          <div v-if="loading">
            <div class="skeleton" style="height:72px;margin-bottom:8px"></div>
            <div class="skeleton" style="height:72px;opacity:0.9"></div>
          </div>

          <div v-else>
            <div v-if="filtered.length === 0" class="small" style="padding:20px;text-align:center;color:var(--muted)">No connections yet — add one on the right.</div>
            <div style="display:flex;flex-direction:column;gap:10px">
              <ConnectionCard v-for="c in filtered" :key="c.provider + '-' + c.id" :conn="c" @use="useConnection" @delete="(id)=>removeConnection(c.provider,id)" />
            </div>
          </div>
        </section>

      <aside class="card">
        <h3 style="margin:0 0 12px 0">Add / Test Connection</h3>
        <div class="provider-toggle" style="margin-bottom:12px">
          <button class="btn btn-ghost" :class="provider==='gcp' ? 'active' : ''" @click="provider='gcp'">GCP</button>
          <button class="btn btn-ghost" :class="provider==='aws' ? 'active' : ''" @click="provider='aws'">AWS</button>
        </div>

        <form @submit.prevent="onSave" style="display:flex;flex-direction:column;gap:10px">
          <input v-model="form.name" placeholder="Connection name" />
          <input v-model="form.bucket" placeholder="Bucket name" />
          <textarea v-model="form.credentials" rows="6" :placeholder="provider === 'gcp' ? 'Service account JSON' : awsPlaceholder"></textarea>
          <div style="display:flex;gap:10px">
            <button type="button" class="btn btn-primary" @click="onTest" :disabled="testing">{{ testing ? 'Testing…' : 'Test' }}</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</button>
          </div>
          <div v-if="error" class="notice" style="background:rgba(255,0,0,0.04);color:var(--danger)">{{ error }}</div>
          <div v-if="notice" class="notice">{{ notice }}</div>
        </form>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Header from './components/Header.vue'

const connections = ref([])
const loading = ref(false)
const testing = ref(false)
const saving = ref(false)
const notice = ref('')
const provider = ref('gcp')
const form = ref({ name: '', bucket: '', credentials: '' })
const awsPlaceholder = '{"access_key_id":"...","secret_access_key":"...","region":"..."}'

const activeButtonStyle = { flex: '1', padding: '0.45rem', borderRadius: '6px', border: '1px solid #e6e6e6', background: '#111827', color: '#fff' }
const inactiveButtonStyle = { flex: '1', padding: '0.45rem', borderRadius: '6px', border: '1px solid #e6e6e6', background: 'transparent' }

async function fetchConnections() {
  loading.value = true
  try {
    const [g1, g2] = await Promise.all([
      fetch('/api/gcp/connections').then(r => r.ok ? r.json() : []),
      fetch('/api/aws/connections').then(r => r.ok ? r.json() : [])
    ])
    const gcpList = (g1 || []).map(c => ({ ...c, provider: 'gcp' }))
    const awsList = (g2 || []).map(c => ({ ...c, provider: 'aws' }))
    connections.value = [...gcpList, ...awsList]
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(fetchConnections)

function useConnection(c) {
  provider.value = c.provider
  form.value.name = c.name
  form.value.bucket = c.bucket
  form.value.credentials = c.credentials
  notice.value = 'Loaded connection into form.'
}

defineExpose({})

async function testConnection() {
  testing.value = true
  notice.value = ''
  try {
    const endpoint = provider.value === 'gcp' ? '/api/gcp/test' : '/api/aws/test'
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ bucket: form.value.bucket, credentials: form.value.credentials })
    })
    if (!res.ok) {
      const txt = await res.text()
      notice.value = 'Test failed: ' + txt
    } else {
      notice.value = 'Connection test succeeded.'
    }
  } catch (err) {
    notice.value = 'Error: ' + err.message
  } finally { testing.value = false }
}

async function saveConnection() {
  saving.value = true
  notice.value = ''
  try {
    const endpoint = provider.value === 'gcp' ? '/api/gcp/connection' : '/api/aws/connection'
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })
    if (!res.ok) {
      const txt = await res.text()
      notice.value = 'Save failed: ' + txt
    } else {
      notice.value = 'Connection saved.'
      form.value = { name: '', bucket: '', credentials: '' }
      await fetchConnections()
    }
  } catch (err) {
    notice.value = 'Error: ' + err.message
  } finally { saving.value = false }
}

async function removeConnection(p, id) {
  try {
    const endpoint = p === 'gcp' ? `/api/gcp/connection/${id}` : `/api/aws/connection/${id}`
    const res = await fetch(endpoint, { method: 'DELETE' })
    if (res.ok) {
      await fetchConnections()
    }
  } catch (err) {
    console.error(err)
  }
}
</script>
