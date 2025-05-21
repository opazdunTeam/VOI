package dto

// RegisterRequest представляет запрос на регистрацию нового пользователя
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginRequest представляет запрос на аутентификацию пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse представляет ответ при успешной аутентификации
type AuthResponse struct {
	User UserDTO `json:"user"`
}

// UserDTO представляет публичную информацию о пользователе
type UserDTO struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

// LogoutRequest представляет запрос на выход из системы
type LogoutRequest struct {
	All bool `json:"all"` // Выйти из всех сессий
}

// ErrorResponse представляет структуру ошибки в API
type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// UpdateProfileRequest представляет запрос на обновление профиля пользователя
type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required"`
}

// ChangePasswordRequest представляет запрос на смену пароля
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}
