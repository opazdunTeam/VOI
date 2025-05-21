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
          <template v-if="!isEditing">
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
              v-if="!content.isPublished"
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
          <template v-else>
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

      <!-- Информация о контенте -->
      <div class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <div class="flex justify-between mb-6">
          <div>
            <div class="flex items-center mb-2">
              <span v-if="content.isGenerated" class="bg-indigo-100 text-indigo-800 text-xs px-2 py-1 rounded mr-2">
                Сгенерировано
              </span>
              <span :class="[
                'text-xs px-2 py-1 rounded',
                content.isPublished 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-gray-100 text-gray-800'
              ]">
                {{ content.isPublished ? 'Опубликовано' : 'Черновик' }}
              </span>
            </div>
            <h2 v-if="!isEditing" class="text-2xl font-semibold text-gray-900">{{ content.title }}</h2>
            <input 
              v-else
              v-model="editedContent.title"
              class="w-full text-2xl font-semibold px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent mb-2"
              placeholder="Заголовок"
            />
            <p class="text-gray-500 text-sm">
              Создано: {{ content.createdAt }} 
              <span v-if="content.updatedAt !== content.createdAt">
                · Обновлено: {{ content.updatedAt }}
              </span>
            </p>
          </div>
          <div>
            <Button 
              variant="outline"
              size="sm"
              @click="regenerateContent"
              :is-loading="isRegenerating"
              v-if="content.isGenerated && !isEditing"
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
          <div v-html="content.content"></div>
        </div>
        <div v-else>
          <textarea
            v-model="editedContent.content"
            rows="20"
            class="w-full px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent font-mono"
            placeholder="Содержимое контента в HTML-формате"
          ></textarea>
        </div>
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
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'

const route = useRoute()
const router = useRouter()

// Состояния
const isLoading = ref(true)
const isEditing = ref(false)
const isSaving = ref(false)
const isPublishing = ref(false)
const isRegenerating = ref(false)
const isDeleting = ref(false)
const showDeleteConfirm = ref(false)

// Мок-данные для примера
const content = reactive({
  id: '1',
  title: 'Будущее удаленной работы в технологических компаниях',
  content: `
    <h1>Будущее удаленной работы в технологических компаниях</h1>
    <p>Удаленная работа полностью изменила функционирование технологических компаний. От распределенных команд до виртуальных рабочих пространств, мы наблюдаем революцию в том, как организуется работа.</p>
    <h2>Преимущества удаленной работы</h2>
    <p>Удаленная работа предлагает множество преимуществ как для сотрудников, так и для работодателей. Сотрудники получают гибкость в организации своего времени и рабочего пространства, что может привести к лучшему балансу между работой и личной жизнью.</p>
    <p>Для компаний удаленная работа означает возможность привлекать таланты со всего мира, не ограничиваясь географическим местоположением офиса. Это особенно важно в условиях глобальной конкуренции за квалифицированных специалистов.</p>
    <h2>Вызовы удаленной работы</h2>
    <p>Несмотря на преимущества, удаленная работа также создает ряд вызовов. Одним из основных является поддержание эффективной коммуникации и командного духа в виртуальной среде.</p>
    <p>Технологические компании разрабатывают и внедряют новые инструменты для решения этих проблем, от платформ для видеоконференций до сложных систем управления проектами и асинхронной коммуникации.</p>
    <h2>Гибридные модели работы</h2>
    <p>В будущем, вероятно, наиболее популярными станут гибридные модели работы, сочетающие элементы как удаленной, так и офисной работы. Эти модели позволяют использовать преимущества обоих подходов, адаптируясь к потребностям как компаний, так и сотрудников.</p>
    <blockquote>
      <p>Будущее работы не ограничивается выбором между удаленной и офисной работой. Речь идет о создании гибких, адаптивных моделей, которые работают для каждой конкретной компании и ее команды.</p>
    </blockquote>
    <h2>Заключение</h2>
    <p>Технологические компании находятся на переднем крае изменений в организации работы. От внедрения новых инструментов до переосмысления корпоративной культуры, эти изменения формируют будущее работы не только в технологическом секторе, но и в мире в целом.</p>
  `,
  createdAt: '10.06.2023',
  updatedAt: '12.06.2023',
  isPublished: false,
  isGenerated: true
})

// Копия для редактирования
const editedContent = reactive({
  title: '',
  content: ''
})

// Загрузка контента при монтировании компонента
onMounted(async () => {
  const contentId = route.params.id
  
  try {
    // В реальном приложении здесь был бы запрос к API для получения контента
    // В этом примере мы уже используем мок-данные
    
    // Имитация загрузки
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // В реальном приложении здесь мы бы обновили состояние content
    isLoading.value = false
  } catch (error) {
    console.error('Ошибка при загрузке контента:', error)
    isLoading.value = false
  }
})

// Методы
const saveContent = async () => {
  isSaving.value = true
  
  try {
    // В реальном приложении здесь был бы запрос к API для сохранения контента
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Обновляем данные
    content.title = editedContent.title
    content.content = editedContent.content
    content.updatedAt = new Date().toLocaleDateString()
    
    // Выходим из режима редактирования
    isEditing.value = false
  } catch (error) {
    console.error('Ошибка при сохранении контента:', error)
  } finally {
    isSaving.value = false
  }
}

const cancelEditing = () => {
  // Сбрасываем изменения и выходим из режима редактирования
  isEditing.value = false
}

const publishContent = async () => {
  isPublishing.value = true
  
  try {
    // В реальном приложении здесь был бы запрос к API для публикации контента
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Обновляем статус
    content.isPublished = true
    content.updatedAt = new Date().toLocaleDateString()
  } catch (error) {
    console.error('Ошибка при публикации контента:', error)
  } finally {
    isPublishing.value = false
  }
}

const regenerateContent = async () => {
  isRegenerating.value = true
  
  try {
    // В реальном приложении здесь был бы запрос к API для регенерации контента
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Обновляем контент (в реальном приложении новый контент пришел бы с сервера)
    content.content = `
      <h1>${content.title}</h1>
      <p>Это обновленный контент после регенерации. В реальном приложении здесь будет новый вариант сгенерированного контента с помощью API.</p>
      <h2>Новый подзаголовок 1</h2>
      <p>Регенерированный текст для первого раздела, который отличается от предыдущего варианта.</p>
      <ul>
        <li>Новый пункт 1</li>
        <li>Новый пункт 2</li>
        <li>Новый пункт 3</li>
      </ul>
      <h2>Новый подзаголовок 2</h2>
      <p>Дополнительная информация для второго раздела, которая отличается от предыдущего варианта.</p>
      <blockquote>
        <p>Обновленная цитата в новой версии контента после регенерации.</p>
      </blockquote>
      <p>Новый заключительный абзац, подводящий итоги всего вышесказанного.</p>
    `
    content.updatedAt = new Date().toLocaleDateString()
  } catch (error) {
    console.error('Ошибка при регенерации контента:', error)
  } finally {
    isRegenerating.value = false
  }
}

const deleteContent = async () => {
  isDeleting.value = true
  
  try {
    // В реальном приложении здесь был бы запрос к API для удаления контента
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Перенаправляем пользователя на список контента
    router.push('/content')
  } catch (error) {
    console.error('Ошибка при удалении контента:', error)
    showDeleteConfirm.value = false
  } finally {
    isDeleting.value = false
  }
}

// При входе в режим редактирования копируем данные
watch(isEditing, (newValue) => {
  if (newValue) {
    editedContent.title = content.title
    editedContent.content = content.content
  }
})
</script> 