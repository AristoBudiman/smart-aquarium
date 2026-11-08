<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import PageFeedback from '../components/PageFeedback.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { getApiErrorMessage, getApiErrorStatus } from '../api/client'
import { useNodesStore } from '../stores/nodes'
import { formatNumber, formatWib } from '../utils/format'
import type { ActuatorCommand } from '../types'

const route = useRoute()
const store = useNodesStore()
const nodeId = computed(() => typeof route.params.id === 'string' ? route.params.id : '')
const loading = ref(true)
const error = ref('')
const actionError = ref('')
const waitingForState = computed(() => store.pendingActuators[nodeId.value] === true)
const detail = computed(() => store.details[nodeId.value] ?? null)

watch(nodeId, async (id) => {
  loading.value = true
  error.value = ''
  try {
    await store.loadNode(id)
  } catch (cause) {
    error.value = getApiErrorMessage(cause, 'Detail node tidak dapat dimuat.')
  } finally {
    loading.value = false
  }
}, { immediate: true })

const temperatureProgress = computed(() => Math.min(100, Math.max(0, ((detail.value?.temp ?? 0) / 40) * 100)))
const turbidityProgress = computed(() => Math.min(100, Math.max(0, ((detail.value?.ntu ?? 0) / 200) * 100)))

async function toggleMode() {
  if (!detail.value) return
  actionError.value = ''
  try {
    await store.changeMode(detail.value.node_id, detail.value.mode === 'auto' ? 'manual' : 'auto')
  } catch (cause) {
    actionError.value = getApiErrorMessage(cause, 'Mode tidak dapat diubah.')
  }
}

async function toggleActuator(key: keyof ActuatorCommand) {
  if (!detail.value) return
  actionError.value = ''
  const current: ActuatorCommand = {
    heater: detail.value.heater === 1 ? 1 : 0,
    fan: detail.value.fan === 1 ? 1 : 0,
    pump: detail.value.pump === 1 ? 1 : 0,
  }
  current[key] = current[key] === 1 ? 0 : 1
  try {
    await store.changeActuators(detail.value.node_id, current)
  } catch (cause) {
    actionError.value = getApiErrorStatus(cause) === 409
      ? 'Node harus berada di mode Manual untuk mengubah aktuator.'
      : getApiErrorMessage(cause, 'Perintah aktuator gagal dikirim.')
  }
}
</script>

