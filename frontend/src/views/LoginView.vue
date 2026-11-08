<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/nodes'
    void router.replace(redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '/nodes')
  } catch {
    // Pesan kegagalan sudah disimpan oleh auth store.
  }
}
</script>

<template>
  <main class="login-page">
    <div class="login-decoration decoration-one"></div>
    <div class="login-decoration decoration-two"></div>
    <section class="login-panel">
      <RouterLink class="brand login-brand" to="/login">
        <span class="brand-mark">A</span>
        <span>
          <strong>AquaSense</strong>
          <small>SMART AQUARIUM</small>
        </span>
      </RouterLink>
      <p class="eyebrow">DASHBOARD MONITORING</p>
      <h1>Selamat datang</h1>
      <p class="muted login-intro">Masuk untuk memantau kondisi akuarium Anda.</p>

      <form class="login-form" @submit.prevent="submit">
        <label for="username">Username</label>
        <input id="username" v-model="username" autocomplete="username" placeholder="Masukkan username" required>
        <label for="password">Kata sandi</label>
        <input id="password" v-model="password" type="password" autocomplete="current-password" placeholder="Masukkan kata sandi" required>
        <p v-if="validationError || auth.error" class="form-error" role="alert">
          {{ validationError || auth.error }}
        </p>
        <button class="button button-primary login-submit" type="submit" :disabled="auth.submitting">
          {{ auth.submitting ? 'Memproses...' : 'Masuk ke dashboard' }}
          <span aria-hidden="true">→</span>
        </button>
      </form>

      <p class="login-footnote"><span class="secure-dot"></span> Koneksi lokal dengan autentikasi JWT</p>
    </section>
    <footer class="login-footer">AI Smart Aquarium <span>·</span> Sistem pemantauan kualitas air</footer>
  </main>
</template>
