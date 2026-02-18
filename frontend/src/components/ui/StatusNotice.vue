<template>
  <transition name="notice-fade">
    <div v-if="message" :class="['status-notice', `status-notice--${type}`]" role="alert">
      <span class="status-notice__icon">{{ icons[type] }}</span>
      <span class="status-notice__text">{{ message }}</span>
    </div>
  </transition>
</template>

<script setup>
defineProps({
  message: { type: String, default: '' },
  type: { type: String, default: 'info' }, // 'info' | 'success' | 'error'
})

const icons = { info: 'ℹ️', success: '✓', error: '✕' }
</script>
