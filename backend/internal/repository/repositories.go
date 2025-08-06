package repository

import (
	"gorm.io/gorm"
	"verilog-oj/backend/internal/services"
)

// Repositories 包含所有Repository的结构体
type Repositories struct {
	UserRepository       services.UserRepository
	ProblemRepository    services.ProblemRepository
	SubmissionRepository services.SubmissionRepository
	ForumRepository      services.ForumRepository
	NewsRepository       services.NewsRepository
}

// NewRepositories 创建Repositories实例
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository:       NewUserRepository(db),
		ProblemRepository:    NewProblemRepository(db),
		SubmissionRepository: NewSubmissionRepository(db),
		ForumRepository:      NewForumRepository(db),
		NewsRepository:       NewNewsRepository(db),
	}
}