<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import PageFeedback from '../components/PageFeedback.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { getNodeLogs } from '../api/nodes'
import { getApiErrorMessage } from '../api/client'
import { useNodesStore } from '../stores/nodes'
import { formatWib } from '../utils/format'
import type { NodeLogs } from '../types'

const route = useRoute()
const store = useNodesStore()
const nodeId = computed(() => typeof route.params.id === 'string' ? route.params.id : '')
const nodeName = computed(() => store.details[nodeId.value]?.name ?? nodeId.value)
const logs = ref<NodeLogs>({ actuators: [], alerts: [] })
const loading = ref(true)
const error = ref('')
const alerts = computed(() => [
  ...(store.realtimeAlerts[nodeId.value] ?? []),
  ...logs.value.alerts,
])

watch(nodeId, async (id) => {
  loading.value = true
  error.value = ''
  try {
    logs.value = await getNodeLogs(id)
  } catch (cause) {
    error.value = getApiErrorMessage(cause, 'Log dan notifikasi tidak dapat dimuat.')
  } finally {
    loading.value = false
  }
}, { immediate: true })
</script>

<template>
  <section class="page">
    <div class="breadcrumb"><RouterLink :to="{ name: 'node-detail', params: { id: nodeId } }">{{ nodeName }}</RouterLink><span>/</span><span>Log &amp; notifikasi</span></div>
    <div class="page-heading">
      <p class="eyebrow">AKTIVITAS NODE · {{ nodeId }}</p>
      <h1>Log &amp; notifikasi</h1>
      <p class="muted">Riwayat perubahan aktuator dan peringatan sistem.</p>
    </div>

    <PageFeedback :loading="loading" :error="error">
      <div class="logs-grid">
        <article class="panel log-panel">
          <div class="panel-heading"><div><p class="eyebrow">RIWAYAT PERANGKAT</p><h2>Aktivitas aktuator</h2></div><span class="count-pill">{{ logs.actuators.length }}</span></div>
          <div v-if="!logs.actuators.length" class="empty-inline">Belum ada riwayat aktuator.</div>
          <div v-else class="timeline">
            <div v-for="entry in logs.actuators" :key="entry.id" class="timeline-entry">
              <span class="timeline-marker actuator-marker">↻</span>
              <div class="timeline-content">
                <div class="timeline-title"><strong>Perubahan aktuator</strong><span class="source-pill">Mode {{ entry.source }}</span></div>
                <p>Pemanas {{ entry.heater ? 'aktif' : 'nonaktif' }} · Kipas {{ entry.fan ? 'aktif' : 'nonaktif' }} · Pompa {{ entry.pump ? 'aktif' : 'nonaktif' }}</p>
                <time>{{ formatWib(entry.ts) }} WIB</time>
              </div>
            </div>
          </div>
        </article>

        <article class="panel log-panel">
          <div class="panel-heading"><div><p class="eyebrow">PUSAT PERINGATAN</p><h2>Notifikasi</h2></div><span class="count-pill">{{ alerts.length }}</span></div>
          <div v-if="!alerts.length" class="empty-inline">Tidak ada notifikasi untuk ditampilkan.</div>
          <div v-else class="timeline">
            <div v-for="alert in alerts" :key="alert.id" class="timeline-entry">
              <span class="timeline-marker alert-marker">!</span>
              <div class="timeline-content">
                <div class="timeline-title"><strong>{{ alert.message }}</strong><StatusBadge :status="alert.level" /></div>
                <time>{{ formatWib(alert.ts) }} WIB</time>
              </div>
            </div>
          </div>
        </article>
      </div>
    </PageFeedback>
  </section>
</template>
