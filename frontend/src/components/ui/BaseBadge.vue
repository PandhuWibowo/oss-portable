<template>
  <span :class="['base-badge', `base-badge--${provider}`]">
    {{ provider.toUpperCase() }}
  </span>
</template>

<script setup>
defineProps({
  provider: { type: String, required: true }, // 'gcp' | 'aws'
})
</script>
