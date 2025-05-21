# Руководство-план по реализации базовой версии платформы **VOY** (MVP)
## Оглавление
1. Введение и стратегия
2. Послойный план разработки (Road-map)
3. Архитектурная карта сервисов и маршрутов
4. Модели данных и схема БД (черновик)
5. Руководство для фронтенд-команды
6. Руководство для бэкенд-команды
7. Настройка окружений и конфигурации

---

## 1. Введение и стратегия

| Роль | Цель в MVP |
|------|------------|
| Бизнес | Проверить гипотезу персонализированной генерации контента с минимальными затратами |
| Техлид | Обеспечить надёжный фундамент, позволяющий расширять функциональность слоями |
| Разработчики | Следовать гайдлайнам, минимизировать «технический креатив» и расхождения |

Ключевой принцип — **vertical slice**: каждая задача приводит к полезному, демонстрируемому результату (end-to-end). Слои реализуются последовательно, но с минимальными блокировками между командами.

---

## 2. Послойный план разработки (Road-map)

### Слой 0. Подготовительный (⏱ 1 день)
1. Настроить базовые проекты для фронтенда и бэкенда.
2. Настроить локальную среду разработки: PostgreSQL для разработки.
3. Настроить Git Flow + шаблоны PR + базовые линтеры.
4. Создать шаблоны `.env` файлов и систему управления конфигурацией.

### Слой 1. Инфраструктурный (⏱ 2 дня)
1. **Backend**
   - Шаблон микросервиса на Go + Gin (`cmd`, `internal`, `pkg`).
   - Библиотека общих компонентов: логгер (zap), middleware (CORS, Cookie+JWT, Session), utils.
   - Настройка Swagger для документации API.
   - Настройка GORM для работы с БД.
   - Настройка автосоздания БД и базовой миграции при старте.
   - Загрузка и валидация конфигурации из `.env` файлов.
   - Реализация разделения сред разработки и продакшн.
2. **Frontend**
   - Базовое приложение (Vue 3).
   - Конфигурация маршрутизации (Vue Router).
   - Shared UI Library (основные компоненты).
   - Настройка переключения между dev/prod API через env переменные.

### Слой 2. Данные (⏱ 2 дня)
1. ER-модель: `users`, `voice_profiles`, `notes`, `posts`, `sessions`.
2. Миграции через GORM AutoMigrate.
3. Тестовые данные для разработки.
4. Скрипты для автоматического восстановления БД в dev-окружении.

### Слой 3. Бэкенд сервисы (⏱ 6 дней)
| Приоритет | Сервис | Функция | Эндпоинты v1 |
|-----------|--------|---------|--------------|
| P0 | **Auth** | Регистрация/вход, Cookie+JWT, управление сессиями | `POST /auth/register`, `POST /auth/login`, `POST /auth/logout`, `GET /auth/me` |
| P0 | **Profile** | Voice DNA CRUD | `GET /profile`, `PUT /profile` |
| P0 | **Orchestrator** | Управление генерацией | `POST /posts/generate`, `GET /posts`, `GET /posts/:id` |
| P1 | **Speech** | Транскрибация аудио (Whisper) | `POST /speech/transcribe` |
| P1 | **Generator** | Интеграция с LLM | `POST /internal/generate` (internal) |

> Сервисы P0 необходимы для «сквозного» сценария текста → пост. P1 можно симулировать (мок) до готовности.

### Слой 4. Фронтенд (⏱ 6 дней)
1. **Ключевые страницы**
   - `/login` / `/register`
   - `/dashboard` (лендинг после логина)
   - `/onboarding` (мастер Voice DNA)
   - `/posts/new` (создание поста)
   - `/posts/:id/edit` (редактирование поста)
   - `/voice-dna` (настройка профиля)
2. **Основные компоненты**
   - `NoteInput.vue` (текст/аудио)
   - `PostGenerator.vue` (progress bar, state machine)
   - `PostCard.vue` + `PostEditor.vue`
3. Сервис-слой (Axios): `authService` (с поддержкой куки), `postService`, `profileService`.
4. Глобальное состояние (Pinia stores): `useUserStore`, `usePostStore`, `useProfileStore`.

