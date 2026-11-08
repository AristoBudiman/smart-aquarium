import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

window.addEventListener('auth:unauthorized', () => {
  useAuthStore(pinia).clear()
  void router.replace({ name: 'login' })
})

app.mount('#app')
