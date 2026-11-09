<script setup lang="ts">
import { computed } from 'vue'
import {
  AlertOctagon,
  AlertTriangle,
  CheckCircle2,
  Clock,
  Info,
  WifiOff,
} from 'lucide-vue-next'
import type { StatusClass } from '../types'

const props = defineProps<{
  status: StatusClass | null
  offline?: boolean
  size?: 'md' | 'lg'
}>()

const badgeClass = computed(() => {
  if (props.offline) return 'status-offline'
  return `status-${(props.status ?? 'unknown').toLowerCase()}`
})

const badgeText = computed(() => {
  if (props.offline) return 'OFFLINE'
  return props.status ?? 'BELUM ADA DATA'
})

const iconComponent = computed(() => {
  if (props.offline) return WifiOff
  switch (props.status) {
    case 'NORMAL':
      return CheckCircle2
    case 'CAUTION':
      return Info
    case 'WARNING':
      return AlertTriangle
    case 'SEVERE':
      return AlertOctagon
    default:
      return Clock
  }
})
</script>

<template>
  <span class="status-badge" :class="[badgeClass, size === 'lg' ? 'status-badge-lg' : '']">
    <component :is="iconComponent" :size="size === 'lg' ? 18 : 15" :stroke-width="2.2" aria-hidden="true" />
    <span>{{ badgeText }}</span>
  </span>
</template>
