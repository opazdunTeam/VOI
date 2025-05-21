package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "backend/orchestrator-service/docs"

	"backend/orchestrator-service/internal/handlers"
	"backend/orchestrator-service/internal/models"
	"backend/orchestrator-service/internal/repository"
	"backend/orchestrator-service/internal/service"
	"backend/pkg/config"
	"backend/pkg/db"
	"backend/pkg/logger"
	"backend/pkg/middleware/cors"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           VOY API - Orchestrator Service
// @version         1.0
// @description     Сервис оркестрации для платформы VOY

// @host      localhost:8082
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате "Bearer {token}"

func main() {
	// Загрузка конфигурации из .env файлов
	cfg, err := config.LoadConfigWithServiceID("orchestrator")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	log, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Настройка режима Gin
	if cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Инициализация подключения к базе данных
	database, err := db.NewDatabase(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Автоматическая миграция моделей
	if err := database.AutoMigrate(&models.Note{}, &models.Post{}); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Инициализация роутера Gin
	router := gin.New()

	// Использование middleware
	router.Use(gin.Recovery())
	router.Use(cors.Middleware(cfg))

	// Настройка маршрутов
	setupRoutes(router, cfg, database, log)

	// Настройка Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Создание HTTP сервера
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		log.Info("Starting server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}

// setupRoutes настраивает маршруты API
func setupRoutes(router *gin.Engine, cfg *config.Config, database *db.Database, log *logger.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API группа
	api := router.Group("/api/v1")
	{
		// Инициализация репозиториев
		postRepo := repository.NewPostRepository(database.DB, log)

		// Инициализация сервисов
		postService := service.NewPostService(postRepo, cfg, log)

		// Инициализация обработчиков
		postHandlers := handlers.NewPostHandlers(postService, cfg, log)

		// Настройка маршрутов для постов
		handlers.SetupPostRoutes(api, postHandlers, cfg, database.DB)
	}
}
