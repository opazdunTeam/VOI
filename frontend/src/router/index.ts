import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Маршруты
const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false, guestOnly: true }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false, guestOnly: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/RegisterView.vue'),
    meta: { requiresAuth: false, guestOnly: true }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/content',
    name: 'content-list',
    component: () => import('@/views/ContentListView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/content/new',
    name: 'content-new',
    component: () => import('@/views/ContentNewView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/content/:id',
    name: 'content-view',
    component: () => import('@/views/ContentView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/voice-dna',
    name: 'voice-dna',
    component: () => import('@/views/VoiceDNAView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue')
  }
]

// Создание маршрутизатора
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Защита маршрутов
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Список страниц для гостей, где не нужно проверять авторизацию
  const guestPages = ['login', 'register', 'home']
  const isGuestPage = guestPages.includes(to.name as string)
  
  // Проверяем авторизацию только если это не страница для гостей
  if (!authStore.isAuthenticated && !isGuestPage) {
    try {
      await authStore.getCurrentUser()
    } catch (error) {
      // Ошибка при получении пользователя обрабатывается в хранилище,
      // здесь ничего не делаем, так как isAuthenticated станет false
    }
  }
  
  // Проверяем требования маршрута
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // Перенаправляем неавторизованных пользователей на страницу входа
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.guestOnly && authStore.isAuthenticated) {
    // Перенаправляем авторизованных пользователей на дашборд, если они пытаются открыть страницу для гостей
    next({ name: 'dashboard' })
  } else {
    // Разрешаем переход
    next()
  }
})

export default router
