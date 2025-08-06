package services

import (
	"errors"
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubmissionRepository Mock 提交仓储
type MockSubmissionRepository struct {
	mock.Mock
}

func (m *MockSubmissionRepository) Create(submission *domain.Submission) error {
	args := m.Called(submission)
	return args.Error(0)
}

func (m *MockSubmissionRepository) GetByID(id uint) (*domain.Submission, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) List(page, limit int, userID, problemID uint, status string) ([]domain.Submission, int64, error) {
	args := m.Called(page, limit, userID, problemID, status)
	return args.Get(0).([]domain.Submission), args.Get(1).(int64), args.Error(2)
}

func (m *MockSubmissionRepository) UpdateStatus(id uint, status string, score int, runTime, memory int, errorMessage string, passedTests, totalTests int) error {
	args := m.Called(id, status, score, runTime, memory, errorMessage, passedTests, totalTests)
	return args.Error(0)
}

func (m *MockSubmissionRepository) CountAcceptedByUser(userID, problemID uint) (int64, error) {
	args := m.Called(userID, problemID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSubmissionRepository) GetStats(userID uint) (map[string]interface{}, error) {
	args := m.Called(userID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockSubmissionRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestSubmissionService_CreateSubmission 测试创建提交
func TestSubmissionService_CreateSubmission(t *testing.T) {
	tests := []struct {
		name           string
		problemID      uint
		code           string
		language       string
		userID         uint
		mockUser       *domain.User
		mockProblem    *domain.Problem
		userRepoError  error
		problemRepoError error
		createError    error
		expectedError  string
		expectedSubmission bool
	}{
		{
			name:      "成功创建提交",
			problemID: 1,
			code:      "module test(); endmodule",
			language:  "verilog",
			userID:    1,
			mockUser: &domain.User{
				ID:        1,
				Username:  "testuser",
				Solved:    5,
				Submitted: 10,
			},
			mockProblem: &domain.Problem{
				ID:       1,
				Title:    "Test Problem",
				IsPublic: true,
				AuthorID: 2,
			},
			expectedSubmission: true,
		},
		{
			name:      "用户不存在",
			problemID: 1,
			code:      "module test(); endmodule",
			language:  "verilog",
			userID:    1,
			mockUser:  nil,
			expectedError: "用户不存在",
		},
		{
			name:      "题目不存在",
			problemID: 1,
			code:      "module test(); endmodule",
			language:  "verilog",
			userID:    1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			mockProblem:   nil,
			expectedError: "题目不存在",
		},
		{
			name:      "代码为空",
			problemID: 1,
			code:      "",
			language:  "verilog",
			userID:    1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			mockProblem: &domain.Problem{
				ID:       1,
				IsPublic: true,
			},
			expectedError: "代码不能为空",
		},
		{
			name:      "代码长度超过限制",
			problemID: 1,
			code:      string(make([]byte, 100001)),
			language:  "verilog",
			userID:    1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			mockProblem: &domain.Problem{
				ID:       1,
				IsPublic: true,
			},
			expectedError: "代码长度超过限制",
		},
		{
			name:      "私有题目无权限访问",
			problemID: 1,
			code:      "module test(); endmodule",
			language:  "verilog",
			userID:    1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			mockProblem: &domain.Problem{
				ID:       1,
				IsPublic: false,
				AuthorID: 2,
			},
			expectedError: "没有权限访问此题目",
		},
		{
			name:      "私有题目作者可以访问",
			problemID: 1,
			code:      "module test(); endmodule",
			language:  "verilog",
			userID:    1,
			mockUser: &domain.User{
				ID:        1,
				Username:  "testuser",
				Solved:    5,
				Submitted: 10,
			},
			mockProblem: &domain.Problem{
				ID:       1,
				IsPublic: false,
				AuthorID: 1,
			},
			expectedSubmission: true,
		},
		{
			name:      "默认语言设置",
			problemID: 1,
			code:      "module test(); endmodule",
			language:  "",
			userID:    1,
			mockUser: &domain.User{
				ID:        1,
				Username:  "testuser",
				Solved:    5,
				Submitted: 10,
			},
			mockProblem: &domain.Problem{
				ID:       1,
				IsPublic: true,
			},
			expectedSubmission: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockSubmissionRepo := new(MockSubmissionRepository)
			mockProblemRepo := new(MockProblemRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockUserRepo.On("GetByID", tt.userID).Return(tt.mockUser, tt.userRepoError)

			if tt.mockUser != nil {
				mockProblemRepo.On("GetByID", tt.problemID).Return(tt.mockProblem, tt.problemRepoError)
			}

			if tt.mockProblem != nil && tt.code != "" && len(tt.code) <= 100000 && (tt.mockProblem.IsPublic || tt.mockProblem.AuthorID == tt.userID) {
				mockSubmissionRepo.On("Create", mock.AnythingOfType("*domain.Submission")).Return(tt.createError)
				if tt.createError == nil {
					mockProblemRepo.On("UpdateSubmitCount", tt.problemID, 1).Return(nil)
					mockUserRepo.On("UpdateStats", tt.userID, tt.mockUser.Solved, tt.mockUser.Submitted+1).Return(nil)
				}
			}

			// 创建服务
			service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)

			// 执行测试
			result, err := service.CreateSubmission(tt.problemID, tt.code, tt.language, tt.userID)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else if tt.expectedSubmission {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.UserID)
				assert.Equal(t, tt.problemID, result.ProblemID)
				assert.Equal(t, tt.code, result.Code)
				if tt.language == "" {
					assert.Equal(t, "verilog", result.Language)
				} else {
					assert.Equal(t, tt.language, result.Language)
				}
				assert.Equal(t, "pending", result.Status)
				assert.Equal(t, 0, result.Score)
			}

			// 验证 Mock 调用
			mockUserRepo.AssertExpectations(t)
			mockProblemRepo.AssertExpectations(t)
			mockSubmissionRepo.AssertExpectations(t)
		})
	}
}

// TestSubmissionService_GetSubmission 测试获取提交详情
func TestSubmissionService_GetSubmission(t *testing.T) {
	tests := []struct {
		name           string
		id             uint
		mockSubmission *domain.Submission
		repoError      error
		expectedError  string
	}{
		{
			name: "成功获取提交",
			id:   1,
			mockSubmission: &domain.Submission{
				ID:        1,
				UserID:    1,
				ProblemID: 1,
				Code:      "module test(); endmodule",
				Language:  "verilog",
				Status:    "accepted",
			},
		},
		{
			name:           "提交不存在",
			id:             999,
			mockSubmission: nil,
			expectedError:  "提交记录不存在",
		},
		{
			name:          "仓储错误",
			id:            1,
			repoError:     errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockSubmissionRepo := new(MockSubmissionRepository)
			mockProblemRepo := new(MockProblemRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockSubmissionRepo.On("GetByID", tt.id).Return(tt.mockSubmission, tt.repoError)

			// 创建服务
			service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)

			// 执行测试
			result, err := service.GetSubmission(tt.id)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockSubmission.ID, result.ID)
				assert.Equal(t, tt.mockSubmission.UserID, result.UserID)
				assert.Equal(t, tt.mockSubmission.ProblemID, result.ProblemID)
			}

			// 验证 Mock 调用
			mockSubmissionRepo.AssertExpectations(t)
		})
	}
}

