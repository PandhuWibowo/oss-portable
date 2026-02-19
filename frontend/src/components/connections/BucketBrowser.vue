<template>
  <div class="browser-view" ref="rootEl"
       @dragover.prevent="onDragOver"
       @dragleave.self="onDragLeave"
       @drop.prevent="onDrop">

    <!-- ── Header ──────────────────────────────────────────────── -->
    <div class="browser-hd">
      <div class="browser-hd__left">
        <div class="browser-prov-icon" :class="`browser-prov-icon--${conn.provider}`">
          {{ conn.provider === 'gcp' ? 'G' : 'A' }}
        </div>
        <div style="min-width:0">
          <div class="browser-conn-name">
            {{ conn.name }}
            <BaseBadge :provider="conn.provider" />
          </div>
          <div class="browser-conn-bucket">{{ conn.bucket }}</div>

          <!-- Breadcrumbs -->
          <div class="breadcrumbs" v-if="currentPrefix">
            <button class="bread-item" @click="navigateTo('')">root</button>
            <template v-for="(crumb, i) in breadcrumbs" :key="i">
              <span class="bread-sep">/</span>
              <button
                class="bread-item"
                :style="i === breadcrumbs.length - 1 ? 'color:var(--text-2);cursor:default' : ''"
                @click="i < breadcrumbs.length - 1 && navigateTo(crumb.prefix)"
              >{{ crumb.label }}</button>
            </template>
          </div>
        </div>
      </div>

      <div class="browser-hd__actions">
        <!-- Stats toggle -->
        <button
          class="icon-btn"
          :style="showStats ? 'background:var(--accent-bg);color:var(--accent);border-color:var(--accent-ring)' : ''"
          @click="toggleStats"
          title="Bucket stats"
        >
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>
          </svg>
        </button>

        <!-- Refresh -->
        <button class="icon-btn" :disabled="loading" @click="refresh" title="Refresh (R)">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
               :style="loading ? 'animation:spin .6s linear infinite' : ''">
            <polyline points="23 4 23 10 17 10"/>
            <polyline points="1 20 1 14 7 14"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
        </button>

        <!-- Delete connection -->
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

    <!-- ── Stats bar ────────────────────────────────────────────── -->
    <transition name="slide-down">
      <div v-if="showStats" class="stats-bar">
        <template v-if="statsLoading">
          <div class="stat-item">
            <div class="skeleton-item" style="width:60px;height:22px;border-radius:4px"></div>
            <span class="stat-lbl">loading…</span>
          </div>
        </template>
        <template v-else-if="stats">
          <div class="stat-item">
            <span class="stat-val">{{ stats.object_count.toLocaleString() }}</span>
            <span class="stat-lbl">{{ stats.truncated ? 'objects (est.)' : 'objects' }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-val">{{ formatSize(stats.total_size) }}</span>
            <span class="stat-lbl">total size</span>
          </div>
        </template>
        <template v-else-if="statsError">
          <span style="font-size:12px;color:var(--muted)">{{ statsError }}</span>
        </template>
      </div>
    </transition>

    <!-- ── Toolbar ──────────────────────────────────────────────── -->
    <div class="browser-toolbar">
      <div class="search-field">
        <span class="search-field__icon">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
        </span>
        <input
          ref="searchInput"
          class="search-field__input"
          v-model="searchQuery"
          placeholder="Search files… (/)"
          @keydown.escape.stop="searchQuery = ''"
        />
        <button v-if="searchQuery" class="search-field__clear" @click="searchQuery = ''" aria-label="Clear search">×</button>
      </div>

      <div class="toolbar-spacer"></div>

      <label class="upload-label" title="Upload files">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/>
          <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/>
        </svg>
        Upload
        <input type="file" multiple style="display:none" @change="onFileInput" />
      </label>
    </div>

    <!-- ── Upload progress ──────────────────────────────────────── -->
    <transition name="slide-down">
      <div v-if="uploading" class="upload-progress">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="animation:spin .6s linear infinite">
          <line x1="12" y1="2" x2="12" y2="6"/><line x1="12" y1="18" x2="12" y2="22"/>
          <line x1="4.93" y1="4.93" x2="7.76" y2="7.76"/><line x1="16.24" y1="16.24" x2="19.07" y2="19.07"/>
          <line x1="2" y1="12" x2="6" y2="12"/><line x1="18" y1="12" x2="22" y2="12"/>
          <line x1="4.93" y1="19.07" x2="7.76" y2="16.24"/><line x1="16.24" y1="7.76" x2="19.07" y2="4.93"/>
        </svg>
        <span>Uploading {{ uploadingCount }} file{{ uploadingCount !== 1 ? 's' : '' }}…</span>
        <div class="progress-bar"><div class="progress-fill" style="width:100%"></div></div>
      </div>
    </transition>

    <!-- ── Upload error ─────────────────────────────────────────── -->
    <transition name="slide-down">
      <div v-if="uploadError" class="upload-progress" style="background:var(--danger-bg);border-color:var(--danger);color:var(--danger)">
        <span>{{ uploadError }}</span>
        <button @click="uploadError=''" style="margin-left:auto;background:none;border:none;cursor:pointer;color:inherit;font-size:16px">×</button>
      </div>
    </transition>

    <!-- ── Body ─────────────────────────────────────────────────── -->
    <div class="browser-body" ref="bodyEl" style="position:relative;overflow-y:auto;flex:1">

      <!-- Drag overlay -->
      <div v-if="isDragging" class="drop-overlay">
        <div class="drop-overlay__inner">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-bottom:8px">
            <polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/>
            <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/>
          </svg>
          Drop to upload here
        </div>
      </div>

      <!-- Loading -->
      <table v-if="loading" class="file-table">
        <thead>
          <tr>
            <th>Name</th><th>Size</th><th>Modified</th><th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="i in 10" :key="i">
            <td colspan="4">
              <div class="skeleton-item" :style="`height:16px;border-radius:3px;width:${30 + (i * 37 % 55)}%`"></div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Browse error -->
      <div v-else-if="browseError" class="empty-state">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.4;margin-bottom:6px">
          <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <p style="font-size:13px;color:var(--text-2);max-width:320px">{{ browseError }}</p>
        <button class="base-btn base-btn--ghost" @click="refresh" style="font-size:12px;padding:6px 12px;margin-top:4px">Retry</button>
      </div>

      <!-- Empty -->
      <div v-else-if="filteredEntries.length === 0" class="empty-state">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.4;margin-bottom:6px">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
        </svg>
        <p>{{ searchQuery ? 'No files match your search.' : 'This folder is empty.' }}</p>
        <button v-if="searchQuery" class="base-btn base-btn--ghost" @click="searchQuery = ''" style="font-size:12px;padding:6px 12px;margin-top:4px">Clear search</button>
      </div>

      <!-- File table -->
      <table v-else class="file-table">
        <thead>
          <tr>
            <th class="sortable" @click="cycleSort('name')">
              Name
              <span class="sort-indicator" :class="{ active: sortKey === 'name' }">
                {{ sortKey === 'name' ? (sortDir === 'asc' ? '↑' : '↓') : '↕' }}
              </span>
            </th>
            <th class="sortable" @click="cycleSort('size')">
              Size
              <span class="sort-indicator" :class="{ active: sortKey === 'size' }">
                {{ sortKey === 'size' ? (sortDir === 'asc' ? '↑' : '↓') : '↕' }}
              </span>
            </th>
            <th class="sortable" @click="cycleSort('date')">
              Modified
              <span class="sort-indicator" :class="{ active: sortKey === 'date' }">
                {{ sortKey === 'date' ? (sortDir === 'asc' ? '↑' : '↓') : '↕' }}
              </span>
            </th>
            <th style="width:120px"></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="entry in filteredEntries"
            :key="entry.name"
            :class="{ 'is-dir': entry.type === 'dir' }"
          >
            <td>
              <div
                class="file-name"
                :style="entry.type === 'dir' ? 'cursor:pointer' : ''"
                @click="entry.type === 'dir' ? navigateTo(entry.name) : null"
              >
                <!-- Folder icon -->
                <svg v-if="entry.type === 'dir'" class="file-icon" width="13" height="13" viewBox="0 0 24 24" fill="currentColor" stroke="none" style="color:var(--aws)">
                  <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" opacity=".8"/>
                </svg>
                <!-- File icon -->
                <svg v-else class="file-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
                  <polyline points="13 2 13 9 20 9"/>
                </svg>
                {{ entry.display }}
              </div>
            </td>
            <td class="file-size">{{ entry.type === 'dir' ? '—' : formatSize(entry.size) }}</td>
            <td class="file-date">{{ entry.type === 'dir' ? '—' : formatDate(entry.updated) }}</td>
            <td class="file-actions">
              <!-- Copy path -->
              <button class="row-btn" @click.stop="copyPath(entry)" title="Copy path">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
              </button>
              <template v-if="entry.type === 'file'">
                <!-- Download -->
                <button class="row-btn" @click.stop="download(entry)" title="Download">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="7 10 12 15 17 10"/>
                    <line x1="12" y1="15" x2="12" y2="3"/>
                  </svg>
                </button>
                <!-- Preview -->
                <button class="row-btn" @click.stop="openPreview(entry)" title="Preview">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                </button>
                <!-- Delete -->
                <button class="row-btn danger" @click.stop="confirmDelete(entry)" title="Delete">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                    <path d="M10 11v6"/><path d="M14 11v6"/>
                    <path d="M9 6V4h6v2"/>
                  </svg>
                </button>
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ── Preview panel ────────────────────────────────────────── -->
    <transition name="slide-right">
      <div v-if="previewEntry" class="preview-panel">
        <div class="preview-hd">
          <span class="preview-hd__name">{{ previewEntry.display }}</span>
          <button class="preview-close" @click="closePreview" aria-label="Close preview">×</button>
        </div>

        <div class="preview-body">
          <!-- Loading -->
          <div v-if="previewLoading" class="preview-unsupported">
            <div class="base-btn__spinner" style="width:20px;height:20px;border-width:2px"></div>
            <span style="margin-top:10px;font-size:12px;color:var(--muted)">Loading preview…</span>
          </div>
          <!-- Image -->
          <img v-else-if="isImage(previewEntry) && previewUrl" :src="previewUrl" class="preview-img" @error="previewLoadError = true" />
          <!-- Text / JSON -->
          <pre v-else-if="isText(previewEntry) && previewContent" class="preview-text">{{ previewContent }}</pre>
          <!-- Error or unsupported -->
          <div v-else class="preview-unsupported">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.3">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
              <polyline points="13 2 13 9 20 9"/>
            </svg>
            <p>{{ previewLoadError ? 'Failed to load preview.' : 'No preview for this file type.' }}</p>
            <button class="base-btn base-btn--ghost" @click="download(previewEntry)" style="font-size:12px;padding:6px 12px;margin-top:4px">Download file</button>
          </div>
        </div>

        <div class="preview-ft">
          <span class="preview-meta">
            {{ formatSize(previewEntry.size) }}
            <template v-if="previewEntry.updated"> · {{ formatDate(previewEntry.updated) }}</template>
            <template v-if="previewEntry.content_type"> · {{ previewEntry.content_type }}</template>
          </span>
          <button class="base-btn base-btn--ghost" @click="download(previewEntry)" style="font-size:12px;padding:5px 10px">
            Download
          </button>
        </div>
      </div>
    </transition>

  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import BaseBadge from '../ui/BaseBadge.vue'
import { useConnections } from '../../composables/useConnections.js'

const props = defineProps({
  conn: { type: Object, required: true },
})

defineEmits(['delete'])

const {
  browseObjects,
  getDownloadURL,
  deleteObject,
  uploadObjects,
  getBucketStats,
} = useConnections()

// ── State ──────────────────────────────────────────────────────
const currentPrefix  = ref('')
const entries        = ref([])
const loading        = ref(false)
const browseError    = ref('')

const searchQuery    = ref('')
const sortKey        = ref('name')   // 'name' | 'size' | 'date'
const sortDir        = ref('asc')    // 'asc' | 'desc'

const stats          = ref(null)
const statsLoading   = ref(false)
const statsError     = ref('')
const showStats      = ref(false)
const statsLoaded    = ref(false)

const uploading      = ref(false)
const uploadingCount = ref(0)
const uploadError    = ref('')
const isDragging     = ref(false)

const previewEntry   = ref(null)
const previewUrl     = ref('')
const previewContent = ref('')
const previewLoading = ref(false)
const previewLoadError = ref(false)

const searchInput    = ref(null)
const bodyEl         = ref(null)

// ── Breadcrumbs ─────────────────────────────────────────────────
const breadcrumbs = computed(() => {
  if (!currentPrefix.value) return []
  const parts = currentPrefix.value.replace(/\/$/, '').split('/')
  return parts.map((label, i) => ({
    label,
    prefix: parts.slice(0, i + 1).join('/') + '/',
  }))
})

// ── Filtered + sorted entries ───────────────────────────────────
const filteredEntries = computed(() => {
  let list = entries.value

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(e => e.display.toLowerCase().includes(q))
  }

  const dirs  = list.filter(e => e.type === 'dir')
  const files = list.filter(e => e.type === 'file')

  const sortFn = (a, b) => {
    let va, vb
    if (sortKey.value === 'size') {
      va = a.size ?? 0; vb = b.size ?? 0
    } else if (sortKey.value === 'date') {
      va = a.updated ? new Date(a.updated).getTime() : 0
      vb = b.updated ? new Date(b.updated).getTime() : 0
    } else {
      va = a.display.toLowerCase(); vb = b.display.toLowerCase()
    }
    if (va < vb) return sortDir.value === 'asc' ? -1 : 1
    if (va > vb) return sortDir.value === 'asc' ? 1 : -1
    return 0
  }

  return [
    ...dirs.sort(sortFn),
    ...files.sort(sortFn),
  ]
})

