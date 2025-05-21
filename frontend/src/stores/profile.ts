import { defineStore } from 'pinia'
import { ref } from 'vue'
import { profileApi } from '@/api/client'

export interface ProfileData {
  userId: number
  dnaData: string | null
  updatedAt: string
  createdAt: string
}

export const useProfileStore = defineStore('profile', () => {
  const profile = ref<ProfileData | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

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
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Обновление профиля пользователя
  const updateProfile = async (data: { dnaData: string }) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await profileApi.put<ProfileData>('', data)
      profile.value = response.data
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Ошибка при обновлении профиля'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    profile,
    isLoading,
    error,
    getProfile,
    updateProfile
  }
}) 