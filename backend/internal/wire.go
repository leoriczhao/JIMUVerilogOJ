//go:build wireinject
// +build wireinject

package internal

import (
	"verilog-oj/backend/internal/handlers"
	"verilog-oj/backend/internal/repository"
	"verilog-oj/backend/internal/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// InitializeApp 初始化整个应用
func InitializeApp(db *gorm.DB) (*App, error) {
	wire.Build(
		repository.RepositorySet,
		services.ServiceSet,
		handlers.HandlerSet,
		NewApp,
	)
	return &App{}, nil
}

// App 包含所有应用组件的结构体
type App struct {
	Handlers *handlers.Handlers
	Services *services.Services
	Repos    *repository.Repositories
}

// NewApp 创建App实例
func NewApp(
	handlers *handlers.Handlers,
	services *services.Services,
	repos *repository.Repositories,
) *App {
	return &App{
		Handlers: handlers,
		Services: services,
		Repos:    repos,
	}
}
