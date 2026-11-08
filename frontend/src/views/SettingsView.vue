<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import PageFeedback from '../components/PageFeedback.vue'
import { getApiErrorMessage } from '../api/client'
import { saveTelegramChatId } from '../api/nodes'
import { useNodesStore } from '../stores/nodes'

const store = useNodesStore()
const selectedNodeId = ref('')
const nodeName = ref('')
const telegramChatId = ref('')
const loading = ref(true)
const savingName = ref(false)
const savingTelegram = ref(false)
const loadError = ref('')
const actionError = ref('')
const success = ref('')
const selectedNode = computed(() => store.details[selectedNodeId.value] ?? null)

onMounted(async () => {
  try {
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
    const detail = store.details[id] ?? await store.loadNode(id)
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
    await store.renameNode(selectedNodeId.value, nodeName.value)
    success.value = 'Nama node berhasil diperbarui.'
  } catch (cause) {
    actionError.value = getApiErrorMessage(cause, 'Nama node tidak dapat disimpan.')
  } finally {
    savingName.value = false
  }
}

async function saveTelegram() {
  savingTelegram.value = true
  actionError.value = ''
  success.value = ''
  try {
    await saveTelegramChatId(telegramChatId.value)
    success.value = 'Telegram chat ID berhasil disimpan.'
  } catch (cause) {
    actionError.value = getApiErrorMessage(cause, 'Telegram chat ID tidak dapat disimpan.')
  } finally {
    savingTelegram.value = false
  }
}
</script>

<template>
  <section class="page">
    <div class="page-heading">
      <p class="eyebrow">PREFERENSI APLIKASI</p>
      <h1>Pengaturan</h1>
      <p class="muted">Kelola identitas node dan tujuan notifikasi.</p>
    </div>

    <PageFeedback :loading="loading" :error="loadError">
      <div class="settings-grid">
        <article class="panel settings-panel">
          <div class="panel-heading"><div><p class="eyebrow">IDENTITAS PERANGKAT</p><h2>Nama node</h2></div><span class="settings-icon">⌁</span></div>
          <p class="muted settings-description">Nama ini digunakan untuk mengenali akuarium pada dashboard.</p>
          <label for="node-select">Pilih node</label>
          <select id="node-select" v-model="selectedNodeId">
            <option v-for="node in store.nodes" :key="node.node_id" :value="node.node_id">{{ node.name }} · {{ node.node_id }}</option>
          </select>
          <label for="node-name">Nama tampilan</label>
          <input id="node-name" v-model="nodeName" maxlength="100" placeholder="Contoh: Akuarium ruang tamu">
          <p v-if="selectedNode" class="field-hint">ID perangkat: {{ selectedNode.node_id }}</p>
          <button class="button button-primary" type="button" :disabled="savingName || !nodeName.trim()" @click="saveNodeName">
            {{ savingName ? 'Menyimpan...' : 'Simpan nama node' }}
          </button>
        </article>

        <article class="panel settings-panel">
          <div class="panel-heading"><div><p class="eyebrow">NOTIFIKASI TELEGRAM</p><h2>Telegram chat ID</h2></div><span class="settings-icon telegram-icon">✈</span></div>
          <p class="muted settings-description">Masukkan chat ID tujuan untuk menerima peringatan kondisi akuarium.</p>
          <label for="telegram-chat-id">Chat ID</label>
          <input id="telegram-chat-id" v-model="telegramChatId" inputmode="numeric" placeholder="Masukkan Telegram chat ID">
          <p class="field-hint">Nilai chat ID tidak ditampilkan kembali setelah disimpan.</p>
          <button class="button button-primary" type="button" :disabled="savingTelegram || !telegramChatId.trim()" @click="saveTelegram">
            {{ savingTelegram ? 'Menyimpan...' : 'Simpan chat ID' }}
          </button>
        </article>
      </div>
      <p v-if="actionError" class="form-error settings-message" role="alert">{{ actionError }}</p>
      <p v-else-if="success" class="success-message settings-message" role="status">{{ success }}</p>
    </PageFeedback>
  </section>
</template>
