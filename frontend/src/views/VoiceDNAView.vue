<template>
  <DashboardLayout>
    <div>
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Настройка Голоса ДНК</h1>
        <p class="text-gray-600">
          Давайте персонализируем VOY под ваш уникальный стиль и аудиторию
        </p>
      </div>

      <!-- Прогресс создания профиля -->
      <div class="relative mb-8">
        <!-- Линия прогресса -->
        <div class="relative">
          <!-- Фоновая линия -->
          <div style="width: 94%;" class="absolute top-5 inset-x-0 h-1 bg-gray-200"></div>
          
          <!-- Индикаторы шагов -->
          <div class="relative flex justify-between">
            <!-- Шаг 1 -->
            <div class="flex flex-col items-center">
              <div 
                class="w-10 h-10 rounded-full flex items-center justify-center mb-2 cursor-pointer"
                :class="[currentStep >= 1 ? 'bg-indigo-500 text-white' : 'bg-gray-200 text-gray-600']"
                @click="goToStep(1)"
              >
                1
              </div>
              <div class="text-sm text-center" :class="{ 'text-indigo-600 font-medium': currentStep === 1 }">
                О вас
              </div>
            </div>
            
            <!-- Шаг 2 -->
            <div class="flex flex-col items-center">
              <div 
                class="w-10 h-10 rounded-full flex items-center justify-center mb-2 cursor-pointer"
                :class="[currentStep >= 2 ? 'bg-indigo-500 text-white' : 'bg-gray-200 text-gray-600']"
                @click="goToStep(2)"
              >
                2
              </div>
              <div class="text-sm text-center" :class="{ 'text-indigo-600 font-medium': currentStep === 2 }">
                Ваш стиль
              </div>
            </div>
            
            <!-- Шаг 3 -->
            <div class="flex flex-col items-center">
              <div 
                class="w-10 h-10 rounded-full flex items-center justify-center mb-2 cursor-pointer"
                :class="[currentStep >= 3 ? 'bg-indigo-500 text-white' : 'bg-gray-200 text-gray-600']"
                @click="goToStep(3)"
              >
                3
              </div>
              <div class="text-sm text-center" :class="{ 'text-indigo-600 font-medium': currentStep === 3 }">
                Ваша аудитория
              </div>
            </div>
            
            <!-- Шаг 4 -->
            <div class="flex flex-col items-center">
              <div 
                class="w-10 h-10 rounded-full flex items-center justify-center mb-2 cursor-pointer"
                :class="[currentStep >= 4 ? 'bg-indigo-500 text-white' : 'bg-gray-200 text-gray-600']"
                @click="goToStep(4)"
              >
                4
              </div>
              <div class="text-sm text-center" :class="{ 'text-indigo-600 font-medium': currentStep === 4 }">
                Ключевые темы
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Шаг 1: О вас -->
      <div v-if="currentStep === 1" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Расскажите о себе</h2>
        
        <div class="mb-6">
          <p class="text-gray-700 mb-4">
            Эта информация поможет адаптировать контент под вашу профессиональную идентичность.
          </p>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Профессиональная биография</label>
              <textarea 
                v-model="profile.biography"
                rows="5" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 resize-none" 
                placeholder="Например: Я - консультант по цифровому маркетингу с 10-летним опытом помощи SaaS-компаниям в расширении их аудитории..."
              ></textarea>
              <p class="mt-1 text-sm text-gray-500">Ваш опыт, экспертиза и то, что делает вас уникальным</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Отрасль</label>
              <select 
                v-model="profile.industry"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="" disabled selected>Выберите вашу отрасль</option>
                <option value="technology">Технологии и IT</option>
                <option value="marketing">Маркетинг и реклама</option>
                <option value="finance">Финансы и инвестиции</option>
                <option value="healthcare">Здравоохранение</option>
                <option value="education">Образование</option>
                <option value="ecommerce">Электронная коммерция</option>
                <option value="manufacturing">Производство</option>
                <option value="consulting">Консалтинг</option>
                <option value="creative">Креативные индустрии</option>
                <option value="other">Другое</option>
              </select>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end">
          <Button 
            variant="primary" 
            @click="nextStep"
          >
            Далее
          </Button>
        </div>
      </div>

      <!-- Шаг 2: Ваш стиль -->
      <div v-if="currentStep === 2" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Определите ваш стиль письма</h2>
        
        <div class="mb-6">
          <p class="text-gray-700 mb-4">
            Как бы вы описали свой идеальный стиль контента?
          </p>
          
          <div class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Формальность тона</label>
              <div class="relative pt-1">
                <input 
                  v-model="style.formality" 
                  type="range" 
                  min="0" 
                  max="100" 
                  step="1"
                  class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-indigo-600"
                />
                <div class="flex justify-between text-xs text-gray-600 px-2 mt-1">
                  <span>Разговорный</span>
                  <span>Нейтральный</span>
                  <span>Формальный</span>
                </div>
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Характеристики стиля</label>
              <div class="grid grid-cols-2 gap-3">
                <div class="flex items-start">
                  <input 
                    id="conversational" 
                    v-model="style.characteristics" 
                    type="checkbox" 
                    value="conversational"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="conversational" class="ml-2 block text-sm text-gray-700">Разговорный</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="educational" 
                    v-model="style.characteristics" 
                    type="checkbox" 
                    value="educational"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="educational" class="ml-2 block text-sm text-gray-700">Образовательный</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="persuasive" 
                    v-model="style.characteristics" 
                    type="checkbox" 
                    value="persuasive"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="persuasive" class="ml-2 block text-sm text-gray-700">Убедительный</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="inspiring" 
                    v-model="style.characteristics" 
                    type="checkbox" 
                    value="inspiring"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="inspiring" class="ml-2 block text-sm text-gray-700">Вдохновляющий</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="analytical" 
                    v-model="style.characteristics" 
                    type="checkbox" 
                    value="analytical"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="analytical" class="ml-2 block text-sm text-gray-700">Аналитический</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="storytelling" 
                    v-model="style.characteristics" 
                    type="checkbox" 
                    value="storytelling"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="storytelling" class="ml-2 block text-sm text-gray-700">Повествовательный</label>
                </div>
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Примеры текстов (Опционально)</label>
              <textarea 
                v-model="style.examples"
                rows="5" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 resize-none" 
                placeholder="Вставьте абзац вашего текста, который представляет ваш идеальный стиль..."
              ></textarea>
              <p class="mt-1 text-sm text-gray-500">Это поможет VOY лучше понять ваш уникальный голос</p>
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
          >
            Далее
          </Button>
        </div>
      </div>

      <!-- Шаг 3: Ваша аудитория -->
      <div v-if="currentStep === 3" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Определите вашу аудиторию</h2>
        
        <div class="mb-6">
          <p class="text-gray-700 mb-4">
            Для кого вы создаете контент?
          </p>
          
          <div class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Описание аудитории</label>
              <textarea 
                v-model="audience.description"
                rows="5" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 resize-none" 
                placeholder="Например: Владельцы малого бизнеса в технологической отрасли, которые хотят улучшить свой цифровой маркетинг..."
              ></textarea>
            </div>
            
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Возрастной диапазон</label>
                <select 
                  v-model="audience.ageRange"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="" disabled selected>Выберите возрастной диапазон</option>
                  <option value="18-24">18-24</option>
                  <option value="25-34">25-34</option>
                  <option value="35-44">35-44</option>
                  <option value="45-54">45-54</option>
                  <option value="55+">55+</option>
                  <option value="mixed">Смешанный</option>
                </select>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Уровень экспертности</label>
                <select 
                  v-model="audience.expertiseLevel"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="" disabled selected>Выберите уровень экспертности</option>
                  <option value="beginner">Новички</option>
                  <option value="intermediate">Средний уровень</option>
                  <option value="advanced">Продвинутый уровень</option>
                  <option value="expert">Эксперты</option>
                  <option value="mixed">Смешанный</option>
                </select>
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Цели аудитории (Чего они хотят достичь?)</label>
              <div class="grid grid-cols-2 gap-3">
                <div class="flex items-start">
                  <input 
                    id="learn" 
                    v-model="audience.goals" 
                    type="checkbox" 
                    value="learn"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="learn" class="ml-2 block text-sm text-gray-700">Получить новые навыки</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="solve" 
                    v-model="audience.goals" 
                    type="checkbox" 
                    value="solve"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="solve" class="ml-2 block text-sm text-gray-700">Решить проблемы</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="stay_informed" 
                    v-model="audience.goals" 
                    type="checkbox" 
                    value="stay_informed"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="stay_informed" class="ml-2 block text-sm text-gray-700">Быть в курсе</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="inspiration" 
                    v-model="audience.goals" 
                    type="checkbox" 
                    value="inspiration"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="inspiration" class="ml-2 block text-sm text-gray-700">Найти вдохновение</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="make_decisions" 
                    v-model="audience.goals" 
                    type="checkbox" 
                    value="make_decisions"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="make_decisions" class="ml-2 block text-sm text-gray-700">Принять решения</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="connect" 
                    v-model="audience.goals" 
                    type="checkbox" 
                    value="connect"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="connect" class="ml-2 block text-sm text-gray-700">Связаться с другими</label>
                </div>
              </div>
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
          >
            Далее
          </Button>
        </div>
      </div>

      <!-- Шаг 4: Ключевые темы -->
      <div v-if="currentStep === 4" class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Ключевые темы</h2>
        
        <div class="mb-6">
          <p class="text-gray-700 mb-4">
            На каких ключевых темах вы хотите сосредоточиться в своем контенте?
          </p>
          
          <div class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Основные темы</label>
              <textarea 
                v-model="content.mainTopics"
                rows="5" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 resize-none" 
                placeholder="Например: Цифровой маркетинг, SEO, контент-стратегия, социальные сети..."
              ></textarea>
              <p class="mt-1 text-sm text-gray-500">Вводите темы через запятую</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Проблемные точки аудитории</label>
              <textarea 
                v-model="content.painPoints"
                rows="5" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 resize-none" 
                placeholder="Например: Трудности с развитием присутствия в социальных сетях, отсутствие постоянного контента, незнание о чем публиковать..."
              ></textarea>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Цели контента</label>
              <div class="grid grid-cols-2 gap-3">
                <div class="flex items-start">
                  <input 
                    id="educate" 
                    v-model="content.goals" 
                    type="checkbox" 
                    value="educate"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="educate" class="ml-2 block text-sm text-gray-700">Обучать</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="entertain" 
                    v-model="content.goals" 
                    type="checkbox" 
                    value="entertain"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="entertain" class="ml-2 block text-sm text-gray-700">Развлекать</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="inspire" 
                    v-model="content.goals" 
                    type="checkbox" 
                    value="inspire"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="inspire" class="ml-2 block text-sm text-gray-700">Вдохновлять</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="persuade" 
                    v-model="content.goals" 
                    type="checkbox" 
                    value="persuade"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="persuade" class="ml-2 block text-sm text-gray-700">Убеждать</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="awareness" 
                    v-model="content.goals" 
                    type="checkbox" 
                    value="awareness"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="awareness" class="ml-2 block text-sm text-gray-700">Повышать узнаваемость</label>
                </div>
                
                <div class="flex items-start">
                  <input 
                    id="convert" 
                    v-model="content.goals" 
                    type="checkbox" 
                    value="convert"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="convert" class="ml-2 block text-sm text-gray-700">Стимулировать конверсии</label>
                </div>
              </div>
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
            @click="saveVoiceDNA"
            :isLoading="isSaving"
          >
            Сохранить
          </Button>
        </div>
      </div>

      <!-- Подтверждение -->
      <div v-if="currentStep === 5" class="bg-white shadow-sm rounded-lg p-6">
        <div class="text-center mb-6">
          <div class="inline-flex items-center justify-center w-24 h-24 bg-green-100 rounded-full mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          
          <h2 class="text-2xl font-bold text-gray-900 mb-2">Голос ДНК настроен!</h2>
          <p class="text-gray-700 mb-6">
            Ваш уникальный голосовой профиль успешно создан и сохранен.
            Теперь вы можете генерировать контент, который будет звучать в вашем собственном стиле.
          </p>
          
          <div class="bg-indigo-50 p-4 rounded-lg text-left mb-8">
            <h3 class="font-medium text-indigo-800 mb-2">Следующие шаги:</h3>
            <ul class="list-disc pl-5 space-y-1 text-indigo-700">
              <li>Перейдите в раздел "Создать контент", чтобы начать использовать ваш голосовой профиль</li>
              <li>Вы можете обновить профиль в любое время, вернувшись на эту страницу</li>
            </ul>
          </div>
        </div>
        
        <div class="flex justify-between">
          <Button 
            variant="outline" 
            @click="router.push('/dashboard')"
          >
            Вернуться в дашборд
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
import { useToast } from 'vue-toast-notification'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'
import { useProfileStore } from '@/stores/profile'

