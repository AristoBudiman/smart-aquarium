<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import {
  Activity,
  Bell,
  Droplets,
  Fan,
  Flame,
  History,
  ShieldAlert,
  ShieldCheck,
  Sliders,
} from 'lucide-vue-next'
import PageFeedback from '../components/PageFeedback.vue'
import PageHeader from '../components/PageHeader.vue'
import SensorGauge from '../components/SensorGauge.vue'
import StatusBadge from '../components/StatusBadge.vue'
import ToggleSwitch from '../components/ToggleSwitch.vue'
import { getApiErrorMessage, getApiErrorStatus } from '../api/client'
import { useNodesStore } from '../stores/nodes'
import { useToastStore } from '../stores/toast'
import { formatWib } from '../utils/format'
import type { ActuatorCommand } from '../types'

const route = useRoute()
const store = useNodesStore()
const toast = useToastStore()
const nodeId = computed(() => (typeof route.params.id === 'string' ? route.params.id : ''))
const loading = ref(true)
const error = ref('')
const actionError = ref('')
const waitingForState = computed(() => store.pendingActuators[nodeId.value] === true)
const detail = computed(() => store.details[nodeId.value] ?? null)

watch(
  nodeId,
  async (id) => {
    loading.value = true
    error.value = ''
    try {
      await store.loadNode(id)
    } catch (cause) {
      error.value = getApiErrorMessage(cause, 'Detail node tidak dapat dimuat.')
    } finally {
      loading.value = false
    }
  },
  { immediate: true },
)

async function toggleMode() {
  if (!detail.value) return
  actionError.value = ''
  const nextMode = detail.value.mode === 'auto' ? 'manual' : 'auto'
  try {
    await store.changeMode(detail.value.node_id, nextMode)
    toast.success(`Mode kendali diubah ke ${nextMode === 'auto' ? 'Otomatis' : 'Manual'}.`)
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
    actionError.value =
      getApiErrorStatus(cause) === 409
        ? 'Node harus berada di mode Manual untuk mengubah aktuator.'
        : getApiErrorMessage(cause, 'Perintah aktuator gagal dikirim.')
  }
}

const actuatorItems = computed(() => [
  { key: 'heater' as const, label: 'Pemanas air', value: detail.value?.heater ?? null, icon: Flame },
  { key: 'fan' as const, label: 'Kipas pendingin', value: detail.value?.fan ?? null, icon: Fan },
  { key: 'pump' as const, label: 'Pompa filter', value: detail.value?.pump ?? null, icon: Droplets },
])
</script>

