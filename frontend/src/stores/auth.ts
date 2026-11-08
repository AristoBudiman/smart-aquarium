import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginRequest } from '../api/auth'
import { getApiErrorMessage } from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') ?? '')
  const submitting = ref(false)
  const error = ref('')

  async function login(username: string, password: string) {
    submitting.value = true
    error.value = ''
    try {
      const response = await loginRequest(username, password)
      token.value = response.token
      localStorage.setItem('token', response.token)
    } catch (cause) {
      error.value = getApiErrorMessage(cause, 'Login gagal. Silakan coba lagi.')
      throw cause
    } finally {
      submitting.value = false
    }
  }

  function clear() {
    token.value = ''
    localStorage.removeItem('token')
  }

  return { token, submitting, error, login, clear }
})
