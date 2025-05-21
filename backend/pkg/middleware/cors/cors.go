package cors

import (
	"backend/pkg/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config содержит настройки CORS для разных окружений
type Config struct {
	// Разрешенные источники (домены)
	AllowOrigins []string
	// Разрешенные методы HTTP
	AllowMethods []string
	// Разрешенные заголовки
	AllowHeaders []string
	// Заголовки, которые клиенты могут видеть
	ExposeHeaders []string
	// Разрешено ли передавать куки
	AllowCredentials bool
	// Время кэширования preflight запросов
	MaxAge time.Duration
}

// DefaultConfig возвращает настройки CORS по умолчанию
func DefaultConfig() Config {
	return Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// Middleware возвращает middleware для настройки CORS в зависимости от окружения
func Middleware(cfg *config.Config) gin.HandlerFunc {
	corsConfig := DefaultConfig()

	// Настройка AllowOrigins в зависимости от окружения
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

	// Добавляем дополнительные заголовки для работы с куки
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Cookie")
	corsConfig.ExposeHeaders = append(corsConfig.ExposeHeaders, "Set-Cookie")

	return cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowOrigins,
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders,
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
	})
}
