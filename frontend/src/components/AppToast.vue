<script setup lang="ts">
import { AlertTriangle, CheckCircle2, Info, X } from 'lucide-vue-next'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
</script>

<template>
  <Teleport to="body">
    <div class="toast-region" aria-live="polite" aria-atomic="true">
      <TransitionGroup name="toast">
        <div
          v-for="item in toast.items"
          :key="item.id"
          class="toast-card"
          :class="`toast-${item.kind}`"
          role="status"
        >
          <CheckCircle2 v-if="item.kind === 'success'" :size="19" class="toast-icon" aria-hidden="true" />
          <AlertTriangle v-else-if="item.kind === 'error'" :size="19" class="toast-icon" aria-hidden="true" />
          <Info v-else :size="19" class="toast-icon" aria-hidden="true" />
          <span class="toast-text">{{ item.message }}</span>
          <button
            class="toast-close"
            type="button"
            aria-label="Tutup pemberitahuan"
            @click="toast.dismiss(item.id)"
          >
            <X :size="16" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
