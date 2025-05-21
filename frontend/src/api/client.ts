// @ts-nocheck - отключаем проверку типов для этого файла из-за проблем с типами axios
import axios from 'axios'

// Базовая конфигурация для всех API-клиентов
const baseConfig = {
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Для работы с куками
}

/**
 * Создает экземпляр API-клиента для конкретного микросервиса
 * @param servicePath Путь к микросервису (auth, profile, content)
 * @param additionalConfig Дополнительная конфигурация axios
 * @returns Экземпляр axios для работы с микросервисом
 */
export function createApiClient(servicePath, additionalConfig = {}) {
  // Формируем базовый URL для микросервиса
  const baseURL = `/api/v1/${servicePath}`

  // Создаем и возвращаем настроенный экземпляр axios
  const client = axios.create({
    ...baseConfig,
    ...additionalConfig,
    baseURL,
  })

  // Добавляем перехватчики для обработки ошибок
  client.interceptors.response.use(
    (response) => response,
    async (error) => {
      // Общая обработка ошибок для всех сервисов
      const errorMessage = error.response?.data?.message || `Ошибка в сервисе ${servicePath}`
      
      // Можно добавить специфичную обработку для разных статусов
      if (error.response?.status === 401) {
        // Возможно, перенаправление на страницу логина
        console.warn('Необходима авторизация')
      }
      
      // Для разработки
      if (import.meta.env.DEV) {
        console.error(`API Error (${servicePath}):`, error)
      }
      
      return Promise.reject(error)
    }
  )

  return client
}

// Создаем базовый экземпляр axios для auth API
const api = axios.create({
  ...baseConfig,
  baseURL: '/api/v1'
})

// Экспортируем готовые экземпляры для каждого сервиса
export const authApi = createApiClient('auth')
export const profileApi = createApiClient('profile')
export const contentApi = createApiClient('content') 