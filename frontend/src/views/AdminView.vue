<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import {
  Copy,
  Plus,
  RefreshCw,
  ShieldCheck,
  Trash2,
  Unlink,
  User,
} from 'lucide-vue-next'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import EmptyState from '../components/EmptyState.vue'
import PageFeedback from '../components/PageFeedback.vue'
import PageHeader from '../components/PageHeader.vue'
import { getApiErrorMessage, getApiErrorStatus } from '../api/client'
import {
  adminUnclaimNode,
  createAdminNode,
  deleteAdminNode,
  regenerateClaimCode,
} from '../api/nodes'
import { useNodesStore } from '../stores/nodes'
import { useToastStore } from '../stores/toast'
import { formatWib } from '../utils/format'
import type { AdminNode } from '../types'

const store = useNodesStore()
const toast = useToastStore()

const newNodeId = ref('')
const creating = ref(false)
const createError = ref('')
const busyNodeId = ref('')

const unclaimTarget = ref<AdminNode | null>(null)
const deleteTarget = ref<AdminNode | null>(null)
const dialogLoading = ref(false)

const nodeIdRegex = /^(aq-[0-9a-f]{6}|aq-sim[0-9]{3})$/

const onlineCount = computed(() => store.adminNodes.filter((n) => n.online).length)

onMounted(() => {
  void store.loadAdminNodes().catch(() => undefined)
})

async function handleCreateNode() {
  createError.value = ''
  const normalized = newNodeId.value.trim().toLowerCase()
  newNodeId.value = normalized

  if (!nodeIdRegex.test(normalized)) {
    createError.value =
      'Format ID node tidak valid. Gunakan aq- diikuti 6 karakter heksadesimal (contoh: aq-1a2b3c) atau aq-sim diikuti 3 digit (contoh: aq-sim001).'
    return
  }

  creating.value = true
  try {
    const created = await createAdminNode(normalized)
    store.adminNodes.unshift(created)
    newNodeId.value = ''
    toast.success(`Node ${created.node_id} berhasil didaftarkan dengan kode klaim ${created.claim_code}.`)
  } catch (cause) {
    const status = getApiErrorStatus(cause)
    if (status === 409) {
      createError.value = 'ID node tersebut sudah terdaftar di sistem.'
    } else {
      createError.value = getApiErrorMessage(cause, 'Gagal mendaftarkan node baru.')
    }
  } finally {
    creating.value = false
  }
}

async function copyClaimCode(node: AdminNode) {
  try {
    await navigator.clipboard.writeText(node.claim_code)
    toast.success(`Kode klaim ${node.claim_code} (${node.node_id}) disalin ke papan klip.`)
  } catch {
    toast.error('Gagal menyalin kode klaim secara otomatis.')
  }
}

async function handleRegenerateCode(node: AdminNode) {
  busyNodeId.value = node.node_id
  try {
    const updated = await regenerateClaimCode(node.node_id)
    const idx = store.adminNodes.findIndex((item) => item.node_id === node.node_id)
    if (idx >= 0) store.adminNodes[idx] = updated
    toast.success(`Kode klaim baru untuk ${node.node_id}: ${updated.claim_code}`)
  } catch (cause) {
    toast.error(getApiErrorMessage(cause, 'Gagal membuat kode klaim baru.'))
  } finally {
    busyNodeId.value = ''
  }
}

async function confirmAdminUnclaim() {
  if (!unclaimTarget.value) return
  dialogLoading.value = true
  const target = unclaimTarget.value
  try {
    await adminUnclaimNode(target.node_id)
    const idx = store.adminNodes.findIndex((item) => item.node_id === target.node_id)
    if (idx >= 0 && store.adminNodes[idx]) {
      store.adminNodes[idx]!.owner_username = null
    }
    unclaimTarget.value = null
    toast.success(`Kepemilikan ${target.node_id} berhasil dilepas dan mode diatur ke Manual.`)
  } catch (cause) {
    toast.error(getApiErrorMessage(cause, 'Gagal melepas pemilik node.'))
  } finally {
    dialogLoading.value = false
  }
}

async function confirmAdminDelete() {
  if (!deleteTarget.value) return
  dialogLoading.value = true
  const target = deleteTarget.value
  try {
    await deleteAdminNode(target.node_id)
    store.adminNodes = store.adminNodes.filter((item) => item.node_id !== target.node_id)
    deleteTarget.value = null
    toast.success(`Node ${target.node_id} beserta seluruh datanya berhasil dihapus.`)
  } catch (cause) {
    toast.error(getApiErrorMessage(cause, 'Gagal menghapus node.'))
  } finally {
    dialogLoading.value = false
  }
}
</script>

