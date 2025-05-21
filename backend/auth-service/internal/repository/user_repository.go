package repository

import (
	"context"
	"errors"
	"time"

	"backend/auth-service/internal/models"
	"backend/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserRepository интерфейс для работы с пользователями в БД
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint64) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	CreateSession(ctx context.Context, session *models.Session) error
	FindSessionByID(ctx context.Context, id string) (*models.Session, error)
	DeactivateSession(ctx context.Context, userID uint64, sessionID string) error
	DeactivateAllSessions(ctx context.Context, userID uint64) error
	CleanupExpiredSessions(ctx context.Context) error
}

// GormUserRepository реализует UserRepository с использованием GORM
type GormUserRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *gorm.DB, log *logger.Logger) UserRepository {
	return &GormUserRepository{
		db:     db,
		logger: log.Named("user_repository"),
	}
}

// Create создает нового пользователя в БД
func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID ищет пользователя по ID
func (r *GormUserRepository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail ищет пользователя по email
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

// Update обновляет данные пользователя
func (r *GormUserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// CreateSession создает новую сессию пользователя
func (r *GormUserRepository) CreateSession(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// FindSessionByID ищет сессию по ID
func (r *GormUserRepository) FindSessionByID(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // сессия не найдена
		}
		return nil, err
	}
	return &session, nil
}

// DeactivateSession деактивирует конкретную сессию пользователя
func (r *GormUserRepository) DeactivateSession(ctx context.Context, userID uint64, sessionID string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND id = ?", userID, sessionID).
		Update("is_active", false)
	return result.Error
}

// DeactivateAllSessions деактивирует все сессии пользователя
func (r *GormUserRepository) DeactivateAllSessions(ctx context.Context, userID uint64) error {
	result := r.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ?", userID).
		Update("is_active", false)
	return result.Error
}

// CleanupExpiredSessions очищает устаревшие сессии
func (r *GormUserRepository) CleanupExpiredSessions(ctx context.Context) error {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}

	r.logger.Info("Cleaned up expired sessions", zap.Int64("count", result.RowsAffected))
	return nil
}
