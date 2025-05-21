package repository

import (
	"context"
	"errors"

	"backend/pkg/logger"
	"backend/profile-service/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProfileRepository интерфейс для работы с профилями
type ProfileRepository interface {
	GetByUserID(ctx context.Context, userID uint64) (*models.VoiceProfile, error)
	Create(ctx context.Context, profile *models.VoiceProfile) error
	Update(ctx context.Context, profile *models.VoiceProfile) error
}

// profileRepository реализация репозитория профилей
type profileRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewProfileRepository создает новый репозиторий профилей
func NewProfileRepository(db *gorm.DB, logger *logger.Logger) ProfileRepository {
	return &profileRepository{
		db:     db,
		logger: logger.Named("profile_repository"),
	}
}

// GetByUserID получает профиль по ID пользователя
func (r *profileRepository) GetByUserID(ctx context.Context, userID uint64) (*models.VoiceProfile, error) {
	var profile models.VoiceProfile
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Профиль не найден, но это не ошибка
		}
		r.logger.Error("Failed to get profile", zap.Error(result.Error), zap.Uint64("user_id", userID))
		return nil, result.Error
	}
	return &profile, nil
}

// Create создает новый профиль
func (r *profileRepository) Create(ctx context.Context, profile *models.VoiceProfile) error {
	result := r.db.WithContext(ctx).Create(profile)
	if result.Error != nil {
		r.logger.Error("Failed to create profile", zap.Error(result.Error), zap.Uint64("user_id", profile.UserID))
		return result.Error
	}
	return nil
}

// Update обновляет существующий профиль
func (r *profileRepository) Update(ctx context.Context, profile *models.VoiceProfile) error {
	result := r.db.WithContext(ctx).Save(profile)
	if result.Error != nil {
		r.logger.Error("Failed to update profile", zap.Error(result.Error), zap.Uint64("user_id", profile.UserID))
		return result.Error
	}
	return nil
}