<template>
  <section class="page">
    <PageHeader
      eyebrow="ADMINISTRASI PERANGKAT"
      title="Kelola node akuarium"
      subtitle="Daftarkan perangkat baru, cetak kode klaim untuk stiker akuarium, dan kelola kepemilikan akun klien."
    >
      <template #actions>
        <div class="online-count">
          <span class="status-dot dot-online" aria-hidden="true"></span>
          <span>{{ onlineCount }} dari {{ store.adminNodes.length }} perangkat online</span>
        </div>
      </template>
    </PageHeader>

    <article class="panel admin-register-panel">
      <div class="panel-heading">
        <div>
          <p class="eyebrow">PENDAFTARAN PERANGKAT</p>
          <h2>Tambah node baru</h2>
        </div>
        <span class="settings-icon" aria-hidden="true">
          <ShieldCheck :size="20" />
        </span>
      </div>
      <form class="admin-create-form" @submit.prevent="handleCreateNode">
        <div class="admin-create-input-wrap">
          <label for="admin-new-node" class="sr-only">ID Node</label>
          <input
            id="admin-new-node"
            v-model="newNodeId"
            placeholder="Masukkan ID node (contoh: aq-1a2b3c atau aq-sim001)"
            autocomplete="off"
            required
            @input="newNodeId = newNodeId.toLowerCase()"
          />
        </div>
        <button class="button button-primary" type="submit" :disabled="creating">
          <Plus :size="18" aria-hidden="true" />
          <span>{{ creating ? 'Mendaftarkan...' : 'Tambah Node' }}</span>
        </button>
      </form>
      <p class="field-hint">
        Format wajib: <code>aq-xxxxxx</code> (6 heksadesimal huruf kecil) atau <code>aq-simxxx</code> (3 digit angka).
      </p>
      <p v-if="createError" class="form-error" role="alert">{{ createError }}</p>
    </article>

    <PageFeedback :loading="store.loading && !store.adminNodes.length" :error="store.error">
      <EmptyState
        v-if="!store.adminNodes.length"
        :icon="ShieldCheck"
        title="Belum ada node terdaftar"
        description="Daftarkan ID perangkat di atas atau nyalakan perangkat akuarium agar terdaftar otomatis lewat MQTT."
      />

      <template v-else>
        <!-- Tabel Desktop -->
        <div class="panel admin-table-wrap">
          <table class="admin-table">
            <thead>
              <tr>
                <th>Node &amp; Nama</th>
                <th>Status</th>
                <th>Kode Klaim</th>
                <th>Pemilik</th>
                <th>Terakhir Terlihat</th>
                <th>Terdaftar</th>
                <th class="text-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="node in store.adminNodes" :key="node.node_id">
                <td>
                  <div class="admin-node-cell">
                    <strong>{{ node.node_id }}</strong>
                    <span class="muted">{{ node.name }}</span>
                  </div>
                </td>
                <td>
                  <span
                    class="online-pill"
                    :class="node.online ? 'online-pill-active' : 'online-pill-idle'"
                  >
                    <span
                      class="status-dot"
                      :class="node.online ? 'dot-online' : 'dot-offline'"
                    ></span>
                    {{ node.online ? 'Online' : 'Offline' }}
                  </span>
                </td>
                <td>
                  <div class="claim-code-group">
                    <code class="claim-code-badge">{{ node.claim_code }}</code>
                    <button
                      class="button button-quiet button-sm"
                      type="button"
                      :aria-label="`Salin kode klaim ${node.node_id}`"
                      @click="copyClaimCode(node)"
                    >
                      <Copy :size="16" aria-hidden="true" />
                      <span>Salin</span>
                    </button>
                  </div>
                </td>
                <td>
                  <span v-if="node.owner_username" class="owner-badge owner-claimed">
                    <User :size="15" aria-hidden="true" />
                    <span>{{ node.owner_username }}</span>
                  </span>
                  <span v-else class="owner-badge owner-unclaimed">Belum diklaim</span>
                </td>
                <td>{{ formatWib(node.last_seen) }}</td>
                <td>{{ formatWib(node.created_at) }}</td>
                <td>
                  <div class="admin-row-actions">
                    <button
                      class="button button-secondary button-sm"
                      type="button"
                      :disabled="busyNodeId === node.node_id"
                      @click="handleRegenerateCode(node)"
                    >
                      <RefreshCw :size="16" aria-hidden="true" />
                      <span>Buat kode baru</span>
                    </button>
                    <button
                      v-if="node.owner_username"
                      class="button button-quiet button-sm"
                      type="button"
                      :disabled="busyNodeId === node.node_id"
                      @click="unclaimTarget = node"
                    >
                      <Unlink :size="16" aria-hidden="true" />
                      <span>Lepas pemilik</span>
                    </button>
                    <button
                      class="button button-danger-outline button-sm"
                      type="button"
                      :disabled="busyNodeId === node.node_id"
                      @click="deleteTarget = node"
                    >
                      <Trash2 :size="16" aria-hidden="true" />
                      <span>Hapus</span>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Kartu Seluler -->
        <div class="admin-mobile-list">
          <article v-for="node in store.adminNodes" :key="node.node_id" class="panel admin-card">
            <div class="admin-card-top">
              <div>
                <strong class="admin-card-id">{{ node.node_id }}</strong>
                <p class="muted">{{ node.name }}</p>
              </div>
              <span
                class="online-pill"
                :class="node.online ? 'online-pill-active' : 'online-pill-idle'"
              >
                <span class="status-dot" :class="node.online ? 'dot-online' : 'dot-offline'"></span>
                {{ node.online ? 'Online' : 'Offline' }}
              </span>
            </div>

            <div class="admin-card-code-box">
              <div>
                <span class="eyebrow">KODE KLAIM STIKER</span>
                <code class="claim-code-badge">{{ node.claim_code }}</code>
              </div>
              <button class="button button-quiet button-sm" type="button" @click="copyClaimCode(node)">
                <Copy :size="16" aria-hidden="true" />
                <span>Salin</span>
              </button>
            </div>

            <div class="admin-card-meta">
              <div class="class-row">
                <span>Pemilik akun</span>
                <span v-if="node.owner_username" class="owner-badge owner-claimed">
                  <User :size="15" aria-hidden="true" />
                  <span>{{ node.owner_username }}</span>
                </span>
                <span v-else class="owner-badge owner-unclaimed">Belum diklaim</span>
              </div>
              <div class="class-row">
                <span>Terakhir terlihat</span>
                <span>{{ formatWib(node.last_seen) }}</span>
              </div>
              <div class="class-row">
                <span>Waktu terdaftar</span>
                <span>{{ formatWib(node.created_at) }}</span>
              </div>
            </div>

            <div class="admin-card-actions">
              <button
                class="button button-secondary"
                type="button"
                :disabled="busyNodeId === node.node_id"
                @click="handleRegenerateCode(node)"
              >
                <RefreshCw :size="17" aria-hidden="true" />
                <span>Buat kode baru</span>
              </button>
              <button
                v-if="node.owner_username"
                class="button button-quiet"
                type="button"
                :disabled="busyNodeId === node.node_id"
                @click="unclaimTarget = node"
              >
                <Unlink :size="17" aria-hidden="true" />
                <span>Lepas pemilik</span>
              </button>
              <button
                class="button button-danger-outline"
                type="button"
                :disabled="busyNodeId === node.node_id"
                @click="deleteTarget = node"
              >
                <Trash2 :size="17" aria-hidden="true" />
                <span>Hapus node</span>
              </button>
            </div>
          </article>
        </div>
      </template>
    </PageFeedback>

    <ConfirmDialog
      :open="Boolean(unclaimTarget)"
      title="Lepas pemilik akuarium?"
      :description="`Node ${unclaimTarget?.node_id} tidak lagi dimiliki oleh akun '${unclaimTarget?.owner_username}' dan mode perangkat akan dikembalikan ke Manual.`"
      confirm-label="Ya, lepas pemilik"
      :loading="dialogLoading"
      @cancel="unclaimTarget = null"
      @confirm="confirmAdminUnclaim"
    />

    <ConfirmDialog
      :open="Boolean(deleteTarget)"
      title="Hapus node beserta seluruh datanya?"
      :description="`Tindakan ini akan menghapus node ${deleteTarget?.node_id} beserta seluruh riwayat telemetry, prediksi, log aktuator, dan notifikasi secara permanen. Jika perangkat fisik masih menyala dan terhubung ke broker MQTT, perangkat akan terdaftar kembali secara otomatis.`"
      confirm-label="Ya, hapus permanen"
      :danger="true"
      :loading="dialogLoading"
      @cancel="deleteTarget = null"
      @confirm="confirmAdminDelete"
    />
  </section>
</template>
