<template>
  <div class="min-h-screen flex items-center justify-center p-4 bg-gray-50">
    <div class="text-center">
      <div class="mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-24 w-24 mx-auto text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h1 class="text-4xl font-bold text-gray-900 mb-2">404</h1>
      <p class="text-xl text-gray-600 mb-6">Страница не найдена</p>
      <p class="text-gray-500 mb-8">Запрашиваемая страница не существует или была перемещена.</p>
      <router-link to="/dashboard" class="inline-flex items-center justify-center px-6 py-3 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-indigo-600 hover:bg-indigo-700">
        Вернуться на дашборд
      </router-link>
    </div>
  </div>
</template> 