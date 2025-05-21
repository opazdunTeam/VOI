<template>
  <DashboardLayout>
    <div>
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Дашборд</h1>
        <router-link to="/content/new">
          <Button variant="primary">
            Создать новый контент
          </Button>
        </router-link>
      </div>

      <!-- Секция статистики -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <StatsCard title="Всего контента" :value="23" />
        <StatsCard title="Сгенерировано в этом месяце" :value="14" subtitle="↑ 24% с прошлого месяца" />
        <StatsCard title="Статус голосового профиля" value="Активен" />
      </div>

      <!-- Недавний контент -->
      <div class="mb-8">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-900">Недавний контент</h2>
          <router-link to="/content" class="text-sm text-indigo-600 hover:text-indigo-800">
            Просмотреть все
          </router-link>
        </div>
        
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <ContentCard 
            v-for="item in recentContent" 
            :key="item.id"
            :id="item.id"
            :title="item.title"
            :excerpt="item.excerpt"
            :date="item.date"
            :is-generated="item.isGenerated"
          />
        </div>
      </div>

      <!-- Рекомендации -->
      <div>
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Рекомендации</h2>
        <div class="bg-indigo-50 border border-indigo-200 rounded-lg p-6">
          <div class="flex items-start">
            <div class="bg-indigo-100 p-3 rounded-full mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div>
              <h3 class="font-medium text-indigo-800 mb-1">Улучшите свой голосовой профиль</h3>
              <p class="text-indigo-700 mb-4">
                Для более качественной персонализации контента запишите дополнительные образцы голоса.
              </p>
              <router-link to="/voice-dna">
                <Button variant="outline" size="sm">
                  Обновить голосовой профиль
                </Button>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import DashboardLayout from '@/components/layout/DashboardLayout.vue'
import Button from '@/components/ui/Button.vue'
import StatsCard from '@/components/dashboard/StatsCard.vue'
import ContentCard from '@/components/dashboard/ContentCard.vue'

// Мок-данные для примера
const recentContent = ref([
  {
    id: '1',
    title: 'Будущее удаленной работы в технологических компаниях',
    excerpt: 'Исследование тенденций перехода к удаленной работе в технологическом секторе и его влияние на продуктивность и культуру компаний.',
    date: '2 часа назад',
    isGenerated: true
  },
  {
    id: '2',
    title: 'Влияние искусственного интеллекта на творческие профессии',
    excerpt: 'Анализ того, как ИИ может дополнять творческую работу, а не заменять ее, и как креативщики адаптируются к новым технологиям.',
    date: '1 день назад',
    isGenerated: true
  },
  {
    id: '3',
    title: 'Стратегии эффективного управления проектами в 2023 году',
    excerpt: 'Обзор современных подходов к управлению проектами с учетом гибридных команд и возросшей важности цифровых инструментов.',
    date: '3 дня назад',
    isGenerated: false
  },
  {
    id: '4',
    title: 'Эмоциональный интеллект в лидерстве: почему это важно',
    excerpt: 'Исследование роли эмоционального интеллекта в эффективном лидерстве и как развивать эти навыки в современной бизнес-среде.',
    date: '5 дней назад',
    isGenerated: true
  }
])
</script> 