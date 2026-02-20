import { ref } from 'vue'

// Module-level singleton — renders one confirm dialog at a time.
const pending = ref(null) // { title, message, resolve }

export function useConfirm() {
  function confirm(message, title = 'Are you sure?') {
    return new Promise(resolve => {
      pending.value = { title, message, resolve }
    })
  }

  function respond(result) {
    if (pending.value) {
      pending.value.resolve(result)
      pending.value = null
    }
  }

  return { pending, confirm, respond }
}
