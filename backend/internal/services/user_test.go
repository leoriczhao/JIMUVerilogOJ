package services

import (
	"errors"
	"testing"
	"verilog-oj/backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateStats(userID uint, solved, submitted int) error {
	args := m.Called(userID, solved, submitted)
	return args.Error(0)
}

// TestUserService_CreateUser 测试创建用户
func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		email          string
		password       string
		mockSetup      func(*MockUserRepository)
		expectedError  string
		expectedResult bool
	}{
		{
			name:     "成功创建用户",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// 检查用户名不存在
				mockRepo.On("GetByUsername", "testuser").Return(nil, nil)
				// 检查邮箱不存在
				mockRepo.On("GetByEmail", "test@example.com").Return(nil, nil)
				// 创建用户成功
				mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedError:  "",
			expectedResult: true,
		},
		{
			name:     "用户名已存在",
			username: "existinguser",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// 用户名已存在
				existingUser := &domain.User{
					ID:       1,
					Username: "existinguser",
					Email:    "existing@example.com",
				}
				mockRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)
			},
			expectedError:  "用户名已存在",
			expectedResult: false,
		},
		{
			name:     "邮箱已存在",
			username: "newuser",
			email:    "existing@example.com",
			password: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// 用户名不存在
				mockRepo.On("GetByUsername", "newuser").Return(nil, nil)
				// 邮箱已存在
				existingUser := &domain.User{
					ID:       1,
					Username: "existinguser",
					Email:    "existing@example.com",
				}
				mockRepo.On("GetByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedError:  "邮箱已存在",
			expectedResult: false,
		},
		{
			name:     "数据库创建失败",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// 检查用户名不存在
				mockRepo.On("GetByUsername", "testuser").Return(nil, nil)
				// 检查邮箱不存在
				mockRepo.On("GetByEmail", "test@example.com").Return(nil, nil)
				// 创建用户失败
				mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("数据库错误"))
			},
			expectedError:  "数据库错误",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock仓储
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// 创建服务
			userService := NewUserService(mockRepo)

			// 执行测试
			user, err := userService.CreateUser(tt.username, tt.email, tt.password, "student")

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, "student", user.Role)
				assert.Equal(t, 1200, user.Rating)
				assert.True(t, user.IsActive)
				// 验证密码已加密
				assert.NotEqual(t, tt.password, user.Password)
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))
				assert.NoError(t, err)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_GetUserByUsername 测试根据用户名获取用户
func TestUserService_GetUserByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetup     func(*MockUserRepository)
		expectedUser  *domain.User
		expectedError string
	}{
		{
			name:     "成功获取用户",
			username: "testuser",
			mockSetup: func(mockRepo *MockUserRepository) {
				user := &domain.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Role:     "student",
				}
				mockRepo.On("GetByUsername", "testuser").Return(user, nil)
			},
			expectedUser: &domain.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Role:     "student",
			},
			expectedError: "",
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "nonexistent").Return(nil, nil)
			},
			expectedUser:  nil,
			expectedError: "",
		},
		{
			name:     "数据库错误",
			username: "testuser",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByUsername", "testuser").Return(nil, errors.New("数据库连接失败"))
			},
			expectedUser:  nil,
			expectedError: "数据库连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock仓储
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// 创建服务
			userService := NewUserService(mockRepo)

			// 执行测试
			user, err := userService.GetUserByUsername(tt.username)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				if tt.expectedUser == nil {
					assert.Nil(t, user)
				} else {
					assert.NotNil(t, user)
					assert.Equal(t, tt.expectedUser.ID, user.ID)
					assert.Equal(t, tt.expectedUser.Username, user.Username)
					assert.Equal(t, tt.expectedUser.Email, user.Email)
					assert.Equal(t, tt.expectedUser.Role, user.Role)
				}
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_GetUserByID 测试根据ID获取用户
func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		mockSetup     func(*MockUserRepository)
		expectedUser  *domain.User
		expectedError string
	}{
		{
			name:   "成功获取用户",
			userID: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				user := &domain.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Role:     "student",
				}
				mockRepo.On("GetByID", uint(1)).Return(user, nil)
			},
			expectedUser: &domain.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Role:     "student",
			},
			expectedError: "",
		},
		{
			name:   "用户不存在",
			userID: 999,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByID", uint(999)).Return(nil, nil)
			},
			expectedUser:  nil,
			expectedError: "用户不存在",
		},
		{
			name:   "数据库错误",
			userID: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("数据库连接失败"))
			},
			expectedUser:  nil,
			expectedError: "数据库连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock仓储
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// 创建服务
			userService := NewUserService(mockRepo)

			// 执行测试
			user, err := userService.GetUserByID(tt.userID)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Role, user.Role)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_ValidatePassword 测试密码验证
func TestUserService_ValidatePassword(t *testing.T) {
	// 创建一个测试用的加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		user          *domain.User
		password      string
		expectedError bool
	}{
		{
			name: "密码正确",
			user: &domain.User{
				ID:       1,
				Username: "testuser",
				Password: string(hashedPassword),
			},
			password:      "correctpassword",
			expectedError: false,
		},
		{
			name: "密码错误",
			user: &domain.User{
				ID:       1,
				Username: "testuser",
				Password: string(hashedPassword),
			},
			password:      "wrongpassword",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock仓储（这个测试不需要仓储操作）
			mockRepo := new(MockUserRepository)

			// 创建服务
			userService := NewUserService(mockRepo)

			// 执行测试
			err := userService.ValidatePassword(tt.user, tt.password)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUserService_UpdateUser 测试更新用户
func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *domain.User
		mockSetup     func(*MockUserRepository)
		expectedError string
	}{
		{
			name: "成功更新用户",
			user: &domain.User{
				ID:       1,
				Username: "testuser",
				Email:    "updated@example.com",
				Nickname: "Updated Nickname",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "更新失败",
			user: &domain.User{
				ID:       1,
				Username: "testuser",
				Email:    "updated@example.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(errors.New("更新失败"))
			},
			expectedError: "更新失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock仓储
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// 创建服务
			userService := NewUserService(mockRepo)

			// 执行测试
			err := userService.UpdateUser(tt.user)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_UpdateUserStats 测试更新用户统计信息
func TestUserService_UpdateUserStats(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		solved        int
		submitted     int
		mockSetup     func(*MockUserRepository)
		expectedError string
	}{
		{
			name:      "成功更新统计信息",
			userID:    1,
			solved:    10,
			submitted: 20,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("UpdateStats", uint(1), 10, 20).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "更新失败",
			userID:    1,
			solved:    10,
			submitted: 20,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("UpdateStats", uint(1), 10, 20).Return(errors.New("统计更新失败"))
			},
			expectedError: "统计更新失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Mock仓储
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// 创建服务
			userService := NewUserService(mockRepo)

			// 执行测试
			err := userService.UpdateUserStats(tt.userID, tt.solved, tt.submitted)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}
