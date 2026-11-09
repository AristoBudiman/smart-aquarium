<script setup lang="ts">
import { onUnmounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowRight, Thermometer, Waves, Wifi, WifiOff } from 'lucide-vue-next'
import StatusBadge from './StatusBadge.vue'
import { formatNumber, formatWib } from '../utils/format'
import type { NodeSummary } from '../types'

const props = defineProps<{
  node: NodeSummary
  index?: number
}>()

const displayTemp = ref<number | null>(props.node.temp)
const displayNtu = ref<number | null>(props.node.ntu)
let rafId: number | null = null

function prefersReducedMotion(): boolean {
  return typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

watch(
  () => [props.node.temp, props.node.ntu] as const,
  ([nextTemp, nextNtu], prev) => {
    if (rafId !== null) {
      cancelAnimationFrame(rafId)
      rafId = null
    }
    if (!prev || prefersReducedMotion()) {
      displayTemp.value = nextTemp
      displayNtu.value = nextNtu
      return
    }
    const startTemp = displayTemp.value ?? nextTemp ?? 0
    const startNtu = displayNtu.value ?? nextNtu ?? 0
    const duration = 200
    const startTime = performance.now()
    const step = (now: number) => {
      const progress = Math.min(1, (now - startTime) / duration)
      const eased = 1 - Math.pow(1 - progress, 3)
      displayTemp.value = nextTemp === null ? null : startTemp + (nextTemp - startTemp) * eased
      displayNtu.value = nextNtu === null ? null : startNtu + (nextNtu - startNtu) * eased
      if (progress < 1) {
        rafId = requestAnimationFrame(step)
      } else {
        displayTemp.value = nextTemp
        displayNtu.value = nextNtu
        rafId = null
      }
    }
    rafId = requestAnimationFrame(step)
  },
  { immediate: true },
)

onUnmounted(() => {
  if (rafId !== null) cancelAnimationFrame(rafId)
})
</script>

<template>
  <article
    class="node-card stagger-item"
    :style="{ '--stagger-delay': `${Math.min(index ?? 0, 8) * 55}ms` }"
  >
    <div class="node-card-top">
      <div class="node-symbol" :class="{ 'node-symbol-offline': !node.online }" aria-hidden="true">
        <Wifi v-if="node.online" :size="20" :stroke-width="2.1" />
        <WifiOff v-else :size="20" :stroke-width="2.1" />
      </div>
      <div class="node-card-badges">
        <span class="online-pill" :class="node.online ? 'online-pill-active' : 'online-pill-idle'">
          <span class="status-dot" :class="node.online ? 'dot-online' : 'dot-offline'"></span>
          {{ node.online ? 'Online' : 'Offline' }}
        </span>
        <StatusBadge :status="node.status" :offline="!node.online" />
      </div>
    </div>

    <div class="node-card-title-group">
      <h2>{{ node.name }}</h2>
      <div class="node-meta-row">
        <span class="node-id">{{ node.node_id }}</span>
        <span class="mode-chip">{{ node.mode === 'auto' ? 'Otomatis' : 'Manual' }}</span>
      </div>
    </div>

    <div class="node-metrics">
      <div class="metric-box">
        <span class="metric-label">
          <Thermometer :size="15" aria-hidden="true" />
          Suhu air
        </span>
        <strong>
          {{ formatNumber(displayTemp) }}
          <small>°C</small>
        </strong>
      </div>
      <div class="metric-box">
        <span class="metric-label">
          <Waves :size="15" aria-hidden="true" />
          Kekeruhan
        </span>
        <strong>
          {{ formatNumber(displayNtu) }}
          <small>NTU</small>
        </strong>
      </div>
    </div>

    <div class="node-card-bottom">
      <span>{{ node.online ? 'Terakhir diperbarui' : 'Terakhir terlihat' }}</span>
      <time>{{ formatWib(node.last_seen) }}</time>
    </div>

    <RouterLink class="card-link" :to="{ name: 'node-detail', params: { id: node.node_id } }">
      <span>Lihat detail akuarium</span>
      <ArrowRight :size="18" aria-hidden="true" />
    </RouterLink>
  </article>
</template>
