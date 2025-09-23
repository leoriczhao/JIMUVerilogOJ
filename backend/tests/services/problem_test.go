package services

import (
	"errors"
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProblemRepository 模拟题目仓储
type MockProblemRepository struct {
	mock.Mock
}

func (m *MockProblemRepository) Create(problem *domain.Problem) error {
	args := m.Called(problem)
	return args.Error(0)
}

func (m *MockProblemRepository) GetByID(id uint) (*domain.Problem, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Problem), args.Error(1)
}

func (m *MockProblemRepository) List(page, limit int, filters map[string]interface{}) ([]domain.Problem, int64, error) {
	args := m.Called(page, limit, filters)
	return args.Get(0).([]domain.Problem), args.Get(1).(int64), args.Error(2)
}

func (m *MockProblemRepository) Update(problem *domain.Problem) error {
	args := m.Called(problem)
	return args.Error(0)
}

func (m *MockProblemRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProblemRepository) UpdateSubmitCount(id uint, increment int) error {
	args := m.Called(id, increment)
	return args.Error(0)
}

func (m *MockProblemRepository) UpdateAcceptedCount(id uint, increment int) error {
	args := m.Called(id, increment)
	return args.Error(0)
}

func (m *MockProblemRepository) GetTestCases(problemID uint) ([]domain.TestCase, error) {
	args := m.Called(problemID)
	return args.Get(0).([]domain.TestCase), args.Error(1)
}

func (m *MockProblemRepository) CreateTestCase(testCase *domain.TestCase) error {
	args := m.Called(testCase)
	return args.Error(0)
}

func (m *MockProblemRepository) DeleteTestCases(problemID uint) error {
	args := m.Called(problemID)
	return args.Error(0)
}

