import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'dashboard', component: () => import('../views/Dashboard.vue') },
  { path: '/users', name: 'users', component: () => import('../views/Users.vue') },
  { path: '/problems', name: 'problems', component: () => import('../views/Problems.vue') },
  { path: '/submissions', name: 'submissions', component: () => import('../views/Submissions.vue') },
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes,
})

export default router

