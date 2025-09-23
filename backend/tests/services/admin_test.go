package services

import (
	"errors"
	"testing"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepository is a mock for the AdminRepository
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) CountUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAdminRepository) CountProblems() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAdminRepository) CountSubmissions() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// TestAdminService_GetSystemStats tests the GetSystemStats method
func TestAdminService_GetSystemStats(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockAdminRepository)
		expectedStats *services.SystemStats
		expectedError string
	}{
		{
			name: "成功获取系统统计",
			mockSetup: func(mockRepo *MockAdminRepository) {
				mockRepo.On("CountUsers").Return(int64(100), nil)
				mockRepo.On("CountProblems").Return(int64(50), nil)
				mockRepo.On("CountSubmissions").Return(int64(200), nil)
			},
			expectedStats: &services.SystemStats{
				Users:        100,
				Problems:     50,
				Submissions:  200,
				JudgesOnline: 1, // As per current implementation
			},
			expectedError: "",
		},
		{
			name: "获取用户数失败",
			mockSetup: func(mockRepo *MockAdminRepository) {
				mockRepo.On("CountUsers").Return(int64(0), errors.New("数据库错误"))
			},
			expectedStats: nil,
			expectedError: "数据库错误",
		},
		{
			name: "获取问题数失败",
			mockSetup: func(mockRepo *MockAdminRepository) {
				mockRepo.On("CountUsers").Return(int64(100), nil)
				mockRepo.On("CountProblems").Return(int64(0), errors.New("数据库错误"))
			},
			expectedStats: nil,
			expectedError: "数据库错误",
		},
		{
			name: "获取提交数失败",
			mockSetup: func(mockRepo *MockAdminRepository) {
				mockRepo.On("CountUsers").Return(int64(100), nil)
				mockRepo.On("CountProblems").Return(int64(50), nil)
				mockRepo.On("CountSubmissions").Return(int64(0), errors.New("数据库错误"))
			},
			expectedStats: nil,
			expectedError: "数据库错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminRepository)
			tt.mockSetup(mockRepo)

			adminService := services.NewAdminService(mockRepo)
			stats, err := adminService.GetSystemStats()

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStats, stats)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
