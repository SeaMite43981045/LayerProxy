import LoginView from '@/views/LoginView.vue'
import SetupView from '@/views/SetupView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      name: 'Login',
      path: '/login',
      component: LoginView,
    },
    {
      name: 'Setup',
      path: '/setup',
      component: SetupView,
    },
  ],
})

export default router
