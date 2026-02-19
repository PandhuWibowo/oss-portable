<template>
  <aside class="sidebar">
    <!-- Brand -->
    <div class="sidebar__brand">
      <div class="brand-icon">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 15a4 4 0 0 0 4 4h9a5 5 0 0 0 1.8-9.7 6 6 0 0 0-11.8-1A4 4 0 0 0 3 15z"/>
        </svg>
      </div>
      <div>
        <div class="brand-name">Anvesa Vestra</div>
        <div class="brand-sub">Cloud storage manager</div>
      </div>
    </div>

    <!-- Body -->
    <div class="sidebar__body">
      <!-- New connection -->
      <button class="btn-new-conn" @click="$emit('new-connection')">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        New Connection
      </button>

      <!-- Search -->
      <div v-if="connections.length > 0" class="sidebar-search">
        <svg class="sidebar-search__icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input
          class="sidebar-search__input"
          v-model="query"
          placeholder="Filter connections…"
        />
      </div>

      <!-- Skeleton while loading -->
      <SkeletonLoader v-if="loading" :count="3" height="42px" />

      <template v-else>
        <!-- Section label -->
        <div v-if="filtered.length" class="section-label">Connections</div>

        <!-- Empty -->
        <div v-if="!filtered.length && !connections.length" class="sidebar-empty">
          No connections yet.<br>Add your first bucket above.
        </div>

        <!-- Items -->
        <div
          v-for="c in filtered"
          :key="c.provider + '-' + c.id"
          class="conn-item"
          :class="{ 'is-active': activeConn?.id === c.id && activeConn?.provider === c.provider }"
          @click="$emit('select', c)"
        >
          <span class="conn-dot" :class="`conn-dot--${c.provider}`"></span>
          <div class="conn-item__body">
            <div class="conn-item__name">{{ c.name }}</div>
            <div class="conn-item__bucket">{{ c.bucket }}</div>
          </div>
          <!-- Edit -->
          <button
            class="conn-item__del"
            @click.stop="$emit('edit', c)"
            title="Edit connection"
          >
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
          </button>
          <!-- Delete -->
          <button
            class="conn-item__del"
            @click.stop="$emit('delete', c.provider, c.id)"
            title="Delete connection"
          >
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
      </template>
    </div>

    <!-- Theme toggle -->
    <div class="sidebar__bottom">
      <button class="theme-btn" @click="toggleTheme">
        <svg v-if="isLight" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="5"/>
          <line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
          <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
          <line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/>
          <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
        </svg>
        <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
        {{ isLight ? 'Light mode' : 'Dark mode' }}
      </button>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue'
import SkeletonLoader from '../ui/SkeletonLoader.vue'
import { useTheme } from '../../composables/useTheme.js'

const props = defineProps({
  connections: { type: Array, default: () => [] },
  loading:     { type: Boolean, default: false },
  activeConn:  { type: Object, default: null },
})

defineEmits(['new-connection', 'select', 'delete'])

const { isLight, toggleTheme } = useTheme()

const query = ref('')

const filtered = computed(() => {
  const q = query.value.toLowerCase().trim()
  if (!q) return props.connections
  return props.connections.filter(c =>
    c.name.toLowerCase().includes(q) || c.bucket.toLowerCase().includes(q)
  )
})
</script>