// ── Navigation ──────────────────────────────────────────────────
function navigateTo(prefix) {
  currentPrefix.value = prefix
  searchQuery.value = ''
  load()
}

function navigateUp() {
  if (!currentPrefix.value) return
  const trimmed = currentPrefix.value.replace(/\/$/, '')
  const parent  = trimmed.includes('/') ? trimmed.slice(0, trimmed.lastIndexOf('/') + 1) : ''
  navigateTo(parent)
}

// ── Data loading ────────────────────────────────────────────────
async function load() {
  loading.value    = true
  browseError.value = ''
  entries.value    = []
  try {
    const result = await browseObjects(
      props.conn.provider,
      props.conn.bucket,
      props.conn.credentials,
      currentPrefix.value,
    )
    entries.value = result.entries ?? []
  } catch (err) {
    browseError.value = err.message
  } finally {
    loading.value = false
  }
}

function refresh() {
  load()
}

// ── Stats ───────────────────────────────────────────────────────
async function loadStats() {
  if (statsLoaded.value) return
  statsLoading.value = true
  statsError.value   = ''
  try {
    stats.value      = await getBucketStats(props.conn.provider, props.conn.bucket, props.conn.credentials)
    statsLoaded.value = true
  } catch (err) {
    statsError.value  = 'Stats unavailable'
  } finally {
    statsLoading.value = false
  }
}

