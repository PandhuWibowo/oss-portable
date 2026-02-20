import { ref } from 'vue'

// Module-level singleton so any component can push toasts.
const toasts = ref([])
let nextId = 1

export function useToast() {
  function push(message, type = 'info', duration = 4000) {
    const id = nextId++
    toasts.value.push({ id, message, type })
    if (duration > 0) {
      setTimeout(() => remove(id), duration)
    }
    return id
  }

  function remove(id) {
    const idx = toasts.value.findIndex(t => t.id === id)
    if (idx !== -1) toasts.value.splice(idx, 1)
  }

  const success = (msg, d) => push(msg, 'success', d)
  const error   = (msg, d) => push(msg, 'error',   d ?? 6000)
  const info    = (msg, d) => push(msg, 'info',     d)

  return { toasts, push, remove, success, error, info }
}
