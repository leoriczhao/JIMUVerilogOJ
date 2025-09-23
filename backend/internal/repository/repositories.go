package repository

import (
	"verilog-oj/backend/internal/services"

	"gorm.io/gorm"
)

// Repositories 包含所有Repository的结构体
type Repositories struct {
	UserRepository       services.UserRepository
	ProblemRepository    services.ProblemRepository
	SubmissionRepository services.SubmissionRepository
	ForumRepository      services.ForumRepository
	NewsRepository       services.NewsRepository
	AdminRepository      services.AdminRepository
}

// NewRepositories 创建Repositories实例
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository:       NewUserRepository(db),
		ProblemRepository:    NewProblemRepository(db),
		SubmissionRepository: NewSubmissionRepository(db),
		ForumRepository:      NewForumRepository(db),
		NewsRepository:       NewNewsRepository(db),
		AdminRepository:      NewAdminRepository(db),
	}
}
