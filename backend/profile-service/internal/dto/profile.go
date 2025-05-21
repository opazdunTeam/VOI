package dto

// ProfileRequest представляет запрос на создание/обновление профиля
type ProfileRequest struct {
	DNAData string `json:"dna_data" binding:"required"`
}

// ProfileResponse представляет ответ с данными профиля
type ProfileResponse struct {
	UserID    uint64 `json:"user_id"`
	DNAData   string `json:"dna_data"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
