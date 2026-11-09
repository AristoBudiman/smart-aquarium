import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMe, login as loginRequest, register as registerRequest } from '../api/auth'
import { getApiErrorMessage } from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') ?? '')
  const username = ref(localStorage.getItem('auth_username') ?? '')
  const isAdmin = ref(localStorage.getItem('auth_is_admin') === '1')
  const telegramLinked = ref(false)
  const initialized = ref(false)
  const submitting = ref(false)
  const error = ref('')
  let restorePromise: Promise<void> | null = null

  async function fetchMe() {
    const me = await getMe()
    username.value = me.username
    isAdmin.value = me.is_admin
    telegramLinked.value = me.telegram_linked
    initialized.value = true
    localStorage.setItem('auth_username', me.username)
    localStorage.setItem('auth_is_admin', me.is_admin ? '1' : '0')
    return me
  }

  async function ensureRestored() {
    if (!token.value) {
      initialized.value = true
      return
    }
    if (initialized.value) return
    if (!restorePromise) {
      restorePromise = fetchMe()
        .then(() => undefined)
        .catch(() => {
          clear()
        })
        .finally(() => {
          restorePromise = null
        })
    }
    await restorePromise
  }

  async function login(user: string, password: string) {
    submitting.value = true
    error.value = ''
    try {
      const response = await loginRequest(user.trim(), password)
      token.value = response.token
      localStorage.setItem('token', response.token)
      await fetchMe()
    } catch (cause) {
      error.value = getApiErrorMessage(cause, 'Login gagal. Periksa kembali username dan kata sandi.')
      throw cause
    } finally {
      submitting.value = false
    }
  }

  async function register(user: string, password: string) {
    submitting.value = true
    error.value = ''
    try {
      const response = await registerRequest(user.trim().toLowerCase(), password)
      token.value = response.token
      localStorage.setItem('token', response.token)
      await fetchMe()
    } catch (cause) {
      error.value = getApiErrorMessage(cause, 'Pendaftaran gagal. Silakan coba lagi.')
      throw cause
    } finally {
      submitting.value = false
    }
  }

  function clear() {
    token.value = ''
    username.value = ''
    isAdmin.value = false
    telegramLinked.value = false
    initialized.value = false
    error.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('auth_username')
    localStorage.removeItem('auth_is_admin')
  }

  return {
    token,
    username,
    isAdmin,
    telegramLinked,
    initialized,
    submitting,
    error,
    fetchMe,
    ensureRestored,
    login,
    register,
    clear,
  }
})
