<template>
  <button
    :class="buttonClasses"
    :type="type"
    :disabled="disabled || isLoading"
    @click="emit('click')"
  >
    <div v-if="isLoading" class="flex items-center justify-center">
      <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <span>Загрузка...</span>
    </div>
    <slot v-else></slot>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'primary' | 'secondary' | 'outline' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  isLoading?: boolean
  disabled?: boolean
  fullWidth?: boolean
  type?: 'button' | 'submit' | 'reset'
  class?: string
}>(), {
  variant: 'primary',
  size: 'md',
  isLoading: false,
  disabled: false,
  fullWidth: false,
  type: 'button',
  class: ''
})

const emit = defineEmits<{
  (e: 'click'): void
}>()

// Базовые классы для всех кнопок
const baseClasses = 'font-medium rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2'

// Классы для разных вариантов кнопок
const variantClasses = {
  primary: 'bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500',
  secondary: 'bg-gray-600 text-white hover:bg-gray-700 focus:ring-gray-500',
  outline: 'border border-indigo-600 text-indigo-600 hover:bg-indigo-50 focus:ring-indigo-500',
  danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500'
}

// Классы для размеров кнопок
const sizeClasses = {
  sm: 'py-1.5 px-3 text-sm',
  md: 'py-2 px-4 text-base',
  lg: 'py-3 px-6 text-lg'
}

// Классы для ширины кнопки
const widthClasses = computed(() => props.fullWidth ? 'w-full' : '')

// Классы для состояния загрузки и отключения
const stateClasses = computed(() => {
  if (props.disabled || props.isLoading) {
    return 'opacity-70 cursor-not-allowed'
  }
  return 'cursor-pointer'
})

// Объединяем все классы
const buttonClasses = computed(() => 
  `${baseClasses} ${variantClasses[props.variant]} ${sizeClasses[props.size]} ${widthClasses.value} ${stateClasses.value} ${props.class}`
)
</script> 