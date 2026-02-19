<template>
  <div class="browser-view">
    <!-- Header -->
    <div class="browser-hd">
      <div class="browser-hd__left">
        <div class="browser-prov-icon" :class="`browser-prov-icon--${conn.provider}`">
          {{ conn.provider === 'gcp' ? 'G' : 'A' }}
        </div>
        <div>
          <div class="browser-conn-name">
            {{ conn.name }}
            <BaseBadge :provider="conn.provider" />
          </div>
          <div class="browser-conn-bucket">{{ conn.bucket }}</div>
        </div>
      </div>

      <div class="browser-hd__actions">
        <button class="icon-btn" :disabled="loading" @click="load" title="Refresh">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :style="loading ? 'animation:spin .6s linear infinite' : ''">
            <polyline points="23 4 23 10 17 10"/>
            <polyline points="1 20 1 14 7 14"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
        </button>
        <button class="icon-btn danger" @click="$emit('delete')" title="Delete connection">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="3 6 5 6 21 6"/>
            <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
            <path d="M10 11v6"/><path d="M14 11v6"/>
            <path d="M9 6V4h6v2"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Meta bar -->
    <div v-if="!loading && !error && objects.length > 0" class="browser-meta">
      {{ objects.length.toLocaleString() }} object{{ objects.length !== 1 ? 's' : '' }}
      <span v-if="truncated">&nbsp;— showing first 1,000</span>
    </div>

    <!-- Body -->
    <div class="browser-body">
      <!-- Loading -->
      <table v-if="loading" class="file-table">
        <thead>
          <tr>
            <th>Name</th><th>Size</th><th>Type</th><th>Modified</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="i in 8" :key="i">
            <td colspan="4">
              <div class="skeleton-item" style="height:18px;border-radius:4px;width:100%"></div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Error -->
      <div v-else-if="error" class="empty-state">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.4;margin-bottom:6px">
          <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <p style="font-size:13px;color:var(--text-2);max-width:300px">{{ error }}</p>
      </div>

      <!-- Empty -->
      <div v-else-if="objects.length === 0" class="empty-state">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.4;margin-bottom:6px">
          <ellipse cx="12" cy="5" rx="9" ry="3"/>
          <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
          <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
        </svg>
        <p>This bucket is empty.</p>
      </div>

      <!-- File table -->
      <table v-else class="file-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Size</th>
            <th v-if="conn.provider === 'gcp'">Type</th>
            <th>Modified</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="obj in objects" :key="obj.name">
            <td>
              <div class="file-name">
                <svg class="file-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
                  <polyline points="13 2 13 9 20 9"/>
                </svg>
                {{ obj.name }}
              </div>
            </td>
            <td class="file-size">{{ formatSize(obj.size) }}</td>
            <td v-if="conn.provider === 'gcp'" class="file-type">{{ obj.content_type || '—' }}</td>
            <td class="file-date">{{ formatDate(obj.updated) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import BaseBadge from '../ui/BaseBadge.vue'
import { useConnections } from '../../composables/useConnections.js'

const props = defineProps({
  conn: { type: Object, required: true },
})

defineEmits(['delete'])

const { listObjects } = useConnections()

const objects  = ref([])
const truncated = ref(false)
const loading  = ref(false)
const error    = ref('')

async function load() {
  loading.value = true
  error.value   = ''
  objects.value  = []
  truncated.value = false
  try {
    const result = await listObjects(props.conn.provider, props.conn.bucket, props.conn.credentials)
    objects.value   = result.objects
    truncated.value = result.truncated
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => props.conn, load)

function formatSize(bytes) {
  if (!bytes) return '—'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
  return (bytes / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1) + ' ' + units[i]
}

function formatDate(iso) {
  if (!iso) return '—'
  try {
    return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  } catch { return '—' }
}
</script>
