<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { KeyRound, Plus, Sparkles, Waves, X } from 'lucide-vue-next'
import EmptyState from '../components/EmptyState.vue'
import NodeCard from '../components/NodeCard.vue'
import PageFeedback from '../components/PageFeedback.vue'
import PageHeader from '../components/PageHeader.vue'
import { getApiErrorMessage, getApiErrorStatus, useMockApi } from '../api/client'
import { useNodesStore } from '../stores/nodes'
import { useToastStore } from '../stores/toast'

const store = useNodesStore()
const toast = useToastStore()
const router = useRouter()

const onlineCount = computed(() => store.nodes.filter((node) => node.online).length)

const claimModalOpen = ref(false)
const claimNodeId = ref('')
const claimCode = ref('')
const claiming = ref(false)
const claimError = ref('')

onMounted(() => {
  void store.loadNodes().catch(() => undefined)
})

function openClaimModal() {
  claimNodeId.value = ''
  claimCode.value = ''
  claimError.value = ''
  claimModalOpen.value = true
}

function closeClaimModal() {
  if (claiming.value) return
  claimModalOpen.value = false
}

function onClaimCodeInput(event: Event) {
  const target = event.target as HTMLInputElement
  claimCode.value = target.value.toUpperCase().trim()
}

function fillMockClaimExample() {
  claimNodeId.value = 'aq-sim003'
  claimCode.value = 'KLIM2026'
  claimError.value = ''
}

async function submitClaim() {
  const id = claimNodeId.value.trim().toLowerCase()
  const code = claimCode.value.trim().toUpperCase()
  claimNodeId.value = id
  claimCode.value = code
  claimError.value = ''

  if (!id || !code) {
    claimError.value = 'ID node dan kode klaim wajib diisi.'
    return
  }
  if (code.length !== 8) {
    claimError.value = 'Kode klaim harus terdiri dari 8 karakter.'
    return
  }

  claiming.value = true
  try {
    const detail = await store.claimNode(id, code)
    claimModalOpen.value = false
    toast.success(`Akuarium ${detail.name} berhasil ditambahkan ke akun Anda.`)
    void router.push({ name: 'node-detail', params: { id: detail.node_id } })
  } catch (cause) {
    const status = getApiErrorStatus(cause)
    if (status === 404) {
      claimError.value = 'ID node atau kode klaim salah. Periksa kembali stiker pada kotak akuarium.'
    } else if (status === 409) {
      claimError.value = 'Akuarium sudah terdaftar di akun lain.'
    } else if (status === 429) {
      claimError.value = 'Terlalu banyak percobaan gagal, silakan coba lagi nanti.'
    } else if (status === 400) {
      claimError.value = getApiErrorMessage(cause, 'Format ID node atau kode klaim tidak sesuai.')
    } else {
      claimError.value = getApiErrorMessage(cause, 'Gagal menambahkan akuarium. Silakan coba lagi.')
    }
  } finally {
    claiming.value = false
  }
}
</script>

<template>
  <section class="page">
    <PageHeader
      eyebrow="RINGKASAN SISTEM"
      title="Daftar akuarium"
      subtitle="Pantau kondisi kualitas air seluruh akuarium yang terhubung ke akun Anda."
    >
      <template #actions>
        <div class="header-actions-group">
          <div v-if="store.nodes.length" class="online-count">
            <span class="status-dot dot-online" aria-hidden="true"></span>
            <span>{{ onlineCount }} dari {{ store.nodes.length }} akuarium online</span>
          </div>
          <button class="button button-primary" type="button" @click="openClaimModal">
            <Plus :size="18" aria-hidden="true" />
            <span>Tambah Akuarium</span>
          </button>
        </div>
      </template>
    </PageHeader>

    <PageFeedback :loading="store.loading && !store.nodes.length" :error="store.error">
      <EmptyState
        v-if="!store.nodes.length"
        :icon="Waves"
        title="Belum ada akuarium terdaftar"
        description="Tambahkan perangkat akuarium Anda menggunakan ID node dan kode klaim 8 karakter yang tertera pada stiker kotak perangkat."
      >
        <template #action>
          <button class="button button-primary" type="button" @click="openClaimModal">
            <Plus :size="18" aria-hidden="true" />
            <span>Tambah Akuarium</span>
          </button>
        </template>
      </EmptyState>

      <div v-else class="node-grid">
        <NodeCard
          v-for="(node, idx) in store.nodes"
          :key="node.node_id"
          :node="node"
          :index="idx"
        />
      </div>
    </PageFeedback>

    <Teleport to="body">
      <Transition name="dialog-fade">
        <div
          v-if="claimModalOpen"
          class="modal-backdrop"
          role="presentation"
          @click.self="closeClaimModal"
        >
          <div
            class="modal-card"
            role="dialog"
            aria-modal="true"
            aria-labelledby="claim-modal-title"
          >
            <div class="modal-header">
              <div class="modal-icon modal-icon-teal" aria-hidden="true">
                <KeyRound :size="22" />
              </div>
              <button
                class="icon-button"
                type="button"
                aria-label="Tutup dialog klaim"
                :disabled="claiming"
                @click="closeClaimModal"
              >
                <X :size="19" />
              </button>
            </div>

            <h2 id="claim-modal-title" class="modal-title">Tambah Akuarium</h2>
            <p class="modal-description">
              Lihat stiker pada kotak akuarium untuk menemukan ID node dan kode klaim 8 karakter.
            </p>

            <div v-if="useMockApi" class="mock-claim-hint">
              <div class="mock-claim-hint-text">
                <Sparkles :size="17" aria-hidden="true" />
                <span>
                  Contoh simulasi siap klaim: ID <code>aq-sim003</code> · Kode
                  <code>KLIM2026</code>
                </span>
              </div>
              <button class="button button-secondary button-sm" type="button" @click="fillMockClaimExample">
                Isi contoh
              </button>
            </div>

            <form class="claim-form" @submit.prevent="submitClaim">
              <label for="claim-node-id">ID Node</label>
              <input
                id="claim-node-id"
                v-model="claimNodeId"
                placeholder="Contoh: aq-1a2b3c"
                autocomplete="off"
                required
              />

              <label for="claim-code-input">Kode Klaim (8 karakter)</label>
              <input
                id="claim-code-input"
                :value="claimCode"
                class="code-input"
                placeholder="Contoh: A1B2C3D4"
                maxlength="8"
                autocomplete="off"
                required
                @input="onClaimCodeInput"
              />
              <p class="field-hint">Lihat stiker pada kotak akuarium.</p>

              <p v-if="claimError" class="form-error" role="alert">{{ claimError }}</p>

              <div class="modal-actions">
                <button
                  class="button button-quiet"
                  type="button"
                  :disabled="claiming"
                  @click="closeClaimModal"
                >
                  Batal
                </button>
                <button class="button button-primary" type="submit" :disabled="claiming">
                  {{ claiming ? 'Mengklaim...' : 'Klaim Akuarium' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </section>
</template>