// TestProblemService_CreateProblem 测试创建题目
func TestProblemService_CreateProblem(t *testing.T) {
	tests := []struct {
		name    string
		problem *domain.Problem
		mockFn  func(*MockProblemRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功创建题目",
			problem: &domain.Problem{
				Title:       "测试题目",
				Description: "这是一个测试题目",
				TimeLimit:   1000,
				MemoryLimit: 256,
			},
			mockFn: func(m *MockProblemRepository) {
				m.On("Create", mock.AnythingOfType("*domain.Problem")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "标题为空",
			problem: &domain.Problem{
				Title:       "",
				Description: "这是一个测试题目",
				TimeLimit:   1000,
				MemoryLimit: 256,
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "题目标题不能为空",
		},
		{
			name: "描述为空",
			problem: &domain.Problem{
				Title:       "测试题目",
				Description: "",
				TimeLimit:   1000,
				MemoryLimit: 256,
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "题目描述不能为空",
		},
		{
			name: "时间限制无效",
			problem: &domain.Problem{
				Title:       "测试题目",
				Description: "这是一个测试题目",
				TimeLimit:   0,
				MemoryLimit: 256,
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "时间限制必须大于0",
		},
		{
			name: "内存限制无效",
			problem: &domain.Problem{
				Title:       "测试题目",
				Description: "这是一个测试题目",
				TimeLimit:   1000,
				MemoryLimit: 0,
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "内存限制必须大于0",
		},
		{
			name: "数据库错误",
			problem: &domain.Problem{
				Title:       "测试题目",
				Description: "这是一个测试题目",
				TimeLimit:   1000,
				MemoryLimit: 256,
			},
			mockFn: func(m *MockProblemRepository) {
				m.On("Create", mock.AnythingOfType("*domain.Problem")).Return(errors.New("数据库错误"))
			},
			wantErr: true,
			errMsg:  "数据库错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProblemRepository)
			tt.mockFn(mockRepo)

			service := services.NewProblemService(mockRepo)
			err := service.CreateProblem(tt.problem)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestProblemService_GetProblem 测试获取题目
func TestProblemService_GetProblem(t *testing.T) {
	tests := []struct {
		name     string
		id       uint
		mockFn   func(*MockProblemRepository)
		wantErr  bool
		errMsg   string
		expected *domain.Problem
	}{
		{
			name: "成功获取题目",
			id:   1,
			mockFn: func(m *MockProblemRepository) {
				problem := &domain.Problem{
					ID:          1,
					Title:       "测试题目",
					Description: "这是一个测试题目",
					TimeLimit:   1000,
					MemoryLimit: 256,
				}
				m.On("GetByID", uint(1)).Return(problem, nil)
			},
			wantErr: false,
			expected: &domain.Problem{
				ID:          1,
				Title:       "测试题目",
				Description: "这是一个测试题目",
				TimeLimit:   1000,
				MemoryLimit: 256,
			},
		},
		{
			name: "题目不存在",
			id:   999,
			mockFn: func(m *MockProblemRepository) {
				m.On("GetByID", uint(999)).Return((*domain.Problem)(nil), nil)
			},
			wantErr: true,
			errMsg:  "题目不存在",
		},
		{
			name: "数据库错误",
			id:   1,
			mockFn: func(m *MockProblemRepository) {
				m.On("GetByID", uint(1)).Return((*domain.Problem)(nil), errors.New("数据库错误"))
			},
			wantErr: true,
			errMsg:  "数据库错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProblemRepository)
			tt.mockFn(mockRepo)

			service := services.NewProblemService(mockRepo)
			result, err := service.GetProblem(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestProblemService_ListProblems 测试获取题目列表
func TestProblemService_ListProblems(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		filters       map[string]interface{}
		mockFn        func(*MockProblemRepository)
		wantErr       bool
		expectedPage  int
		expectedLimit int
	}{
		{
			name:    "成功获取题目列表",
			page:    1,
			limit:   10,
			filters: map[string]interface{}{"difficulty": "easy"},
			mockFn: func(m *MockProblemRepository) {
				problems := []domain.Problem{
					{ID: 1, Title: "题目1"},
					{ID: 2, Title: "题目2"},
				}
				m.On("List", 1, 10, map[string]interface{}{"difficulty": "easy"}).Return(problems, int64(2), nil)
			},
			wantErr:       false,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:    "页码无效时使用默认值",
			page:    0,
			limit:   10,
			filters: nil,
			mockFn: func(m *MockProblemRepository) {
				m.On("List", 1, 10, mock.Anything).Return([]domain.Problem{}, int64(0), nil)
			},
			wantErr:       false,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:    "限制数量无效时使用默认值",
			page:    1,
			limit:   0,
			filters: nil,
			mockFn: func(m *MockProblemRepository) {
				m.On("List", 1, 20, mock.Anything).Return([]domain.Problem{}, int64(0), nil)
			},
			wantErr:       false,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:    "限制数量过大时使用默认值",
			page:    1,
			limit:   200,
			filters: nil,
			mockFn: func(m *MockProblemRepository) {
				m.On("List", 1, 20, mock.Anything).Return([]domain.Problem{}, int64(0), nil)
			},
			wantErr:       false,
			expectedPage:  1,
			expectedLimit: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProblemRepository)
			tt.mockFn(mockRepo)

			service := services.NewProblemService(mockRepo)
			_, _, err := service.ListProblems(tt.page, tt.limit, tt.filters)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestProblemService_UpdateProblem 测试更新题目
func TestProblemService_UpdateProblem(t *testing.T) {
	tests := []struct {
		name    string
		problem *domain.Problem
		mockFn  func(*MockProblemRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功更新题目",
			problem: &domain.Problem{
				ID:          1,
				Title:       "更新后的题目",
				Description: "更新后的描述",
			},
			mockFn: func(m *MockProblemRepository) {
				m.On("Update", mock.AnythingOfType("*domain.Problem")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "标题为空",
			problem: &domain.Problem{
				ID:          1,
				Title:       "",
				Description: "更新后的描述",
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "题目标题不能为空",
		},
		{
			name: "描述为空",
			problem: &domain.Problem{
				ID:          1,
				Title:       "更新后的题目",
				Description: "",
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "题目描述不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProblemRepository)
			tt.mockFn(mockRepo)

			service := services.NewProblemService(mockRepo)
			err := service.UpdateProblem(tt.problem)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestProblemService_DeleteProblem 测试删除题目
func TestProblemService_DeleteProblem(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mockFn  func(*MockProblemRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功删除题目",
			id:   1,
			mockFn: func(m *MockProblemRepository) {
				problem := &domain.Problem{ID: 1, Title: "测试题目"}
				m.On("GetByID", uint(1)).Return(problem, nil)
				m.On("DeleteTestCases", uint(1)).Return(nil)
				m.On("Delete", uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "题目不存在",
			id:   999,
			mockFn: func(m *MockProblemRepository) {
				m.On("GetByID", uint(999)).Return((*domain.Problem)(nil), nil)
			},
			wantErr: true,
			errMsg:  "题目不存在",
		},
		{
			name: "获取题目时数据库错误",
			id:   1,
			mockFn: func(m *MockProblemRepository) {
				m.On("GetByID", uint(1)).Return((*domain.Problem)(nil), errors.New("数据库错误"))
			},
			wantErr: true,
			errMsg:  "数据库错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProblemRepository)
			tt.mockFn(mockRepo)

			service := services.NewProblemService(mockRepo)
			err := service.DeleteProblem(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestProblemService_AddTestCase 测试添加测试用例
func TestProblemService_AddTestCase(t *testing.T) {
	tests := []struct {
		name     string
		testCase *domain.TestCase
		mockFn   func(*MockProblemRepository)
		wantErr  bool
		errMsg   string
	}{
		{
			name: "成功添加测试用例",
			testCase: &domain.TestCase{
				ProblemID: 1,
				Input:     "1 2",
				Output:    "3",
			},
			mockFn: func(m *MockProblemRepository) {
				problem := &domain.Problem{ID: 1, Title: "测试题目"}
				m.On("GetByID", uint(1)).Return(problem, nil)
				m.On("CreateTestCase", mock.AnythingOfType("*domain.TestCase")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "题目ID为空",
			testCase: &domain.TestCase{
				ProblemID: 0,
				Input:     "1 2",
				Output:    "3",
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "题目ID不能为空",
		},
		{
			name: "输入输出同时为空",
			testCase: &domain.TestCase{
				ProblemID: 1,
				Input:     "",
				Output:    "",
			},
			mockFn:  func(m *MockProblemRepository) {},
			wantErr: true,
			errMsg:  "测试用例输入和输出不能同时为空",
		},
		{
			name: "题目不存在",
			testCase: &domain.TestCase{
				ProblemID: 999,
				Input:     "1 2",
				Output:    "3",
			},
			mockFn: func(m *MockProblemRepository) {
				m.On("GetByID", uint(999)).Return((*domain.Problem)(nil), nil)
			},
			wantErr: true,
			errMsg:  "题目不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProblemRepository)
			tt.mockFn(mockRepo)

			service := services.NewProblemService(mockRepo)
			err := service.AddTestCase(tt.testCase)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestProblemService_UpdateProblemStats 测试更新题目统计
func TestProblemService_UpdateProblemStats(t *testing.T) {
	tests := []struct {
		name              string
		problemID         uint
		submitIncrement   int
		acceptedIncrement int
		mockFn            func(*MockProblemRepository)
		wantErr           bool
	}{
		{
			name:              "成功更新提交和通过统计",
			problemID:         1,
			submitIncrement:   1,
			acceptedIncrement: 1,
			mockFn: func(m *MockProblemRepository) {
				m.On("UpdateSubmitCount", uint(1), 1).Return(nil)
				m.On("UpdateAcceptedCount", uint(1), 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:              "只更新提交统计",
			problemID:         1,
			submitIncrement:   1,
			acceptedIncrement: 0,
			mockFn: func(m *MockProblemRepository) {
				m.On("UpdateSubmitCount", uint(1), 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:              "只更新通过统计",
			problemID:         1,
			submitIncrement:   0,
			acceptedIncrement: 1,
			mockFn: func(m *MockProblemRepository) {
				m.On("UpdateAcceptedCount", uint(1), 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:              "都不更新",
			problemID:         1,
			submitIncrement:   0,
			acceptedIncrement: 0,
			mockFn:            func(m *MockProblemRepository) {},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProblemRepository)
			tt.mockFn(mockRepo)

			service := services.NewProblemService(mockRepo)
			err := service.UpdateProblemStats(tt.problemID, tt.submitIncrement, tt.acceptedIncrement)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
