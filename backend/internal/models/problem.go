package models

import (
	"time"

	"gorm.io/gorm"
)

// Problem 题目模型
type Problem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 基本信息
	Title       string `json:"title" gorm:"not null" validate:"required"`
	Description string `json:"description" gorm:"type:text" validate:"required"`
	InputDesc   string `json:"input_desc" gorm:"type:text"`
	OutputDesc  string `json:"output_desc" gorm:"type:text"`

	// 难度和分类
	Difficulty string `json:"difficulty" gorm:"default:Easy"` // Easy, Medium, Hard
	Category   string `json:"category" gorm:"size:50"`
	Tags       string `json:"tags" gorm:"type:text"` // JSON数组字符串

	// 限制条件
	TimeLimit   int `json:"time_limit" gorm:"default:1000"`  // 毫秒
	MemoryLimit int `json:"memory_limit" gorm:"default:128"` // MB

	// 统计信息
	SubmitCount   int `json:"submit_count" gorm:"default:0"`
	AcceptedCount int `json:"accepted_count" gorm:"default:0"`

	// 状态
	IsPublic bool `json:"is_public" gorm:"default:false"`
	AuthorID uint `json:"author_id"`
	Author   User `json:"author" gorm:"foreignKey:AuthorID"`
}

// TestCase 测试用例模型
type TestCase struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	ProblemID uint    `json:"problem_id" gorm:"not null"`
	Problem   Problem `json:"problem" gorm:"foreignKey:ProblemID"`

	Input    string `json:"input" gorm:"type:text"`
	Output   string `json:"output" gorm:"type:text"`
	IsSample bool   `json:"is_sample" gorm:"default:false"`
}
