import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/client'

export interface User {
  id: number
  email: string
  full_name: string
  created_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  full_name: string // используем snake_case для соответствия формату API
  email: string
  password: string
}

export interface AuthResponse {
  user: User
  token?: string // Добавлен опциональный токен в ответе
}

export interface ProfileResponse {
  userId: number
  dnaData: string | null
  updatedAt: string
  createdAt: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => !!user.value)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Регистрация пользователя
  const register = async (data: RegisterRequest) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await authApi.post<AuthResponse>('/register', data)
      user.value = response.data.user
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при регистрации'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Вход пользователя
  const login = async (data: LoginRequest) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await authApi.post<AuthResponse>('/login', data)
      user.value = response.data.user 
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при входе'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Выход пользователя
  const logout = async () => {
    isLoading.value = true
    try {
      await authApi.post('/logout')
      user.value = null
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при выходе'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Получение текущего пользователя
  const getCurrentUser = async () => {
    isLoading.value = true
    error.value = null
    try {
      const response = await authApi.get<User>('/me')
      user.value = response.data
      return response.data
    } catch (err: any) {
      // Если 401, не показываем ошибку, просто очищаем состояние
      if (err.response?.status === 401) {
        user.value = null
      } else {
        error.value = err.response?.data?.message || 'Ошибка при получении данных пользователя'
      }
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    register,
    login,
    logout,
    getCurrentUser
  }
}) 