<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Bell, CheckCircle2, Clock, Settings, Unlink, User } from 'lucide-vue-next'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import PageFeedback from '../components/PageFeedback.vue'
import PageHeader from '../components/PageHeader.vue'
import { getApiErrorMessage } from '../api/client'
import { saveTelegramChatId } from '../api/nodes'
import { useAuthStore } from '../stores/auth'
import { useNodesStore } from '../stores/nodes'
import { useToastStore } from '../stores/toast'

const auth = useAuthStore()
const store = useNodesStore()
const toast = useToastStore()
const router = useRouter()

const selectedNodeId = ref('')
const nodeName = ref('')
const telegramChatId = ref('')
const loading = ref(true)
const savingName = ref(false)
const savingTelegram = ref(false)
const unclaiming = ref(false)
const confirmUnclaimOpen = ref(false)
const loadError = ref('')
const actionError = ref('')
const success = ref('')

const selectedNode = computed(() => store.details[selectedNodeId.value] ?? null)

onMounted(async () => {
  try {
    await auth.fetchMe().catch(() => undefined)
    if (!store.nodes.length) await store.loadNodes()
    selectedNodeId.value = store.activeNodeId || store.nodes[0]?.node_id || ''
    if (selectedNodeId.value) {
      const detail = await store.loadNode(selectedNodeId.value)
      nodeName.value = detail.name
    }
  } catch (cause) {
    loadError.value = getApiErrorMessage(cause, 'Pengaturan tidak dapat dimuat.')
  } finally {
    loading.value = false
  }
})

watch(selectedNodeId, async (id, previousId) => {
  if (!id || !previousId) return
  loadError.value = ''
  try {
    const detail = store.details[id] ?? (await store.loadNode(id))
    nodeName.value = detail.name
    success.value = ''
  } catch (cause) {
    loadError.value = getApiErrorMessage(cause, 'Detail node tidak dapat dimuat.')
  }
})

async function saveNodeName() {
  if (!selectedNodeId.value || !nodeName.value.trim()) return
  savingName.value = true
  actionError.value = ''
  success.value = ''
  try {
    await store.renameNode(selectedNodeId.value, nodeName.value.trim())
    success.value = 'Nama akuarium berhasil diperbarui.'
    toast.success('Nama akuarium berhasil diperbarui.')
  } catch (cause) {
    actionError.value = getApiErrorMessage(cause, 'Nama node tidak dapat disimpan.')
    toast.error(actionError.value)
  } finally {
    savingName.value = false
  }
}

async function saveTelegram() {
  savingTelegram.value = true
  actionError.value = ''
  success.value = ''
  try {
    await saveTelegramChatId(telegramChatId.value.trim())
    await auth.fetchMe()
    telegramChatId.value = ''
    success.value = 'Telegram chat ID berhasil disimpan.'
    toast.success('Telegram chat ID berhasil disimpan.')
  } catch (cause) {
    actionError.value = getApiErrorMessage(cause, 'Telegram chat ID tidak dapat disimpan.')
    toast.error(actionError.value)
  } finally {
    savingTelegram.value = false
  }
}

async function confirmUnclaim() {
  if (!selectedNodeId.value) return
  unclaiming.value = true
  actionError.value = ''
  try {
    const targetName = selectedNode.value?.name ?? selectedNodeId.value
    await store.unclaimNode(selectedNodeId.value)
    confirmUnclaimOpen.value = false
    toast.success(`Akuarium ${targetName} berhasil dilepas dari akun Anda.`)
    void router.push({ name: 'nodes' })
  } catch (cause) {
    actionError.value = getApiErrorMessage(cause, 'Gagal melepas kepemilikan akuarium.')
    toast.error(actionError.value)
  } finally {
    unclaiming.value = false
  }
}
</script>

