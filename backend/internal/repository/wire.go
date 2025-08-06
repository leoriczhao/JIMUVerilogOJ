//go:build wireinject
// +build wireinject

package repository

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"verilog-oj/backend/internal/services"
)

// RepositorySet 提供所有Repository的Wire集合
var RepositorySet = wire.NewSet(
	// User Repository
	NewUserRepository,

	// Problem Repository
	NewProblemRepository,

	// Submission Repository
	NewSubmissionRepository,

	// Forum Repository (统一的论坛仓储)
	NewForumRepository,

	// News Repository
	NewNewsRepository,

	// Repositories构造函数
	NewRepositories,
)