<template>
  <div class="skeleton-wrap">
    <div
      v-for="i in count"
      :key="i"
      class="skeleton-item"
      :style="{ height: height }"
    ></div>
  </div>
</template>

<script setup>
defineProps({
  count: { type: Number, default: 3 },
  height: { type: String, default: '72px' },
})
</script>
