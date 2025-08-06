package services

import (
	"errors"
	"verilog-oj/backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 创建用户
	Create(user *domain.User) error
	// 根据用户名获取用户
	GetByUsername(username string) (*domain.User, error)
	// 根据ID获取用户
	GetByID(id uint) (*domain.User, error)
	// 根据邮箱获取用户
	GetByEmail(email string) (*domain.User, error)
	// 更新用户信息
	Update(user *domain.User) error
	// 更新用户统计信息
	UpdateStats(userID uint, solved, submitted int) error
}

// UserService 用户服务
type UserService struct {
	userRepo UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, email, password string) (*domain.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "student",
		IsActive: true,
		Rating:   1200,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(username)
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *domain.User) error {
	return s.userRepo.Update(user)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(email)
}

// UpdateUserStats 更新用户统计信息
func (s *UserService) UpdateUserStats(userID uint, solved, submitted int) error {
	return s.userRepo.UpdateStats(userID, solved, submitted)
}

// ValidatePassword 验证密码
func (s *UserService) ValidatePassword(user *domain.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}