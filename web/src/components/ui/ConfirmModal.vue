<template>
  <teleport to="body">
    <transition name="modal-fade">
      <div v-if="pending" class="modal-backdrop" @mousedown.self="respond(false)">
        <div class="modal" role="dialog" :aria-label="pending.title">
          <div class="modal-hd">
            <span class="modal-title">{{ pending.title }}</span>
          </div>
          <div class="modal-bd">
            <p class="modal-msg">{{ pending.message }}</p>
          </div>
          <div class="modal-ft">
            <button class="base-btn base-btn--ghost" @click="respond(false)">Cancel</button>
            <button class="base-btn base-btn--danger" @click="respond(true)" style="border:1px solid var(--danger)">
              Confirm
            </button>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { useConfirm } from '../../composables/useConfirm.js'
const { pending, respond } = useConfirm()
</script>