function toggleStats() {
  showStats.value = !showStats.value
  if (showStats.value && !statsLoaded.value) loadStats()
}

// ── Sort ────────────────────────────────────────────────────────
function cycleSort(key) {
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortDir.value = 'asc'
  }
}

// ── Copy path ───────────────────────────────────────────────────
function copyPath(entry) {
  const path = `gs://${props.conn.bucket}/${entry.name}`
  navigator.clipboard?.writeText(path).catch(() => {})
}

// ── Download ────────────────────────────────────────────────────
async function download(entry) {
  try {
    const url = await getDownloadURL(
      props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name,
    )
    const a = document.createElement('a')
    a.href = url
    a.download = entry.display
    a.target = '_blank'
    a.rel = 'noopener'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
  } catch (err) {
    uploadError.value = 'Download failed: ' + err.message
  }
}

// ── Delete ──────────────────────────────────────────────────────
async function confirmDelete(entry) {
  if (!window.confirm(`Delete "${entry.display}"? This cannot be undone.`)) return
  try {
    await deleteObject(props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name)
    if (previewEntry.value?.name === entry.name) closePreview()
    await load()
    if (statsLoaded.value) { statsLoaded.value = false; loadStats() }
  } catch (err) {
    uploadError.value = 'Delete failed: ' + err.message
  }
}

