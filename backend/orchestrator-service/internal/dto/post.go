package dto

// NoteRequest представляет запрос на создание заметки
type NoteRequest struct {
	Text   string `json:"text" binding:"required"`
	Source string `json:"source" binding:"required,oneof=text voice"` // voice|text
}

// GeneratePostRequest представляет запрос на генерацию поста
type GeneratePostRequest struct {
	NoteID uint64 `json:"note_id" binding:"required"`
}

// PostResponse представляет ответ с данными поста
type PostResponse struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	NoteID    uint64 `json:"note_id"`
	ContentMD string `json:"content_md"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PostListResponse представляет ответ со списком постов
type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// NoteResponse представляет ответ с данными заметки
type NoteResponse struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	OriginalText string `json:"original_text"`
	Source       string `json:"source"`
	CreatedAt    string `json:"created_at"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
