<template>
  <div class="conn-card">
    <div class="conn-left">
      <div class="avatar">{{ providerAbbrev }}</div>
      <div>
        <div style="display:flex;gap:8px;align-items:center">
          <div class="name">{{ conn.name }}</div>
          <div class="badge">{{ conn.provider.toUpperCase() }}</div>
        </div>
        <div class="meta">{{ conn.bucket }}</div>
        <div class="small" v-if="showCreated">{{ formattedDate }}</div>
      </div>
    </div>

    <div class="actions">
      <button class="icon-btn" @click="copyBucket" title="Copy bucket name">📋</button>
      <button class="icon-btn" @click="toggleReveal" :title="revealed ? 'Hide credentials' : 'Show credentials'">{{ revealed ? '🔓' : '🔒' }}</button>
      <button class="btn btn-ghost" @click="$emit('use', conn)">Use</button>
      <button class="btn btn-ghost" @click="$emit('delete', conn.id)">Delete</button>
    </div>

    <transition name="fade">
      <div v-if="revealed" style="width:100%;margin-top:10px;grid-column:1/-1">
        <pre style="white-space:pre-wrap;background:rgba(0,0,0,0.03);padding:10px;border-radius:8px;font-family:var(--mono);font-size:12px">{{ conn.credentials }}</pre>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
const props = defineProps({ conn: Object, showCreated: { type: Boolean, default: false } })
const emits = defineEmits(['use','delete'])
const revealed = ref(false)
function toggleReveal(){ revealed.value = !revealed.value }
function copyBucket(){ navigator.clipboard?.writeText(props.conn.bucket) }
const providerAbbrev = computed(()=> props.conn.provider.slice(0,1).toUpperCase())
const formattedDate = computed(()=> new Date(props.conn.created_at).toLocaleString())
</script>
