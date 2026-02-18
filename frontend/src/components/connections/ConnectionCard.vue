<template>
  <div class="conn-card" :class="{ 'conn-card--revealed': revealed }">
    <div class="conn-card__main">
      <div class="conn-card__avatar" :class="`conn-card__avatar--${conn.provider}`">
        {{ conn.provider === 'gcp' ? 'G' : 'A' }}
      </div>
      <div class="conn-card__info">
        <div class="conn-card__name-row">
          <span class="conn-card__name">{{ conn.name }}</span>
          <BaseBadge :provider="conn.provider" />
        </div>
        <div class="conn-card__bucket">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="opacity:0.5">
            <ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
          </svg>
          {{ conn.bucket }}
        </div>
        <div v-if="conn.created_at" class="conn-card__date">{{ formattedDate }}</div>
      </div>
    </div>

    <div class="conn-card__actions">
      <button class="icon-btn" @click="copyBucket" :title="copied ? 'Copied!' : 'Copy bucket name'">
        <svg v-if="!copied" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
        </svg>
        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="20 6 9 17 4 12"/>
        </svg>
      </button>
      <button class="icon-btn" @click="toggleReveal" :title="revealed ? 'Hide credentials' : 'Show credentials'">
        <svg v-if="!revealed" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>
        </svg>
        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
          <line x1="1" y1="1" x2="23" y2="23"/>
        </svg>
      </button>
      <BaseButton variant="ghost" @click="$emit('use', conn)">Use</BaseButton>
      <BaseButton variant="danger" @click="$emit('delete', conn.id)">Delete</BaseButton>
    </div>

    <transition name="slide-down">
      <div v-if="revealed" class="conn-card__credentials">
        <pre>{{ conn.credentials }}</pre>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import BaseBadge from '../ui/BaseBadge.vue'
import BaseButton from '../ui/BaseButton.vue'

const props = defineProps({ conn: { type: Object, required: true } })
defineEmits(['use', 'delete'])

const revealed = ref(false)
const copied = ref(false)

function toggleReveal() { revealed.value = !revealed.value }

async function copyBucket() {
  await navigator.clipboard?.writeText(props.conn.bucket)
  copied.value = true
  setTimeout(() => { copied.value = false }, 1500)
}

const formattedDate = computed(() => {
  try { return new Date(props.conn.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }) }
  catch { return '' }
})
</script>
