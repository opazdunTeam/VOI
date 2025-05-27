package repository

import (
	"context"
	"errors"

	"backend/orchestrator-service/internal/models"
	"backend/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PostRepository интерфейс для работы с постами
type PostRepository interface {
	CreateNote(ctx context.Context, note *models.Note) error
	GetNoteByID(ctx context.Context, noteID uint64) (*models.Note, error)
	GetNotes(ctx context.Context, userID uint64, page, size int) ([]*models.Note, int64, error)
	CreatePost(ctx context.Context, post *models.Post) error
	UpdatePost(ctx context.Context, post *models.Post) error
	GetPostByID(ctx context.Context, postID uint64) (*models.Post, error)
	GetPostsByUserID(ctx context.Context, userID uint64, page, size int) ([]*models.Post, int64, error)
	DeletePost(ctx context.Context, postID, userID uint64) error
}

// postRepository реализация репозитория постов
type postRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewPostRepository создает новый репозиторий постов
func NewPostRepository(db *gorm.DB, logger *logger.Logger) PostRepository {
	return &postRepository{
		db:     db,
		logger: logger.Named("post_repository"),
	}
}

// CreateNote создает новую заметку
func (r *postRepository) CreateNote(ctx context.Context, note *models.Note) error {
	result := r.db.WithContext(ctx).Create(note)
	if result.Error != nil {
		r.logger.Error("Failed to create note", zap.Error(result.Error), zap.Uint64("user_id", note.UserID))
		return result.Error
	}
	return nil
}

// GetNoteByID получает заметку по ID
func (r *postRepository) GetNoteByID(ctx context.Context, noteID uint64) (*models.Note, error) {
	var note models.Note
	result := r.db.WithContext(ctx).First(&note, noteID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Заметка не найдена, но это не ошибка
		}
		r.logger.Error("Failed to get note", zap.Error(result.Error), zap.Uint64("note_id", noteID))
		return nil, result.Error
	}
	return &note, nil
}

// CreatePost создает новый пост
func (r *postRepository) CreatePost(ctx context.Context, post *models.Post) error {
	result := r.db.WithContext(ctx).Create(post)
	if result.Error != nil {
		r.logger.Error("Failed to create post", zap.Error(result.Error), zap.Uint64("user_id", post.UserID))
		return result.Error
	}
	return nil
}

// UpdatePost обновляет существующий пост
func (r *postRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	result := r.db.WithContext(ctx).Save(post)
	if result.Error != nil {
		r.logger.Error("Failed to update post", zap.Error(result.Error), zap.Uint64("post_id", post.ID))
		return result.Error
	}
	return nil
}

// GetPostByID получает пост по ID
func (r *postRepository) GetPostByID(ctx context.Context, postID uint64) (*models.Post, error) {
	var post models.Post
	result := r.db.WithContext(ctx).First(&post, postID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Пост не найден, но это не ошибка
		}
		r.logger.Error("Failed to get post", zap.Error(result.Error), zap.Uint64("post_id", postID))
		return nil, result.Error
	}
	return &post, nil
}

// GetPostsByUserID получает список постов пользователя
func (r *postRepository) GetPostsByUserID(ctx context.Context, userID uint64, page, size int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	offset := (page - 1) * size

	// Получаем общее количество
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count posts", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, 0, err
	}

	// Получаем посты с пагинацией
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&posts)

	if result.Error != nil {
		r.logger.Error("Failed to get posts", zap.Error(result.Error), zap.Uint64("user_id", userID))
		return nil, 0, result.Error
	}

	return posts, total, nil
}

// DeletePost удаляет пост
func (r *postRepository) DeletePost(ctx context.Context, postID, userID uint64) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", postID, userID).Delete(&models.Post{})
	if result.Error != nil {
		r.logger.Error("Failed to delete post", zap.Error(result.Error), zap.Uint64("post_id", postID))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("post not found or not owned by user")
	}

	return nil
}

// GetNotes получает список заметок пользователя
func (r *postRepository) GetNotes(ctx context.Context, userID uint64, page, size int) ([]*models.Note, int64, error) {
	var notes []*models.Note
	var total int64

	offset := (page - 1) * size

	// Получаем общее количество
	if err := r.db.WithContext(ctx).Model(&models.Note{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count notes", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, 0, err
	}

	// Получаем заметки с пагинацией
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&notes)

	if result.Error != nil {
		r.logger.Error("Failed to get notes", zap.Error(result.Error), zap.Uint64("user_id", userID))
		return nil, 0, result.Error
	}

	return notes, total, nil
}
