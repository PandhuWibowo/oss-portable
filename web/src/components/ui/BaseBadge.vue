<template>
  <span :class="['base-badge', `base-badge--${provider}`]">
    {{ LABELS[provider] ?? provider.toUpperCase() }}
  </span>
</template>

<script setup>
defineProps({
  provider: { type: String, required: true },
})

const LABELS = {
  gcp:     'GCS',
  aws:     'S3',
  huawei:  'OBS',
  alibaba: 'OSS',
  azure:   'Azure',
}
</script>
