<template>
  <div class="app-root">
    <AppHeader />
    <main class="app-main">
      <AppLayout>
        <template #list>
          <ConnectionList
            :connections="connections"
            :loading="loading"
            @use="handleUse"
            @delete="removeConnection"
          />
        </template>
        <template #form>
          <AddConnectionForm
            ref="formRef"
            :testing="testing"
            :saving="saving"
            :error="error"
            :notice="notice"
            @test="testConnection"
            @save="handleSave"
          />
        </template>
      </AppLayout>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppHeader from './components/layout/AppHeader.vue'
import AppLayout from './components/layout/AppLayout.vue'
import ConnectionList from './components/connections/ConnectionList.vue'
import AddConnectionForm from './components/connections/AddConnectionForm.vue'
import { useConnections } from './composables/useConnections.js'

const {
  connections, loading, testing, saving, error, notice,
  fetchConnections, testConnection, saveConnection, removeConnection,
} = useConnections()

const formRef = ref(null)

onMounted(fetchConnections)

function handleUse(conn) {
  formRef.value?.loadConnection(conn)
  notice.value = 'Connection loaded into form.'
}

async function handleSave(provider, form, resolve) {
  const success = await saveConnection(provider, form)
  resolve?.(success)
}
</script>
