package service

import (
	"context"
	"errors"
	"time"

	"backend/pkg/config"
	"backend/pkg/logger"
	"backend/profile-service/internal/dto"
	"backend/profile-service/internal/models"
	"backend/profile-service/internal/repository"

	"go.uber.org/zap"
)

// ProfileService интерфейс для работы с профилями
type ProfileService interface {
	GetProfile(ctx context.Context, userID uint64) (*dto.ProfileResponse, error)
	CreateProfile(ctx context.Context, userID uint64, req *dto.ProfileRequest) (*dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req *dto.ProfileRequest) (*dto.ProfileResponse, error)
}

// profileService реализация сервиса профилей
type profileService struct {
	repo   repository.ProfileRepository
	config *config.Config
	logger *logger.Logger
}

// NewProfileService создает новый сервис профилей
func NewProfileService(repo repository.ProfileRepository, cfg *config.Config, logger *logger.Logger) ProfileService {
	return &profileService{
		repo:   repo,
		config: cfg,
		logger: logger.Named("profile_service"),
	}
}

// GetProfile получает профиль пользователя
func (s *profileService) GetProfile(ctx context.Context, userID uint64) (*dto.ProfileResponse, error) {
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get profile", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	if profile == nil {
		return nil, errors.New("profile not found")
	}

	return &dto.ProfileResponse{
		UserID:    profile.UserID,
		DNAData:   profile.DNAData,
		UpdatedAt: profile.UpdatedAt.Format(time.RFC3339),
		CreatedAt: profile.CreatedAt.Format(time.RFC3339),
	}, nil
}

// CreateProfile создает новый профиль пользователя
func (s *profileService) CreateProfile(ctx context.Context, userID uint64, req *dto.ProfileRequest) (*dto.ProfileResponse, error) {
	// Создаем новый профиль
	profile := &models.VoiceProfile{
		UserID:  userID,
		DNAData: req.DNAData,
	}

	if err := s.repo.Create(ctx, profile); err != nil {
		s.logger.Error("Failed to create profile", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	return &dto.ProfileResponse{
		UserID:    profile.UserID,
		DNAData:   profile.DNAData,
		UpdatedAt: profile.UpdatedAt.Format(time.RFC3339),
		CreatedAt: profile.CreatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateProfile обновляет существующий профиль пользователя
func (s *profileService) UpdateProfile(ctx context.Context, userID uint64, req *dto.ProfileRequest) (*dto.ProfileResponse, error) {
	// Проверяем, существует ли профиль
	existingProfile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get profile for update", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	if existingProfile == nil {
		return nil, errors.New("profile not found")
	}

	// Обновляем профиль
	existingProfile.DNAData = req.DNAData

	if err := s.repo.Update(ctx, existingProfile); err != nil {
		s.logger.Error("Failed to update profile", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	return &dto.ProfileResponse{
		UserID:    existingProfile.UserID,
		DNAData:   existingProfile.DNAData,
		UpdatedAt: existingProfile.UpdatedAt.Format(time.RFC3339),
		CreatedAt: existingProfile.CreatedAt.Format(time.RFC3339),
	}, nil
}
