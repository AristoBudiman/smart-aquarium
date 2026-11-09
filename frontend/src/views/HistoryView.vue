<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Thermometer, Waves } from 'lucide-vue-next'
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
import PageHeader from '../components/PageHeader.vue'
import { getNodeTelemetry } from '../api/nodes'
import { getApiErrorMessage } from '../api/client'
import { useNodesStore } from '../stores/nodes'
import { useThemeStore } from '../stores/theme'
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
const theme = useThemeStore()
const nodeId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))
const range = ref<RangeKey>('24h')
const points = ref<TelemetryPoint[]>([])
const loading = ref(true)
const error = ref('')
const nodeName = computed(() => store.details[nodeId.value]?.name ?? nodeId.value)

async function loadTelemetry() {
  loading.value = true
  error.value = ''
  try {
    if (!store.details[nodeId.value]) {
      await store.loadNode(nodeId.value).catch(() => undefined)
    }
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

watch(
  [nodeId, range],
  () => {
    void loadTelemetry()
  },
  { immediate: true },
)

function readCssVar(name: string): string {
  if (typeof document === 'undefined') return ''
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

const chartOption = computed<EChartsOption>(() => {
  // Baca ulang variabel CSS saat tema terang/gelap berubah
  void theme.isDark
  const seriesTemp = readCssVar('--chart-series-temp')
  const seriesNtu = readCssVar('--chart-series-ntu')
  const axisLabelColor = readCssVar('--chart-label')
  const axisLineColor = readCssVar('--chart-axis')
  const splitLineColor = readCssVar('--chart-grid')
  const tooltipBg = readCssVar('--chart-tooltip-bg')
  const tooltipBorder = readCssVar('--chart-tooltip-border')
  const tooltipText = readCssVar('--chart-tooltip-text')

  return {
    backgroundColor: 'transparent',
    color: [seriesTemp, seriesNtu],
    tooltip: {
      trigger: 'axis',
      backgroundColor: tooltipBg,
      borderColor: tooltipBorder,
      textStyle: { color: tooltipText, fontSize: 13 },
    },
    legend: {
      top: 0,
      right: 4,
      itemWidth: 12,
      itemHeight: 12,
      textStyle: { color: axisLabelColor, fontSize: 13 },
    },
    grid: { left: 48, right: 48, top: 44, bottom: 44 },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: points.value.map((point) =>
        new Intl.DateTimeFormat('id-ID', {
          timeZone: 'Asia/Jakarta',
          hour: range.value === '1h' || range.value === '24h' ? '2-digit' : undefined,
          minute: range.value === '1h' ? '2-digit' : undefined,
          day: range.value === '1h' || range.value === '24h' ? undefined : '2-digit',
          month: range.value === '1h' || range.value === '24h' ? undefined : 'short',
        }).format(new Date(point.ts)),
      ),
      axisLine: { lineStyle: { color: axisLineColor } },
      axisLabel: { color: axisLabelColor, fontSize: 13, hideOverlap: true },
    },
    yAxis: [
      {
        type: 'value',
        name: '°C',
        min: 0,
        max: 40,
        nameTextStyle: { color: axisLabelColor, fontSize: 13 },
        axisLabel: { color: axisLabelColor, fontSize: 13 },
        splitLine: { lineStyle: { color: splitLineColor } },
      },
      {
        type: 'value',
        name: 'NTU',
        min: 0,
        max: 200,
        nameTextStyle: { color: axisLabelColor, fontSize: 13 },
        axisLabel: { color: axisLabelColor, fontSize: 13 },
        splitLine: { show: false },
      },
    ],
    dataZoom: [{ type: 'inside', zoomOnMouseWheel: false }],
    series: [
      {
        name: 'Suhu (°C)',
        type: 'line',
        smooth: true,
        showSymbol: false,
        lineStyle: { width: 3 },
        areaStyle: { opacity: theme.isDark ? 0.16 : 0.09 },
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
  }
})
</script>

<template>
  <section class="page">
    <PageHeader
      eyebrow="ANALISIS DATA"
      title="Grafik riwayat"
      :subtitle="`Perubahan suhu dan kekeruhan · ${nodeName}`"
    >
      <template #breadcrumb>
        <RouterLink to="/nodes">Daftar akuarium</RouterLink>
        <span>/</span>
        <RouterLink :to="{ name: 'node-detail', params: { id: nodeId } }">{{ nodeName }}</RouterLink>
        <span>/</span>
        <span>Grafik riwayat</span>
      </template>
      <template #actions>
        <div class="range-tabs" role="group" aria-label="Rentang waktu grafik">
          <button
            v-for="item in ranges"
            :key="item.value"
            :class="{ selected: range === item.value }"
            type="button"
            @click="range = item.value"
          >
            {{ item.label }}
          </button>
        </div>
      </template>
    </PageHeader>

    <PageFeedback :loading="loading" :error="error">
      <article class="panel chart-panel">
        <div class="chart-summary">
          <div class="chart-stat">
            <span class="chart-stat-icon temp-stat" aria-hidden="true">
              <Thermometer :size="18" />
            </span>
            <div>
              <span>Suhu terbaru</span>
              <strong>{{ formatNumber(points.at(-1)?.temp) }} °C</strong>
            </div>
          </div>
          <div class="chart-stat">
            <span class="chart-stat-icon ntu-stat" aria-hidden="true">
              <Waves :size="18" />
            </span>
            <div>
              <span>Kekeruhan terbaru</span>
              <strong>{{ formatNumber(points.at(-1)?.ntu) }} NTU</strong>
            </div>
          </div>
          <span v-if="points.length" class="chart-updated">
            Diperbarui {{ formatWib(points.at(-1)?.ts) }} WIB
          </span>
        </div>

        <div v-if="!points.length" class="empty-inline">
          Belum ada data telemetry pada rentang waktu ini.
        </div>
        <VChart v-else class="history-chart" :option="chartOption" autoresize />
      </article>
    </PageFeedback>
  </section>
</template>
