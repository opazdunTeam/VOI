<template>
  <router-link
    :to="to"
    class="flex items-center px-4 py-3 rounded-lg mb-1 transition-colors"
    :class="active 
      ? 'bg-white/20 text-white font-medium' 
      : 'text-white/80 hover:bg-white/10'"
  >
    <span class="mr-3" v-html="icon"></span>
    <span>{{ label }}</span>
  </router-link>
</template>

<script setup lang="ts">
defineProps<{
  to: string
  icon: string
  label: string
  active: boolean
}>()
</script> 