<template>
  <section class="page">
    <PageFeedback :loading="loading" :error="error">
      <template v-if="detail">
        <PageHeader
          :eyebrow="detail.node_id"
          :title="detail.name"
          :subtitle="`Pembaruan terakhir ${formatWib(detail.last_seen)}`"
        >
          <template #breadcrumb>
            <RouterLink to="/nodes">Daftar akuarium</RouterLink>
            <span>/</span>
            <span>Detail akuarium</span>
          </template>
          <template #actions>
            <StatusBadge :status="detail.status" :offline="!detail.online" size="lg" />
          </template>
        </PageHeader>

        <div class="detail-grid">
          <SensorGauge
            eyebrow="SENSOR SUHU"
            title="Suhu air"
            :value="detail.temp"
            unit="°C"
            :max="40"
            min-label="0 °C"
            range-label="Rentang 0–40 °C (Acuan 27,48 °C)"
            max-label="40 °C"
            class-label="Kelas suhu"
            :status-class="detail.class_temp"
            variant="temp"
          />

          <SensorGauge
            eyebrow="SENSOR KEKERUHAN"
            title="Kekeruhan air"
            :value="detail.ntu"
            unit="NTU"
            :max="200"
            min-label="0 NTU"
            range-label="Rentang 0–200 NTU"
            max-label="200 NTU"
            class-label="Kelas kekeruhan"
            :status-class="detail.class_turbidity"
            variant="water"
          />

          <article class="panel overall-panel">
            <div class="panel-heading">
              <div>
                <p class="eyebrow">KLASIFIKASI MODEL</p>
                <h2>Kondisi keseluruhan</h2>
              </div>
              <span class="sensor-icon teal-icon" aria-hidden="true">
                <Activity :size="20" :stroke-width="2.1" />
              </span>
            </div>

            <div class="overall-status">
              <StatusBadge :status="detail.status" :offline="!detail.online" size="lg" />
              <span>Status utama keputusan sistem</span>
            </div>

            <div class="class-row">
              <span>Kelas model overall</span>
              <StatusBadge :status="detail.class_overall" />
            </div>
            <div class="class-row">
              <span>Koneksi perangkat</span>
              <span class="inline-status" :class="detail.online ? 'text-online' : 'text-offline'">
                <span class="status-dot" :class="detail.online ? 'dot-online' : 'dot-offline'"></span>
                {{ detail.online ? 'Online' : 'Offline' }}
              </span>
            </div>
            <div class="class-row">
              <span>Mode kendali</span>
              <span class="mode-label">{{ detail.mode === 'auto' ? 'Otomatis (AI)' : 'Manual' }}</span>
            </div>
          </article>

          <article class="panel control-panel">
            <div class="panel-heading">
              <div>
                <p class="eyebrow">KENDALI PERANGKAT</p>
                <h2>Mode dan aktuator</h2>
              </div>
              <ToggleSwitch
                :model-value="detail.mode === 'auto'"
                :label="`Mode ${detail.mode === 'auto' ? 'otomatis' : 'manual'}`"
                @toggle="toggleMode"
              />
            </div>

            <div class="mode-card">
              <div class="mode-card-head">
                <Sliders :size="17" aria-hidden="true" />
                <strong>Mode {{ detail.mode === 'auto' ? 'Otomatis' : 'Manual' }}</strong>
              </div>
              <span>
                {{
                  detail.mode === 'auto'
                    ? 'Sistem mengatur pemanas, kipas, dan pompa secara otomatis berdasarkan hasil klasifikasi Random Forest.'
                    : 'Anda memegang kendali penuh untuk menyalakan atau mematikan relay secara manual.'
                }}
              </span>
            </div>

            <p v-if="waitingForState" class="pending-message" role="status">
              Menunggu konfirmasi perubahan status relay dari perangkat...
            </p>
            <p v-if="actionError" class="form-error" role="alert">{{ actionError }}</p>

            <div class="actuator-list">
              <div v-for="item in actuatorItems" :key="item.key" class="actuator-row">
                <div class="actuator-info">
                  <span class="actuator-icon" :class="{ 'actuator-icon-active': item.value === 1 }" aria-hidden="true">
                    <component :is="item.icon" :size="18" />
                  </span>
                  <div>
                    <strong>{{ item.label }}</strong>
                    <span :class="item.value === 1 ? 'text-online' : 'muted'">
                      {{ item.value === null ? 'Belum diketahui' : item.value === 1 ? 'Aktif (ON)' : 'Nonaktif (OFF)' }}
                    </span>
                  </div>
                </div>
                <ToggleSwitch
                  :model-value="item.value === 1"
                  :label="`${item.label} ${item.value === 1 ? 'aktif' : 'nonaktif'}`"
                  :disabled="detail.mode === 'auto' || waitingForState || !detail.online"
                  @toggle="toggleActuator(item.key)"
                />
              </div>
            </div>

            <p class="failsafe-note" :class="{ 'failsafe-active': detail.failsafe }">
              <ShieldAlert v-if="detail.failsafe" :size="18" aria-hidden="true" />
              <ShieldCheck v-else :size="18" aria-hidden="true" />
              <span>
                Proteksi fail-safe lokal:
                <strong>{{ detail.failsafe ? 'AKTIF' : 'Tidak aktif (Normal)' }}</strong>
              </span>
            </p>
          </article>
        </div>

        <div class="detail-links">
          <RouterLink class="button button-secondary" :to="{ name: 'node-history', params: { id: detail.node_id } }">
            <History :size="18" aria-hidden="true" />
            <span>Lihat grafik riwayat</span>
          </RouterLink>
          <RouterLink class="button button-secondary" :to="{ name: 'node-logs', params: { id: detail.node_id } }">
            <Bell :size="18" aria-hidden="true" />
            <span>Log &amp; notifikasi</span>
          </RouterLink>
        </div>
      </template>
    </PageFeedback>
  </section>
</template>
