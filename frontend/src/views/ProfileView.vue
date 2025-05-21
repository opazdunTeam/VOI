<template>
  <DashboardLayout>
    <div>
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Профиль пользователя</h1>
        <p class="text-gray-600">
          Управляйте своими персональными данными и настройками.
        </p>
      </div>

      <div class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <div class="mb-6">
          <div class="flex items-center mb-6">
            <div class="bg-indigo-100 w-20 h-20 rounded-full flex items-center justify-center mr-4">
              <span class="text-3xl font-semibold text-indigo-700">
                {{ getUserInitial }}
              </span>
            </div>
            <div>
              <h2 class="text-xl font-semibold text-gray-900">{{ userData?.full_name || 'Загрузка...' }}</h2>
              <p class="text-gray-600">{{ userData?.email || 'Загрузка...' }}</p>
            </div>
          </div>

          <div class="border-t border-gray-200 pt-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Личные данные</h3>
            
            <form @submit.prevent="updateUserProfile" class="space-y-4">
              <div>
                <label for="full_name" class="block text-sm font-medium text-gray-700 mb-1">
                  Полное имя
                </label>
                <input
                  type="text"
                  id="full_name"
                  v-model="formData.full_name"
                  class="w-full px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                  placeholder="Ваше полное имя"
                  required
                />
              </div>
              
              <div>
                <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  type="email"
                  id="email"
                  v-model="formData.email"
                  class="w-full px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                  placeholder="Ваш email"
                  required
                  disabled
                />
                <p class="mt-1 text-xs text-gray-500">Email изменить нельзя.</p>
              </div>
              
              <div class="flex justify-end">
                <Button
                  type="submit"
                  :is-loading="isUpdatingProfile"
                  :disabled="!formDataChanged"
                >
                  Сохранить изменения
                </Button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <div class="bg-white shadow-sm rounded-lg p-6 mb-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Голосовой профиль</h2>
        
        <div class="mb-6">
          <div class="flex items-center justify-between mb-4">
            <div>
              <p class="text-gray-700">
                Статус: 
                <span v-if="profileData?.dnaData" class="text-green-600 font-medium">Активен</span>
                <span v-else class="text-red-600 font-medium">Не настроен</span>
              </p>
              
              <p v-if="profileData?.dnaData" class="text-sm text-gray-500">
                Обновлен: {{ new Date(profileData.updatedAt).toLocaleDateString() }}
              </p>
            </div>
            
            <router-link to="/voice-dna">
              <Button>
                {{ profileData?.dnaData ? 'Обновить профиль' : 'Создать профиль' }}
              </Button>
            </router-link>
          </div>
          
          <div v-if="profileData?.dnaData" class="bg-gray-50 p-4 rounded-lg">
            <h3 class="text-sm font-medium text-gray-900 mb-2">Информация о профиле:</h3>
            <div class="space-y-2">
              <div class="flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-gray-700">Голосовые образцы: 5 записей</span>
              </div>
              <div class="flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-gray-700">Качество: Хорошее</span>
              </div>
            </div>
          </div>
          
          <div v-else class="bg-indigo-50 p-4 rounded-lg">
            <p class="text-indigo-700">
              Создайте голосовой профиль, чтобы генерировать контент в вашем собственном стиле.
              Это займет всего несколько минут.
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white shadow-sm rounded-lg p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Безопасность</h2>
        
        <div class="mb-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Смена пароля</h3>
          
          <form @submit.prevent="changePassword" class="space-y-4">
            <div>
              <label for="currentPassword" class="block text-sm font-medium text-gray-700 mb-1">
                Текущий пароль
              </label>
              <input
                type="password"
                id="currentPassword"
                v-model="passwordData.currentPassword"
                class="w-full px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                placeholder="••••••••"
                required
              />
            </div>
            
            <div>
              <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">
                Новый пароль
              </label>
              <input
                type="password"
                id="newPassword"
                v-model="passwordData.newPassword"
                class="w-full px-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                placeholder="••••••••"
                required
                minlength="8"
              />
              <p class="mt-1 text-xs text-gray-500">Минимум 8 символов</p>
            </div>
            
            <div>
              <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">
                Подтверждение пароля
              </label>
              <input
                type="password"
                id="confirmPassword"
                v-model="passwordData.confirmPassword"
                class="w-full px-4 py-2 rounded-lg border focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                :class="passwordError ? 'border-red-500' : 'border-gray-300'"
                placeholder="••••••••"
                required
              />
              <p v-if="passwordError" class="mt-1 text-xs text-red-500">{{ passwordError }}</p>
            </div>
            
            <div class="flex justify-end">
              <Button
                type="submit"
                :is-loading="isChangingPassword"
                :disabled="!passwordDataFilled"
              >
                Изменить пароль
              </Button>
            </div>
          </form>
        </div>
        
        <div class="border-t border-gray-200 pt-6">
          <div class="flex justify-between items-center">
            <div>
              <h3 class="text-lg font-medium text-red-600">Удаление аккаунта</h3>
              <p class="text-gray-600 text-sm">
                Удаление аккаунта приведет к безвозвратной потере всех ваших данных.
              </p>
            </div>
            <Button
              variant="danger"
              @click="showDeleteConfirm = true"
            >
              Удалить аккаунт
            </Button>
          </div>
        </div>
      </div>
      
      <!-- Диалог подтверждения удаления аккаунта -->
      <div v-if="showDeleteConfirm" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
        <div class="bg-white rounded-lg p-6 max-w-md w-full">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Подтверждение удаления аккаунта</h3>
          <p class="text-gray-700 mb-6">
            Вы уверены, что хотите удалить свой аккаунт? Это действие нельзя отменить,
            и все ваши данные будут безвозвратно удалены.
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
              :is-loading="isDeletingAccount"
              @click="deleteAccount"
            >
              Подтвердить удаление
            </Button>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useProfileStore, type ProfileData } from '@/stores/profile'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'
