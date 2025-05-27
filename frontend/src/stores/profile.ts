import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { profileApi } from '@/api/client'
import { useToast } from 'vue-toast-notification'

export interface ProfileData {
  user_id: number
  dna_data: string | null
  updated_at: string
  created_at: string
}

// Интерфейс для статуса голосового профиля
export interface VoiceProfileStatus {
  isActive: boolean
  lastUpdated: string | null
  message: string
}

export const useProfileStore = defineStore('profile', () => {
  const profile = ref<ProfileData | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const toast = useToast()

  // Вычисляемые свойства
  const hasVoiceProfile = computed(() => {
    if (!profile.value?.dna_data) {
      return false
    }
    
    try {
      // Пытаемся распарсить JSON, чтобы убедиться, что это не пустой объект
      const dnaData = JSON.parse(profile.value.dna_data)
      return dnaData && 
             Object.keys(dnaData).length > 0 && 
             dnaData.profile && 
             dnaData.style
    } catch (e) {
      console.error('Ошибка при парсинге dna_data:', e)
      return false
    }
  })

  // Получение статуса голосового профиля
  const getVoiceProfileStatus = computed((): VoiceProfileStatus => {
    if (!profile.value) {
      return {
        isActive: false,
        lastUpdated: null,
        message: 'Профиль не загружен'
      }
    }

    // Проверяем наличие профиля с той же логикой, что и в hasVoiceProfile
    if (!profile.value.dna_data) {
      return {
        isActive: false,
        lastUpdated: null,
        message: 'Голосовой профиль не настроен'
      }
    }
    
    try {
      const dnaData = JSON.parse(profile.value.dna_data)
      if (!dnaData || Object.keys(dnaData).length === 0 || !dnaData.profile || !dnaData.style) {
        return {
          isActive: false,
          lastUpdated: null,
          message: 'Голосовой профиль не настроен'
        }
      }
    } catch (e) {
      return {
        isActive: false,
        lastUpdated: null,
        message: 'Голосовой профиль повреждён'
      }
    }

    return {
      isActive: true,
      lastUpdated: profile.value.updated_at,
      message: 'Голосовой профиль активен'
    }
  })

  // Получение профиля пользователя
  const getProfile = async () => {
    isLoading.value = true
    error.value = null
    try {
      const response = await profileApi.get<ProfileData>('')
      profile.value = response.data
      return response.data
    } catch (err: any) {
      // Если профиль не найден (404), это не ошибка
      if (err.response?.status === 404) {
        profile.value = null
        return null
      }
      error.value = err.response?.data?.error || 'Ошибка при получении профиля'
      console.error('Ошибка при получении профиля:', err)
      
      // Показываем уведомление только если это не 404
      toast.error(error.value || 'Ошибка при получении профиля', { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Обновление профиля пользователя
  const updateProfile = async (data: { dna_data: string }) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await profileApi.put<ProfileData>('', data)
      profile.value = response.data
      
      // Показываем уведомление об успехе
      toast.success('Голосовой профиль успешно обновлен', { duration: 3000 })
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Ошибка при обновлении профиля'
      console.error('Ошибка при обновлении профиля:', err)
      toast.error(error.value || 'Ошибка при обновлении профиля', { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    profile,
    isLoading,
    error,
    hasVoiceProfile,
    getVoiceProfileStatus,
    getProfile,
    updateProfile
  }
}) 