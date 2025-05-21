<template>
  <div class="border border-gray-200 rounded-lg p-6 bg-white">
    <div class="flex justify-between mb-2">
      <router-link :to="`/content/${id}`" class="text-lg font-semibold text-gray-900 hover:text-indigo-600">
        {{ title }}
      </router-link>
      <span v-if="isGenerated" class="bg-indigo-100 text-indigo-800 text-xs px-2 py-1 rounded">
        Сгенерировано
      </span>
    </div>
    <p class="text-gray-600 mb-4">{{ excerpt }}</p>
    <div class="flex justify-between items-center">
      <span class="text-sm text-gray-500">
        {{ date }}
      </span>
      <router-link :to="`/content/${id}`" class="text-sm text-indigo-600 hover:text-indigo-800">
        Просмотреть
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  id: string
  title: string
  excerpt: string
  date: string
  isGenerated?: boolean
}>()
</script> 