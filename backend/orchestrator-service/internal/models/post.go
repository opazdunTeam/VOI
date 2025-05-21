package models

import (
	"time"

	"gorm.io/gorm"
)

// Note представляет исходную заметку пользователя
type Note struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64         `gorm:"index;not null" json:"user_id"`
	OriginalText string         `gorm:"type:text;not null" json:"original_text"`
	Source       string         `gorm:"type:varchar(20);not null;default:'text'" json:"source"` // voice|text
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Post представляет сгенерированный пост
type Post struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64         `gorm:"index;not null" json:"user_id"`
	NoteID    uint64         `gorm:"index;not null" json:"note_id"`
	ContentMD string         `gorm:"type:text;not null" json:"content_md"`
	Status    string         `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // draft|ready
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate автоматически вызывается перед созданием записи
func (n *Note) BeforeCreate(tx *gorm.DB) (err error) {
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate автоматически вызывается перед обновлением записи
func (n *Note) BeforeUpdate(tx *gorm.DB) (err error) {
	n.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate автоматически вызывается перед созданием записи
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate автоматически вызывается перед обновлением записи
func (p *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return nil
}
