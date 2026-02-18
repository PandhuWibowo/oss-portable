<template>
  <div style="font-family: system-ui, Arial, sans-serif; padding: 2rem;">
    <h1>OSS Portable</h1>
    <p>{{ message }}</p>
    <button @click="callApi">Call backend</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const message = ref('Frontend ready')

async function callApi() {
  try {
    const res = await fetch('/api/hello')
    const data = await res.json()
    message.value = data.message || JSON.stringify(data)
  } catch (err) {
    message.value = 'Error: ' + err.message
  }
}
</script>
