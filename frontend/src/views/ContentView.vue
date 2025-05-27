<template>
  <DashboardLayout>
    <div>
      <div class="flex items-center justify-between mb-8">
        <div class="flex items-center">
          <router-link to="/content" class="text-indigo-600 hover:text-indigo-800 mr-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </router-link>
          <h1 class="text-3xl font-bold text-gray-900">{{ isEditing ? 'Редактирование контента' : 'Просмотр контента' }}</h1>
        </div>
        <div class="flex items-center space-x-4">
          <template v-if="!isEditing && contentStore.currentContent">
            <Button 
              variant="outline" 
              @click="isEditing = true"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
              </svg>
              Редактировать
            </Button>
            <Button 
              variant="primary"
              @click="publishContent"
              :is-loading="isPublishing"
              v-if="contentStore.currentContent && !contentStore.currentContent.is_published"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Опубликовать
            </Button>
            <Button 
              variant="danger"
              @click="showDeleteConfirm = true"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              Удалить
            </Button>
          </template>
          <template v-else-if="isEditing">
            <Button 
              variant="outline" 
              @click="cancelEditing"
            >
              Отмена
            </Button>
            <Button 
              variant="primary"
              @click="saveContent"
              :is-loading="isSaving"
            >
              Сохранить
            </Button>
          </template>
        </div>
      </div>

      <!-- Индикатор загрузки -->
      <div v-if="contentStore.isLoading" class="text-center py-12">
        <svg class="animate-spin h-10 w-10 text-indigo-600 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="mt-4 text-gray-600">Загрузка контента...</p>
      </div>

      <!-- Информация о контенте -->
      <div v-else-if="contentStore.currentContent" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <div class="flex justify-between mb-6">
          <div>
            <div class="flex items-center mb-2">
              <span v-if="contentStore.currentContent.is_generated" class="bg-indigo-100 text-indigo-800 text-xs px-2 py-1 rounded mr-2">
                Сгенерировано
              </span>
              <span :class="[
                'text-xs px-2 py-1 rounded',
                contentStore.currentContent.is_published 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-gray-100 text-gray-800'
              ]">
                {{ contentStore.currentContent.is_published ? 'Опубликовано' : 'Черновик' }}
              </span>
            </div>
            <h2 v-if="!isEditing" class="text-2xl font-semibold text-gray-900">{{ contentStore.currentContent.title }}</h2>
            <input 
              v-else
              v-model="editedContent.title"
              class="w-full text-2xl font-semibold px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent mb-2"
              placeholder="Заголовок"
            />
            <p class="text-gray-500 text-sm">
              Создано: {{ formatDate(contentStore.currentContent.created_at) }}
              <span v-if="contentStore.currentContent.updated_at !== contentStore.currentContent.created_at">
                · Обновлено: {{ formatDate(contentStore.currentContent.updated_at) }}
              </span>
            </p>
          </div>
          <div>
            <Button 
              variant="outline"
              size="sm"
              @click="regenerateContent"
              :is-loading="isRegenerating"
              v-if="contentStore.currentContent.is_generated && !isEditing"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Сгенерировать заново
            </Button>
          </div>
        </div>

        <!-- Содержимое контента -->
        <div v-if="!isEditing" class="prose prose-indigo max-w-none">
          <div v-html="markdownToHtml(contentStore.currentContent?.content_md || '')"></div>
        </div>
        <div v-else>
          <textarea
            v-model="editedContent.content_md"
            rows="20"
            class="w-full px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent font-mono resize-none"
            placeholder="Содержимое контента в Markdown-формате"
          ></textarea>
        </div>
      </div>
      
      <!-- Сообщение об ошибке -->
      <div v-else-if="contentStore.error" class="bg-white shadow-sm rounded-lg p-6 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-red-500 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">Ошибка при загрузке контента</h3>
        <p class="text-gray-600 mb-6">{{ contentStore.error }}</p>
        <router-link to="/content">
          <Button variant="primary">
            Вернуться к списку
          </Button>
        </router-link>
      </div>
      
      <!-- Диалог подтверждения удаления -->
      <div v-if="showDeleteConfirm" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white rounded-lg p-6 max-w-md w-full">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Подтверждение удаления</h3>
          <p class="text-gray-700 mb-6">
            Вы уверены, что хотите удалить этот контент? Это действие нельзя отменить.
          </p>
          <div class="flex justify-end space-x-4">
            <Button
              variant="outline"
              @click="showDeleteConfirm = false"
            >
              Отмена
            </Button>
            <Button
              variant="danger"
              :is-loading="isDeleting"
              @click="deleteContent"
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
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useContentStore } from '@/stores/content'
import type { Content } from '@/stores/content'
import { useToast } from 'vue-toast-notification'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'
import { marked } from 'marked' // Правильный импорт библиотеки marked

