//go:build wireinject
// +build wireinject

package services

import (
	"github.com/google/wire"
)

// ServiceSet 提供所有Service的Wire集合
var ServiceSet = wire.NewSet(
	// User Service
	NewUserService,

	// Problem Service
	NewProblemService,

	// Submission Service
	NewSubmissionService,

	// Forum Service
	NewForumService,

	// News Service
	NewNewsService,

	// Services构造函数
	NewServices,
)