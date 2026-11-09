<script setup lang="ts">
import { AlertTriangle } from 'lucide-vue-next'

defineProps<{
  loading?: boolean
  error?: string
}>()
</script>

<template>
  <div v-if="loading" class="skeleton-grid" role="status" aria-label="Memuat data">
    <div v-for="n in 3" :key="n" class="skeleton-card">
      <div class="skeleton-line skeleton-badge"></div>
      <div class="skeleton-line skeleton-title"></div>
      <div class="skeleton-line skeleton-sub"></div>
      <div class="skeleton-metrics">
        <div class="skeleton-block"></div>
        <div class="skeleton-block"></div>
      </div>
      <div class="skeleton-line skeleton-footer"></div>
    </div>
  </div>
  <div v-else-if="error" class="feedback-card feedback-error" role="alert">
    <div class="feedback-error-head">
      <AlertTriangle :size="20" aria-hidden="true" />
      <strong>Data belum dapat ditampilkan</strong>
    </div>
    <span>{{ error }}</span>
  </div>
  <slot v-else />
</template>
