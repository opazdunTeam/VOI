import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  // Настройка прокси для API запросов
  server: {
    proxy: {
      '/api/v1/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/api/v1/profile': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
      '/api/v1/content': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      }
    },
    // Расширенные настройки HMR для устранения проблем с WebSocket
    hmr: {
      overlay: false, // Отключаем оверлей ошибок
      protocol: 'ws', // Явно указываем протокол WebSocket
      host: 'localhost', // Явно указываем хост
      port: 5173, // Совпадает с портом сервера разработки
      clientPort: 5173, // Порт на стороне клиента (для прокси-случаев)
      timeout: 10000, // Увеличиваем тайм-аут для соединения
    },
    // Дополнительные настройки сервера
    host: true, // Слушать все сетевые интерфейсы (0.0.0.0)
    strictPort: true, // Строгое использование указанного порта
    port: 5173, // Порт Vite сервера
  },
  // Оптимизации для развертывания
  build: {
    // Минимизация кода
    minify: 'terser',
    // Разделение кода
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          tailwind: ['tailwindcss']
        }
      }
    },
    // Предварительная загрузка модулей
    modulePreload: true,
  },
  // Оптимизация разработки
  optimizeDeps: {
    include: ['vue', 'vue-router', 'pinia', 'axios'],
  }
})
