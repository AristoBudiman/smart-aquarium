<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'
import { Thermometer, Waves } from 'lucide-vue-next'
import StatusBadge from './StatusBadge.vue'
import { formatNumber } from '../utils/format'
import type { StatusClass } from '../types'

const props = defineProps<{
  eyebrow: string
  title: string
  value: number | null
  unit: string
  max: number
  minLabel: string
  rangeLabel: string
  maxLabel: string
  classLabel: string
  statusClass: StatusClass | null
  variant: 'temp' | 'water'
}>()

const animatedValue = ref<number | null>(props.value)
let rafId: number | null = null

function prefersReducedMotion(): boolean {
  return typeof window !== 'undefined' && window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

watch(
  () => props.value,
  (target, previous) => {
    if (rafId !== null) {
      cancelAnimationFrame(rafId)
      rafId = null
    }
    if (target === null || previous === null || previous === undefined || prefersReducedMotion()) {
      animatedValue.value = target
      return
    }
    const start = animatedValue.value ?? previous
    const diff = target - start
    if (Math.abs(diff) < 0.05) {
      animatedValue.value = target
      return
    }
    const duration = 220
    const startTime = performance.now()
    const step = (now: number) => {
      const progress = Math.min(1, (now - startTime) / duration)
      const eased = 1 - Math.pow(1 - progress, 3)
      animatedValue.value = start + diff * eased
      if (progress < 1) {
        rafId = requestAnimationFrame(step)
      } else {
        animatedValue.value = target
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

const progressRatio = computed(() => {
  if (props.value === null || !Number.isFinite(props.value)) return 0
  return Math.min(1, Math.max(0, props.value / props.max))
})

const radius = 64
const circumference = 2 * Math.PI * radius
const arcFraction = 0.78
const arcLength = circumference * arcFraction
const dashOffset = computed(() => arcLength * (1 - progressRatio.value))
</script>

<template>
  <article class="panel gauge-panel">
    <div class="panel-heading">
      <div>
        <p class="eyebrow">{{ eyebrow }}</p>
        <h2>{{ title }}</h2>
      </div>
      <span class="sensor-icon" :class="variant === 'temp' ? 'warm' : 'cool'" aria-hidden="true">
        <Thermometer v-if="variant === 'temp'" :size="20" :stroke-width="2.1" />
        <Waves v-else :size="20" :stroke-width="2.1" />
      </span>
    </div>

    <div class="gauge-wrap">
      <div class="gauge-svg-box">
        <svg class="gauge-svg" viewBox="0 0 160 160" aria-hidden="true">
          <circle
            class="gauge-track"
            cx="80"
            cy="80"
            :r="radius"
            fill="none"
            stroke-width="12"
            stroke-linecap="round"
            :stroke-dasharray="`${arcLength} ${circumference}`"
            transform="rotate(130 80 80)"
          />
          <circle
            class="gauge-bar"
            :class="variant === 'temp' ? 'gauge-bar-temp' : 'gauge-bar-water'"
            cx="80"
            cy="80"
            :r="radius"
            fill="none"
            stroke-width="12"
            stroke-linecap="round"
            :stroke-dasharray="`${arcLength} ${circumference}`"
            :stroke-dashoffset="dashOffset"
            transform="rotate(130 80 80)"
          />
        </svg>
        <div class="gauge-center">
          <strong>{{ formatNumber(animatedValue) }}</strong>
          <span>{{ unit }}</span>
        </div>
      </div>
    </div>

    <div class="gauge-scale">
      <span>{{ minLabel }}</span>
      <span>{{ rangeLabel }}</span>
      <span>{{ maxLabel }}</span>
    </div>

    <div class="class-row">
      <span>{{ classLabel }}</span>
      <StatusBadge :status="statusClass" />
    </div>
  </article>
</template>