### Слой 5. Интеграция (⏱ 2 дня)
1. Ручное тестирование основных потоков.
2. Исправление критичных ошибок.

### Слой 6. Релиз (⏱ 1 день)
1. Настройка простого хостинга для фронтенда.
2. Размещение бэкенда на простом хостинге (типа Heroku, Render).
3. Базовое руководство для пользователей.

> **Итого MVP:** ~18 человеко-дней при слаженной работе.

---

## 3. Архитектурная карта сервисов и маршрутов

### 3.1 Backend REST API (public)
```http
POST   /api/v1/auth/register         # создать пользователя
POST   /api/v1/auth/login            # получить пары токенов в куки + создать сессию
POST   /api/v1/auth/logout           # удалить куки + инвалидировать сессию 
GET    /api/v1/auth/me               # профиль текущего пользователя

GET    /api/v1/profile               # получить Voice DNA
PUT    /api/v1/profile               # обновить Voice DNA

POST   /api/v1/posts/generate        # инициировать генерацию
GET    /api/v1/posts                 # список постов (пагинация)
GET    /api/v1/posts/:id             # детальный просмотр
PUT    /api/v1/posts/:id             # редактировать
DELETE /api/v1/posts/:id             # удалить

POST   /api/v1/speech/transcribe     # вернуть текст заметки
```

### 3.2 Внутренние маршруты (service-to-service)
```http
POST /internal/generate              # Generator ⇐ Orchestrator
```
> Защищается API-ключами или внутренней сетью.

### 3.3 Frontend маршрутизация (Vue Router)
| Path | Компонент | Auth? |
|------|-----------|-------|
| `/login` | `LoginPage` | ✖ |
| `/register` | `RegisterPage` | ✖ |
| `/onboarding` | `OnboardingWizard` | ✔ |
| `/dashboard` | `Dashboard` | ✔ |
| `/posts/new` | `NoteInput` | ✔ |
| `/posts/:id/edit` | `PostEditor` | ✔ |
| `/voice-dna` | `ProfileSettings` | ✔ |

---

## 4. Модели данных и схема БД (черновик)
```mermaid
erDiagram
    users {
      bigint id PK
      text   email  "uniq"
      text   password_hash
      text   full_name
      timestamptz created_at
    }
    sessions {
      bigint id PK
      bigint user_id FK "users.id" 
      text   token "jwt token hash"
      text   user_agent
      text   ip_address
      bool   is_active
      timestamptz expires_at
      timestamptz created_at
    }
    voice_profiles {
      bigint user_id PK FK "users.id"
      jsonb  dna_data
      timestamptz updated_at
    }
    notes {
      bigint id PK
      bigint user_id FK "users.id"
      text   original_text
      text   source  # voice|text
      timestamptz created_at
    }
    posts {
      bigint id PK
      bigint user_id FK "users.id"
      bigint note_id FK "notes.id"
      text   content_md
      text   status   # draft|ready
      timestamptz created_at
    }
```
> Все модели GORM с автоматическим созданием через AutoMigrate.

---

## 5. Руководство для фронтенд-команды
1. **Запуск локально**
   ```bash
   cd voy-frontend
   npm i && npm run dev
   ```
2. **Создание нового UI компонента**: следовать структуре в `/components`.
3. **Обращение к API**: через `apiClient` (Axios, с поддержкой credentials: 'include' для работы с куки).
4. **Code style**: Prettier + ESLint `@vue/eslint-config-typescript`.
5. **Окружения**:
   ```js
   // .env.development
   VITE_API_BASE_URL=http://localhost:8080/api/v1
   
   // .env.production
   VITE_API_BASE_URL=/api/v1
   ```

---

## 6. Руководство для бэкенд-команды
1. **Шаблон сервиса**: использовать `template-service` как отправную точку.
2. **Middlewares**: `RequestID`, `Recovery`, `Logger`, `SessionAuth` (куки + JWT).
3. **ORM**: GORM для работы с базой данных и автомиграций.
4. **Документация API**: Swagger (swaggo) для автоматического генерирования документации.
5. **Сессии**: JWT токены хранятся в secure httpOnly cookie с проверкой в БД.
6. **Автосоздание БД**: при старте в dev-режиме проверять и создавать БД при необходимости.

