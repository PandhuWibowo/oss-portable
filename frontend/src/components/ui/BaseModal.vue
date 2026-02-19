<template>
  <teleport to="body">
    <transition name="modal-fade">
      <div v-if="open" class="modal-backdrop" @mousedown.self="$emit('update:open', false)">
        <div class="modal" role="dialog" :aria-label="title" @keydown.escape.stop="$emit('update:open', false)">
          <div class="modal-hd">
            <span class="modal-title">{{ title }}</span>
            <button class="modal-close" @click="$emit('update:open', false)" aria-label="Close">×</button>
          </div>
          <div class="modal-bd">
            <slot />
          </div>
          <div v-if="$slots.footer" class="modal-ft">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
defineProps({
  open:  { type: Boolean, required: true },
  title: { type: String,  default: '' },
})
defineEmits(['update:open'])
</script>
