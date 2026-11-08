import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/',
      component: () => import('../components/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/nodes' },
        { path: 'nodes', name: 'nodes', component: () => import('../views/NodeListView.vue') },
        { path: 'nodes/:id', name: 'node-detail', component: () => import('../views/NodeDetailView.vue') },
        { path: 'nodes/:id/history', name: 'node-history', component: () => import('../views/HistoryView.vue') },
        { path: 'nodes/:id/logs', name: 'node-logs', component: () => import('../views/LogsView.vue') },
        { path: 'settings', name: 'settings', component: () => import('../views/SettingsView.vue') },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const hasToken = Boolean(localStorage.getItem('token'))
  if (to.meta.requiresAuth && !hasToken) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && hasToken) return { name: 'nodes' }
  return true
})

export default router
