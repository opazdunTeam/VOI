package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/auth-service/internal/dto"
	"backend/auth-service/internal/models"
	"backend/auth-service/internal/repository"
	"backend/pkg/config"
	"backend/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidLogin    = errors.New("invalid email or password")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidToken    = errors.New("invalid or expired token")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrInternalServer  = errors.New("internal server error")
	ErrInvalidPassword = errors.New("invalid current password")
)

// AuthService предоставляет методы для аутентификации
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, string, error)
	Login(ctx context.Context, req dto.LoginRequest, userAgent, ipAddress string) (*dto.AuthResponse, string, error)
	Logout(ctx context.Context, userID uint64, sessionID string, all bool) error
	GetUserProfile(ctx context.Context, userID uint64) (*dto.UserDTO, error)
	UpdateProfile(ctx context.Context, userID uint64, req dto.UpdateProfileRequest) (*dto.UserDTO, error)
	ChangePassword(ctx context.Context, userID uint64, req dto.ChangePasswordRequest) error
	CleanupExpiredSessions(ctx context.Context) error
}

// AuthServiceImpl реализует AuthService
type AuthServiceImpl struct {
	userRepo repository.UserRepository
	config   *config.Config
	logger   *logger.Logger
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(userRepo repository.UserRepository, config *config.Config, logger *logger.Logger) AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
		config:   config,
		logger:   logger.Named("auth_service"),
	}
}

// Register регистрирует нового пользователя
func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, string, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Error checking existing user",
			zap.Error(err),
			zap.String("email", req.Email))
		return nil, "", ErrInternalServer
	}

	if existingUser != nil {
		return nil, "", ErrUserExists
	}

	// Создаем нового пользователя
	user := &models.User{
		Email:    req.Email,
		FullName: req.FullName,
	}

	// Устанавливаем пароль
	if err := user.SetPassword(req.Password); err != nil {
		s.logger.Error("Error hashing password", zap.Error(err))
		return nil, "", ErrInternalServer
	}

	// Сохраняем пользователя
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Error creating user", zap.Error(err))
		return nil, "", ErrInternalServer
	}

	// Создаем сессию
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(s.config.Auth.TokenExpiration)

	// Генерируем JWT token
	token, err := s.generateJWT(user.ID, user.Email, sessionID, expiresAt)
	if err != nil {
		s.logger.Error("Error generating JWT", zap.Error(err))
		return nil, "", ErrInternalServer
	}

	// Сохраняем сессию
	session := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     token,
		UserAgent: "Registration", // Можно добавить User-Agent из контекста, если нужно
		IPAddress: "0.0.0.0",      // Можно добавить IP из контекста, если нужно
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		s.logger.Error("Error creating session", zap.Error(err))
		return nil, "", ErrInternalServer
	}

	// Возвращаем DTO пользователя и токен
	return &dto.AuthResponse{
		User: dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, token, nil
}

// Login выполняет аутентификацию пользователя
func (s *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest, userAgent, ipAddress string) (*dto.AuthResponse, string, error) {
	// Ищем пользователя по email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Error finding user by email", zap.Error(err), zap.String("email", req.Email))
		return nil, "", ErrInternalServer
	}

	if user == nil || !user.CheckPassword(req.Password) {
		return nil, "", ErrInvalidLogin
	}

	// Создаем сессию
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(s.config.Auth.TokenExpiration)

	// Генерируем JWT token
	token, err := s.generateJWT(user.ID, user.Email, sessionID, expiresAt)
	if err != nil {
		s.logger.Error("Error generating JWT", zap.Error(err))
		return nil, "", ErrInternalServer
	}

	// Сохраняем сессию
	session := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     token, // В реальности здесь должен быть хеш токена
		UserAgent: userAgent,
		IPAddress: ipAddress,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		s.logger.Error("Error creating session", zap.Error(err))
		return nil, "", ErrInternalServer
	}

	return &dto.AuthResponse{
		User: dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, token, nil
}

// Logout выполняет выход пользователя
func (s *AuthServiceImpl) Logout(ctx context.Context, userID uint64, sessionID string, all bool) error {
	var err error
	if all {
		err = s.userRepo.DeactivateAllSessions(ctx, userID)
	} else {
		err = s.userRepo.DeactivateSession(ctx, userID, sessionID)
	}

	if err != nil {
		s.logger.Error("Error during logout", zap.Error(err), zap.Uint64("userID", userID))
		return ErrInternalServer
	}

	return nil
}

// GetUserProfile получает профиль пользователя
func (s *AuthServiceImpl) GetUserProfile(ctx context.Context, userID uint64) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("Error finding user", zap.Error(err), zap.Uint64("userID", userID))
		return nil, ErrInternalServer
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return &dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateProfile обновляет профиль пользователя
func (s *AuthServiceImpl) UpdateProfile(ctx context.Context, userID uint64, req dto.UpdateProfileRequest) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("Error finding user", zap.Error(err), zap.Uint64("userID", userID))
		return nil, ErrInternalServer
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	user.FullName = req.FullName

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("Error updating user", zap.Error(err), zap.Uint64("userID", userID))
		return nil, ErrInternalServer
	}

	return &dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// ChangePassword меняет пароль пользователя
func (s *AuthServiceImpl) ChangePassword(ctx context.Context, userID uint64, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("Error finding user", zap.Error(err), zap.Uint64("userID", userID))
		return ErrInternalServer
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Проверяем текущий пароль
	if !user.CheckPassword(req.CurrentPassword) {
		return ErrInvalidPassword
	}

	// Устанавливаем новый пароль
	if err := user.SetPassword(req.NewPassword); err != nil {
		s.logger.Error("Error setting new password", zap.Error(err), zap.Uint64("userID", userID))
		return ErrInternalServer
	}

	// Обновляем пользователя в БД
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("Error updating user password", zap.Error(err), zap.Uint64("userID", userID))
		return ErrInternalServer
	}

	return nil
}

// CleanupExpiredSessions очищает просроченные сессии
func (s *AuthServiceImpl) CleanupExpiredSessions(ctx context.Context) error {
	if err := s.userRepo.CleanupExpiredSessions(ctx); err != nil {
		s.logger.Error("Error cleaning up expired sessions", zap.Error(err))
		return ErrInternalServer
	}
	return nil
}

// generateJWT генерирует JWT токен
func (s *AuthServiceImpl) generateJWT(userID uint64, email, sessionID string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"email":      email,
		"session_id": sessionID,
		"exp":        expiresAt.Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return signedToken, nil
}
