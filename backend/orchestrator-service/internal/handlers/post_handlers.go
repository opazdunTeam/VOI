package handlers

import (
	"net/http"
	"strconv"

	"backend/orchestrator-service/internal/dto"
	"backend/orchestrator-service/internal/service"
	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PostHandlers обработчики для постов
type PostHandlers struct {
	service service.PostService
	config  *config.Config
	logger  *logger.Logger
}

// NewPostHandlers создает новые обработчики постов
func NewPostHandlers(service service.PostService, cfg *config.Config, logger *logger.Logger) *PostHandlers {
	return &PostHandlers{
		service: service,
		config:  cfg,
		logger:  logger.Named("post_handlers"),
	}
}

// SetupPostRoutes настраивает маршруты для постов
func SetupPostRoutes(router *gin.RouterGroup, handlers *PostHandlers, cfg *config.Config, db *gorm.DB) {
	// Защищенные маршруты
	posts := router.Group("/posts")
	posts.Use(middleware.AuthMiddleware(cfg, db))
	{
		posts.POST("/generate", handlers.GeneratePost)
		posts.GET("", handlers.GetPosts)
		posts.GET("/:id", handlers.GetPost)
		posts.PUT("/:id", handlers.UpdatePost)
		posts.DELETE("/:id", handlers.DeletePost)

		// Эндпоинты для заметок
		posts.POST("/notes", handlers.CreateNote)
		posts.GET("/notes", handlers.GetNotes)
	}
}

// GeneratePost генерирует новый пост
// @Summary Сгенерировать пост
// @Description Генерирует новый пост на основе заметки
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GeneratePostRequest true "Данные для генерации"
// @Success 201 {object} dto.PostResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts/generate [post]
func (h *PostHandlers) GeneratePost(c *gin.Context) {
	// Получаем ID пользователя из контекста (установлен middleware)
	userClaims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.GeneratePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	post, err := h.service.GeneratePost(c.Request.Context(), userClaims.UserID, &req)
	if err != nil {
		if err.Error() == "note not found" || err.Error() == "note not found or not owned by user" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		h.logger.Error("Failed to generate post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to generate post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPosts получает список постов пользователя
// @Summary Получить список постов
// @Description Получает список постов текущего пользователя
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы"
// @Param size query int false "Размер страницы"
// @Success 200 {object} dto.PostListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts [get]
func (h *PostHandlers) GetPosts(c *gin.Context) {
	// Получаем ID пользователя из контекста (установлен middleware)
	userClaims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	posts, err := h.service.GetPosts(c.Request.Context(), userClaims.UserID, page, size)
	if err != nil {
		h.logger.Error("Failed to get posts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost получает пост по ID
// @Summary Получить пост
// @Description Получает пост по ID
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Success 200 {object} dto.PostResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts/{id} [get]
func (h *PostHandlers) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid post ID"})
		return
	}

	post, err := h.service.GetPost(c.Request.Context(), postID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "post not found"})
			return
		}
		h.logger.Error("Failed to get post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get post"})
		return
	}

	// Проверяем, что пост принадлежит пользователю
	userClaims, exists := middleware.GetCurrentUser(c)
	if !exists || post.UserID != userClaims.UserID {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost удаляет пост
// @Summary Удалить пост
// @Description Удаляет пост по ID
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Success 204 "No Content"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts/{id} [delete]
func (h *PostHandlers) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid post ID"})
		return
	}

	// Получаем ID пользователя из контекста (установлен middleware)
	userClaims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	err = h.service.DeletePost(c.Request.Context(), postID, userClaims.UserID)
	if err != nil {
		if err.Error() == "post not found or not owned by user" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "post not found"})
			return
		}
		h.logger.Error("Failed to delete post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to delete post"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdatePost обновляет пост
// @Summary Обновить пост
// @Description Обновляет пост по ID
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Param request body dto.UpdatePostRequest true "Данные для обновления поста"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts/{id} [put]
func (h *PostHandlers) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid post ID"})
		return
	}

	// Получаем ID пользователя из контекста (установлен middleware)
	userClaims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	post, err := h.service.UpdatePost(c.Request.Context(), postID, userClaims.UserID, &req)
	if err != nil {
		if err.Error() == "post not found" || err.Error() == "post not found or not owned by user" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "post not found"})
			return
		}
		h.logger.Error("Failed to update post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// CreateNote создает новую заметку
// @Summary Создать заметку
// @Description Создает новую заметку
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.NoteRequest true "Данные для создания заметки"
// @Success 201 {object} dto.NoteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts/notes [post]
func (h *PostHandlers) CreateNote(c *gin.Context) {
	// Получаем ID пользователя из контекста (установлен middleware)
	userClaims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.NoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	note, err := h.service.CreateNote(c.Request.Context(), userClaims.UserID, &req)
	if err != nil {
		h.logger.Error("Failed to create note", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetNotes получает список заметок пользователя
// @Summary Получить список заметок
// @Description Получает список заметок текущего пользователя
// @Tags notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы"
// @Param size query int false "Размер страницы"
// @Success 200 {object} dto.NoteListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts/notes [get]
func (h *PostHandlers) GetNotes(c *gin.Context) {
	// Получаем ID пользователя из контекста (установлен middleware)
	userClaims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	notes, err := h.service.GetNotes(c.Request.Context(), userClaims.UserID, page, size)
	if err != nil {
		h.logger.Error("Failed to get notes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}
