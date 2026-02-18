<template>
  <aside class="panel">
    <h2 class="panel__title" style="margin-bottom:4px">Add Connection</h2>
    <p class="panel__subtitle" style="margin-bottom:16px">Test and save a new cloud bucket</p>

    <!-- Provider Toggle -->
    <div class="provider-tabs">
      <button
        class="provider-tab"
        :class="{ 'provider-tab--active': provider === 'gcp' }"
        @click="provider = 'gcp'"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 15a4 4 0 0 0 4 4h9a5 5 0 0 0 1.8-9.7 6 6 0 0 0-11.8-1A4 4 0 0 0 3 15z"/>
        </svg>
        GCP
      </button>
      <button
        class="provider-tab"
        :class="{ 'provider-tab--active': provider === 'aws' }"
        @click="provider = 'aws'"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
        </svg>
        AWS S3
      </button>
    </div>

    <form @submit.prevent="handleSave" class="conn-form">
      <div class="form-group">
        <label class="form-label">Connection name</label>
        <BaseInput v-model="form.name" placeholder="e.g. Production GCS" />
      </div>
      <div class="form-group">
        <label class="form-label">Bucket name</label>
        <BaseInput v-model="form.bucket" placeholder="my-bucket-name" />
      </div>
      <div class="form-group">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:5px">
          <label class="form-label">
            {{ provider === 'gcp' ? 'Service account JSON' : 'AWS credentials (JSON)' }}
            <span v-if="provider === 'gcp'" class="form-label-optional">(optional for public buckets)</span>
          </label>
        </div>
        <textarea
          class="base-textarea"
          v-model="form.credentials"
          rows="7"
          :placeholder="provider === 'gcp' ? gcpPlaceholder : awsPlaceholder"
        ></textarea>
        <p v-if="provider === 'gcp'" class="form-hint">
          Leave empty to connect to a public GCS bucket without authentication.
        </p>
      </div>

      <StatusNotice :message="error" type="error" />
      <StatusNotice :message="notice" type="success" />

      <div class="form-actions">
        <BaseButton type="button" variant="ghost" :loading="testing" @click="handleTest">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="9 11 12 14 22 4"/><path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
          </svg>
          {{ testing ? 'Testing…' : 'Test' }}
        </BaseButton>
        <BaseButton type="submit" variant="primary" :loading="saving">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/>
          </svg>
          {{ saving ? 'Saving…' : 'Save' }}
        </BaseButton>
      </div>
    </form>
  </aside>
</template>

<script setup>
import { ref } from 'vue'
import BaseInput from '../ui/BaseInput.vue'
import BaseButton from '../ui/BaseButton.vue'
import StatusNotice from '../ui/StatusNotice.vue'

const props = defineProps({
  testing: { type: Boolean, default: false },
  saving: { type: Boolean, default: false },
  error: { type: String, default: '' },
  notice: { type: String, default: '' },
})
const emit = defineEmits(['test', 'save', 'load'])

const provider = ref('gcp')
const form = ref({ name: '', bucket: '', credentials: '' })

const gcpPlaceholder = `{\n  "type": "service_account",\n  "project_id": "...",\n  ...\n}`
const awsPlaceholder = `{\n  "access_key_id": "...",\n  "secret_access_key": "...",\n  "region": "us-east-1"\n}`

function handleTest() {
  emit('test', provider.value, form.value.bucket, form.value.credentials)
}

async function handleSave() {
  const success = await new Promise(resolve => emit('save', provider.value, { ...form.value }, resolve))
  if (success) {
    form.value = { name: '', bucket: '', credentials: '' }
  }
}

function loadConnection(conn) {
  provider.value = conn.provider
  form.value.name = conn.name
  form.value.bucket = conn.bucket
  form.value.credentials = conn.credentials
}

defineExpose({ loadConnection })
</script>
