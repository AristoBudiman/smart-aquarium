<script setup lang="ts">
import type { StatusClass } from '../types'

defineProps<{
  status: StatusClass | null
  offline?: boolean
}>()
</script>

<template>
  <span
    class="status-badge"
    :class="offline ? 'status-offline' : `status-${(status ?? 'unknown').toLowerCase()}`"
  >
    <i></i>{{ offline ? 'OFFLINE' : status ?? 'BELUM ADA DATA' }}
  </span>
</template>
