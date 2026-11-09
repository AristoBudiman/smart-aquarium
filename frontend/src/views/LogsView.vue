<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Bell, RefreshCw } from 'lucide-vue-next'
import PageFeedback from '../components/PageFeedback.vue'
import PageHeader from '../components/PageHeader.vue'
import Paginator from '../components/Paginator.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { getNodeLogs } from '../api/nodes'
import { getApiErrorMessage } from '../api/client'
import { useNodesStore } from '../stores/nodes'
import { formatWib } from '../utils/format'
import type { NodeLogs } from '../types'

const PAGE_SIZE = 10

const route = useRoute()
const store = useNodesStore()
const nodeId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))
const nodeName = computed(() => store.details[nodeId.value]?.name ?? nodeId.value)
const logs = ref<NodeLogs>({ actuators: [], alerts: [] })
const loading = ref(true)
const error = ref('')

const actuatorPage = ref(1)
const alertPage = ref(1)
let alertRefreshTimer: number | undefined

const actuatorTotalPages = computed(() =>
  Math.max(1, Math.ceil(logs.value.actuators.length / PAGE_SIZE)),
)
const alertTotalPages = computed(() =>
  Math.max(1, Math.ceil(logs.value.alerts.length / PAGE_SIZE)),
)

const pagedActuators = computed(() => {
  const start = (actuatorPage.value - 1) * PAGE_SIZE
  return logs.value.actuators.slice(start, start + PAGE_SIZE)
})

const pagedAlerts = computed(() => {
  const start = (alertPage.value - 1) * PAGE_SIZE
  return logs.value.alerts.slice(start, start + PAGE_SIZE)
})

async function fetchLogs(id: string, silent = false) {
  if (!id) return
  if (!silent) {
    loading.value = true
    error.value = ''
  }
  try {
    if (!store.details[id]) {
      await store.loadNode(id).catch(() => undefined)
    }
    logs.value = await getNodeLogs(id)
    if (actuatorPage.value > actuatorTotalPages.value) {
      actuatorPage.value = actuatorTotalPages.value
    }
    if (alertPage.value > alertTotalPages.value) {
      alertPage.value = alertTotalPages.value
    }
  } catch (cause) {
    if (!silent) {
      error.value = getApiErrorMessage(cause, 'Log dan notifikasi tidak dapat dimuat.')
    }
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

watch(
  nodeId,
  (id) => {
    actuatorPage.value = 1
    alertPage.value = 1
    void fetchLogs(id)
  },
  { immediate: true },
)

// Saat event WebSocket "alert" masuk, jangan menambahkan item secara lokal.
// Ambil ulang GET /api/nodes/{id}/logs dengan debounce ~500 ms dan gunakan id dari API sebagai key.
watch(
  () => store.alertVersion[nodeId.value],
  (nextVer, prevVer) => {
    if (!nextVer || nextVer === prevVer) return
    if (alertRefreshTimer !== undefined) {
      window.clearTimeout(alertRefreshTimer)
    }
    alertRefreshTimer = window.setTimeout(() => {
      void fetchLogs(nodeId.value, true)
    }, 500)
  },
)

onUnmounted(() => {
  if (alertRefreshTimer !== undefined) {
    window.clearTimeout(alertRefreshTimer)
  }
})
</script>

<template>
  <section class="page">
    <PageHeader
      :eyebrow="`AKTIVITAS NODE · ${nodeId}`"
      title="Log & notifikasi"
      subtitle="Riwayat perubahan aktuator dan peringatan kualitas air."
    >
      <template #breadcrumb>
        <RouterLink to="/nodes">Daftar akuarium</RouterLink>
        <span>/</span>
        <RouterLink :to="{ name: 'node-detail', params: { id: nodeId } }">{{ nodeName }}</RouterLink>
        <span>/</span>
        <span>Log &amp; notifikasi</span>
      </template>
    </PageHeader>

    <PageFeedback :loading="loading" :error="error">
      <div class="logs-grid">
        <article class="panel log-panel">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">RIWAYAT PERANGKAT</p>
              <h2>Aktivitas aktuator</h2>
            </div>
            <span class="count-pill" :aria-label="`Total ${logs.actuators.length} riwayat aktuator`">
              {{ logs.actuators.length }}
            </span>
          </div>

          <div v-if="!logs.actuators.length" class="empty-inline">
            Belum ada riwayat perubahan aktuator.
          </div>
          <template v-else>
            <div class="timeline">
              <div v-for="entry in pagedActuators" :key="entry.id" class="timeline-entry">
                <span class="timeline-marker actuator-marker" aria-hidden="true">
                  <RefreshCw :size="16" />
                </span>
                <div class="timeline-content">
                  <div class="timeline-title">
                    <strong>Perubahan status relay</strong>
                    <span class="source-pill">Mode {{ entry.source }}</span>
                  </div>
                  <p>
                    Pemanas {{ entry.heater ? 'ON' : 'OFF' }} · Kipas {{ entry.fan ? 'ON' : 'OFF' }} ·
                    Pompa {{ entry.pump ? 'ON' : 'OFF' }}
                  </p>
                  <time>{{ formatWib(entry.ts) }} WIB</time>
                </div>
              </div>
            </div>
            <Paginator
              v-model:page="actuatorPage"
              :total-pages="actuatorTotalPages"
              :total-items="logs.actuators.length"
            />
          </template>
        </article>

        <article class="panel log-panel">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">PUSAT PERINGATAN</p>
              <h2>Notifikasi</h2>
            </div>
            <span class="count-pill" :aria-label="`Total ${logs.alerts.length} notifikasi`">
              {{ logs.alerts.length }}
            </span>
          </div>

          <div v-if="!logs.alerts.length" class="empty-inline">
            Tidak ada notifikasi peringatan untuk ditampilkan.
          </div>
          <template v-else>
            <div class="timeline">
              <div v-for="alert in pagedAlerts" :key="alert.id" class="timeline-entry">
                <span class="timeline-marker alert-marker" aria-hidden="true">
                  <Bell :size="16" />
                </span>
                <div class="timeline-content">
                  <div class="timeline-title">
                    <strong>{{ alert.message }}</strong>
                    <StatusBadge :status="alert.level" />
                  </div>
                  <time>{{ formatWib(alert.ts) }} WIB</time>
                </div>
              </div>
            </div>
            <Paginator
              v-model:page="alertPage"
              :total-pages="alertTotalPages"
              :total-items="logs.alerts.length"
            />
          </template>
        </article>
      </div>
    </PageFeedback>
  </section>
</template>