---

## 7. Настройка окружений и конфигурации

### 7.1 Структура конфигурации бэкенда
```go
type Config struct {
    App struct {
        Name string `env:"APP_NAME" envDefault:"voy-api"`
        Port int    `env:"PORT" envDefault:"8080"`
        Mode string `env:"APP_MODE" envDefault:"development"` // development или production
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
        JWTSecret        string        `env:"JWT_SECRET" envDefault:"super-secret-key-change-in-production"`
        TokenExpiration  time.Duration `env:"TOKEN_EXPIRATION" envDefault:"24h"`
        CookieDomain     string        `env:"COOKIE_DOMAIN"`
        CookieSecure     bool          `env:"COOKIE_SECURE" envDefault:"false"`
        SessionsCleanup  time.Duration `env:"SESSIONS_CLEANUP" envDefault:"1h"`
    }
    
    Services struct {
        WhisperAPIKey string `env:"WHISPER_API_KEY"`
        LLMAPIURL     string `env:"LLM_API_URL" envDefault:"https://api.openai.com/v1"`
        LLMAPIKey     string `env:"LLM_API_KEY"` 
    }
}
```

### 7.2 Загрузка конфигурации
```go
import "github.com/caarlos0/env/v6"

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}
```

### 7.3 Автосоздание БД (пример для dev-режима)
```go
func EnsureDatabaseExists(config *Config) error {
    // Только для dev-режима
    if config.App.Mode != "development" {
        return nil
    }
    
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
        config.Database.Host, config.Database.Port, 
        config.Database.User, config.Database.Password, config.Database.SSLMode)
    
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return err
    }
    defer db.Close()
    
    // Проверяем существование БД
    var exists bool
    query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", config.Database.Name)
    err = db.QueryRow(query).Scan(&exists)
    if err != nil && err != sql.ErrNoRows {
        return err
    }
    
    // Создаем БД если не существует
    if !exists {
        _, err = db.Exec("CREATE DATABASE " + config.Database.Name)
        if err != nil {
            return err
        }
        log.Printf("База данных %s создана", config.Database.Name)
    }
    
    return nil
}
```

### 7.4 Шаблоны .env файлов

**.env.example (для разработчика)**
```
# Основные настройки
APP_NAME=voy-api
PORT=8080
APP_MODE=development

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_NAME=voy_db
DB_USER=postgres
DB_PASS=postgres
DB_SSL_MODE=disable

# Аутентификация
JWT_SECRET=development-jwt-secret-key-change-me
TOKEN_EXPIRATION=24h
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false
SESSIONS_CLEANUP=1h

# Внешние сервисы
WHISPER_API_KEY=your-whisper-api-key
LLM_API_URL=https://api.openai.com/v1
LLM_API_KEY=your-openai-api-key
```

**.env.production.example (для деплоя)**
```
# Основные настройки
APP_NAME=voy-api
PORT=80
APP_MODE=production

# База данных
DB_HOST=db.example.com
DB_PORT=5432
DB_NAME=voy_production
DB_USER=voy_user
DB_PASS=strong-password
DB_SSL_MODE=require

# Аутентификация
JWT_SECRET=very-secure-random-string-min-32-chars
TOKEN_EXPIRATION=12h
COOKIE_DOMAIN=api.example.com
COOKIE_SECURE=true
SESSIONS_CLEANUP=6h

# Внешние сервисы
WHISPER_API_KEY=actual-whisper-api-key
LLM_API_URL=https://api.openai.com/v1
LLM_API_KEY=actual-openai-api-key
```

> **Важно**: Никогда не коммитить реальные `.env` файлы в репозиторий. Использовать `.env.example` как шаблон.

> Документ обновляется по мере эволюции архитектуры. Все PR, влияющие на публичные контракты, должны включать изменение данного файла. 