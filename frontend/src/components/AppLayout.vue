<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { connectRealtime } from '../api/ws'
import { useAuthStore } from '../stores/auth'
import { useNodesStore } from '../stores/nodes'

const auth = useAuthStore()
const nodes = useNodesStore()
const router = useRouter()
let disconnect: (() => void) | undefined

onMounted(() => {
  disconnect = connectRealtime(nodes.applyRealtimeEvent)
})

onUnmounted(() => disconnect?.())

function logout() {
  auth.clear()
  void router.replace({ name: 'login' })
}
</script>

<template>
  <div class="app-shell">
    <header class="topbar">
      <RouterLink class="brand" to="/nodes" aria-label="AI Smart Aquarium, beranda">
        <span class="brand-mark">A</span>
        <span>
          <strong>AquaSense</strong>
          <small>SMART AQUARIUM</small>
        </span>
      </RouterLink>
      <div class="topbar-actions">
        <span class="connection-label"><i></i> Data simulasi aktif</span>
        <button class="button button-quiet logout-button" type="button" @click="logout">Keluar</button>
      </div>
    </header>

    <div class="workspace">
      <aside class="sidebar">
        <p class="nav-caption">MENU UTAMA</p>
        <nav class="side-nav" aria-label="Navigasi utama">
          <RouterLink to="/nodes" class="nav-link" active-class="is-active">
            <span class="nav-icon">▦</span> Daftar node
          </RouterLink>
          <RouterLink to="/settings" class="nav-link" active-class="is-active">
            <span class="nav-icon">⚙</span> Pengaturan
          </RouterLink>
        </nav>
        <div class="sidebar-note">
          <span class="note-icon">i</span>
          <p>Data saat ini berasal dari simulator lokal.</p>
        </div>
      </aside>

      <main class="main-content">
        <RouterView />
      </main>
    </div>

    <nav class="mobile-nav" aria-label="Navigasi seluler">
      <RouterLink to="/nodes" class="mobile-nav-link" active-class="is-active">
        <span>▦</span>Node
      </RouterLink>
      <RouterLink to="/settings" class="mobile-nav-link" active-class="is-active">
        <span>⚙</span>Pengaturan
      </RouterLink>
      <button class="mobile-nav-link" type="button" @click="logout">
        <span>↪</span>Keluar
      </button>
    </nav>
  </div>
</template>
