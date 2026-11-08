<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import PageFeedback from '../components/PageFeedback.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { useNodesStore } from '../stores/nodes'
import { formatNumber, formatWib } from '../utils/format'

const store = useNodesStore()
const onlineCount = computed(() => store.nodes.filter((node) => node.online).length)

onMounted(() => {
  void store.loadNodes().catch(() => undefined)
})
</script>

<template>
  <section class="page">
    <div class="page-heading heading-with-action">
      <div>
        <p class="eyebrow">RINGKASAN SISTEM</p>
        <h1>Daftar node</h1>
        <p class="muted">Pantau kondisi seluruh akuarium yang terhubung.</p>
      </div>
      <div class="online-count"><span></span>{{ onlineCount }} dari {{ store.nodes.length }} node online</div>
    </div>

    <PageFeedback :loading="store.loading && !store.nodes.length" :error="store.error">
      <div v-if="!store.nodes.length" class="empty-state">
        <span class="empty-icon">⌁</span>
        <h2>Belum ada node</h2>
        <p>Node akan muncul otomatis setelah perangkat terhubung.</p>
      </div>
      <div v-else class="node-grid">
        <article v-for="node in store.nodes" :key="node.node_id" class="node-card">
          <div class="node-card-top">
            <div class="node-symbol" :class="{ 'node-symbol-offline': !node.online }">⌁</div>
            <StatusBadge :status="node.status" :offline="!node.online" />
          </div>
          <h2>{{ node.name }}</h2>
          <p class="node-id">{{ node.node_id }}</p>
          <div class="node-metrics">
            <div>
              <span class="metric-label">Suhu air</span>
              <strong>{{ formatNumber(node.temp) }}<small>°C</small></strong>
            </div>
            <div>
              <span class="metric-label">Kekeruhan</span>
              <strong>{{ formatNumber(node.ntu) }}<small> NTU</small></strong>
            </div>
          </div>
          <div class="node-card-bottom">
            <span>{{ node.online ? 'Terakhir diperbarui' : 'Terakhir terlihat' }}</span>
            <time>{{ formatWib(node.last_seen) }}</time>
          </div>
          <RouterLink class="card-link" :to="{ name: 'node-detail', params: { id: node.node_id } }">
            Lihat detail <span aria-hidden="true">→</span>
          </RouterLink>
        </article>
      </div>
    </PageFeedback>
  </section>
</template>