// ── Upload ──────────────────────────────────────────────────────
async function handleUpload(files) {
  if (!files || files.length === 0) return
  uploading.value      = true
  uploadingCount.value = files.length
  uploadError.value    = ''
  try {
    await uploadObjects(
      props.conn.provider, props.conn.bucket, props.conn.credentials, currentPrefix.value, files,
    )
    await load()
    if (statsLoaded.value) { statsLoaded.value = false; loadStats() }
  } catch (err) {
    uploadError.value = 'Upload failed: ' + err.message
  } finally {
    uploading.value = false
  }
}

function onFileInput(e) {
  handleUpload(e.target.files)
  e.target.value = ''
}

// ── Drag-and-drop ───────────────────────────────────────────────
let dragCounter = 0
function onDragOver(e) {
  if (!e.dataTransfer?.types.includes('Files')) return
  dragCounter++
  isDragging.value = true
}
function onDragLeave() {
  dragCounter--
  if (dragCounter <= 0) { dragCounter = 0; isDragging.value = false }
}
function onDrop(e) {
  dragCounter = 0
  isDragging.value = false
  handleUpload(e.dataTransfer?.files)
}

// ── Preview ─────────────────────────────────────────────────────
function isImage(entry) {
  if (!entry) return false
  const ct = (entry.content_type || '').toLowerCase()
  if (ct.startsWith('image/')) return true
  const ext = entry.display.split('.').pop().toLowerCase()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg', 'ico', 'bmp'].includes(ext)
}

