<template>
  <button
    :class="['base-btn', `base-btn--${variant}`, { 'base-btn--loading': loading }]"
    :disabled="disabled || loading"
    v-bind="$attrs"
  >
    <span v-if="loading" class="base-btn__spinner" aria-hidden="true"></span>
    <slot />
  </button>
</template>

<script setup>
defineProps({
  variant: { type: String, default: 'primary' }, // primary | ghost | danger | soft
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
})
</script>
