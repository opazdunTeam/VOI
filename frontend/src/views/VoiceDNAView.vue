<template>
  <DashboardLayout>
    <div>
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Настройка голосового профиля</h1>
        <p class="text-gray-600">
          Создайте ваш уникальный голосовой профиль для персонализации генерируемого контента.
        </p>
      </div>

      <!-- Прогресс создания профиля -->
      <div class="relative mb-8">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-700">Прогресс создания профиля</span>
          <span class="text-sm font-medium text-indigo-600">{{ currentStep }} из {{ totalSteps }}</span>
        </div>
        <div class="overflow-hidden h-2 rounded-full bg-gray-200">
          <div 
            class="h-full bg-indigo-600 rounded-full transition-all duration-500" 
            :style="{ width: `${(currentStep / totalSteps) * 100}%` }"
          ></div>
        </div>
      </div>

      <!-- Шаг 1: Инструкция -->
      <div v-if="currentStep === 1" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Добро пожаловать в настройку голосового профиля</h2>
        
        <div class="mb-6">
          <p class="text-gray-700 mb-4">
            Голосовой профиль (Voice DNA) - это уникальный образец вашего голоса и стиля речи, который позволяет 
            нашей системе генерировать контент, максимально приближенный к вашей манере общения.
          </p>
          
          <p class="text-gray-700 mb-4">
            Процесс состоит из нескольких простых шагов:
          </p>
          
          <ul class="list-disc pl-5 space-y-2 mb-4">
            <li class="text-gray-700">Вам нужно будет записать несколько образцов своего голоса</li>
            <li class="text-gray-700">Наша система проанализирует ваш голос и стиль речи</li>
            <li class="text-gray-700">На основе анализа будет создан ваш уникальный голосовой профиль</li>
            <li class="text-gray-700">Весь процесс займет около 5-10 минут</li>
          </ul>
          
          <div class="bg-indigo-50 p-4 rounded-lg">
            <p class="text-indigo-700">
              <strong>Важно:</strong> Для наилучших результатов рекомендуем использовать качественный микрофон и 
              находиться в тихом помещении во время записи.
            </p>
          </div>
        </div>
        
        <div class="flex justify-between">
          <Button 
            variant="outline" 
            @click="router.push('/dashboard')"
          >
            Отмена
          </Button>
          <Button 
            variant="primary" 
            @click="nextStep"
          >
            Начать создание профиля
          </Button>
        </div>
      </div>

      <!-- Шаг 2: Проверка микрофона -->
      <div v-if="currentStep === 2" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Проверка микрофона</h2>
        
        <div class="mb-6">
          <p class="text-gray-700 mb-4">
            Перед началом записи давайте убедимся, что ваш микрофон работает корректно.
            Нажмите кнопку "Проверить микрофон" и произнесите несколько слов.
          </p>
          
          <div class="flex flex-col items-center space-y-4 mb-6">
            <div 
              class="w-32 h-32 rounded-full flex items-center justify-center"
              :class="[
                isRecording ? 'bg-red-100 animate-pulse' : 'bg-gray-100'
              ]"
            >
              <svg 
                xmlns="http://www.w3.org/2000/svg" 
                class="h-16 w-16" 
                :class="{ 'text-red-500': isRecording, 'text-gray-400': !isRecording }"
                fill="none" 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
              </svg>
            </div>
            
            <Button 
              v-if="!isRecording" 
              @click="startRecording" 
              variant="primary"
            >
              Проверить микрофон
            </Button>
            <Button 
              v-else 
              @click="stopRecording" 
              variant="danger"
            >
              Остановить запись
            </Button>
          </div>
          
          <div v-if="recordingResult" class="mb-6">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Результат проверки:</h3>
            <div :class="[
              'p-4 rounded-lg',
              recordingResult === 'success' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
            ]">
              <p v-if="recordingResult === 'success'">
                <strong>Отлично!</strong> Ваш микрофон работает корректно.
                Можно переходить к следующему шагу.
              </p>
              <p v-else>
                <strong>Проблема с микрофоном.</strong> Убедитесь, что микрофон подключен и 
                в настройках браузера разрешен доступ к микрофону.
              </p>
            </div>
          </div>
        </div>
        
        <div class="flex justify-between">
          <Button 
            variant="outline" 
            @click="prevStep"
          >
            Назад
          </Button>
          <Button 
            variant="primary" 
            @click="nextStep"
            :disabled="!recordingResult || recordingResult === 'error'"
          >
            Продолжить
          </Button>
        </div>
      </div>

      <!-- Шаг 3: Запись образцов голоса -->
      <div v-if="currentStep === 3" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Запись образцов голоса</h2>
        
        <div class="mb-6">
          <p class="text-gray-700 mb-4">
            Сейчас вам будет предложено записать несколько фраз. Старайтесь говорить естественно,
            своим обычным голосом и в привычном для вас темпе.
          </p>
          
          <div v-if="currentPrompt" class="border border-gray-200 rounded-lg p-4 mb-6">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Фраза {{ currentPromptIndex + 1 }} из {{ voicePrompts.length }}:</h3>
            <p class="text-gray-800 text-lg mb-4">{{ currentPrompt }}</p>
            
            <div class="flex flex-col items-center space-y-4">
              <div 
                class="w-16 h-16 rounded-full flex items-center justify-center"
                :class="[
                  isRecording ? 'bg-red-100 animate-pulse' : 'bg-gray-100'
                ]"
              >
                <svg 
                  xmlns="http://www.w3.org/2000/svg" 
                  class="h-8 w-8" 
                  :class="{ 'text-red-500': isRecording, 'text-gray-400': !isRecording }"
                  fill="none" 
                  viewBox="0 0 24 24" 
                  stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
                </svg>
              </div>
              
              <Button 
                v-if="!isRecording" 
                @click="startPromptRecording" 
                variant="primary"
                size="sm"
              >
                Начать запись
              </Button>
              <Button 
                v-else 
                @click="stopPromptRecording" 
                variant="danger"
                size="sm"
              >
                Остановить запись
              </Button>
            </div>
          </div>
          
          <div v-if="completedPrompts.length > 0" class="mb-6">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Записанные фразы:</h3>
            <ul class="space-y-2">
              <li 
                v-for="(prompt, index) in completedPrompts" 
                :key="index"
                class="flex items-center text-gray-700"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                {{ prompt }}
              </li>
            </ul>
          </div>
          
          <div v-if="isProcessingVoice" class="mb-6">
            <div class="flex items-center justify-center p-6 bg-indigo-50 rounded-lg">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-indigo-700">Обработка голосовых данных...</span>
            </div>
          </div>
        </div>
        
        <div class="flex justify-between">
          <Button 
            variant="outline" 
            @click="prevStep"
            :disabled="isProcessingVoice"
          >
            Назад
          </Button>
          <Button 
            variant="primary" 
            @click="processVoiceData"
            :disabled="completedPrompts.length < voicePrompts.length || isProcessingVoice"
          >
            Создать голосовой профиль
          </Button>
        </div>
      </div>

      <!-- Шаг 4: Подтверждение и завершение -->
      <div v-if="currentStep === 4" class="bg-white shadow-sm rounded-lg p-6">
        <div class="text-center mb-6">
          <div class="inline-flex items-center justify-center w-24 h-24 bg-green-100 rounded-full mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          
          <h2 class="text-2xl font-bold text-gray-900 mb-2">Голосовой профиль создан!</h2>
          <p class="text-gray-700 mb-6">
            Ваш уникальный голосовой профиль успешно создан и сохранен.
            Теперь вы можете генерировать контент, который будет звучать в вашем собственном стиле.
          </p>
          
          <div class="bg-indigo-50 p-4 rounded-lg text-left mb-8">
            <h3 class="font-medium text-indigo-800 mb-2">Следующие шаги:</h3>
            <ul class="list-disc pl-5 space-y-1 text-indigo-700">
              <li>Перейдите в раздел "Создать контент", чтобы начать использовать ваш голосовой профиль</li>
              <li>Вы можете обновить профиль в любое время, вернувшись на эту страницу</li>
              <li>Чем больше образцов голоса вы предоставите, тем точнее будет результат</li>
            </ul>
          </div>
        </div>
        
        <div class="flex justify-between">
          <Button 
            variant="outline" 
            @click="router.push('/profile')"
          >
            Перейти в профиль
          </Button>
          <Button 
            variant="primary" 
            @click="router.push('/content/new')"
          >
            Создать контент
          </Button>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'
