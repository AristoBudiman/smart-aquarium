<script setup lang="ts">
defineProps<{
  loading?: boolean
  error?: string
}>()
</script>

<template>
  <div v-if="loading" class="feedback-card" role="status">
    <span class="loader"></span>
    <span>Memuat data...</span>
  </div>
  <div v-else-if="error" class="feedback-card feedback-error" role="alert">
    <strong>Data belum dapat ditampilkan</strong>
    <span>{{ error }}</span>
  </div>
  <slot v-else />
</template>
