<template>
  <aside class="w-64 bg-gradient-to-b from-indigo-800 to-indigo-600 text-white h-screen fixed left-0 top-0 overflow-y-auto">
    <div class="p-6">
      <router-link to="/" class="flex items-center justify-center mb-8">
        <h1 class="text-3xl font-bold">VOY</h1>
      </router-link>

      <nav class="mt-8">
        <nav-item 
          to="/dashboard" 
          :icon="DashboardIcon" 
          label="Дашборд" 
          :active="route.path === '/dashboard'" 
        />
        <nav-item 
          to="/content" 
          :icon="ContentIcon" 
          label="Контент" 
          :active="route.path.startsWith('/content')" 
        />
        <nav-item 
          to="/voice-dna" 
          :icon="VoiceIcon" 
          label="Голос ДНК" 
          :active="route.path.startsWith('/voice-dna')" 
        />
        <nav-item 
          to="/profile" 
          :icon="ProfileIcon" 
          label="Профиль" 
          :active="route.path.startsWith('/profile')" 
        />
      </nav>
    </div>

    <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-white/10">
      <div class="space-y-4">
        <router-link to="/profile" class="flex items-center hover:bg-white/10 p-2 rounded-lg transition-colors">
          <div class="bg-indigo-900 w-10 h-10 rounded-full flex items-center justify-center mr-3">
            <span class="text-lg font-semibold">
              {{ getUserInitial }}
            </span>
          </div>
          <div class="flex flex-col">
            <span class="text-sm font-medium">
              {{ isLoading ? 'Загрузка...' : userData?.full_name || 'Пользователь' }}
            </span>
            <span class="text-xs text-white/70">
              {{ isLoading ? '' : userData?.email || 'user@example.com' }}
            </span>
          </div>
        </router-link>

        <button 
          @click="handleLogout" 
          class="w-full flex items-center text-white/70 hover:text-white hover:bg-white/10 p-2 rounded-lg transition-colors"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Выйти
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavItem from './NavItem.vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// Состояние
const isLoading = ref(true)
const userData = ref<any>(null)

// Иконки для меню (используем компоненты SVG)
const DashboardIcon = `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
</svg>`

const ContentIcon = `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
</svg>`

const VoiceIcon = `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
</svg>`

const ProfileIcon = `<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
</svg>`

// Получение пользователя
onMounted(async () => {
  try {
    await authStore.getCurrentUser()
    userData.value = authStore.user
  } catch (error) {
    console.error('Ошибка при получении пользователя:', error)
  } finally {
    isLoading.value = false
  }
})

// Получаем первую букву имени пользователя
const getUserInitial = computed(() => {
  if (isLoading.value || !userData.value?.full_name) return '?'
  return userData.value.full_name.charAt(0).toUpperCase()
})

const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('Ошибка при выходе:', error)
  }
}
</script> 