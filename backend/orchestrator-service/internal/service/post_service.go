package service

import (
	"context"
	"errors"
	"time"

	"backend/orchestrator-service/internal/dto"
	"backend/orchestrator-service/internal/models"
	"backend/orchestrator-service/internal/repository"
	"backend/pkg/config"
	"backend/pkg/logger"

	"go.uber.org/zap"
)

// PostService интерфейс для работы с постами
type PostService interface {
	CreateNote(ctx context.Context, userID uint64, req *dto.NoteRequest) (*dto.NoteResponse, error)
	GeneratePost(ctx context.Context, userID uint64, req *dto.GeneratePostRequest) (*dto.PostResponse, error)
	GetPost(ctx context.Context, postID uint64) (*dto.PostResponse, error)
	GetPosts(ctx context.Context, userID uint64, page, size int) (*dto.PostListResponse, error)
	DeletePost(ctx context.Context, postID, userID uint64) error
}

// postService реализация сервиса постов
type postService struct {
	repo   repository.PostRepository
	config *config.Config
	logger *logger.Logger
}

// NewPostService создает новый сервис постов
func NewPostService(repo repository.PostRepository, cfg *config.Config, logger *logger.Logger) PostService {
	return &postService{
		repo:   repo,
		config: cfg,
		logger: logger.Named("post_service"),
	}
}

// CreateNote создает новую заметку
func (s *postService) CreateNote(ctx context.Context, userID uint64, req *dto.NoteRequest) (*dto.NoteResponse, error) {
	note := &models.Note{
		UserID:       userID,
		OriginalText: req.Text,
		Source:       req.Source,
	}

	if err := s.repo.CreateNote(ctx, note); err != nil {
		s.logger.Error("Failed to create note", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	return &dto.NoteResponse{
		ID:           note.ID,
		UserID:       note.UserID,
		OriginalText: note.OriginalText,
		Source:       note.Source,
		CreatedAt:    note.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GeneratePost генерирует новый пост на основе заметки
func (s *postService) GeneratePost(ctx context.Context, userID uint64, req *dto.GeneratePostRequest) (*dto.PostResponse, error) {
	// Получаем заметку
	note, err := s.repo.GetNoteByID(ctx, req.NoteID)
	if err != nil {
		s.logger.Error("Failed to get note", zap.Error(err), zap.Uint64("note_id", req.NoteID))
		return nil, err
	}

	if note == nil {
		return nil, errors.New("note not found")
	}

	// Проверяем, что заметка принадлежит пользователю
	if note.UserID != userID {
		return nil, errors.New("note not found or not owned by user")
	}

	// TODO: Вызвать генератор контента через API

	// Создаем пост (пока с заглушкой)
	post := &models.Post{
		UserID:    userID,
		NoteID:    note.ID,
		ContentMD: "# Сгенерированный пост\n\nЭто заглушка для сгенерированного поста на основе заметки: " + note.OriginalText,
		Status:    "draft",
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		s.logger.Error("Failed to create post", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	return &dto.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		NoteID:    post.NoteID,
		ContentMD: post.ContentMD,
		Status:    post.Status,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetPost получает пост по ID
func (s *postService) GetPost(ctx context.Context, postID uint64) (*dto.PostResponse, error) {
	post, err := s.repo.GetPostByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", zap.Error(err), zap.Uint64("post_id", postID))
		return nil, err
	}

	if post == nil {
		return nil, errors.New("post not found")
	}

	return &dto.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		NoteID:    post.NoteID,
		ContentMD: post.ContentMD,
		Status:    post.Status,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetPosts получает список постов пользователя
func (s *postService) GetPosts(ctx context.Context, userID uint64, page, size int) (*dto.PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	posts, total, err := s.repo.GetPostsByUserID(ctx, userID, page, size)
	if err != nil {
		s.logger.Error("Failed to get posts", zap.Error(err), zap.Uint64("user_id", userID))
		return nil, err
	}

	response := &dto.PostListResponse{
		Posts: make([]dto.PostResponse, 0, len(posts)),
		Total: total,
		Page:  page,
		Size:  size,
	}

	for _, post := range posts {
		response.Posts = append(response.Posts, dto.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			NoteID:    post.NoteID,
			ContentMD: post.ContentMD,
			Status:    post.Status,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

// DeletePost удаляет пост
func (s *postService) DeletePost(ctx context.Context, postID, userID uint64) error {
	return s.repo.DeletePost(ctx, postID, userID)
}
