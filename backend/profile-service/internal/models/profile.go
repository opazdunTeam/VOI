package models

import (
	"time"

	"gorm.io/gorm"
)

// VoiceProfile представляет профиль "голоса" пользователя (Voice DNA)
type VoiceProfile struct {
	UserID    uint64         `gorm:"primaryKey;uniqueIndex" json:"user_id"`
	DNAData   string         `gorm:"type:jsonb;not null" json:"dna_data"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate автоматически вызывается перед созданием записи
func (p *VoiceProfile) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate автоматически вызывается перед обновлением записи
func (p *VoiceProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return nil
}
