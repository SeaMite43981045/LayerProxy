import AppLayout from '@/components/AppLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import LoginView from '@/views/LoginView.vue'
import LoggingView from '@/views/LoggingView.vue'
import SettingsView from '@/views/SettingsView.vue'
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
    {
      path: '/',
      component: AppLayout,
      redirect: { name: 'Dashboard' },
      children: [
        {
          name: 'Dashboard',
          path: 'dashboard',
          component: DashboardView,
        },
        {
          name: 'Logging',
          path: 'logging',
          component: LoggingView,
        },
        {
          name: 'Settings',
          path: 'settings',
          component: SettingsView,
        },
      ],
    },
  ],
})

export default router
