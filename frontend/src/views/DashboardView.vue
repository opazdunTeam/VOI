<template>
  <DashboardLayout>
    <div>
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Дашборд</h1>
        <p class="text-gray-600">
          Добро пожаловать в ваш личный кабинет. Здесь вы можете управлять контентом и настройками.
        </p>
      </div>
      
      <!-- Предупреждение об отсутствии голосового профиля -->
      <VoiceProfileAlert 
        message="Для полноценного использования системы рекомендуем настроить ваш голосовой профиль." 
      />

      <!-- Статистика и основная информация -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <!-- Карточка статистики контента -->
        <div class="bg-white shadow-sm rounded-lg p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Ваш контент</h2>
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-indigo-50 p-4 rounded-lg text-center">
              <span class="text-2xl font-bold text-indigo-600 block">{{ contentStats.total || 0 }}</span>
              <span class="text-sm text-gray-600">Всего</span>
            </div>
            <div class="bg-green-50 p-4 rounded-lg text-center">
              <span class="text-2xl font-bold text-green-600 block">{{ contentStats.published || 0 }}</span>
              <span class="text-sm text-gray-600">Опубликовано</span>
            </div>
          </div>
          <div class="mt-4 text-center">
            <router-link to="/content" class="text-indigo-600 hover:text-indigo-800 text-sm font-medium">
              Управление контентом →
            </router-link>
          </div>
        </div>
        
        <!-- Карточка голосового профиля -->
        <div class="bg-white shadow-sm rounded-lg p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Голосовой профиль</h2>
          <div class="flex items-center justify-center mb-4">
            <div class="w-12 h-12 rounded-full bg-indigo-100 flex items-center justify-center">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
              </svg>
            </div>
          </div>
          <div class="text-center">
            <p class="text-gray-700 mb-3">
              {{ hasVoiceProfile ? 'Ваш голосовой профиль настроен' : 'Голосовой профиль не настроен' }}
            </p>
            <router-link to="/voice-dna" class="text-indigo-600 hover:text-indigo-800 text-sm font-medium">
              {{ hasVoiceProfile ? 'Обновить профиль' : 'Настроить профиль' }} →
            </router-link>
          </div>
        </div>
        
        <!-- Карточка создания контента -->
        <div class="bg-white shadow-sm rounded-lg p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Быстрые действия</h2>
          <div class="space-y-3">
            <router-link to="/content/new" class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition">
              <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
              </div>
              <span class="text-gray-700">Создать новый контент</span>
            </router-link>
            
            <router-link to="/profile" class="flex items-center p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition">
              <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center mr-3">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <span class="text-gray-700">Управление профилем</span>
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- Последний контент -->
      <div class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-xl font-semibold text-gray-900">Последний контент</h2>
          <router-link to="/content" class="text-indigo-600 hover:text-indigo-800 text-sm font-medium">
            Посмотреть все →
          </router-link>
        </div>
        
        <div v-if="recentContent.length === 0" class="text-center py-8">
          <p class="text-gray-500">У вас пока нет созданного контента</p>
          <router-link to="/content/new" class="mt-2 inline-block text-indigo-600 hover:text-indigo-800">
            Создать первый контент
          </router-link>
        </div>
        
        <div v-else class="space-y-4">
          <div v-for="item in recentContent" :key="item.id" class="border-b border-gray-100 pb-4 last:border-0 last:pb-0">
            <div class="flex justify-between items-start">
              <div>
                <router-link :to="`/content/${item.id}`" class="text-lg font-medium text-gray-900 hover:text-indigo-600">
                  {{ item.title || 'Без заголовка' }}
                </router-link>
                <p class="text-sm text-gray-500 mt-1">{{ new Date(item.created_at).toLocaleDateString() }}</p>
              </div>
              <div class="bg-gray-100 text-gray-700 text-xs px-2 py-1 rounded">
                {{ item.is_published ? 'Опубликован' : 'Черновик' }}
              </div>
            </div>
            <p class="text-gray-700 mt-2 line-clamp-2">{{ item.excerpt }}</p>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import { useContentStore, type Content } from '@/stores/content'
import { useProfileStore } from '@/stores/profile'
import VoiceProfileAlert from '@/components/ui/VoiceProfileAlert.vue'

const contentStore = useContentStore()
const profileStore = useProfileStore()

// Состояния
const recentContent = ref<Content[]>([])
const contentStats = ref({
  total: 0,
  published: 0,
  drafts: 0
})
const hasVoiceProfile = ref(false)
const isLoading = ref(true)

// Получение данных при загрузке компонента
onMounted(async () => {
  isLoading.value = true
  
  try {
    // Получаем данные контента
    let content = []
    try {
      content = await contentStore.getContentList()
      
      // Фильтруем для последних 3 элементов
      recentContent.value = content.slice(0, 3)
      
      // Считаем статистику
      contentStats.value = {
        total: content.length,
        published: content.filter(item => item.is_published).length,
        drafts: content.filter(item => !item.is_published).length
      }
    } catch (contentError) {
      console.error('Ошибка при получении списка контента:', contentError)
      // Уведомление об ошибке только в консоли
    }
    
    // Проверяем наличие голосового профиля
    try {
      const profile = await profileStore.getProfile()
      hasVoiceProfile.value = profileStore.hasVoiceProfile
    } catch (profileError) {
      console.error('Ошибка при получении профиля:', profileError)
      hasVoiceProfile.value = false
    }
  } catch (error) {
    console.error('Ошибка при получении данных:', error)
  } finally {
    isLoading.value = false
  }
})
</script> 