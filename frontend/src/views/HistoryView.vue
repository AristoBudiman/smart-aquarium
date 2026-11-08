<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import {
  DataZoomComponent,
  GridComponent,
  LegendComponent,
  TooltipComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'
import type { EChartsOption } from 'echarts'
import PageFeedback from '../components/PageFeedback.vue'
import { getNodeTelemetry } from '../api/nodes'
import { getApiErrorMessage } from '../api/client'
import { useNodesStore } from '../stores/nodes'
import { formatNumber, formatWib } from '../utils/format'
import type { TelemetryPoint } from '../types'

use([LineChart, GridComponent, LegendComponent, TooltipComponent, DataZoomComponent, CanvasRenderer])

type RangeKey = '1h' | '24h' | '7d' | '30d'
const ranges: { value: RangeKey; label: string; milliseconds: number }[] = [
  { value: '1h', label: '1 jam', milliseconds: 60 * 60_000 },
  { value: '24h', label: '24 jam', milliseconds: 24 * 60 * 60_000 },
  { value: '7d', label: '7 hari', milliseconds: 7 * 86400_000 },
  { value: '30d', label: '30 hari', milliseconds: 30 * 86400_000 },
]

const route = useRoute()
const store = useNodesStore()
const nodeId = computed(() => typeof route.params.id === 'string' ? route.params.id : '')
const range = ref<RangeKey>('24h')
const points = ref<TelemetryPoint[]>([])
const loading = ref(true)
const error = ref('')
const nodeName = computed(() => store.details[nodeId.value]?.name ?? nodeId.value)

async function loadTelemetry() {
  loading.value = true
  error.value = ''
  try {
    const selected = ranges.find((item) => item.value === range.value) ?? ranges[1]!
    const to = new Date()
    const from = new Date(to.getTime() - selected.milliseconds)
    points.value = await getNodeTelemetry(nodeId.value, from.toISOString(), to.toISOString())
  } catch (cause) {
    error.value = getApiErrorMessage(cause, 'Grafik riwayat tidak dapat dimuat.')
  } finally {
    loading.value = false
  }
}

watch([nodeId, range], () => { void loadTelemetry() }, { immediate: true })

const chartOption = computed<EChartsOption>(() => ({
  color: ['#0f9f8f', '#f4a340'],
  tooltip: { trigger: 'axis' },
  legend: { top: 0, right: 4, itemWidth: 10, itemHeight: 10, textStyle: { color: '#667085' } },
  grid: { left: 48, right: 48, top: 42, bottom: 42 },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: points.value.map((point) => new Intl.DateTimeFormat('id-ID', {
      timeZone: 'Asia/Jakarta',
      hour: range.value === '1h' || range.value === '24h' ? '2-digit' : undefined,
      minute: range.value === '1h' ? '2-digit' : undefined,
      day: range.value === '1h' || range.value === '24h' ? undefined : '2-digit',
      month: range.value === '1h' || range.value === '24h' ? undefined : 'short',
    }).format(new Date(point.ts))),
    axisLine: { lineStyle: { color: '#e6ebf0' } },
    axisLabel: { color: '#84909d', hideOverlap: true },
  },
  yAxis: [
    { type: 'value', name: '°C', min: 0, max: 40, splitLine: { lineStyle: { color: '#f0f3f6' } } },
    { type: 'value', name: 'NTU', min: 0, max: 200, splitLine: { show: false } },
  ],
  dataZoom: [{ type: 'inside', zoomOnMouseWheel: false }],
  series: [
    {
      name: 'Suhu (°C)',
      type: 'line',
      smooth: true,
      showSymbol: false,
      lineStyle: { width: 3 },
      areaStyle: { opacity: 0.08 },
      data: points.value.map((point) => point.temp),
    },
    {
      name: 'Kekeruhan (NTU)',
      type: 'line',
      yAxisIndex: 1,
      smooth: true,
      showSymbol: false,
      lineStyle: { width: 3 },
      data: points.value.map((point) => point.ntu),
    },
  ],
}))
</script>

<template>
  <section class="page">
    <div class="breadcrumb"><RouterLink :to="{ name: 'node-detail', params: { id: nodeId } }">{{ nodeName }}</RouterLink><span>/</span><span>Grafik riwayat</span></div>
    <div class="page-heading heading-with-action">
      <div>
        <p class="eyebrow">ANALISIS DATA</p>
        <h1>Grafik riwayat</h1>
        <p class="muted">Perubahan suhu dan kekeruhan · {{ nodeName }}</p>
      </div>
      <div class="range-tabs" role="group" aria-label="Rentang waktu grafik">
        <button v-for="item in ranges" :key="item.value" :class="{ selected: range === item.value }" type="button" @click="range = item.value">
          {{ item.label }}
        </button>
      </div>
    </div>

    <PageFeedback :loading="loading" :error="error">
      <article class="panel chart-panel">
        <div class="chart-summary">
          <div><span class="chart-dot temp-dot"></span><span>Suhu terbaru</span><strong>{{ formatNumber(points.at(-1)?.temp) }} °C</strong></div>
          <div><span class="chart-dot ntu-dot"></span><span>Kekeruhan terbaru</span><strong>{{ formatNumber(points.at(-1)?.ntu) }} NTU</strong></div>
          <span v-if="points.length" class="chart-updated">Diperbarui {{ formatWib(points.at(-1)?.ts) }} WIB</span>
        </div>
        <div v-if="!points.length" class="empty-inline">Belum ada data telemetry pada rentang waktu ini.</div>
        <VChart v-else class="history-chart" :option="chartOption" autoresize />
      </article>
    </PageFeedback>
  </section>
</template>
