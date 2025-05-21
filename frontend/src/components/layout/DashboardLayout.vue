<template>
  <div class="flex min-h-screen bg-gray-50">
    <Sidebar />
    
    <div class="ml-64 w-full">
      <main class="py-6 px-8 w-full">
        <slot></slot>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
</script> 