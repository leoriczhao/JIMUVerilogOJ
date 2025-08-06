package repository

import (
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/services"

	"gorm.io/gorm"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) services.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *domain.User) error {
	// 将domain.User转换为models.User
	modelUser := UserDomainToModel(user)
	
	err := r.db.Create(modelUser).Error
	if err != nil {
		return err
	}
	
	// 更新ID和时间戳
	user.ID = modelUser.ID
	user.CreatedAt = modelUser.CreatedAt
	user.UpdatedAt = modelUser.UpdatedAt
	
	return nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var modelUser models.User
	err := r.db.Where("username = ?", username).First(&modelUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	// 转换为domain.User
	return UserModelToDomain(&modelUser), nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	var modelUser models.User
	err := r.db.First(&modelUser, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	// 转换为domain.User
	return UserModelToDomain(&modelUser), nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var modelUser models.User
	err := r.db.Where("email = ?", email).First(&modelUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	// 转换为domain.User
	return UserModelToDomain(&modelUser), nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *domain.User) error {
	// 将domain.User转换为models.User
	modelUser := UserDomainToModel(user)
	
	err := r.db.Save(modelUser).Error
	if err != nil {
		return err
	}
	
	// 更新时间戳
	user.UpdatedAt = modelUser.UpdatedAt
	
	return nil
}

// UpdateStats 更新用户统计信息
func (r *UserRepository) UpdateStats(userID uint, solved, submitted int) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"solved":    gorm.Expr("solved + ?", solved),
		"submitted": gorm.Expr("submitted + ?", submitted),
	}).Error
}