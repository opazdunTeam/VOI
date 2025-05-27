import { defineStore } from 'pinia'
import { ref } from 'vue'
import { contentApi, generatorApi, profileApi } from '@/api/client'
import { useToast } from 'vue-toast-notification'
import { useProfileStore } from '@/stores/profile'

// Интерфейсы для типизации данных
export interface Content {
  id: string | number
  user_id: number
  note_id?: number
  content_md: string // контент в формате Markdown
  status: string // Статус публикации
  created_at: string
  updated_at: string
  
  // Дополнительные поля для фронтенда
  title?: string // Заголовок, который можно извлечь из content_md
  excerpt?: string // Краткое описание, которое можно извлечь из content_md
  is_generated?: boolean // Сгенерирован ли контент, можно определить по наличию note_id
  is_published?: boolean // Опубликован ли контент, можно определить по полю status
}

export interface ContentCreateRequest {
  title: string
  content?: string
  note_id?: number
  status?: string // draft, published
}

export interface GenerateContentRequest {
  note_id?: number
  prompt?: string
  voice_dna?: string | null
  max_length?: number
  temperature?: number
  include_images?: boolean
  use_voice_profile?: boolean
}

export interface GenerateContentResponse {
  content_md: string
  status: string
  id?: number
}

// Ответ от функции generateContent
export interface GenerateResponse {
  content: string
  content_md?: string
  status: string
  post_id?: number
}

