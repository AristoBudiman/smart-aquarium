<script setup lang="ts">
import { AlertTriangle, ShieldAlert, X } from 'lucide-vue-next'

defineProps<{
  open: boolean
  title: string
  description: string
  confirmLabel?: string
  cancelLabel?: string
  danger?: boolean
  loading?: boolean
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <Teleport to="body">
    <Transition name="dialog-fade">
      <div v-if="open" class="modal-backdrop" role="presentation" @click.self="emit('cancel')">
        <div
          class="modal-card"
          role="alertdialog"
          aria-modal="true"
          :aria-labelledby="'confirm-dialog-title'"
        >
          <div class="modal-header">
            <div class="modal-icon" :class="danger ? 'modal-icon-danger' : 'modal-icon-warn'" aria-hidden="true">
              <ShieldAlert v-if="danger" :size="22" />
              <AlertTriangle v-else :size="22" />
            </div>
            <button
              class="icon-button"
              type="button"
              aria-label="Tutup dialog"
              :disabled="loading"
              @click="emit('cancel')"
            >
              <X :size="19" />
            </button>
          </div>
          <h2 id="confirm-dialog-title" class="modal-title">{{ title }}</h2>
          <p class="modal-description">{{ description }}</p>
          <div class="modal-actions">
            <button
              class="button button-quiet"
              type="button"
              :disabled="loading"
              @click="emit('cancel')"
            >
              {{ cancelLabel ?? 'Batal' }}
            </button>
            <button
              class="button"
              :class="danger ? 'button-danger' : 'button-primary'"
              type="button"
              :disabled="loading"
              @click="emit('confirm')"
            >
              {{ loading ? 'Memproses...' : (confirmLabel ?? 'Konfirmasi') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
