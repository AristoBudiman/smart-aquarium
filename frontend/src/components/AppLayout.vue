<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import {
  Droplets,
  Info,
  LayoutGrid,
  LogOut,
  Settings,
  ShieldCheck,
  User,
} from 'lucide-vue-next'
import { useMockApi } from '../api/client'
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
  nodes.reset()
  void router.replace({ name: 'login' })
}
</script>

<template>
  <div class="app-shell">
    <header class="topbar">
      <RouterLink
        class="brand"
        :to="auth.isAdmin ? '/admin' : '/nodes'"
        aria-label="AquaSense Smart Aquarium, beranda"
      >
        <span class="brand-mark" aria-hidden="true">
          <Droplets :size="22" :stroke-width="2.2" />
        </span>
        <span>
          <strong>AquaSense</strong>
          <small>SMART AQUARIUM</small>
        </span>
      </RouterLink>

      <div class="topbar-actions">
        <span v-if="useMockApi" class="connection-label">
          <span class="status-dot dot-online" aria-hidden="true"></span>
          Mode simulasi aktif
        </span>
        <span v-if="auth.username" class="user-badge">
          <ShieldCheck v-if="auth.isAdmin" :size="17" aria-hidden="true" />
          <User v-else :size="17" aria-hidden="true" />
          <span>{{ auth.username }}</span>
        </span>
        <button class="button button-quiet logout-button" type="button" @click="logout">
          <LogOut :size="18" aria-hidden="true" />
          <span>Keluar</span>
        </button>
      </div>
    </header>

    <div class="workspace">
      <aside class="sidebar">
        <p class="nav-caption">{{ auth.isAdmin ? 'PANEL ADMIN' : 'MENU UTAMA' }}</p>
        <nav class="side-nav" aria-label="Navigasi utama">
          <template v-if="auth.isAdmin">
            <RouterLink to="/admin" class="nav-link" active-class="is-active">
              <ShieldCheck :size="19" aria-hidden="true" />
              <span>Kelola Node</span>
            </RouterLink>
          </template>
          <template v-else>
            <RouterLink to="/nodes" class="nav-link" active-class="is-active">
              <LayoutGrid :size="19" aria-hidden="true" />
              <span>Daftar Node</span>
            </RouterLink>
            <RouterLink to="/settings" class="nav-link" active-class="is-active">
              <Settings :size="19" aria-hidden="true" />
              <span>Pengaturan</span>
            </RouterLink>
          </template>
        </nav>

        <div v-if="useMockApi" class="sidebar-note">
          <Info :size="18" class="note-icon-svg" aria-hidden="true" />
          <p>Data saat ini berasal dari mode simulasi lokal.</p>
        </div>
      </aside>

      <main class="main-content">
        <RouterView v-slot="{ Component }">
          <Transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>

    <nav class="mobile-nav" aria-label="Navigasi seluler">
      <template v-if="auth.isAdmin">
        <RouterLink to="/admin" class="mobile-nav-link" active-class="is-active">
          <ShieldCheck :size="20" aria-hidden="true" />
          <span>Kelola Node</span>
        </RouterLink>
      </template>
      <template v-else>
        <RouterLink to="/nodes" class="mobile-nav-link" active-class="is-active">
          <LayoutGrid :size="20" aria-hidden="true" />
          <span>Node</span>
        </RouterLink>
        <RouterLink to="/settings" class="mobile-nav-link" active-class="is-active">
          <Settings :size="20" aria-hidden="true" />
          <span>Pengaturan</span>
        </RouterLink>
      </template>
      <button class="mobile-nav-link" type="button" @click="logout">
        <LogOut :size="20" aria-hidden="true" />
        <span>Keluar</span>
      </button>
    </nav>
  </div>
</template>
