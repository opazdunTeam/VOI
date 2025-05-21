package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"backend/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Константы для работы с JWT и куки
const (
	authCookieName  = "voy_auth"
	bearerPrefix    = "Bearer "
	contextUserKey  = "user"
	contextTokenKey = "token"
)

// UserClaims представляет данные, хранимые в JWT-токене
type UserClaims struct {
	UserID    uint64 `json:"user_id"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware создает middleware для проверки авторизации
func AuthMiddleware(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из cookie
		tokenString, err := getTokenFromRequest(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		// Проверяем валидность токена
		claims := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Auth.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token"})
			c.Abort()
			return
		}

		// Проверяем существование сессии в БД
		var session struct {
			IsActive bool
		}
		result := db.Table("session").
			Select("is_active").
			Where("user_id = ? AND id = ? AND is_active = ? AND expires_at > ?",
				claims.UserID, claims.SessionID, true, time.Now()).
			Scan(&session)

		if result.Error != nil || !session.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: session expired or revoked"})
			c.Abort()
			return
		}

		// Устанавливаем данные пользователя в контекст
		c.Set(contextUserKey, claims)
		c.Set(contextTokenKey, tokenString)
		c.Next()
	}
}

// getTokenFromRequest извлекает JWT токен из запроса (куки или заголовка)
func getTokenFromRequest(c *gin.Context) (string, error) {
	// Сначала пробуем получить из куки
	if cookie, err := c.Cookie(authCookieName); err == nil {
		return cookie, nil
	}

	// Если нет в куки, пробуем из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header or cookie is missing")
	}

	// Проверяем префикс Bearer
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	// Извлекаем токен
	return authHeader[len(bearerPrefix):], nil
}

// GetCurrentUser извлекает пользователя из контекста
func GetCurrentUser(c *gin.Context) (*UserClaims, bool) {
	user, exists := c.Get(contextUserKey)
	if !exists {
		return nil, false
	}

	claims, ok := user.(*UserClaims)
	return claims, ok
}

// GetCurrentToken извлекает токен из контекста
func GetCurrentToken(c *gin.Context) (string, bool) {
	token, exists := c.Get(contextTokenKey)
	if !exists {
		return "", false
	}

	tokenStr, ok := token.(string)
	return tokenStr, ok
}
