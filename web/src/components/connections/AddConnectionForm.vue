<template>
  <div class="form-view">
    <div class="form-view__hd">
      <h1 class="form-view__title">{{ editConn ? 'Edit Connection' : 'New Connection' }}</h1>
      <p class="form-view__sub">
        {{ editConn ? 'Update the connection details below.' : 'Test and save a cloud bucket to start browsing files.' }}
      </p>
    </div>

    <!-- Provider toggle (disabled in edit mode) -->
    <div class="provider-tabs">
      <button
        class="provider-tab"
        :class="{ 'provider-tab--active': provider === 'gcp' }"
        :disabled="!!editConn"
        @click="!editConn && (provider = 'gcp')"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 15a4 4 0 0 0 4 4h9a5 5 0 0 0 1.8-9.7 6 6 0 0 0-11.8-1A4 4 0 0 0 3 15z"/>
        </svg>
        Google Cloud Storage
      </button>
      <button
        class="provider-tab"
        :class="{ 'provider-tab--active': provider === 'aws' }"
        :disabled="!!editConn"
        @click="!editConn && (provider = 'aws')"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
          <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
        </svg>
        AWS S3
      </button>
    </div>

    <form @submit.prevent="handleSave" class="conn-form">
      <div class="form-group">
        <label class="form-label">Name</label>
        <BaseInput v-model="form.name" placeholder="e.g. Production GCS" />
      </div>

      <div class="form-group">
        <label class="form-label">Bucket</label>
        <BaseInput v-model="form.bucket" :placeholder="provider === 'gcp' ? 'my-bucket-name' : 'my-s3-bucket'" />
      </div>

      <div class="form-group">
        <label class="form-label">
          {{ provider === 'gcp' ? 'Service account JSON' : 'Credentials JSON' }}
          <span class="form-label-optional" v-if="provider === 'gcp'">(optional for public buckets)</span>
        </label>
        <textarea
          class="base-textarea"
          v-model="form.credentials"
          rows="6"
          :placeholder="provider === 'gcp' ? gcpPlaceholder : awsPlaceholder"
        ></textarea>
        <p v-if="provider === 'gcp'" class="form-hint">
          Leave empty to connect to a publicly accessible GCS bucket.
        </p>
        <p v-else class="form-hint">
          For Cloudflare R2 or MinIO, include an <code style="font-family:var(--mono);font-size:11px">"endpoint"</code> key pointing to your custom S3-compatible URL.
        </p>
      </div>

      <StatusNotice :message="error"  type="error"   />
      <StatusNotice :message="notice" type="success" />

      <div class="form-actions">
        <BaseButton type="button" variant="ghost" :loading="testing" @click="handleTest">
          {{ testing ? 'Testing…' : 'Test connection' }}
        </BaseButton>
        <BaseButton type="submit" variant="primary" :loading="saving">
          {{ saving ? 'Saving…' : (editConn ? 'Update' : 'Save') }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import BaseInput    from '../ui/BaseInput.vue'
import BaseButton   from '../ui/BaseButton.vue'
import StatusNotice from '../ui/StatusNotice.vue'

const props = defineProps({
  testing:  { type: Boolean, default: false },
  saving:   { type: Boolean, default: false },
  error:    { type: String,  default: '' },
  notice:   { type: String,  default: '' },
  editConn: { type: Object,  default: null }, // null = create mode
})

const emit = defineEmits(['test', 'save'])

const provider = ref(props.editConn?.provider ?? 'gcp')
const form     = ref({
  name:        props.editConn?.name        ?? '',
  bucket:      props.editConn?.bucket      ?? '',
  credentials: props.editConn?.credentials ?? '',
})

// Re-sync when editConn changes (e.g. switching which connection to edit)
watch(() => props.editConn, conn => {
  provider.value = conn?.provider ?? 'gcp'
  form.value = {
    name:        conn?.name        ?? '',
    bucket:      conn?.bucket      ?? '',
    credentials: conn?.credentials ?? '',
  }
})

const gcpPlaceholder = `{\n  "type": "service_account",\n  "project_id": "...",\n  ...\n}`
const awsPlaceholder = `{\n  "access_key_id": "...",\n  "secret_access_key": "...",\n  "region": "us-east-1",\n  "endpoint": "https://...r2.cloudflarestorage.com"  ← optional, for R2/MinIO\n}`

function handleTest() {
  emit('test', provider.value, form.value.bucket, form.value.credentials)
}

async function handleSave() {
  const success = await new Promise(resolve =>
    emit('save', provider.value, { ...form.value }, resolve, props.editConn?.id ?? null)
  )
  if (success && !props.editConn) {
    // Only reset form on create; keep values visible on edit
    form.value = { name: '', bucket: '', credentials: '' }
  }
}
</script>
