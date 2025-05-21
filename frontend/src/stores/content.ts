import { defineStore } from 'pinia'
import { ref } from 'vue'
import { contentApi } from '@/api/client'

export interface Post {
  id: number
  title: string
  content: string
  userId: number
  createdAt: string
  updatedAt: string
}

export const useContentStore = defineStore('content', () => {
  const posts = ref<Post[]>([])
  const currentPost = ref<Post | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Получение списка постов
  const getPosts = async () => {
    isLoading.value = true
    error.value = null
    try {
      const response = await contentApi.get<Post[]>('/posts')
      posts.value = response.data
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при получении списка постов'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Получение поста по ID
  const getPost = async (id: number) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await contentApi.get<Post>(`/posts/${id}`)
      currentPost.value = response.data
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при получении поста'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Создание нового поста
  const createPost = async (data: { title: string; content: string }) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await contentApi.post<Post>('/posts', data)
      posts.value = [response.data, ...posts.value]
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при создании поста'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Обновление поста
  const updatePost = async (id: number, data: { title: string; content: string }) => {
    isLoading.value = true
    error.value = null
    try {
      const response = await contentApi.put<Post>(`/posts/${id}`, data)
      const index = posts.value.findIndex(p => p.id === id)
      if (index !== -1) {
        posts.value[index] = response.data
      }
      currentPost.value = response.data
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при обновлении поста'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Удаление поста
  const deletePost = async (id: number) => {
    isLoading.value = true
    error.value = null
    try {
      await contentApi.delete(`/posts/${id}`)
      posts.value = posts.value.filter(p => p.id !== id)
      if (currentPost.value?.id === id) {
        currentPost.value = null
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Ошибка при удалении поста'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    posts,
    currentPost,
    isLoading,
    error,
    getPosts,
    getPost,
    createPost,
    updatePost,
    deletePost
  }
}) 