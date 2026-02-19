<template>
  <div class="form-view">
    <!-- Header -->
    <div class="form-view__hd">
      <h1 class="form-view__title">New Connection</h1>
      <p class="form-view__sub">Test and save a cloud bucket to start browsing files.</p>
    </div>

    <!-- Provider toggle -->
    <div class="provider-tabs">
      <button
        class="provider-tab"
        :class="{ 'provider-tab--active': provider === 'gcp' }"
        @click="provider = 'gcp'"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 15a4 4 0 0 0 4 4h9a5 5 0 0 0 1.8-9.7 6 6 0 0 0-11.8-1A4 4 0 0 0 3 15z"/>
        </svg>
        Google Cloud Storage
      </button>
      <button
        class="provider-tab"
        :class="{ 'provider-tab--active': provider === 'aws' }"
        @click="provider = 'aws'"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <ellipse cx="12" cy="5" rx="9" ry="3"/>
          <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
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
      </div>

      <StatusNotice :message="error"  type="error"   />
      <StatusNotice :message="notice" type="success" />

      <div class="form-actions">
        <BaseButton type="button" variant="ghost" :loading="testing" @click="handleTest">
          {{ testing ? 'Testing…' : 'Test connection' }}
        </BaseButton>
        <BaseButton type="submit" variant="primary" :loading="saving">
          {{ saving ? 'Saving…' : 'Save' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import BaseInput   from '../ui/BaseInput.vue'
import BaseButton  from '../ui/BaseButton.vue'
import StatusNotice from '../ui/StatusNotice.vue'

defineProps({
  testing: { type: Boolean, default: false },
  saving:  { type: Boolean, default: false },
  error:   { type: String,  default: '' },
  notice:  { type: String,  default: '' },
})

const emit = defineEmits(['test', 'save'])

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
</script>
