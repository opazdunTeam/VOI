package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User представляет информацию о пользователе
type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null;type:varchar(255)" json:"email"`
	PasswordHash string         `gorm:"not null;type:varchar(255)" json:"-"`
	FullName     string         `gorm:"type:varchar(255)" json:"full_name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Связи с другими моделями
	Sessions []Session `gorm:"foreignKey:UserID" json:"-"`
}

// SetPassword хеширует и устанавливает пароль пользователя
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword проверяет правильность пароля
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Session представляет информацию о сессии пользователя
type Session struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    uint64         `gorm:"index;not null" json:"user_id"`
	Token     string         `gorm:"type:varchar(255);not null" json:"-"` // Хеш токена
	UserAgent string         `gorm:"type:varchar(255)" json:"user_agent"`
	IPAddress string         `gorm:"type:varchar(45)" json:"ip_address"` // IPv6 может быть до 45 символов
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Связи с другими моделями
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate автоматически вызывается перед созданием записи
func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	// Обновлять UpdatedAt каждый раз при сохранении
	s.UpdatedAt = time.Now()
	return nil
}

// VoiceProfile представляет профиль "голоса" пользователя (Voice DNA)
type VoiceProfile struct {
	UserID    uint64         `gorm:"primaryKey" json:"user_id"`
	DNAData   string         `gorm:"type:jsonb;not null" json:"dna_data"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Связи с другими моделями
	User User `gorm:"foreignKey:UserID" json:"-"`
}
