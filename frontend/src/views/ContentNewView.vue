<template>
  <DashboardLayout>
    <div>
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Создание нового контента</h1>
        <p class="text-gray-600">
          Введите тему или краткое описание, и мы сгенерируем контент в вашем собственном стиле с помощью вашего голосового профиля.
        </p>
      </div>

      <div class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <div class="mb-6">
          <label for="title" class="block text-sm font-medium text-gray-700 mb-2">
            Заголовок
          </label>
          <input
            type="text"
            id="title"
            v-model="formData.title"
            class="w-full px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
            placeholder="Укажите заголовок для вашего контента"
            required
          />
        </div>

        <div class="mb-6">
          <div class="flex justify-between mb-2">
            <label for="description" class="block text-sm font-medium text-gray-700">
              Описание или заметки
            </label>
            <span class="text-sm text-gray-500">
              {{ formData.description.length }} / 500
            </span>
          </div>
          <textarea
            id="description"
            v-model="formData.description"
            rows="5"
            class="w-full px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
            placeholder="Опишите, о чем должен быть контент, или вставьте свои заметки"
            maxlength="500"
            required
          ></textarea>
        </div>

        <div class="mb-8">
          <div class="flex items-center justify-between mb-2">
            <label class="block text-sm font-medium text-gray-700">
              Или запишите голосовые заметки
            </label>
            <div v-if="recordingStatus === 'ready'" class="text-sm text-gray-500">
              Готово к записи
            </div>
            <div v-else-if="recordingStatus === 'recording'" class="text-sm text-red-500 animate-pulse">
              Запись... {{ recordingTime }}с
            </div>
            <div v-else-if="recordingStatus === 'processing'" class="text-sm text-indigo-500">
              Обработка...
            </div>
            <div v-else-if="recordingStatus === 'done'" class="text-sm text-green-500">
              Запись завершена!
            </div>
          </div>

          <div class="flex items-center gap-3">
            <button
              v-if="recordingStatus !== 'recording'"
              @click="startRecording"
              class="inline-flex items-center px-4 py-2 bg-indigo-100 text-indigo-800 rounded-lg hover:bg-indigo-200 transition"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
              </svg>
              Начать запись
            </button>
            <button
              v-else
              @click="stopRecording"
              class="inline-flex items-center px-4 py-2 bg-red-100 text-red-800 rounded-lg hover:bg-red-200 transition"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
              </svg>
              Остановить запись
            </button>
            <button
              v-if="formData.audioTranscript"
              @click="clearRecording"
              class="inline-flex items-center px-4 py-2 bg-gray-100 text-gray-800 rounded-lg hover:bg-gray-200 transition"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              Очистить запись
            </button>
          </div>

          <div v-if="formData.audioTranscript" class="mt-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <h3 class="text-sm font-medium text-gray-900 mb-2">Транскрипция:</h3>
            <p class="text-gray-700">{{ formData.audioTranscript }}</p>
          </div>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Дополнительные настройки
          </label>
          <div class="flex flex-col space-y-4">
            <div class="flex items-center">
              <input
                type="checkbox"
                id="useVoiceProfile"
                v-model="formData.useVoiceProfile"
                class="rounded text-indigo-600 focus:ring-indigo-500 h-4 w-4"
              />
              <label for="useVoiceProfile" class="ml-2 text-gray-700">
                Использовать мой голосовой профиль
              </label>
            </div>
            
            <div class="flex items-center">
              <input
                type="checkbox"
                id="includeImages"
                v-model="formData.includeImages"
                class="rounded text-indigo-600 focus:ring-indigo-500 h-4 w-4"
              />
              <label for="includeImages" class="ml-2 text-gray-700">
                Включить тематические изображения
              </label>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between">
          <Button 
            variant="outline" 
            @click="router.push('/dashboard')"
          >
            Отмена
          </Button>
          <Button 
            variant="primary"
            :disabled="isLoading || !formData.title || (!formData.description && !formData.audioTranscript)" 
            :is-loading="isLoading"
            @click="generateContent"
          >
            Сгенерировать контент
          </Button>
        </div>
      </div>

      <!-- Предпросмотр готового контента, если есть -->
      <div v-if="generatedContent" class="bg-white shadow-sm rounded-lg p-6">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-xl font-semibold text-gray-900">Предпросмотр контента</h2>
          <div class="flex items-center space-x-4">
            <Button
              variant="outline"
              size="sm"
              @click="regenerateContent"
              :is-loading="isRegenerating"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Сгенерировать заново
            </Button>
            <Button
              variant="primary"
              size="sm"
              @click="saveContent"
              :is-loading="isSaving"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              Сохранить
            </Button>
          </div>
        </div>
        
        <div class="prose prose-indigo max-w-none">
          <div v-html="generatedContent"></div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'

const router = useRouter()