<template>
  <section class="page">
    <div class="breadcrumb"><RouterLink to="/nodes">Daftar node</RouterLink><span>/</span><span>Detail node</span></div>
    <PageFeedback :loading="loading" :error="error">
      <template v-if="detail">
        <div class="page-heading heading-with-action">
          <div>
            <p class="eyebrow">{{ detail.node_id }}</p>
            <h1>{{ detail.name }}</h1>
            <p class="muted">Pembaruan terakhir {{ formatWib(detail.last_seen) }}</p>
          </div>
          <StatusBadge :status="detail.status" :offline="!detail.online" />
        </div>

        <div class="detail-grid">
          <article class="panel gauge-panel">
            <div class="panel-heading"><div><p class="eyebrow">SENSOR SUHU</p><h2>Suhu air</h2></div><span class="sensor-icon warm">°</span></div>
            <div class="gauge-wrap">
              <div class="gauge gauge-temp" :style="{ '--gauge-value': `${temperatureProgress}%` }">
                <div class="gauge-center"><strong>{{ formatNumber(detail.temp) }}</strong><span>°C</span></div>
              </div>
            </div>
            <div class="gauge-scale"><span>0 °C</span><span>Rentang 0–40 °C</span><span>40 °C</span></div>
            <div class="class-row"><span>Kelas suhu</span><StatusBadge :status="detail.class_temp" /></div>
          </article>

          <article class="panel gauge-panel">
            <div class="panel-heading"><div><p class="eyebrow">SENSOR KEKERUHAN</p><h2>Kekeruhan air</h2></div><span class="sensor-icon cool">≈</span></div>
            <div class="gauge-wrap">
              <div class="gauge gauge-water" :style="{ '--gauge-value': `${turbidityProgress}%` }">
                <div class="gauge-center"><strong>{{ formatNumber(detail.ntu) }}</strong><span>NTU</span></div>
              </div>
            </div>
            <div class="gauge-scale"><span>0 NTU</span><span>Rentang 0–200 NTU</span><span>200 NTU</span></div>
            <div class="class-row"><span>Kelas kekeruhan</span><StatusBadge :status="detail.class_turbidity" /></div>
          </article>

          <article class="panel overall-panel">
            <div><p class="eyebrow">KLASIFIKASI MODEL</p><h2>Kondisi keseluruhan</h2></div>
            <div class="overall-status"><StatusBadge :status="detail.status" /><span>Status utama</span></div>
            <div class="class-row"><span>Kelas overall</span><StatusBadge :status="detail.class_overall" /></div>
            <div class="class-row"><span>Koneksi node</span><span :class="detail.online ? 'text-online' : 'text-offline'">{{ detail.online ? 'Online' : 'Offline' }}</span></div>
            <div class="class-row"><span>Mode kendali</span><span class="mode-label">{{ detail.mode === 'auto' ? 'Otomatis' : 'Manual' }}</span></div>
          </article>

          <article class="panel control-panel">
            <div class="panel-heading">
              <div><p class="eyebrow">KENDALI PERANGKAT</p><h2>Mode dan aktuator</h2></div>
              <button
                class="switch"
                :class="{ 'switch-on': detail.mode === 'auto' }"
                role="switch"
                :aria-checked="detail.mode === 'auto'"
                :aria-label="`Mode ${detail.mode === 'auto' ? 'otomatis' : 'manual'}`"
                type="button"
                @click="toggleMode"
              ><span></span></button>
            </div>
            <div class="mode-card"><strong>Mode {{ detail.mode === 'auto' ? 'Otomatis' : 'Manual' }}</strong><span>{{ detail.mode === 'auto' ? 'Sistem mengatur aktuator berdasarkan kondisi air.' : 'Atur aktuator secara manual.' }}</span></div>
            <p v-if="waitingForState" class="pending-message" role="status">Menunggu konfirmasi status aktuator dari node...</p>
            <p v-if="actionError" class="form-error" role="alert">{{ actionError }}</p>
            <div class="actuator-list">
              <div v-for="item in [
                { key: 'heater', label: 'Pemanas', value: detail.heater },
                { key: 'fan', label: 'Kipas', value: detail.fan },
                { key: 'pump', label: 'Pompa filter', value: detail.pump },
              ]" :key="item.key" class="actuator-row">
                <div><strong>{{ item.label }}</strong><span :class="item.value === 1 ? 'text-online' : 'muted'">{{ item.value === null ? 'Belum diketahui' : item.value === 1 ? 'Aktif' : 'Nonaktif' }}</span></div>
                <button
                  class="switch"
                  :class="{ 'switch-on': item.value === 1 }"
                  role="switch"
                  :aria-checked="item.value === 1"
                  :aria-label="`${item.label} ${item.value === 1 ? 'aktif' : 'nonaktif'}`"
                  :disabled="detail.mode === 'auto' || waitingForState || !detail.online"
                  type="button"
                  @click="toggleActuator(item.key as keyof ActuatorCommand)"
                ><span></span></button>
              </div>
            </div>
            <p class="failsafe-note" :class="{ 'failsafe-active': detail.failsafe }">
              <span></span> Fail-safe {{ detail.failsafe ? 'aktif' : 'tidak aktif' }}
            </p>
          </article>
        </div>

        <div class="detail-links">
          <RouterLink class="button button-secondary" :to="{ name: 'node-history', params: { id: detail.node_id } }">Lihat grafik riwayat <span>→</span></RouterLink>
          <RouterLink class="button button-secondary" :to="{ name: 'node-logs', params: { id: detail.node_id } }">Log &amp; notifikasi <span>→</span></RouterLink>
        </div>
      </template>
    </PageFeedback>
  </section>
</template>
