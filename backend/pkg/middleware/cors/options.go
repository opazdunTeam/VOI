package cors

import (
	"backend/pkg/config"
	"time"
)

// Option - функциональный опцион для настройки CORS
type Option func(*Config)

// WithOrigins добавляет дополнительные разрешенные источники
func WithOrigins(origins ...string) Option {
	return func(c *Config) {
		c.AllowOrigins = append(c.AllowOrigins, origins...)
	}
}

// WithMethods устанавливает разрешенные HTTP методы
func WithMethods(methods ...string) Option {
	return func(c *Config) {
		c.AllowMethods = methods
	}
}

// WithHeaders устанавливает разрешенные заголовки
func WithHeaders(headers ...string) Option {
	return func(c *Config) {
		c.AllowHeaders = headers
	}
}

// WithExposeHeaders устанавливает заголовки, видимые клиенту
func WithExposeHeaders(headers ...string) Option {
	return func(c *Config) {
		c.ExposeHeaders = headers
	}
}

// WithCredentials устанавливает разрешение на отправку куки
func WithCredentials(allow bool) Option {
	return func(c *Config) {
		c.AllowCredentials = allow
	}
}

// WithMaxAge устанавливает время кэширования preflight запросов
func WithMaxAge(duration time.Duration) Option {
	return func(c *Config) {
		c.MaxAge = duration
	}
}

// CustomMiddleware создает middleware с дополнительными настройками
func CustomMiddleware(cfg *config.Config, options ...Option) (Config, error) {
	// Получаем базовую конфигурацию в зависимости от окружения
	corsConfig := DefaultConfig()

	// Базовая настройка AllowOrigins в зависимости от окружения
	if cfg.IsDevelopment() {
		corsConfig.AllowOrigins = []string{
			"http://localhost:3000", // Фронтенд на React/Vue в dev режиме
			"http://localhost:5173", // Фронтенд на Vite в dev режиме
			"http://127.0.0.1:3000", // Альтернативные локальные адреса
			"http://127.0.0.1:5173",
		}
	} else if cfg.IsProduction() {
		corsConfig.AllowOrigins = []string{
			"https://voy.io",     // Основной домен
			"https://www.voy.io", // C www
			"https://app.voy.io", // Поддомен для приложения
		}
	}

	// Применяем дополнительные опции
	for _, option := range options {
		option(&corsConfig)
	}

	return corsConfig, nil
}

// EnvironmentSpecificConfig возвращает конфигурацию CORS с учетом переменных окружения
func EnvironmentSpecificConfig(cfg *config.Config, options ...Option) Config {
	// Создание базовой конфигурации
	corsConfig, _ := CustomMiddleware(cfg, options...)

	// Добавление дополнительных источников из переменных окружения, если они есть
	// Например, для тестирования или временных доменов
	if envOrigins := cfg.GetCORSAllowedOrigins(); len(envOrigins) > 0 {
		corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, envOrigins...)
	}

	return corsConfig
}