// Состояния
const isLoading = ref(false)
const isRegenerating = ref(false)
const isSaving = ref(false)
const recordingStatus = ref<'ready' | 'recording' | 'processing' | 'done'>('ready')
const recordingTime = ref(0)
const recordingInterval = ref<number | null>(null)
const generatedContent = ref<string | null>(null)

// Данные формы
const formData = reactive({
  title: '',
  description: '',
  audioTranscript: '',
  useVoiceProfile: true,
  includeImages: false
})

// Функции для записи голоса
const startRecording = () => {
  recordingStatus.value = 'recording'
  recordingTime.value = 0
  
  // Интервал для обновления счетчика времени
  recordingInterval.value = window.setInterval(() => {
    recordingTime.value++
  }, 1000)
  
  // Здесь должна быть реальная логика записи с использованием Web Audio API
  // В этом примере мы просто имитируем запись
}

const stopRecording = () => {
  recordingStatus.value = 'processing'
  
  if (recordingInterval.value) {
    clearInterval(recordingInterval.value)
    recordingInterval.value = null
  }
  
  // Здесь должна быть логика остановки записи и отправки аудио на сервер для транскрипции
  // В этом примере мы имитируем процесс с задержкой
  
  setTimeout(() => {
    recordingStatus.value = 'done'
    formData.audioTranscript = 'Это пример транскрипции голосовой записи. В реальном приложении здесь будет текст, полученный из вашей голосовой записи с помощью API транскрипции.'
  }, 2000)
}

const clearRecording = () => {
  formData.audioTranscript = ''
  recordingStatus.value = 'ready'
}

// Генерация контента
const generateContent = async () => {
  isLoading.value = true
  
  try {
    // Здесь должен быть запрос к API для генерации контента
    // В этом примере мы имитируем ответ от сервера с задержкой
    
    await new Promise(resolve => setTimeout(resolve, 3000))
    
    // Пример сгенерированного содержимого
    generatedContent.value = `
      <h1>${formData.title}</h1>
      <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl vel tincidunt luctus, nunc nibh aliquam nunc, eget tincidunt nisl nisi vel nisl. Sed euismod, nisl vel tincidunt luctus, nunc nibh aliquam nunc, eget tincidunt nisl nisi vel nisl.</p>
      <h2>Подзаголовок 1</h2>
      <p>Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi. Nulla facilisi.</p>
      <ul>
        <li>Пункт 1</li>
        <li>Пункт 2</li>
        <li>Пункт 3</li>
      </ul>
      <h2>Подзаголовок 2</h2>
      <p>Donec auctor, nisl eget aliquam luctus, nisl nisl aliquam nisl, eget aliquam nisl nisl eget nisl. Donec auctor, nisl eget aliquam luctus, nisl nisl aliquam nisl, eget aliquam nisl nisl eget nisl.</p>
      <blockquote>
        <p>Цитата или важная мысль из сгенерированного контента.</p>
      </blockquote>
      <p>Заключительный абзац, подводящий итоги всего вышесказанного.</p>
    `
  } catch (error) {
    console.error('Ошибка при генерации контента:', error)
  } finally {
    isLoading.value = false
  }
}

// Регенерация контента
const regenerateContent = async () => {
  isRegenerating.value = true
  
  try {
    // Здесь должен быть запрос к API для регенерации контента
    await new Promise(resolve => setTimeout(resolve, 3000))
    
    // Обновляем сгенерированный контент с небольшими изменениями
    generatedContent.value = `
      <h1>${formData.title}</h1>
      <p>Это обновленный контент после регенерации. В реальном приложении здесь будет новый вариант сгенерированного контента с помощью API.</p>
      <h2>Новый подзаголовок 1</h2>
      <p>Обновленный текст для первого раздела, который отличается от предыдущего варианта.</p>
      <ul>
        <li>Обновленный пункт 1</li>
        <li>Обновленный пункт 2</li>
        <li>Новый пункт 3</li>
      </ul>
      <h2>Новый подзаголовок 2</h2>
      <p>Дополнительная информация для второго раздела, которая отличается от предыдущего варианта.</p>
      <blockquote>
        <p>Обновленная цитата в новой версии контента.</p>
      </blockquote>
      <p>Новый заключительный абзац, подводящий итоги всего вышесказанного.</p>
    `
  } catch (error) {
    console.error('Ошибка при регенерации контента:', error)
  } finally {
    isRegenerating.value = false
  }
}

// Сохранение контента
const saveContent = async () => {
  isSaving.value = true
  
  try {
    // Здесь должен быть запрос к API для сохранения контента
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // После успешного сохранения перенаправляем на страницу с контентом
    router.push('/content/1') // В реальном приложении здесь будет ID созданного контента
  } catch (error) {
    console.error('Ошибка при сохранении контента:', error)
  } finally {
    isSaving.value = false
  }
}
</script> 