<template>
  <section class="page">
    <PageHeader
      eyebrow="PREFERENSI AKUN & AKUARIUM"
      title="Pengaturan"
      subtitle="Kelola profil akun, notifikasi Telegram, dan akuarium Anda."
    />

    <PageFeedback :loading="loading" :error="loadError">
      <article class="panel account-summary-panel">
        <div class="account-summary-item">
          <span class="settings-icon" aria-hidden="true">
            <User :size="20" />
          </span>
          <div>
            <p class="eyebrow">AKUN LOGIN</p>
            <strong>{{ auth.username || 'Klien' }}</strong>
          </div>
        </div>
        <div class="account-summary-item">
          <span
            class="telegram-status-pill"
            :class="auth.telegramLinked ? 'telegram-linked' : 'telegram-unlinked'"
          >
            <CheckCircle2 v-if="auth.telegramLinked" :size="17" aria-hidden="true" />
            <Clock v-else :size="17" aria-hidden="true" />
            <span>Status Telegram: {{ auth.telegramLinked ? 'Terhubung' : 'Belum terhubung' }}</span>
          </span>
        </div>
      </article>

      <div class="settings-grid">
        <article class="panel settings-panel">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">IDENTITAS AKUARIUM</p>
              <h2>Nama &amp; kepemilikan akuarium</h2>
            </div>
            <span class="settings-icon" aria-hidden="true">
              <Settings :size="20" />
            </span>
          </div>

          <p class="muted settings-description">
            Ubah nama tampilan akuarium pada dashboard atau lepas akuarium dari akun Anda.
          </p>

          <template v-if="store.nodes.length">
            <label for="node-select">Pilih akuarium</label>
            <select id="node-select" v-model="selectedNodeId">
              <option v-for="node in store.nodes" :key="node.node_id" :value="node.node_id">
                {{ node.name }} · {{ node.node_id }}
              </option>
            </select>

            <label for="node-name">Nama tampilan</label>
            <input
              id="node-name"
              v-model="nodeName"
              maxlength="100"
              placeholder="Contoh: Akuarium Ruang Tamu"
            />
            <p v-if="selectedNode" class="field-hint">ID perangkat: {{ selectedNode.node_id }}</p>

            <div class="settings-actions-row">
              <button
                class="button button-primary"
                type="button"
                :disabled="savingName || !nodeName.trim()"
                @click="saveNodeName"
              >
                {{ savingName ? 'Menyimpan...' : 'Simpan nama node' }}
              </button>

              <button
                class="button button-danger-outline"
                type="button"
                :disabled="!selectedNodeId || unclaiming"
                @click="confirmUnclaimOpen = true"
              >
                <Unlink :size="18" aria-hidden="true" />
                <span>Lepas Akuarium</span>
              </button>
            </div>
          </template>
          <p v-else class="muted">
            Anda belum memiliki akuarium terdaftar. Tambahkan akuarium melalui halaman Daftar Node.
          </p>
        </article>

        <article class="panel settings-panel">
          <div class="panel-heading">
            <div>
              <p class="eyebrow">NOTIFIKASI TELEGRAM</p>
              <h2>Telegram Chat ID</h2>
            </div>
            <span class="settings-icon telegram-icon" aria-hidden="true">
              <Bell :size="20" />
            </span>
          </div>

          <p class="muted settings-description">
            Masukkan Chat ID Telegram untuk menerima peringatan saat kondisi air membutuhkan perhatian.
          </p>

          <label for="telegram-chat-id">Chat ID tujuan</label>
          <input
            id="telegram-chat-id"
            v-model="telegramChatId"
            inputmode="numeric"
            placeholder="Contoh: 123456789"
          />
          <p class="field-hint">
            Demi privasi, nilai Chat ID tidak ditampilkan kembali setelah disimpan.
          </p>

          <button
            class="button button-primary"
            type="button"
            :disabled="savingTelegram || !telegramChatId.trim()"
            @click="saveTelegram"
          >
            {{ savingTelegram ? 'Menyimpan...' : 'Simpan Chat ID' }}
          </button>
        </article>
      </div>

      <p v-if="actionError" class="form-error settings-message" role="alert">{{ actionError }}</p>
      <p v-else-if="success" class="success-message settings-message" role="status">{{ success }}</p>
    </PageFeedback>

    <ConfirmDialog
      :open="confirmUnclaimOpen"
      title="Lepas akuarium dari akun?"
      :description="`Akuarium ${selectedNode?.name ?? selectedNodeId} (${selectedNodeId}) akan dilepas dari akun Anda dan dikembalikan ke mode Manual. Anda dapat mengklaimnya kembali kapan saja dengan kode klaim.`"
      confirm-label="Ya, lepas akuarium"
      :danger="true"
      :loading="unclaiming"
      @cancel="confirmUnclaimOpen = false"
      @confirm="confirmUnclaim"
    />
  </section>
</template>
