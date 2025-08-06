//go:build wireinject
// +build wireinject

package handlers

import (
	"github.com/google/wire"
	"verilog-oj/backend/internal/services"
)

// HandlerSet 提供所有Handler的Wire集合
var HandlerSet = wire.NewSet(
	// User Handler
	NewUserHandler,

	// Problem Handler
	NewProblemHandler,

	// Submission Handler
	NewSubmissionHandler,

	// Forum Handler
	NewForumHandler,

	// News Handler
	NewNewsHandler,

	// Handlers构造函数
	NewHandlers,
)