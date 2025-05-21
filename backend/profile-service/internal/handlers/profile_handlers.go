package handlers

import (
	"net/http"

	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/middleware"
	"backend/profile-service/internal/dto"
	"backend/profile-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProfileHandlers обработчики для профилей
type ProfileHandlers struct {
	service service.ProfileService
	config  *config.Config
	logger  *logger.Logger
}

// NewProfileHandlers создает новые обработчики профилей
func NewProfileHandlers(service service.ProfileService, cfg *config.Config, logger *logger.Logger) *ProfileHandlers {
	return &ProfileHandlers{
		service: service,
		config:  cfg,
		logger:  logger.Named("profile_handlers"),
	}
}

// SetupProfileRoutes настраивает маршруты для профилей
func SetupProfileRoutes(router *gin.RouterGroup, handlers *ProfileHandlers, cfg *config.Config, db *gorm.DB) {
	// Защищенные маршруты
	profile := router.Group("/profile")
	profile.Use(middleware.AuthMiddleware(cfg, db))
	{
		profile.GET("", handlers.GetProfile)
		profile.PUT("", handlers.UpdateProfile)
	}
}

// GetProfile получает профиль текущего пользователя
// @Summary Получить профиль
// @Description Получает профиль текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ProfileResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /profile [get]
func (h *ProfileHandlers) GetProfile(c *gin.Context) {
	// Получаем данные пользователя из контекста (установлен middleware)
	claims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), claims.UserID)
	if err != nil {
		if err.Error() == "profile not found" {
			// Если профиль не найден, создаем новый с пустыми данными
			emptyProfile := &dto.ProfileRequest{
				DNAData: "{}", // Пустой JSON объект
			}
			profile, err = h.service.CreateProfile(c.Request.Context(), claims.UserID, emptyProfile)
			if err != nil {
				h.logger.Error("Failed to create profile", zap.Error(err))
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create profile"})
				return
			}
			c.JSON(http.StatusCreated, profile)
			return
		}
		h.logger.Error("Failed to get profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile обновляет профиль текущего пользователя
// @Summary Обновить профиль
// @Description Обновляет профиль текущего пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ProfileRequest true "Данные профиля"
// @Success 200 {object} dto.ProfileResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /profile [put]
func (h *ProfileHandlers) UpdateProfile(c *gin.Context) {
	// Получаем данные пользователя из контекста (установлен middleware)
	claims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	// Проверяем, существует ли профиль
	_, err := h.service.GetProfile(c.Request.Context(), claims.UserID)
	if err != nil {
		if err.Error() == "profile not found" {
			// Профиль не существует, создаем новый
			profile, err := h.service.CreateProfile(c.Request.Context(), claims.UserID, &req)
			if err != nil {
				h.logger.Error("Failed to create profile", zap.Error(err))
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create profile"})
				return
			}
			c.JSON(http.StatusCreated, profile)
			return
		}
		h.logger.Error("Failed to check profile existence", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	// Профиль существует, обновляем его
	profile, err := h.service.UpdateProfile(c.Request.Context(), claims.UserID, &req)
	if err != nil {
		h.logger.Error("Failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
