package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения
type Config struct {
	App struct {
		Name      string `env:"APP_NAME" envDefault:"voy-api"`
		Port      int    `env:"PORT" envDefault:"8080"`
		Mode      string `env:"APP_MODE" envDefault:"development"` // development или production
		ServiceID string
	}

	Database struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     int    `env:"DB_PORT" envDefault:"5432"`
		Name     string `env:"DB_NAME" envDefault:"voy_db"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASS" envDefault:"postgres"`
		SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
	}

	Auth struct {
		JWTSecret       string        `env:"JWT_SECRET" envDefault:"super-secret-key-change-in-production"`
		TokenExpiration time.Duration `env:"TOKEN_EXPIRATION" envDefault:"24h"`
		CookieDomain    string        `env:"COOKIE_DOMAIN"`
		CookieSecure    bool          `env:"COOKIE_SECURE" envDefault:"false"`
		SessionsCleanup time.Duration `env:"SESSIONS_CLEANUP" envDefault:"1h"`
	}

	CORS struct {
		// Дополнительные разрешенные источники (домены), разделенные запятыми
		AllowedOrigins string `env:"CORS_ALLOWED_ORIGINS"`
	}

	Services struct {
		WhisperAPIKey string `env:"WHISPER_API_KEY"`
		LLMAPIURL     string `env:"LLM_API_URL" envDefault:"https://api.openai.com/v1"`
		LLMAPIKey     string `env:"LLM_API_KEY"`
	}
}

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() (*Config, error) {
	return LoadConfigWithServiceID("")
}

// LoadConfigWithServiceID загружает конфигурацию для конкретного сервиса
func LoadConfigWithServiceID(serviceID string) (*Config, error) {
	// Загружаем .env файл
	workDir, _ := os.Getwd()

	// Пытаемся загрузить .env из текущей директории
	_ = godotenv.Load(filepath.Join(workDir, ".env"))

	// Пытаемся загрузить .env из родительской директории (если вызов из вложенных папок)
	_ = godotenv.Load(filepath.Join(workDir, "../.env"))
	_ = godotenv.Load(filepath.Join(workDir, "../../.env"))

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Устанавливаем ID сервиса
	cfg.App.ServiceID = serviceID

	// Если указан serviceID, используем соответствующий порт
	if serviceID != "" {
		// Определяем порт на основе serviceID
		switch strings.ToLower(serviceID) {
		case "auth":
			envPort := os.Getenv("AUTH_PORT")
			if envPort != "" {
				// Используем переменную окружения, если задана
				cfg.App.Port = parseInt(envPort, 8080)
			} else {
				cfg.App.Port = 8080
			}
		case "profile":
			envPort := os.Getenv("PROFILE_PORT")
			if envPort != "" {
				cfg.App.Port = parseInt(envPort, 8081)
			} else {
				cfg.App.Port = 8081
			}
		case "orchestrator":
			envPort := os.Getenv("ORCHESTRATOR_PORT")
			if envPort != "" {
				cfg.App.Port = parseInt(envPort, 8082)
			} else {
				cfg.App.Port = 8082
			}
		}
	}

	return cfg, nil
}

// parseInt преобразует строку в int с проверкой ошибок
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	var value int
	_, err := fmt.Sscanf(s, "%d", &value)
	if err != nil {
		return defaultValue
	}

	return value
}

// IsDevelopment проверяет, запущено ли приложение в режиме разработки
func (c *Config) IsDevelopment() bool {
	return c.App.Mode == "development"
}

// IsProduction проверяет, запущено ли приложение в режиме продакшена
func (c *Config) IsProduction() bool {
	return c.App.Mode == "production"
}

// GetCORSAllowedOrigins возвращает список разрешенных источников из переменной окружения
func (c *Config) GetCORSAllowedOrigins() []string {
	if c.CORS.AllowedOrigins == "" {
		return []string{}
	}

	// Разделяем строку по запятым
	origins := strings.Split(c.CORS.AllowedOrigins, ",")

	// Убираем лишние пробелы
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return origins
}

// GetDSN возвращает строку подключения к базе данных
func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + string(c.Database.Port) +
		" user=" + c.Database.User +
		" dbname=" + c.Database.Name +
		" password=" + c.Database.Password +
		" sslmode=" + c.Database.SSLMode
}