// TestSubmissionService_ListSubmissions 测试获取提交列表
func TestSubmissionService_ListSubmissions(t *testing.T) {
	tests := []struct {
		name            string
		page            int
		limit           int
		userID          uint
		problemID       uint
		status          string
		mockSubmissions []domain.Submission
		mockTotal       int64
		repoError       error
		expectedPage    int
		expectedLimit   int
		expectedError   string
	}{
		{
			name:      "成功获取提交列表",
			page:      1,
			limit:     10,
			userID:    1,
			problemID: 1,
			status:    "accepted",
			mockSubmissions: []domain.Submission{
				{
					ID:        1,
					UserID:    1,
					ProblemID: 1,
					Status:    "accepted",
				},
			},
			mockTotal:     1,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "页码修正",
			page:          0,
			limit:         10,
			expectedPage:  1,
			expectedLimit: 10,
			mockSubmissions: []domain.Submission{},
			mockTotal:       0,
		},
		{
			name:            "限制修正",
			page:            1,
			limit:           0,
			expectedPage:    1,
			expectedLimit:   20,
			mockSubmissions: []domain.Submission{},
			mockTotal:       0,
		},
		{
			name:            "限制超过最大值",
			page:            1,
			limit:           150,
			expectedPage:    1,
			expectedLimit:   20,
			mockSubmissions: []domain.Submission{},
			mockTotal:       0,
		},
		{
			name:          "仓储错误",
			page:          1,
			limit:         10,
			repoError:     errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockSubmissionRepo := new(MockSubmissionRepository)
			mockProblemRepo := new(MockProblemRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			// 模拟服务层的参数处理逻辑
			expectedPage := tt.page
			expectedLimit := tt.limit
			if expectedPage <= 0 {
				expectedPage = 1
			}
			if expectedLimit <= 0 || expectedLimit > 100 {
				expectedLimit = 20
			}
			mockSubmissionRepo.On("List", expectedPage, expectedLimit, tt.userID, tt.problemID, tt.status).Return(tt.mockSubmissions, tt.mockTotal, tt.repoError)

			// 创建服务
			service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)

			// 执行测试
			result, err := service.ListSubmissions(tt.page, tt.limit, tt.userID, tt.problemID, tt.status)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockSubmissions, result.Submissions)
				assert.Equal(t, tt.mockTotal, result.Total)
				assert.Equal(t, expectedPage, result.Page)
				assert.Equal(t, expectedLimit, result.Limit)
			}

			// 验证 Mock 调用
			mockSubmissionRepo.AssertExpectations(t)
		})
	}
}

