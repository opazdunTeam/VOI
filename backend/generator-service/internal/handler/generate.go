package handler

import (
	"log" // Добавим log для вывода ошибок
	"net/http"

	"backend/generator-service/internal/dto"
	"backend/generator-service/internal/service"
	"github.com/gin-gonic/gin"
)

// HandleGenerate теперь принимает GeneratorService
func HandleGenerate(c *gin.Context, generatorSvc service.GeneratorService) {
	var req dto.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем логику генерации через сервисный слой
	// Передаем все параметры из req в GenerateContent
	generatedContent, err := generatorSvc.GenerateContent(
		c.Request.Context(), // Передаем контекст из Gin
		req.Prompt,
		req.VoiceDNA,
	)
	if err != nil {
		log.Printf("Error generating content: %v", err) // Логируем ошибку для отладки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content", "details": err.Error()})
		return
	}

	response := dto.GenerateResponse{
		Content: generatedContent,
		Status:  "success",
	}

	c.JSON(http.StatusOK, response)
}