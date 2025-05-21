package handlers

import (
	"net/http"

	"backend/auth-service/internal/dto"
	"backend/auth-service/internal/service"
	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Константы для работы с куками
const (
	authCookieName = "voy_auth"
)

// AuthHandlers содержит обработчики для аутентификации
type AuthHandlers struct {
	authService service.AuthService
	config      *config.Config
	logger      *logger.Logger
}

// NewAuthHandlers создает новый экземпляр AuthHandlers
func NewAuthHandlers(authService service.AuthService, cfg *config.Config, log *logger.Logger) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		config:      cfg,
		logger:      log.Named("auth_handlers"),
	}
}

// Register обрабатывает регистрацию нового пользователя
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.AuthResponse "User created"
// @Failure 400 {object} dto.ErrorResponse "Invalid input"
// @Failure 409 {object} dto.ErrorResponse "User already exists"
// @Failure 500 {object} dto.ErrorResponse "Server error"
// @Router /auth/register [post]
func (h *AuthHandlers) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: []string{err.Error()},
		})
		return
	}

	response, token, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch err {
		case service.ErrUserExists:
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "User already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	// Устанавливаем JWT в куки
	cookieMaxAge := int(h.config.Auth.TokenExpiration.Seconds())
	c.SetCookie(
		authCookieName,
		token,
		cookieMaxAge,
		"/",
		h.config.Auth.CookieDomain,
		h.config.Auth.CookieSecure,
		true, // HttpOnly
	)

	c.JSON(http.StatusCreated, response)
}

// Login обрабатывает вход пользователя
// @Summary Login user
// @Description Authenticate a user and return token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login data"
// @Success 200 {object} dto.AuthResponse "Successful login"
// @Failure 400 {object} dto.ErrorResponse "Invalid input"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Server error"
// @Router /auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: []string{err.Error()},
		})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	response, token, err := h.authService.Login(c.Request.Context(), req, userAgent, ipAddress)
	if err != nil {
		switch err {
		case service.ErrInvalidLogin:
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	// Устанавливаем JWT в куки
	cookieMaxAge := int(h.config.Auth.TokenExpiration.Seconds())
	c.SetCookie(
		authCookieName,
		token,
		cookieMaxAge,
		"/",
		h.config.Auth.CookieDomain,
		h.config.Auth.CookieSecure,
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, response)
}

// Logout обрабатывает выход пользователя
// @Summary Logout user
// @Description Logout user and invalidate token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest false "Logout options"
// @Success 200 {string} string "Logged out successfully"
// @Failure 401 {object} dto.ErrorResponse "Not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Server error"
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandlers) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Если тело запроса отсутствует, просто устанавливаем значения по умолчанию
		req = dto.LogoutRequest{All: false}
	}

	// Получаем пользователя из контекста (установлен middleware)
	claims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Not authenticated",
		})
		return
	}

	err := h.authService.Logout(c.Request.Context(), claims.UserID, claims.SessionID, req.All)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	// Удаляем куки
	c.SetCookie(
		authCookieName,
		"",
		-1,
		"/",
		h.config.Auth.CookieDomain,
		h.config.Auth.CookieSecure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetProfile получает профиль пользователя
// @Summary Get user profile
// @Description Get current user profile
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserDTO "User profile"
// @Failure 401 {object} dto.ErrorResponse "Not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Server error"
// @Security BearerAuth
// @Router /auth/me [get]
func (h *AuthHandlers) GetProfile(c *gin.Context) {
	claims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Not authenticated",
		})
		return
	}

	profile, err := h.authService.GetUserProfile(c.Request.Context(), claims.UserID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile обновляет профиль пользователя
// @Summary Update user profile
// @Description Update current user profile
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} dto.UserDTO "Updated user profile"
// @Failure 400 {object} dto.ErrorResponse "Invalid input"
// @Failure 401 {object} dto.ErrorResponse "Not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Server error"
// @Security BearerAuth
// @Router /auth/profile [put]
func (h *AuthHandlers) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: []string{err.Error()},
		})
		return
	}

	claims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Not authenticated",
		})
		return
	}

	profile, err := h.authService.UpdateProfile(c.Request.Context(), claims.UserID, req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ChangePassword обрабатывает запрос на изменение пароля
// @Summary Change password
// @Description Change user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} dto.ErrorResponse "Invalid input"
// @Failure 401 {object} dto.ErrorResponse "Not authenticated"
// @Failure 403 {object} dto.ErrorResponse "Invalid current password"
// @Failure 500 {object} dto.ErrorResponse "Server error"
// @Security BearerAuth
// @Router /auth/password [put]
func (h *AuthHandlers) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: []string{err.Error()},
		})
		return
	}

	claims, exists := middleware.GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Not authenticated",
		})
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), claims.UserID, req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "User not found",
			})
		case service.ErrInvalidPassword:
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "Invalid current password",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// SetupAuthRoutes настраивает маршруты для аутентификации
func SetupAuthRoutes(router *gin.RouterGroup, handlers *AuthHandlers, cfg *config.Config, db interface{}) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)

		// Защищенные маршруты с middleware
		secured := auth.Group("")
		secured.Use(middleware.AuthMiddleware(cfg, db.(*gorm.DB)))
		{
			secured.POST("/logout", handlers.Logout)
			secured.GET("/me", handlers.GetProfile)
			secured.PUT("/profile", handlers.UpdateProfile)
			secured.PUT("/password", handlers.ChangePassword)
		}
	}
}
