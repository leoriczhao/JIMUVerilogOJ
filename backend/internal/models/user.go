package models

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	Username string `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=20"`
	Email    string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password string `json:"-" gorm:"not null" validate:"required,min=6"`
	
	// 用户信息
	Nickname    string    `json:"nickname" gorm:"size:50"`
	Avatar      string    `json:"avatar"`
	School      string    `json:"school" gorm:"size:100"`
	StudentID   string    `json:"student_id" gorm:"size:20"`
	Role        string    `json:"role" gorm:"default:student"` // student, teacher, admin
	
	// 统计信息
	Solved      int `json:"solved" gorm:"default:0"`
	Submitted   int `json:"submitted" gorm:"default:0"`
	Rating      int `json:"rating" gorm:"default:1200"`
	
	// 状态
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	LastLoginAt time.Time `json:"last_login_at"`
} 