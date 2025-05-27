<template>
  <div v-if="!hasProfile && !isLoading" class="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-6">
    <div class="flex">
      <div class="flex-shrink-0">
        <svg class="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
      </div>
      <div class="ml-3">
        <p class="text-sm text-yellow-700">
          {{ message || 'Голосовой профиль не настроен. Вы не сможете генерировать контент в вашем собственном стиле.' }}
          <router-link v-if="showLink" to="/voice-dna" class="font-medium underline text-yellow-700 hover:text-yellow-600">
            Настроить голосовой профиль
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, defineProps } from 'vue'
import { useProfileStore } from '@/stores/profile'

const props = defineProps({
  message: {
    type: String,
    default: ''
  },
  showLink: {
    type: Boolean,
    default: true
  }
})

const profileStore = useProfileStore()
const hasProfile = ref(false)
const isLoading = ref(true)

// Функция для проверки наличия и валидности голосового профиля
const checkVoiceProfile = () => {
  const profile = profileStore.profile
  
  // Проверяем, что dna_data существует и содержит реальные данные (не пустой объект)
  if (profile?.dna_data) {
    try {
      // Пытаемся распарсить JSON, чтобы убедиться, что это не пустой объект
      const dnaData = JSON.parse(profile.dna_data)
      hasProfile.value = dnaData && 
                        Object.keys(dnaData).length > 0 && 
                        dnaData.profile && 
                        dnaData.style
    } catch (e) {
      console.error('Ошибка при парсинге dna_data:', e)
      hasProfile.value = false
    }
  } else {
    hasProfile.value = false
  }
}

// Следим за изменениями профиля
watch(() => profileStore.profile, (newProfile) => {
  checkVoiceProfile()
}, { deep: true })

onMounted(async () => {
  isLoading.value = true
  try {
    await profileStore.getProfile()
    checkVoiceProfile()
  } catch (error) {
    console.error('Ошибка при получении данных профиля:', error)
    hasProfile.value = false
  } finally {
    isLoading.value = false
  }
})
</script> 