import { useAuthApi } from '@/services/auth'

interface UserData {
  id: number
  full_name: string
  email: string
}

const router = useRouter()
const authStore = useAuthStore()
const profileStore = useProfileStore()
const authApi = useAuthApi()

const userData = ref<UserData | null>(null)
const profileData = ref<ProfileData | null>(null)
const isUpdatingProfile = ref(false)
const isChangingPassword = ref(false)
const isDeletingAccount = ref(false)
const showDeleteConfirm = ref(false)
const passwordError = ref('')

const formData = reactive({
  full_name: '',
  email: ''
})

// Форма для изменения пароля
const passwordData = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Вычисляемые свойства
const getUserInitial = computed(() => {
  if (!userData.value?.full_name) return '?'
  return userData.value.full_name.charAt(0).toUpperCase()
})

const formDataChanged = computed(() => {
  if (!userData.value) return false
  return formData.full_name !== userData.value.full_name
})

const passwordDataFilled = computed(() => {
  return (
    passwordData.currentPassword.length > 0 &&
    passwordData.newPassword.length >= 8 &&
    passwordData.confirmPassword.length >= 8
  )
})

// Получение данных при монтировании компонента
onMounted(async () => {
  try {
    // Получаем данные пользователя
    userData.value = await authStore.getCurrentUser()
    
    if (!userData.value) {
      throw new Error('Не удалось получить данные пользователя')
    }
    
    // Заполняем форму
    formData.full_name = userData.value.full_name
    formData.email = userData.value.email
    
    // Получаем данные профиля
    const profile = await profileStore.getProfile()
    profileData.value = profile
  } catch (error) {
    console.error('Ошибка при получении данных:', error)
  }
})

// Методы
const updateUserProfile = async () => {
  if (!formDataChanged.value) return
  
  isUpdatingProfile.value = true
  
  try {
    const updatedUser = await authApi.updateProfile({
      full_name: formData.full_name
    })
    
    // Обновляем локальные данные
    if (userData.value) {
      userData.value.full_name = updatedUser.data.full_name
    }
    
    // Показываем уведомление об успехе
    alert('Профиль успешно обновлен')
  } catch (error) {
    console.error('Ошибка при обновлении профиля:', error)
    alert('Ошибка при обновлении профиля')
  } finally {
    isUpdatingProfile.value = false
  }
}

const changePassword = async () => {
  // Сбрасываем ошибку
  passwordError.value = ''
  
  // Проверяем совпадение паролей
  if (passwordData.newPassword !== passwordData.confirmPassword) {
    passwordError.value = 'Пароли не совпадают'
    return
  }
  
  isChangingPassword.value = true
  
  try {
    await authApi.changePassword({
      current_password: passwordData.currentPassword,
      new_password: passwordData.newPassword
    })
    
    // Очищаем форму
    passwordData.currentPassword = ''
    passwordData.newPassword = ''
    passwordData.confirmPassword = ''
    
    // Показываем уведомление об успехе
    alert('Пароль успешно изменен')
  } catch (error: any) {
    console.error('Ошибка при изменении пароля:', error)
    
    if (error.response?.status === 403) {
      passwordError.value = 'Неверный текущий пароль'
    } else {
      alert('Ошибка при изменении пароля')
    }
  } finally {
    isChangingPassword.value = false
  }
}

const deleteAccount = async () => {
  isDeletingAccount.value = true
  
  try {
    // Здесь должен быть код для удаления аккаунта через API
    // В этом примере просто имитируем успешное удаление
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // Выходим из системы
    await authStore.logout()
    
    // Перенаправляем на страницу входа
    router.push('/login')
  } catch (error) {
    console.error('Ошибка при удалении аккаунта:', error)
    showDeleteConfirm.value = false
  } finally {
    isDeletingAccount.value = false
  }
}
</script> 