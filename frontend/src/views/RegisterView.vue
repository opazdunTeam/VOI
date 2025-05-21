<template>
  <div class="flex min-h-screen bg-gray-100 w-full">
    <!-- Левая панель с информацией о продукте (такая же, как на странице логина) -->
    <div class="hidden md:block md:w-1/2 bg-gradient-to-b from-indigo-800 to-indigo-600 text-white p-10 flex flex-col">
      <h1 class="text-4xl font-bold mb-4">VOY</h1>
      <h2 class="text-2xl font-medium mb-8">Voice Of You</h2>
      
      <div class="flex-grow flex flex-col justify-center">
        <p class="text-xl mb-8 max-w-md">
          Преобразуйте свои голосовые заметки и идеи в персонализированный, увлекательный контент с помощью нашей платформы на базе ИИ.
        </p>
        
        <div class="space-y-6">
          <div class="flex items-start">
            <div class="bg-white/20 p-2 rounded-lg mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-lg">Голосовой профиль</h3>
              <p class="text-white/70">Уникальный ДНК-профиль вашего голоса и стиля</p>
            </div>
          </div>
          
          <div class="flex items-start">
            <div class="bg-white/20 p-2 rounded-lg mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h12a2 2 0 012 2v12a2 2 0 01-2 2H6a2 2 0 01-2-2V6z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-lg">Персонализация</h3>
              <p class="text-white/70">Индивидуальный контент, отражающий ваш стиль</p>
            </div>
          </div>
          
          <div class="flex items-start">
            <div class="bg-white/20 p-2 rounded-lg mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-lg">Генерация контента</h3>
              <p class="text-white/70">Мгновенное преобразование идей в готовый контент</p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Правая панель с формой регистрации -->
    <div class="w-full md:w-1/2 flex items-center justify-center px-6 py-12">
      <div class="w-full max-w-md">
        <div class="text-center mb-10">
          <h2 class="text-3xl font-bold text-gray-900 mb-2">Создайте аккаунт</h2>
          <p class="text-gray-600">Присоединяйтесь к VOY и начните создавать персонализированный контент</p>
        </div>
        
        <!-- Показываем ошибку только если она связана с регистрацией -->
        <div v-if="registerError" class="mb-6 p-4 bg-red-50 border-l-4 border-red-500 text-red-700 rounded">
          <p>{{ registerError }}</p>
        </div>
        
        <form @submit.prevent="handleSubmit">
          <div class="mb-6">
            <label for="fullName" class="block text-sm font-medium text-gray-700 mb-2">
              Полное имя
            </label>
            <input
              type="text"
              id="fullName"
              v-model="formData.fullName"
              class="w-full px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="Иван Иванов"
              required
            />
          </div>
          
          <div class="mb-6">
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              type="email"
              id="email"
              v-model="formData.email"
              class="w-full px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="ваша@почта.com"
              required
            />
          </div>
          
          <div class="mb-6">
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              Пароль
            </label>
            <input
              type="password"
              id="password"
              v-model="formData.password"
              class="w-full px-4 py-3 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="••••••••"
              required
              minlength="8"
            />
            <p class="mt-1 text-xs text-gray-500">
              Минимум 8 символов
            </p>
          </div>
          
          <div class="mb-8">
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              Подтверждение пароля
            </label>
            <input
              type="password"
              id="confirmPassword"
              v-model="formData.confirmPassword"
              class="w-full px-4 py-3 rounded-lg border focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              :class="passwordError ? 'border-red-500' : 'border-gray-300'"
              placeholder="••••••••"
              required
            />
            <p v-if="passwordError" class="mt-1 text-xs text-red-500">{{ passwordError }}</p>
          </div>
          
          <button
            type="submit"
            :disabled="authStore.isLoading"
            class="w-full py-3 px-4 bg-indigo-600 text-white rounded-lg font-medium transition-colors hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
            :class="{ 'opacity-70 cursor-not-allowed': authStore.isLoading }"
          >
            {{ authStore.isLoading ? 'Регистрация...' : 'Зарегистрироваться' }}
          </button>
        </form>
        
        <div class="mt-8 text-center">
          <p class="text-gray-600">
            Уже есть аккаунт? 
            <router-link to="/login" class="text-primary font-medium ml-1 hover:text-indigo-800">
              Войти
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const passwordError = ref('')
const registerError = ref('')

// Очищаем общую ошибку хранилища при монтировании
onMounted(() => {
  // Очистить ошибку, которая могла возникнуть при проверке авторизации
  authStore.error = null
})

const formData = reactive({
  fullName: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const handleSubmit = async () => {
  // Сбрасываем ошибку пароля
  passwordError.value = ''
  registerError.value = ''
  
  // Проверка совпадения паролей
  if (formData.password !== formData.confirmPassword) {
    passwordError.value = 'Пароли не совпадают'
    return
  }
  
  try {
    await authStore.register({
      full_name: formData.fullName, // используем snake_case для API
      email: formData.email,
      password: formData.password
    })
    
    // После успешной регистрации перенаправляем на дашборд
    router.push('/dashboard')
  } catch (error) {
    console.error('Registration error:', error)
    registerError.value = authStore.error || 'Ошибка при регистрации'
  }
}
</script> 