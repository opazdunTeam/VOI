<template>
  <div class="bg-white p-6 rounded-lg shadow-sm">
    <h3 class="text-gray-500 text-sm font-medium">{{ title }}</h3>
    <p class="mt-2 text-3xl font-semibold">{{ value }}</p>
    <p v-if="subtitle" class="text-gray-500 text-sm mt-1">{{ subtitle }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  value: string | number
  subtitle?: string
}>()
</script> 