const route = useRoute()
const router = useRouter()
const contentStore = useContentStore()
const toast = useToast()

// Состояния
const isEditing = ref(false)
const isSaving = ref(false)
const isPublishing = ref(false)
const isRegenerating = ref(false)
const isDeleting = ref(false)
const showDeleteConfirm = ref(false)

// Получаем ID контента из URL
const contentId = route.params.id as string

// Редактируемые данные
const editedContent = reactive({
  title: '',
  content_md: ''
})

// Функция для преобразования Markdown в HTML
const markdownToHtml = (markdown: string) => {
  try {
    return marked(markdown)
  } catch (error) {
    console.error('Ошибка при преобразовании Markdown в HTML:', error)
    return markdown // Возвращаем исходный текст в случае ошибки
  }
}

// Загружаем контент при монтировании компонента
onMounted(async () => {
  try {
    await contentStore.getContentById(contentId)
    
    // Инициализируем форму редактирования данными из хранилища
    if (contentStore.currentContent) {
      editedContent.title = contentStore.currentContent.title || ''
      editedContent.content_md = contentStore.currentContent.content_md || ''
    }
  } catch (error) {
    console.error('Ошибка при загрузке контента:', error)
    router.push('/content')
  }
})

// Наблюдаем за изменением currentContent в хранилище
watch(() => contentStore.currentContent, (newContent) => {
  if (newContent) {
    editedContent.title = newContent.title || ''
    editedContent.content_md = newContent.content_md || ''
  }
})

// Форматирование даты
const formatDate = (dateString: string) => {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('ru-RU', { 
      day: '2-digit', 
      month: '2-digit', 
      year: 'numeric' 
    })
  } catch (e) {
    return dateString
  }
}

// Обработчики действий
const cancelEditing = () => {
  // Сбрасываем изменения и возвращаемся в режим просмотра
  if (contentStore.currentContent) {
    editedContent.title = contentStore.currentContent.title || ''
    editedContent.content_md = contentStore.currentContent.content_md || ''
  }
  isEditing.value = false
}

const saveContent = async () => {
  if (!editedContent.title.trim()) {
    toast.error('Заголовок не может быть пустым', { duration: 3000 })
    return
  }
  
  if (!editedContent.content_md.trim()) {
    toast.error('Содержимое не может быть пустым', { duration: 3000 })
    return
  }
  
  isSaving.value = true
  
  try {
    await contentStore.updateContent(contentId, {
      title: editedContent.title,
      content_md: editedContent.content_md
    })
    
    toast.success('Контент успешно сохранен', { duration: 3000 })
    isEditing.value = false
  } catch (error) {
    console.error('Ошибка при сохранении контента:', error)
    // Ошибка уже будет показана через перехватчик в хранилище
  } finally {
    isSaving.value = false
  }
}

const publishContent = async () => {
  isPublishing.value = true
  
  try {
    await contentStore.publishContent(contentId)
    toast.success('Контент успешно опубликован', { duration: 3000 })
  } catch (error) {
    console.error('Ошибка при публикации контента:', error)
    // Ошибка уже будет показана через перехватчик в хранилище
  } finally {
    isPublishing.value = false
  }
}

const regenerateContent = async () => {
  if (!contentStore.currentContent) return
  
  isRegenerating.value = true
  
  try {
    // Используем текущий заголовок и содержимое как основу для регенерации
    const response = await contentStore.generateContent({
      prompt: contentStore.currentContent.title || '',
      use_voice_profile: true
    })
    
    // Обновляем только содержимое, сохраняя заголовок
    await contentStore.updateContent(contentId, {
      content_md: response.content_md || response.content || ''
    })
    
    toast.success('Контент успешно обновлен', { duration: 3000 })
  } catch (error) {
    console.error('Ошибка при регенерации контента:', error)
    // Ошибка уже будет показана через перехватчик в хранилище
  } finally {
    isRegenerating.value = false
  }
}

const deleteContent = async () => {
  isDeleting.value = true
  
  try {
    await contentStore.deleteContent(contentId)
    toast.success('Контент успешно удален', { duration: 3000 })
    router.push('/content')
  } catch (error) {
    console.error('Ошибка при удалении контента:', error)
    // Ошибка уже будет показана через перехватчик в хранилище
  } finally {
    isDeleting.value = false
    showDeleteConfirm.value = false
  }
}
</script> 