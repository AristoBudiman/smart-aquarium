<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { ArrowRight, Droplets, ShieldCheck } from 'lucide-vue-next'
import ThemeSwitcher from '../components/ThemeSwitcher.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const username = ref('')
const password = ref('')
const validationError = ref('')

async function submit() {
  validationError.value = ''
  if (!username.value.trim() || !password.value) {
    validationError.value = 'Username dan kata sandi wajib diisi.'
    return
  }
  try {
    await auth.login(username.value, password.value)
    if (auth.isAdmin) {
      void router.replace({ name: 'admin' })
      return
    }
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/nodes'
    const safeRedirect =
      redirect.startsWith('/') && !redirect.startsWith('//') && !redirect.startsWith('/admin')
        ? redirect
        : '/nodes'
    void router.replace(safeRedirect)
  } catch {
    // Error disimpan di store auth.
  }
}
</script>

<template>
  <main class="login-page">
    <div class="login-topbar">
      <ThemeSwitcher />
    </div>
    <div class="login-decoration decoration-one" aria-hidden="true"></div>
    <div class="login-decoration decoration-two" aria-hidden="true"></div>
    <section class="login-panel">
      <RouterLink class="brand login-brand" to="/login">
        <span class="brand-mark" aria-hidden="true">
          <Droplets :size="22" :stroke-width="2.2" />
        </span>
        <span>
          <strong>AquaSense</strong>
          <small>SMART AQUARIUM</small>
        </span>
      </RouterLink>
      <p class="eyebrow">DASHBOARD PEMANTAUAN</p>
      <h1>Selamat datang kembali</h1>
      <p class="muted login-intro">Masuk untuk memantau kualitas air dan kendali akuarium Anda.</p>

      <form class="login-form" @submit.prevent="submit">
        <label for="username">Username</label>
        <input
          id="username"
          v-model="username"
          autocomplete="username"
          placeholder="Masukkan username"
          required
        />

        <label for="password">Kata sandi</label>
        <input
          id="password"
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="Masukkan kata sandi"
          required
        />

        <p v-if="validationError || auth.error" class="form-error" role="alert">
          {{ validationError || auth.error }}
        </p>

        <button class="button button-primary login-submit" type="submit" :disabled="auth.submitting">
          <span>{{ auth.submitting ? 'Memproses...' : 'Masuk ke dashboard' }}</span>
          <ArrowRight :size="18" aria-hidden="true" />
        </button>
      </form>

      <p class="auth-switch-text">
        Belum punya akun klien?
        <RouterLink to="/register" class="auth-switch-link">Daftar akun baru</RouterLink>
      </p>

      <p class="login-footnote">
        <ShieldCheck :size="16" class="secure-icon" aria-hidden="true" />
        <span>Koneksi terenkripsi dengan autentikasi JWT</span>
      </p>
    </section>
    <footer class="login-footer">AquaSense · Sistem pemantauan kualitas air cerdas</footer>
  </main>
</template>
