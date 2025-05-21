<template>
  <DashboardLayout>
    <div>
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Мой контент</h1>
        <router-link to="/content/new">
          <Button variant="primary">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Создать новый контент
          </Button>
        </router-link>
      </div>

      <!-- Фильтры и поиск -->
      <div class="bg-white shadow-sm rounded-lg p-4 mb-8">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div class="flex md:items-center gap-4 flex-col md:flex-row">
            <div>
              <label for="contentType" class="block text-sm font-medium text-gray-700 mb-1">
                Тип контента
              </label>
              <select
                id="contentType"
                v-model="filters.contentType"
                class="w-full md:w-auto px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              >
                <option value="all">Все</option>
                <option value="generated">Сгенерированные</option>
                <option value="manual">Созданные вручную</option>
              </select>
            </div>
            <div>
              <label for="status" class="block text-sm font-medium text-gray-700 mb-1">
                Статус
              </label>
              <select
                id="status"
                v-model="filters.status"
                class="w-full md:w-auto px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              >
                <option value="all">Все</option>
                <option value="published">Опубликованные</option>
                <option value="draft">Черновики</option>
              </select>
            </div>
          </div>
          <div class="flex-grow md:max-w-md">
            <label for="search" class="block text-sm font-medium text-gray-700 mb-1">
              Поиск
            </label>
            <div class="relative rounded-md shadow-sm">
              <input
                type="text"
                id="search"
                v-model="filters.search"
                class="focus:ring-indigo-500 focus:border-indigo-500 block w-full px-4 py-2 pr-10 border-gray-300 rounded-md"
                placeholder="Поиск по заголовку..."
              />
              <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Список контента -->
      <div v-if="isLoading" class="text-center py-12">
        <svg class="animate-spin h-10 w-10 text-indigo-600 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="mt-4 text-gray-600">Загрузка контента...</p>
      </div>

      <div v-else-if="filteredContent.length === 0" class="bg-white shadow-sm rounded-lg py-12 px-4 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-gray-400 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">Контент не найден</h3>
        <p class="text-gray-600 mb-6">
          {{ 
            filters.search 
              ? 'Не удалось найти контент по вашему запросу. Попробуйте изменить параметры поиска.' 
              : 'У вас пока нет созданного контента. Создайте новый, чтобы начать работу.'
          }}
        </p>
        <router-link to="/content/new">
          <Button variant="primary">
            Создать контент
          </Button>
        </router-link>
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="item in filteredContent" 
          :key="item.id"
          class="bg-white shadow-sm rounded-lg p-6 transition-shadow hover:shadow-md"
        >
          <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <div class="flex items-center mb-2">
                <span v-if="item.isGenerated" class="bg-indigo-100 text-indigo-800 text-xs px-2 py-1 rounded mr-2">
                  Сгенерировано
                </span>
                <span :class="[
                  'text-xs px-2 py-1 rounded',
                  item.isPublished 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-gray-100 text-gray-800'
                ]">
                  {{ item.isPublished ? 'Опубликовано' : 'Черновик' }}
                </span>
              </div>
              <router-link :to="`/content/${item.id}`" class="text-xl font-semibold text-gray-900 hover:text-indigo-600">
                {{ item.title }}
              </router-link>
              <p class="text-gray-600 mt-2 line-clamp-2">{{ item.excerpt }}</p>
              <p class="text-gray-500 text-sm mt-2">
                Создано: {{ item.createdAt }} 
                <span v-if="item.updatedAt !== item.createdAt">
                  · Обновлено: {{ item.updatedAt }}
                </span>
              </p>
            </div>
            <div class="flex items-center space-x-3">
              <router-link :to="`/content/${item.id}`" class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                Просмотреть
              </router-link>
              <button 
                @click="deleteContent(item.id)"
                class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-red-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                Удалить
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Пагинация (упрощенная) -->
      <div v-if="filteredContent.length > 0" class="flex justify-center items-center mt-8">
        <button
          :disabled="currentPage === 1"
          @click="currentPage--"
          :class="[
            'inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md mr-2',
            currentPage === 1
              ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
              : 'bg-white text-gray-700 hover:bg-gray-50'
          ]"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <span class="text-gray-700">Страница {{ currentPage }} из {{ totalPages }}</span>
        <button
          :disabled="currentPage === totalPages"
          @click="currentPage++"
          :class="[
            'inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md ml-2',
            currentPage === totalPages
              ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
              : 'bg-white text-gray-700 hover:bg-gray-50'
          ]"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
      
      <!-- Диалог подтверждения удаления -->
      <div v-if="contentToDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white rounded-lg p-6 max-w-md w-full">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Подтверждение удаления</h3>
          <p class="text-gray-700 mb-6">
            Вы уверены, что хотите удалить этот контент? Это действие нельзя отменить.
          </p>
          <div class="flex justify-end space-x-4">
            <Button
              variant="outline"
              @click="contentToDelete = null"
            >
              Отмена
            </Button>
            <Button
              variant="danger"
              :is-loading="isDeleting"
              @click="confirmDelete"
            >
              Удалить контент
            </Button>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'

