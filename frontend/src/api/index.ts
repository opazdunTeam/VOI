// Экспортируем все API-клиенты из client.ts
export * from './client'

// Типы ошибок API
export interface ApiErrorResponse {
  code?: number;
  message?: string;
  details?: string[];
} 