const router = useRouter()
const toast = useToast()
const profileStore = useProfileStore()

// Настройки для шагов
const totalSteps = 5
const currentStep = ref(1)

// Данные профиля
const profile = ref({
  biography: '',
  industry: ''
})

// Данные стиля
const style = ref({
  formality: 50, // По умолчанию нейтральный (середина шкалы от 0 до 100)
  characteristics: [] as string[],
  examples: ''
})

// Данные аудитории
const audience = ref({
  description: '',
  ageRange: '',
  expertiseLevel: '',
  goals: [] as string[]
})

// Данные контента
const content = ref({
  mainTopics: '',
  painPoints: '',
  goals: [] as string[]
})

// Состояние сохранения
const isSaving = ref(false)

// Навигация по шагам
const nextStep = () => {
  if (currentStep.value < totalSteps) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

const goToStep = (step: number) => {
  if (step > 0 && step <= totalSteps) {
    currentStep.value = step
  }
}

// Сохранение данных голосового профиля
const saveVoiceDNA = async () => {
  // Базовая валидация формы
  if (!profile.value.biography.trim()) {
    alert('Пожалуйста, заполните профессиональную биографию')
    return
  }
  
  if (!profile.value.industry) {
    alert('Пожалуйста, выберите отрасль')
    return
  }
  
  isSaving.value = true
  
  try {
    // Предварительная обработка данных формы, чтобы отправить только релевантные данные
    const cleanedStyle = {
      formality: style.value.formality,
      characteristics: [...style.value.characteristics], // Создаем копию массива
      examples: typeof style.value.examples === 'string' ? style.value.examples.substring(0, 5000) : '' // Ограничиваем длину примера
    }
    
    const cleanedContent = {
      goals: [...content.value.goals], // Создаем копию массива
      mainTopics: typeof content.value.mainTopics === 'string' ? content.value.mainTopics.substring(0, 5000) : '',
      painPoints: typeof content.value.painPoints === 'string' ? content.value.painPoints.substring(0, 5000) : ''
    }
    
    const cleanedAudience = {
      goals: [...audience.value.goals], // Создаем копию массива
      ageRange: audience.value.ageRange,
      description: typeof audience.value.description === 'string' ? audience.value.description.substring(0, 5000) : '',
      expertiseLevel: audience.value.expertiseLevel
    }
    
    const cleanedProfile = {
      industry: profile.value.industry,
      biography: typeof profile.value.biography === 'string' ? profile.value.biography.substring(0, 5000) : ''
    }
    
    // Собираем все данные вместе в структурированный объект
    const voiceDNAData = {
      profile: cleanedProfile,
      style: cleanedStyle,
      audience: cleanedAudience,
      content: cleanedContent,
      createdAt: new Date().toISOString()
    }
    
    // Отправляем на сервер через профильное хранилище
    // PUT /api/v1/profile
    await profileStore.updateProfile({ 
      dna_data: JSON.stringify(voiceDNAData)
    })
    
    // Сохранение успешно
    showSuccessMessage('Профиль Голоса ДНК успешно сохранен')
    router.push('/dashboard')
  } catch (error) {
    console.error('Ошибка при сохранении профиля Голоса ДНК:', error)
    showErrorMessage('Ошибка при сохранении профиля')
  } finally {
    isSaving.value = false
  }
}

// Функции для показа уведомлений
const showSuccessMessage = (message: string): void => {
  toast.success(message, {
    duration: 3000,
    position: 'top-right'
  })
}

const showErrorMessage = (message: string): void => {
  toast.error(message, {
    duration: 5000,
    position: 'top-right'
  })
}
</script> 