import { useProfileStore } from '@/stores/profile'

const router = useRouter()
const profileStore = useProfileStore()

// Настройки для шагов
const totalSteps = 4
const currentStep = ref(1)

// Запись голоса
const isRecording = ref(false)
const recordingResult = ref<null | 'success' | 'error'>(null)

// Образцы голоса
const voicePrompts = [
  'Технологии расширяют возможности человека и делают нашу жизнь более комфортной.',
  'Креативность - это способность видеть обычные вещи в необычном свете.',
  'Важно находить баланс между работой и личной жизнью для сохранения эмоционального благополучия.',
  'Постоянное обучение и адаптация к переменам - ключ к успеху в современном мире.',
  'Эффективная коммуникация строится на умении не только говорить, но и слушать.'
]
const currentPromptIndex = ref(0)
const currentPrompt = ref(voicePrompts[0])
const completedPrompts = ref<string[]>([])
const isProcessingVoice = ref(false)

// Навигация по шагам
const nextStep = () => {
  if (currentStep.value < totalSteps) {
    currentStep.value++
    
    // Если переходим на шаг 3, сбрасываем индекс подсказок
    if (currentStep.value === 3) {
      currentPromptIndex.value = 0
      currentPrompt.value = voicePrompts[0]
      completedPrompts.value = []
    }
  }
}

