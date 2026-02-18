<template>
  <section class="panel">
    <div class="panel__header">
      <div>
        <h2 class="panel__title">Connections</h2>
        <p class="panel__subtitle">Manage your saved cloud connections</p>
      </div>
      <div class="search-wrap">
        <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input class="search-input" v-model="query" placeholder="Search name or bucket…" />
      </div>
    </div>

    <SkeletonLoader v-if="loading" :count="3" height="80px" />

    <div v-else>
      <div v-if="filtered.length === 0" class="empty-state">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:0.3;margin-bottom:8px">
          <path d="M3 15a4 4 0 0 0 4 4h9a5 5 0 0 0 1.8-9.7 6 6 0 0 0-11.8-1A4 4 0 0 0 3 15z"/>
        </svg>
        <p>No connections yet — add one on the right.</p>
      </div>
      <div class="conn-list">
        <ConnectionCard
          v-for="c in filtered"
          :key="c.provider + '-' + c.id"
          :conn="c"
          @use="$emit('use', c)"
          @delete="(id) => $emit('delete', c.provider, id)"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed } from 'vue'
import ConnectionCard from './ConnectionCard.vue'
import SkeletonLoader from '../ui/SkeletonLoader.vue'

const props = defineProps({
  connections: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})
defineEmits(['use', 'delete'])

const query = ref('')
const filtered = computed(() => {
  const q = query.value.toLowerCase().trim()
  if (!q) return props.connections
  return props.connections.filter(c =>
    c.name.toLowerCase().includes(q) || c.bucket.toLowerCase().includes(q)
  )
})
</script>
