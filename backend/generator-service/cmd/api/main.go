package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/generator-service/internal/dto"       // Импортируем dto
	"backend/generator-service/internal/handler"   // Импортируем handler
	"backend/generator-service/internal/service"   // Импортируем service
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "backend/generator-service/cmd/api/docs"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)




// @title AI Content Generator API
// @version 1.0
// @description This is an API service for generating AI content using Ollama models.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// --- Настройка GeneratorService с Ollama ---
	// Укажите имя модели Llama, которую вы загрузили в Ollama.
	// Например, "llama3" или "llama3.2:3b"
	ollamaModelName := os.Getenv("OLLAMA_MODEL_NAME")
	if ollamaModelName == "" {
		ollamaModelName = "llama3.2:3b" // Значение по умолчанию, если переменная окружения не установлена
		log.Printf("OLLAMA_MODEL_NAME not set, defaulting to '%s'", ollamaModelName)
	}

	generatorSvc := service.NewGeneratorService(ollamaModelName)
	// --- Конец настройки GeneratorService ---


	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Routes
	api := router.Group("/api/v1")
	{
		// @Summary Generate AI Content
        // @Description Generates content using an Ollama language model, considering user-specific data.
        // @Tags content-generation
        // @Accept json
        // @Produce json
        // @Param request body dto.GenerateRequest true "Content generation request"
        // @Success 200 {object} dto.GenerateResponse "Successfully generated content"
        // @Failure 400 {object} object{error=string} "Bad Request"
        // @Failure 500 {object} object{error=string} "Internal Server Error"
        // @Router /internal/generate [post]
		api.POST("/internal/generate", func(c *gin.Context) {
			handler.HandleGenerate(c, generatorSvc) // Изменено здесь
		})

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}


	// Server setup
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

// handleGenerate теперь будет принимать GeneratorService
// Изменяем сигнатуру функции
func handleGenerate(c *gin.Context, generatorSvc service.GeneratorService) {
	var req dto.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем логику генерации через сервисный слой
	generatedContent, err := generatorSvc.GenerateContent(
		c.Request.Context(), // Передаем контекст из Gin
		req.Prompt,
		req.VoiceDNA,
	)
	if err != nil {
		log.Printf("Error generating content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content", "details": err.Error()})
		return
	}

	response := dto.GenerateResponse{
		Content: generatedContent,
		Status:  "success",
	}

	c.JSON(http.StatusOK, response)
}