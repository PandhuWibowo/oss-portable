<template>
  <header class="header">
    <div>
      <div class="title">OSS Portable</div>
      <div class="subtitle">Clean, minimal cloud bucket dashboard</div>
    </div>
    <nav style="display:flex;gap:12px;align-items:center">
      <a href="/" style="color:var(--muted);text-decoration:none;font-size:13px">API</a>
      <button @click="toggleTheme" :aria-pressed="isLight" style="background:transparent;border:0;color:var(--muted);cursor:pointer">{{ isLight ? '🌞 Light' : '🌙 Dark' }}</button>
    </nav>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const isLight = ref(false)

function applyTheme(light) {
  if (light) {
    document.documentElement.setAttribute('data-theme', 'light')
  } else {
    document.documentElement.removeAttribute('data-theme')
  }
}

function toggleTheme() {
  isLight.value = !isLight.value
  applyTheme(isLight.value)
  try { localStorage.setItem('theme', isLight.value ? 'light' : 'dark') } catch {}
}

onMounted(() => {
  try {
    const saved = localStorage.getItem('theme')
    isLight.value = saved === 'light'
  } catch {
    isLight.value = false
  }
  applyTheme(isLight.value)
})
</script>
