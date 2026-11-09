<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { ArrowRight, Droplets, ShieldCheck } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const auth = useAuthStore()
const toast = useToastStore()
const router = useRouter()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const validationError = ref('')

const usernameRegex = /^[a-z0-9_.-]{3,32}$/

async function submit() {
  validationError.value = ''
  auth.error = ''
  const normalized = username.value.trim().toLowerCase()
  username.value = normalized

  if (!usernameRegex.test(normalized)) {
    validationError.value =
      'Username harus 3–32 karakter dan hanya menggunakan huruf kecil a-z, angka 0-9, _, ., atau -.'
    return
  }
  if (password.value.length < 8) {
    validationError.value = 'Kata sandi minimal 8 karakter.'
    return
  }
  if (password.value !== confirmPassword.value) {
    validationError.value = 'Konfirmasi kata sandi tidak cocok.'
    return
  }

  try {
    await auth.register(normalized, password.value)
    toast.success(`Akun ${auth.username} berhasil dibuat.`)
    void router.replace({ name: 'nodes' })
  } catch {
    // Error API disimpan di store auth.
  }
}
</script>

<template>
  <main class="login-page">
    <div class="login-decoration decoration-one" aria-hidden="true"></div>
    <div class="login-decoration decoration-two" aria-hidden="true"></div>
    <section class="login-panel">
      <RouterLink class="brand login-brand" to="/register">
        <span class="brand-mark" aria-hidden="true">
          <Droplets :size="22" :stroke-width="2.2" />
        </span>
        <span>
          <strong>AquaSense</strong>
          <small>SMART AQUARIUM</small>
        </span>
      </RouterLink>
      <p class="eyebrow">REGISTRASI KLIEN</p>
      <h1>Buat akun baru</h1>
      <p class="muted login-intro">
        Daftarkan akun untuk mengklaim dan memantau perangkat akuarium milik Anda.
      </p>

      <form class="login-form" @submit.prevent="submit">
        <label for="reg-username">Username</label>
        <input
          id="reg-username"
          v-model="username"
          autocomplete="username"
          placeholder="Contoh: budi_aquarium"
          required
          @input="username = username.toLowerCase()"
        />
        <p class="field-hint">3–32 karakter, huruf kecil a-z, angka 0-9, _, ., atau -.</p>

        <label for="reg-password">Kata sandi</label>
        <input
          id="reg-password"
          v-model="password"
          type="password"
          autocomplete="new-password"
          placeholder="Minimal 8 karakter"
          required
        />

        <label for="reg-confirm">Ulangi kata sandi</label>
        <input
          id="reg-confirm"
          v-model="confirmPassword"
          type="password"
          autocomplete="new-password"
          placeholder="Ketik ulang kata sandi"
          required
        />

        <p v-if="validationError || auth.error" class="form-error" role="alert">
          {{ validationError || auth.error }}
        </p>

        <button class="button button-primary login-submit" type="submit" :disabled="auth.submitting">
          <span>{{ auth.submitting ? 'Mendaftarkan...' : 'Daftar & masuk' }}</span>
          <ArrowRight :size="18" aria-hidden="true" />
        </button>
      </form>

      <p class="auth-switch-text">
        Sudah punya akun?
        <RouterLink to="/login" class="auth-switch-link">Masuk di sini</RouterLink>
      </p>

      <p class="login-footnote">
        <ShieldCheck :size="16" class="secure-icon" aria-hidden="true" />
        <span>Akun klien hanya dapat mengakses akuarium yang diklaim</span>
      </p>
    </section>
    <footer class="login-footer">AquaSense · Sistem pemantauan kualitas air cerdas</footer>
  </main>
</template>
