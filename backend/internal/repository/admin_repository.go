package repository

import (
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/services"

	"gorm.io/gorm"
)

// AdminRepositoryImpl 管理端仓储实现
type AdminRepositoryImpl struct {
	db *gorm.DB
}

// NewAdminRepository 创建Admin仓储实例
func NewAdminRepository(db *gorm.DB) services.AdminRepository {
	return &AdminRepositoryImpl{db: db}
}

func (r *AdminRepositoryImpl) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *AdminRepositoryImpl) CountProblems() (int64, error) {
	var count int64
	err := r.db.Model(&models.Problem{}).Count(&count).Error
	return count, err
}

func (r *AdminRepositoryImpl) CountSubmissions() (int64, error) {
	var count int64
	err := r.db.Model(&models.Submission{}).Count(&count).Error
	return count, err
}