function isText(entry) {
  if (!entry) return false
  const ct = (entry.content_type || '').toLowerCase()
  if (ct.startsWith('text/') || ct.includes('json') || ct.includes('xml') || ct.includes('javascript')) return true
  const ext = entry.display.split('.').pop().toLowerCase()
  return ['txt', 'md', 'json', 'yaml', 'yml', 'toml', 'csv', 'xml', 'html', 'js', 'ts', 'py', 'sh', 'log', 'conf', 'ini'].includes(ext)
}

async function openPreview(entry) {
  previewEntry.value    = entry
  previewUrl.value      = ''
  previewContent.value  = ''
  previewLoadError.value = false
  previewLoading.value  = true

  try {
    const url = await getDownloadURL(
      props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name,
    )
    previewUrl.value = url

    if (isText(entry)) {
      const res = await fetch(url)
      if (res.ok) {
        const text = await res.text()
        previewContent.value = text.slice(0, 50_000) // cap at 50k chars
      } else {
        previewLoadError.value = true
      }
    }
  } catch {
    previewLoadError.value = true
  } finally {
    previewLoading.value = false
  }
}

function closePreview() {
  previewEntry.value = null
  previewUrl.value   = ''
  previewContent.value = ''
}

// ── Keyboard shortcuts ──────────────────────────────────────────
function onKeyDown(e) {
  const tag = document.activeElement?.tagName
  const inInput = ['INPUT', 'TEXTAREA'].includes(tag)

  // / → focus search
  if (e.key === '/' && !inInput) {
    e.preventDefault()
    searchInput.value?.focus()
    return
  }
  // R → refresh (not in input)
  if ((e.key === 'r' || e.key === 'R') && !inInput && !e.metaKey && !e.ctrlKey) {
    refresh()
    return
  }
  // Escape → close preview or clear search
  if (e.key === 'Escape') {
    if (previewEntry.value) { closePreview(); return }
    if (searchQuery.value)  { searchQuery.value = ''; return }
  }
  // Backspace → navigate up (not in input)
  if (e.key === 'Backspace' && !inInput && currentPrefix.value) {
    e.preventDefault()
    navigateUp()
  }
}

// ── Watchers ────────────────────────────────────────────────────
watch(() => props.conn, () => {
  currentPrefix.value  = ''
  searchQuery.value    = ''
  stats.value          = null
  statsLoaded.value    = false
  statsError.value     = ''
  previewEntry.value   = null
  load()
})

onMounted(() => {
  load()
  window.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
})

// ── Formatters ──────────────────────────────────────────────────
function formatSize(bytes) {
  if (!bytes && bytes !== 0) return '—'
  if (bytes === 0) return '0 B'
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