// Состояния
const isLoading = ref(true)
const isDeleting = ref(false)
const contentToDelete = ref<string | null>(null)
const currentPage = ref(1)
const totalPages = ref(1)

// Фильтры
const filters = reactive({
  contentType: 'all',
  status: 'all',
  search: ''
})

// Мок-данные для примера
const contentList = ref([
  {
    id: '1',
    title: 'Будущее удаленной работы в технологических компаниях',
    excerpt: 'Исследование тенденций перехода к удаленной работе в технологическом секторе и его влияние на продуктивность и культуру компаний.',
    createdAt: '10.06.2023',
    updatedAt: '12.06.2023',
    isPublished: true,
    isGenerated: true
  },
  {
    id: '2',
    title: 'Влияние искусственного интеллекта на творческие профессии',
    excerpt: 'Анализ того, как ИИ может дополнять творческую работу, а не заменять ее, и как креативщики адаптируются к новым технологиям.',
    createdAt: '01.06.2023',
    updatedAt: '01.06.2023',
    isPublished: true,
    isGenerated: true
  },
  {
    id: '3',
    title: 'Стратегии эффективного управления проектами в 2023 году',
    excerpt: 'Обзор современных подходов к управлению проектами с учетом гибридных команд и возросшей важности цифровых инструментов.',
    createdAt: '20.05.2023',
    updatedAt: '25.05.2023',
    isPublished: false,
    isGenerated: false
  },
  {
    id: '4',
    title: 'Эмоциональный интеллект в лидерстве: почему это важно',
    excerpt: 'Исследование роли эмоционального интеллекта в эффективном лидерстве и как развивать эти навыки в современной бизнес-среде.',
    createdAt: '15.05.2023',
    updatedAt: '16.05.2023',
    isPublished: true,
    isGenerated: true
  },
  {
    id: '5',
    title: 'Устойчивое развитие и корпоративная ответственность',
    excerpt: 'Как компании могут внедрять принципы устойчивого развития в свои бизнес-модели и какие выгоды это приносит.',
    createdAt: '05.05.2023',
    updatedAt: '07.05.2023',
    isPublished: false,
    isGenerated: true
  }
])

// Фильтрация контента
const filteredContent = computed(() => {
  return contentList.value.filter(item => {
    // Фильтр по типу контента
    if (filters.contentType === 'generated' && !item.isGenerated) return false
    if (filters.contentType === 'manual' && item.isGenerated) return false
    
    // Фильтр по статусу
    if (filters.status === 'published' && !item.isPublished) return false
    if (filters.status === 'draft' && item.isPublished) return false
    
    // Поиск по заголовку
    if (filters.search && !item.title.toLowerCase().includes(filters.search.toLowerCase())) return false
    
    return true
  })
})

// Загрузка данных при монтировании компонента
onMounted(async () => {
  try {
    // В реальном приложении здесь был бы запрос к API для получения списка контента
    // В этом примере мы уже используем мок-данные
    
    // Имитация загрузки
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Устанавливаем количество страниц
    totalPages.value = Math.ceil(contentList.value.length / 10) // 10 элементов на страницу
    
    isLoading.value = false
  } catch (error) {
    console.error('Ошибка при загрузке списка контента:', error)
    isLoading.value = false
  }
})

// Удаление контента
const deleteContent = (id: string) => {
  contentToDelete.value = id
}

const confirmDelete = async () => {
  if (!contentToDelete.value) return
  
  isDeleting.value = true
  
  try {
    // В реальном приложении здесь был бы запрос к API для удаления контента
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Удаляем элемент из списка
    const index = contentList.value.findIndex(item => item.id === contentToDelete.value)
    if (index !== -1) {
      contentList.value.splice(index, 1)
    }
    
    // Обновляем количество страниц
    totalPages.value = Math.ceil(contentList.value.length / 10)
    
    // Закрываем диалог
    contentToDelete.value = null
  } catch (error) {
    console.error('Ошибка при удалении контента:', error)
  } finally {
    isDeleting.value = false
  }
}
</script> 