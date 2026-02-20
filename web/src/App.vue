<template>
  <div class="shell">
    <!-- Sidebar -->
    <AppSidebar
      :connections="connections"
      :loading="loading"
      :activeConn="activeConn"
      @new-connection="startNew"
      @select="handleSelect"
      @edit="handleEdit"
      @delete="handleDelete"
    />

    <!-- Main area -->
    <main class="main">

      <!-- Welcome -->
      <div v-if="mode === 'welcome'" class="welcome">
        <div class="welcome__icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 15a4 4 0 0 0 4 4h9a5 5 0 0 0 1.8-9.7 6 6 0 0 0-11.8-1A4 4 0 0 0 3 15z"/>
          </svg>
        </div>
        <p class="welcome__title">
          {{ connections.length ? 'Select a connection' : 'No connections yet' }}
        </p>
        <p class="welcome__sub">
          {{ connections.length
            ? 'Choose a connection from the sidebar to browse its files.'
            : 'Add your first cloud bucket to start browsing files.' }}
        </p>
        <BaseButton v-if="!connections.length" variant="primary" @click="startNew" style="margin-top:8px">
          New Connection
        </BaseButton>
      </div>

      <!-- New connection form -->
      <AddConnectionForm
        v-else-if="mode === 'form'"
        :testing="testing"
        :saving="saving"
        :error="error"
        :notice="notice"
        @test="testConnection"
        @save="handleSave"
      />

      <!-- Edit connection form -->
      <AddConnectionForm
        v-else-if="mode === 'edit' && editingConn"
        :testing="testing"
        :saving="saving"
        :error="error"
        :notice="notice"
        :editConn="editingConn"
        @test="testConnection"
        @save="handleUpdate"
      />

      <!-- Bucket browser -->
      <BucketBrowser
        v-else-if="mode === 'browse' && activeConn"
        :conn="activeConn"
        @delete="handleDelete(activeConn.provider, activeConn.id)"
      />

    </main>

    <!-- Global overlays -->
    <ToastContainer />
    <ConfirmModal />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import AppSidebar        from './components/layout/AppHeader.vue'
import AddConnectionForm from './components/connections/AddConnectionForm.vue'
import BucketBrowser     from './components/connections/BucketBrowser.vue'
import BaseButton        from './components/ui/BaseButton.vue'
import ToastContainer    from './components/ui/ToastContainer.vue'
import ConfirmModal      from './components/ui/ConfirmModal.vue'
import { useConnections } from './composables/useConnections.js'
import { useToast }       from './composables/useToast.js'

const {
  connections, loading, testing, saving, error, notice,
  fetchConnections, testConnection, saveConnection, updateConnection,
  removeConnection, clearMessages,
} = useConnections()

const toast = useToast()

const mode        = ref('welcome') // 'welcome' | 'form' | 'edit' | 'browse'
const activeConn  = ref(null)
const editingConn = ref(null)

onMounted(() => {
  fetchConnections()
  window.addEventListener('keydown', onAppKeyDown)
})
onUnmounted(() => window.removeEventListener('keydown', onAppKeyDown))

function onAppKeyDown(e) {
  const inInput = ['INPUT', 'TEXTAREA'].includes(document.activeElement?.tagName)
  if ((e.key === 'n' || e.key === 'N') && !inInput && !e.metaKey && !e.ctrlKey) {
    startNew()
  }
}

// ── Navigation ────────────────────────────────────────────────

function handleSelect(conn) {
  activeConn.value  = conn
  editingConn.value = null
  mode.value        = 'browse'
  clearMessages()
}

function startNew() {
  editingConn.value = null
  mode.value        = 'form'
  clearMessages()
}

function handleEdit(conn) {
  editingConn.value = conn
  mode.value        = 'edit'
  clearMessages()
}

// ── Save / Update ─────────────────────────────────────────────

async function handleSave(provider, form, resolve) {
  const success = await saveConnection(provider, form)
  if (success) {
    const saved = connections.value.find(
      c => c.provider === provider && c.name === form.name
    )
    if (saved) {
      activeConn.value  = saved
      editingConn.value = null
      mode.value        = 'browse'
    }
    toast.success('Connection saved.')
  }
  resolve?.(success)
}

async function handleUpdate(provider, form, resolve, id) {
  const success = await updateConnection(provider, id, form)
  if (success) {
    // Refresh activeConn if we just edited it
    if (activeConn.value?.id === id && activeConn.value?.provider === provider) {
      const updated = connections.value.find(
        c => c.provider === provider && c.id === id
      )
      if (updated) activeConn.value = updated
    }
    editingConn.value = null
    mode.value        = activeConn.value ? 'browse' : 'welcome'
    toast.success('Connection updated.')
  }
  resolve?.(success)
}

// ── Delete ────────────────────────────────────────────────────

function handleDelete(provider, id) {
  if (activeConn.value?.id === id && activeConn.value?.provider === provider) {
    activeConn.value = null
    mode.value       = 'welcome'
  }
  if (editingConn.value?.id === id && editingConn.value?.provider === provider) {
    editingConn.value = null
    mode.value        = 'welcome'
  }
  removeConnection(provider, id)
}
</script>
