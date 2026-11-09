import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import { useNodesStore } from './stores/nodes'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

window.addEventListener('auth:unauthorized', () => {
  useAuthStore(pinia).clear()
  useNodesStore(pinia).reset()
  void router.replace({ name: 'login' })
})

app.mount('#app')
