package models

import (
	"time"

	"gorm.io/gorm"
)

// News represents a news article in the system.
type News struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Use json:"-" to hide from API output

	Title   string `json:"title" gorm:"not null;size:255" validate:"required,min=3,max=100"`
	Content string `json:"content" gorm:"type:text" validate:"required"`
	Summary string `json:"summary" gorm:"type:text"`

	AuthorID uint `json:"author_id" gorm:"not null"`
	Author   User `json:"author" gorm:"foreignKey:AuthorID"`

	IsPublished bool `json:"is_published" gorm:"default:false"`
	IsFeatured  bool `json:"is_featured" gorm:"default:false"`

	Category  string `json:"category" gorm:"size:50"`
	Tags      string `json:"tags" gorm:"type:text"`
	ViewCount int    `json:"view_count" gorm:"default:0"`
}