// TestSubmissionService_UpdateSubmissionStatus 测试更新提交状态
func TestSubmissionService_UpdateSubmissionStatus(t *testing.T) {
	tests := []struct {
		name           string
		id             uint
		status         string
		score          int
		runTime        int
		memory         int
		errorMessage   string
		passedTests    int
		totalTests     int
		mockSubmission *domain.Submission
		mockUser       *domain.User
		acceptedCount  int64
		updateError    error
		getError       error
		expectedError  string
	}{
		{
			name:        "成功更新为pending状态",
			id:          1,
			status:      "pending",
			score:       0,
			runTime:     0,
			memory:      0,
			errorMessage: "",
			passedTests: 0,
			totalTests:  10,
		},
		{
			name:        "成功更新为accepted状态-首次通过",
			id:          1,
			status:      "accepted",
			score:       100,
			runTime:     500,
			memory:      1024,
			errorMessage: "",
			passedTests: 10,
			totalTests:  10,
			mockSubmission: &domain.Submission{
				ID:        1,
				UserID:    1,
				ProblemID: 1,
			},
			mockUser: &domain.User{
				ID:        1,
				Solved:    5,
				Submitted: 10,
			},
			acceptedCount: 1,
		},
		{
			name:        "成功更新为accepted状态-非首次通过",
			id:          1,
			status:      "accepted",
			score:       100,
			runTime:     500,
			memory:      1024,
			errorMessage: "",
			passedTests: 10,
			totalTests:  10,
			mockSubmission: &domain.Submission{
				ID:        1,
				UserID:    1,
				ProblemID: 1,
			},
			acceptedCount: 2,
		},
		{
			name:          "更新状态失败",
			id:            1,
			status:        "accepted",
			updateError:   errors.New("update failed"),
			expectedError: "update failed",
		},
		{
			name:        "accepted状态但获取提交失败",
			id:          1,
			status:      "accepted",
			getError:    errors.New("get submission failed"),
			expectedError: "get submission failed",
		},
		{
			name:           "accepted状态但提交不存在",
			id:             1,
			status:         "accepted",
			mockSubmission: nil,
			expectedError:  "提交记录不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockSubmissionRepo := new(MockSubmissionRepository)
			mockProblemRepo := new(MockProblemRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockSubmissionRepo.On("UpdateStatus", tt.id, tt.status, tt.score, tt.runTime, tt.memory, tt.errorMessage, tt.passedTests, tt.totalTests).Return(tt.updateError)

			if tt.updateError == nil && tt.status == "accepted" {
				mockSubmissionRepo.On("GetByID", tt.id).Return(tt.mockSubmission, tt.getError)
				if tt.getError == nil && tt.mockSubmission != nil {
					mockSubmissionRepo.On("CountAcceptedByUser", tt.mockSubmission.UserID, tt.mockSubmission.ProblemID).Return(tt.acceptedCount, nil)
					if tt.acceptedCount == 1 {
						mockUserRepo.On("GetByID", tt.mockSubmission.UserID).Return(tt.mockUser, nil)
						if tt.mockUser != nil {
							mockUserRepo.On("UpdateStats", tt.mockSubmission.UserID, tt.mockUser.Solved+1, tt.mockUser.Submitted).Return(nil)
						}
						mockProblemRepo.On("UpdateAcceptedCount", tt.mockSubmission.ProblemID, 1).Return(nil)
					}
				}
			}

			// 创建服务
			service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)

			// 执行测试
			err := service.UpdateSubmissionStatus(tt.id, tt.status, tt.score, tt.runTime, tt.memory, tt.errorMessage, tt.passedTests, tt.totalTests)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockSubmissionRepo.AssertExpectations(t)
			mockProblemRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestSubmissionService_GetUserSubmissions 测试获取用户提交记录
func TestSubmissionService_GetUserSubmissions(t *testing.T) {
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockProblemRepo := new(MockProblemRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uint(1)
	page := 1
	limit := 10
	mockSubmissions := []domain.Submission{
		{ID: 1, UserID: userID, ProblemID: 1},
	}
	mockTotal := int64(1)

	mockSubmissionRepo.On("List", page, limit, userID, uint(0), "").Return(mockSubmissions, mockTotal, nil)

	service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)
	result, err := service.GetUserSubmissions(userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockSubmissions, result.Submissions)
	assert.Equal(t, mockTotal, result.Total)
	mockSubmissionRepo.AssertExpectations(t)
}

// TestSubmissionService_GetProblemSubmissions 测试获取题目提交记录
func TestSubmissionService_GetProblemSubmissions(t *testing.T) {
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockProblemRepo := new(MockProblemRepository)
	mockUserRepo := new(MockUserRepository)

	problemID := uint(1)
	page := 1
	limit := 10
	mockSubmissions := []domain.Submission{
		{ID: 1, UserID: 1, ProblemID: problemID},
	}
	mockTotal := int64(1)

	mockSubmissionRepo.On("List", page, limit, uint(0), problemID, "").Return(mockSubmissions, mockTotal, nil)

	service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)
	result, err := service.GetProblemSubmissions(problemID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockSubmissions, result.Submissions)
	assert.Equal(t, mockTotal, result.Total)
	mockSubmissionRepo.AssertExpectations(t)
}

// TestSubmissionService_GetSubmissionStats 测试获取提交统计信息
func TestSubmissionService_GetSubmissionStats(t *testing.T) {
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockProblemRepo := new(MockProblemRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uint(1)
	mockStats := map[string]interface{}{
		"total_submitted": 10,
		"total_accepted":  5,
		"acceptance_rate":  0.5,
	}

	mockSubmissionRepo.On("GetStats", userID).Return(mockStats, nil)

	service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)
	result, err := service.GetSubmissionStats(userID)

	assert.NoError(t, err)
	assert.Equal(t, mockStats, result)
	mockSubmissionRepo.AssertExpectations(t)
}

// TestSubmissionService_ValidateSubmissionAccess 测试验证提交访问权限
func TestSubmissionService_ValidateSubmissionAccess(t *testing.T) {
	tests := []struct {
		name           string
		submissionID   uint
		userID         uint
		mockSubmission *domain.Submission
		mockUser       *domain.User
		submissionError error
		userError      error
		expectedError  string
	}{
		{
			name:         "提交者本人访问",
			submissionID: 1,
			userID:       1,
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
			mockUser: &domain.User{
				ID:   1,
				Role: "student",
			},
		},
		{
			name:         "管理员访问",
			submissionID: 1,
			userID:       2,
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
			mockUser: &domain.User{
				ID:   2,
				Role: "admin",
			},
		},
		{
			name:         "教师访问",
			submissionID: 1,
			userID:       2,
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
			mockUser: &domain.User{
				ID:   2,
				Role: "teacher",
			},
		},
		{
			name:         "无权限访问",
			submissionID: 1,
			userID:       2,
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
			mockUser: &domain.User{
				ID:   2,
				Role: "student",
			},
			expectedError: "没有权限访问此提交记录",
		},
		{
			name:            "提交不存在",
			submissionID:    999,
			userID:          1,
			mockSubmission:  nil,
			expectedError:   "提交记录不存在",
		},
		{
			name:         "用户不存在",
			submissionID: 1,
			userID:       999,
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
			mockUser:      nil,
			expectedError: "用户不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockSubmissionRepo := new(MockSubmissionRepository)
			mockProblemRepo := new(MockProblemRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockSubmissionRepo.On("GetByID", tt.submissionID).Return(tt.mockSubmission, tt.submissionError)
			if tt.mockSubmission != nil {
				mockUserRepo.On("GetByID", tt.userID).Return(tt.mockUser, tt.userError)
			}

			// 创建服务
			service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)

			// 执行测试
			err := service.ValidateSubmissionAccess(tt.submissionID, tt.userID)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockSubmissionRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestSubmissionService_DeleteSubmission 测试删除提交记录
func TestSubmissionService_DeleteSubmission(t *testing.T) {
	tests := []struct {
		name           string
		id             uint
		userID         uint
		userRole       string
		mockSubmission *domain.Submission
		getError       error
		deleteError    error
		expectedError  string
	}{
		{
			name:     "提交者删除自己的提交",
			id:       1,
			userID:   1,
			userRole: "student",
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
		},
		{
			name:     "管理员删除提交",
			id:       1,
			userID:   2,
			userRole: "admin",
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
		},
		{
			name:     "无权限删除",
			id:       1,
			userID:   2,
			userRole: "student",
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
			expectedError: "没有权限删除此提交记录",
		},
		{
			name:          "获取提交失败",
			id:            1,
			userID:        1,
			userRole:      "student",
			getError:      errors.New("get submission failed"),
			expectedError: "get submission failed",
		},
		{
			name:     "删除失败",
			id:       1,
			userID:   1,
			userRole: "student",
			mockSubmission: &domain.Submission{
				ID:     1,
				UserID: 1,
			},
			deleteError:   errors.New("delete failed"),
			expectedError: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockSubmissionRepo := new(MockSubmissionRepository)
			mockProblemRepo := new(MockProblemRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockSubmissionRepo.On("GetByID", tt.id).Return(tt.mockSubmission, tt.getError)
			if tt.getError == nil && tt.mockSubmission != nil && (tt.mockSubmission.UserID == tt.userID || tt.userRole == "admin") {
				mockSubmissionRepo.On("SoftDelete", tt.id).Return(tt.deleteError)
			}

			// 创建服务
			service := services.NewSubmissionService(mockSubmissionRepo, mockProblemRepo, mockUserRepo)

			// 执行测试
			err := service.DeleteSubmission(tt.id, tt.userID, tt.userRole)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockSubmissionRepo.AssertExpectations(t)
		})
	}
}