const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

// Проверка микрофона
const startRecording = () => {
  isRecording.value = true
  
  // Здесь должен быть код для запуска записи с микрофона
  // Это упрощенная версия для демонстрации
  
  // Имитация проверки микрофона с результатом через 2 секунды
  setTimeout(() => {
    stopRecording()
  }, 2000)
}

const stopRecording = () => {
  isRecording.value = false
  
  // В реальном приложении здесь бы проверялось качество записи
  // Для демонстрации просто устанавливаем успешный результат
  recordingResult.value = 'success'
}

// Запись образцов голоса
const startPromptRecording = () => {
  isRecording.value = true
  
  // В реальном приложении здесь был бы код для начала записи
  
  // Имитация записи в течение 3 секунд
  setTimeout(() => {
    stopPromptRecording()
  }, 3000)
}

const stopPromptRecording = () => {
  isRecording.value = false
  
  // Добавляем текущую подсказку в список завершенных
  completedPrompts.value.push(currentPrompt.value)
  
  // Переходим к следующей подсказке, если она есть
  if (currentPromptIndex.value < voicePrompts.length - 1) {
    currentPromptIndex.value++
    currentPrompt.value = voicePrompts[currentPromptIndex.value]
  } else {
    // Если все подсказки записаны, устанавливаем currentPrompt в пустую строку
    currentPrompt.value = ''
  }
}

// Обработка голосовых данных
const processVoiceData = async () => {
  isProcessingVoice.value = true
  
  try {
    // В реальном приложении здесь был бы запрос к API для обработки голосовых данных
    // Имитация обработки голосовых данных
    await new Promise(resolve => setTimeout(resolve, 3000))
    
    // Обновление профиля пользователя
    await profileStore.updateProfile({ dnaData: JSON.stringify({ voiceProfile: 'v1', status: 'active' }) })
    
    // Переход к следующему шагу
    nextStep()
  } catch (error) {
    console.error('Ошибка при обработке голосовых данных:', error)
  } finally {
    isProcessingVoice.value = false
  }
}
</script> 