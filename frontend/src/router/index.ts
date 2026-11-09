import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/',
      component: () => import('../components/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/nodes' },
        {
          path: 'admin',
          name: 'admin',
          component: () => import('../views/AdminView.vue'),
          meta: { requiresAdmin: true },
        },
        {
          path: 'nodes',
          name: 'nodes',
          component: () => import('../views/NodeListView.vue'),
          meta: { requiresClient: true },
        },
        {
          path: 'nodes/:id',
          name: 'node-detail',
          component: () => import('../views/NodeDetailView.vue'),
          meta: { requiresClient: true },
        },
        {
          path: 'nodes/:id/history',
          name: 'node-history',
          component: () => import('../views/HistoryView.vue'),
          meta: { requiresClient: true },
        },
        {
          path: 'nodes/:id/logs',
          name: 'node-logs',
          component: () => import('../views/LogsView.vue'),
          meta: { requiresClient: true },
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('../views/SettingsView.vue'),
          meta: { requiresClient: true },
        },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (auth.token) {
    await auth.ensureRestored()
  }

  if (to.meta.requiresAuth && !auth.token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && auth.token) {
    return { name: auth.isAdmin ? 'admin' : 'nodes' }
  }

  if (auth.token && auth.isAdmin && to.name !== 'admin') {
    return { name: 'admin' }
  }

  if (auth.token && !auth.isAdmin && to.meta.requiresAdmin) {
    return { name: 'nodes' }
  }

  return true
})

export default router
