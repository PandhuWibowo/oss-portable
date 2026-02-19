<template>
  <div class="shell">
    <!-- Sidebar -->
    <AppSidebar
      :connections="connections"
      :loading="loading"
      :activeConn="activeConn"
      @new-connection="startNew"
      @select="handleSelect"
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

      <!-- Add connection form -->
      <AddConnectionForm
        v-else-if="mode === 'form'"
        :testing="testing"
        :saving="saving"
        :error="error"
        :notice="notice"
        @test="testConnection"
        @save="handleSave"
      />

      <!-- Bucket browser -->
      <BucketBrowser
        v-else-if="mode === 'browse' && activeConn"
        :conn="activeConn"
        @delete="handleDelete(activeConn.provider, activeConn.id)"
      />

    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppSidebar       from './components/layout/AppHeader.vue'
import AddConnectionForm from './components/connections/AddConnectionForm.vue'
import BucketBrowser    from './components/connections/BucketBrowser.vue'
import BaseButton       from './components/ui/BaseButton.vue'
import { useConnections } from './composables/useConnections.js'

const {
  connections, loading, testing, saving, error, notice,
  fetchConnections, testConnection, saveConnection, removeConnection, clearMessages,
} = useConnections()

const mode       = ref('welcome') // 'welcome' | 'form' | 'browse'
const activeConn = ref(null)

onMounted(fetchConnections)

function handleSelect(conn) {
  activeConn.value = conn
  mode.value = 'browse'
  clearMessages()
}

function startNew() {
  activeConn.value = null
  mode.value = 'form'
  clearMessages()
}

async function handleSave(provider, form, resolve) {
  const success = await saveConnection(provider, form)
  if (success) {
    // Auto-select the newly saved connection (it's first in the list)
    const saved = connections.value.find(
      c => c.provider === provider && c.name === form.name
    )
    if (saved) {
      activeConn.value = saved
      mode.value = 'browse'
    }
  }
  resolve?.(success)
}

function handleDelete(provider, id) {
  if (activeConn.value?.id === id && activeConn.value?.provider === provider) {
    activeConn.value = null
    mode.value = 'welcome'
  }
  removeConnection(provider, id)
}
</script>