export const useContentStore = defineStore('content', () => {
  const toast = useToast()
  const profileStore = useProfileStore()
  
  // Состояния
  const contentList = ref<Content[]>([])
  const currentContent = ref<Content | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Получение списка контента
  const getContentList = async () => {
    isLoading.value = true
    error.value = null
    try {
      // Используем пустой путь, так как базовый URL уже содержит /api/v1/posts
      const response = await contentApi.get<any>('')
      
      // Преобразуем данные от API в формат для фронтенда
      const posts = response.data.posts || response.data || []
      contentList.value = posts.map(transformPostData)
      
      return contentList.value
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Ошибка при получении списка контента'
      error.value = errorMsg
      toast.error(errorMsg, { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Получение контента по ID
  const getContentById = async (id: string) => {
    isLoading.value = true
    error.value = null
    try {
      // Используем только ID, так как базовый URL уже содержит /api/v1/posts
      const response = await contentApi.get<any>(`/${id}`)
      
      // Преобразуем данные от API в формат для фронтенда
      currentContent.value = transformPostData(response.data)
      
      return currentContent.value
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Ошибка при получении контента'
      error.value = errorMsg
      toast.error(errorMsg, { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Функция для преобразования данных поста от API в формат для фронтенда
  const transformPostData = (post: any): Content => {
    // Извлекаем заголовок из Markdown (первая строка h1)
    const titleMatch = post.content_md?.match(/^#\s+(.*)/m)
    const title = titleMatch ? titleMatch[1] : 'Без заголовка'
    
    // Извлекаем отрывок (первый параграф после заголовка)
    const excerptMatch = post.content_md?.match(/^(?!#)(.+)/m)
    const excerpt = excerptMatch ? excerptMatch[1].substring(0, 150) + '...' : 'Без описания'
    
    // Определяем, сгенерирован ли контент (наличие note_id)
    const isGenerated = !!post.note_id
    
    // Определяем, опубликован ли контент (по статусу)
    const isPublished = post.status === 'published'
    
    return {
      ...post,
      title,
      excerpt,
      is_generated: isGenerated,
      is_published: isPublished
    }
  }

  // Обновление контента
  const updateContent = async (id: string, data: Partial<Content>) => {
    isLoading.value = true
    error.value = null
    try {
      // Формируем данные для API
      const apiData: any = {}
      
      // Если есть обновление контента, преобразуем его в формат API
      if (data.content_md) {
        // Преобразуем контент в Markdown, добавляя заголовок
        const title = data.title || currentContent.value?.title || 'Без заголовка'
        apiData.content_md = `# ${title}\n\n${data.content_md}`
      } else if (data.title && currentContent.value) {
        // Если меняется только заголовок, обновляем его в существующем контенте
        const content = currentContent.value.content_md.replace(/^#\s+(.*)$/m, `# ${data.title}`)
        apiData.content_md = content
      }
      
      // Если нужно обновить статус
      if (data.is_published !== undefined) {
        apiData.status = data.is_published ? 'published' : 'draft'
      }
      
      // Используем только ID, так как базовый URL уже содержит /api/v1/posts
      const response = await contentApi.put<any>(`/${id}`, apiData)
      
      // Обновляем данные в хранилище
      const updatedContent = transformPostData(response.data)
      const index = contentList.value.findIndex(c => c.id.toString() === id.toString())
      if (index !== -1) {
        contentList.value[index] = updatedContent
      }
      currentContent.value = updatedContent
      
      return updatedContent
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Ошибка при обновлении контента'
      error.value = errorMsg
      toast.error(errorMsg, { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Публикация контента
  const publishContent = async (id: string) => {
    isLoading.value = true
    error.value = null
    try {
      // Используем только ID, так как базовый URL уже содержит /api/v1/posts
      const response = await contentApi.put<any>(`/${id}`, {
        status: 'published'
      })
      
      // Обновляем данные в хранилище
      const updatedContent = transformPostData(response.data)
      const index = contentList.value.findIndex(c => c.id.toString() === id.toString())
      if (index !== -1) {
        contentList.value[index] = updatedContent
      }
      if (currentContent.value?.id.toString() === id.toString()) {
        currentContent.value = updatedContent
      }
      
      return updatedContent
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Ошибка при публикации контента'
      error.value = errorMsg
      toast.error(errorMsg, { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Создание нового контента
  const createContent = async (data: ContentCreateRequest) => {
    isLoading.value = true
    error.value = null
    try {
      // Формируем данные для API
      const apiData: any = {
        status: data.status || 'draft'
      }
      
      // Если есть контент, преобразуем его в формат API
      if (data.content) {
        // Преобразуем контент в Markdown, добавляя заголовок
        apiData.content_md = `# ${data.title}\n\n${data.content}`
      }
      
      // Если есть ID заметки, добавляем его
      if (data.note_id) {
        apiData.note_id = data.note_id
      }
      
      // Используем пустой путь, так как базовый URL уже содержит /api/v1/posts
      const response = await contentApi.post<any>('', apiData)
      
      // Преобразуем данные и обновляем хранилище
      const newContent = transformPostData(response.data)
      contentList.value = [newContent, ...contentList.value]
      
      return newContent
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Ошибка при создании контента'
      error.value = errorMsg
      toast.error(errorMsg, { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Удаление контента
  const deleteContent = async (id: string) => {
    isLoading.value = true
    error.value = null
    try {
      // Используем только ID, так как базовый URL уже содержит /api/v1/posts
      await contentApi.delete(`/${id}`)
      contentList.value = contentList.value.filter(c => c.id !== id)
      if (currentContent.value?.id === id) {
        currentContent.value = null
      }
      return true
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Ошибка при удалении контента'
      error.value = errorMsg
      toast.error(errorMsg, { duration: 5000 })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Генерация контента с использованием сервиса генерации
  const generateContent = async (request: GenerateContentRequest): Promise<GenerateResponse> => {
    isLoading.value = true;
    error.value = null;
    try {
      // Проверяем наличие действительного профиля для всех запросов
      try {
        if (!profileStore.hasVoiceProfile) {
          throw new Error('Голосовой профиль не настроен');
        }
      } catch (profileError: any) {
        // Если ошибка связана с отсутствием профиля, выдаем понятное сообщение
        if (profileError.response?.status === 404 || 
            profileError.message === 'Голосовой профиль не настроен') {
          throw new Error('Для генерации контента необходимо настроить голосовой профиль');
        }
        // Иначе пробрасываем исходную ошибку
        throw profileError;
      }
      
      let noteId = request.note_id;
      
      // Если у нас есть prompt, но нет note_id, сначала создаем заметку
      if (request.prompt && !noteId) {
        try {
          // Создаем заметку через API заметок
          const noteResponse = await contentApi.post<any>('/notes', {
            text: request.prompt,
            source: 'text' // Источник: text или voice
          });
          
          if (noteResponse.data && noteResponse.data.id) {
            noteId = noteResponse.data.id;
          } else {
            throw new Error('Не удалось создать заметку');
          }
        } catch (noteError) {
          console.error('Ошибка при создании заметки:', noteError);
          throw noteError;
        }
      }
      
      if (!noteId) {
        throw new Error('Необходим ID заметки или текст промпта для генерации контента');
      }
      
      // Отправляем запрос на генерацию поста с учетом настроек голосового профиля
      const response = await contentApi.post<GenerateContentResponse>('/generate', {
        note_id: noteId,
        use_voice_profile: request.use_voice_profile,
        include_images: request.include_images,
        max_length: request.max_length,
        temperature: request.temperature
      });
      
      // Если мы получили готовый пост, возвращаем его
      if (response.data.id) {
        // Получаем полные данные поста
        const postResponse = await contentApi.get<any>(`/${response.data.id}`);
        return {
          content: postResponse.data.content_md,
          status: 'success',
          post_id: postResponse.data.id
        };
      }
      
      // Иначе возвращаем сгенерированный контент
      return {
        content: response.data.content_md || '',
        status: response.data.status || 'success'
      };
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || err.message || 'Ошибка при генерации контента';
      error.value = errorMsg;
      toast.error(errorMsg, { duration: 5000 });
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    contentList,
    currentContent,
    isLoading,
    error,
    getContentList,
    getContentById,
    createContent,
    updateContent,
    publishContent,
    deleteContent,
    generateContent
  }
}) 