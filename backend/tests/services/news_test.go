package services

import (
	"errors"
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNewsRepository Mock 新闻仓储
type MockNewsRepository struct {
	mock.Mock
}

func (m *MockNewsRepository) Create(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsRepository) GetByID(id uint) (*domain.News, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.News), args.Error(1)
}

func (m *MockNewsRepository) List(page, limit int, filters map[string]interface{}) ([]domain.News, int64, error) {
	args := m.Called(page, limit, filters)
	return args.Get(0).([]domain.News), args.Get(1).(int64), args.Error(2)
}

func (m *MockNewsRepository) Update(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNewsRepository) IncrementViewCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestNewsService_CreateNews 测试创建新闻
func TestNewsService_CreateNews(t *testing.T) {
	tests := []struct {
		name          string
		news          *domain.News
		mockUser      *domain.User
		userRepoError error
		createError   error
		expectedError string
	}{
		{
			name: "成功创建新闻",
			news: &domain.News{
				Title:    "测试新闻标题",
				Content:  "这是一条测试新闻内容",
				AuthorID: 1,
				Category: "技术",
				Status:   "published",
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "admin",
				Role:     "admin",
			},
		},
		{
			name: "标题为空",
			news: &domain.News{
				Title:    "",
				Content:  "这是一条测试新闻内容",
				AuthorID: 1,
			},
			expectedError: "新闻标题不能为空",
		},
		{
			name: "内容为空",
			news: &domain.News{
				Title:    "测试新闻标题",
				Content:  "",
				AuthorID: 1,
			},
			expectedError: "新闻内容不能为空",
		},
		{
			name: "作者ID为空",
			news: &domain.News{
				Title:    "测试新闻标题",
				Content:  "这是一条测试新闻内容",
				AuthorID: 0,
			},
			expectedError: "作者ID不能为空",
		},
		{
			name: "作者不存在",
			news: &domain.News{
				Title:    "测试新闻标题",
				Content:  "这是一条测试新闻内容",
				AuthorID: 1,
			},
			mockUser:      nil,
			expectedError: "作者不存在",
		},
		{
			name: "用户仓储错误",
			news: &domain.News{
				Title:    "测试新闻标题",
				Content:  "这是一条测试新闻内容",
				AuthorID: 1,
			},
			userRepoError: errors.New("user repo error"),
			expectedError: "user repo error",
		},
		{
			name: "创建新闻失败",
			news: &domain.News{
				Title:    "测试新闻标题",
				Content:  "这是一条测试新闻内容",
				AuthorID: 1,
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "admin",
			},
			createError:   errors.New("create failed"),
			expectedError: "create failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockNewsRepo := new(MockNewsRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			// 只有当标题、内容和作者ID都不为空时才设置GetByID的Mock期望
			if tt.news.AuthorID != 0 && tt.news.Title != "" && tt.news.Content != "" {
				mockUserRepo.On("GetByID", tt.news.AuthorID).Return(tt.mockUser, tt.userRepoError)
			}

			if tt.userRepoError == nil && tt.mockUser != nil && tt.news.Title != "" && tt.news.Content != "" {
				mockNewsRepo.On("Create", tt.news).Return(tt.createError)
			}

			// 创建服务
			service := services.NewNewsService(mockNewsRepo, mockUserRepo)

			// 执行测试
			err := service.CreateNews(tt.news)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockNewsRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestNewsService_GetNewsList 测试获取新闻列表
func TestNewsService_GetNewsList(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		mockNews      []domain.News
		mockTotal     int64
		repoError     error
		expectedPage  int
		expectedLimit int
		expectedError string
	}{
		{
			name:  "成功获取新闻列表",
			page:  1,
			limit: 10,
			mockNews: []domain.News{
				{
					ID:      1,
					Title:   "新闻1",
					Content: "内容1",
					Status:  "published",
				},
				{
					ID:      2,
					Title:   "新闻2",
					Content: "内容2",
					Status:  "published",
				},
			},
			mockTotal:     2,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "页码修正",
			page:          0,
			limit:         10,
			mockNews:      []domain.News{},
			mockTotal:     0,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "限制修正",
			page:          1,
			limit:         0,
			mockNews:      []domain.News{},
			mockTotal:     0,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "限制超过最大值",
			page:          1,
			limit:         150,
			mockNews:      []domain.News{},
			mockTotal:     0,
			expectedPage:  1,
			expectedLimit: 20,
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
			mockNewsRepo := new(MockNewsRepository)
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

			// 使用 AnythingOfType 匹配 filters 参数
			mockNewsRepo.On("List", expectedPage, expectedLimit, mock.AnythingOfType("map[string]interface {}")).Return(tt.mockNews, tt.mockTotal, tt.repoError)

			// 创建服务
			service := services.NewNewsService(mockNewsRepo, mockUserRepo)

			// 执行测试
			result, total, err := service.GetNewsList(tt.page, tt.limit)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews, result)
				assert.Equal(t, tt.mockTotal, total)
			}

			// 验证 Mock 调用
			mockNewsRepo.AssertExpectations(t)
		})
	}
}

// TestNewsService_GetNews 测试获取新闻详情
func TestNewsService_GetNews(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockNews      *domain.News
		repoError     error
		viewError     error
		expectedError string
	}{
		{
			name: "成功获取新闻",
			id:   1,
			mockNews: &domain.News{
				ID:        1,
				Title:     "测试新闻",
				Content:   "测试内容",
				AuthorID:  1,
				Status:    "published",
				ViewCount: 10,
			},
		},
		{
			name:          "新闻不存在",
			id:            999,
			mockNews:      nil,
			expectedError: "新闻不存在",
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
			mockNewsRepo := new(MockNewsRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockNewsRepo.On("GetByID", tt.id).Return(tt.mockNews, tt.repoError)
			if tt.mockNews != nil && tt.repoError == nil {
				mockNewsRepo.On("IncrementViewCount", tt.id).Return(tt.viewError)
			}

			// 创建服务
			service := services.NewNewsService(mockNewsRepo, mockUserRepo)

			// 执行测试
			result, err := service.GetNews(tt.id)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockNews.ID, result.ID)
				assert.Equal(t, tt.mockNews.Title, result.Title)
				assert.Equal(t, tt.mockNews.Content, result.Content)
			}

			// 验证 Mock 调用
			mockNewsRepo.AssertExpectations(t)
		})
	}
}

// TestNewsService_GetNewsByID 测试根据ID获取新闻详情（兼容方法）
func TestNewsService_GetNewsByID(t *testing.T) {
	mockNewsRepo := new(MockNewsRepository)
	mockUserRepo := new(MockUserRepository)

	newsID := uint(1)
	mockNews := &domain.News{
		ID:      1,
		Title:   "测试新闻",
		Content: "测试内容",
	}

	mockNewsRepo.On("GetByID", newsID).Return(mockNews, nil)
	mockNewsRepo.On("IncrementViewCount", newsID).Return(nil)

	service := services.NewNewsService(mockNewsRepo, mockUserRepo)
	result, err := service.GetNewsByID(newsID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockNews.ID, result.ID)
	assert.Equal(t, mockNews.Title, result.Title)
	mockNewsRepo.AssertExpectations(t)
}

// TestNewsService_UpdateNews 测试更新新闻
func TestNewsService_UpdateNews(t *testing.T) {
	tests := []struct {
		name          string
		news          *domain.News
		updateError   error
		expectedError string
	}{
		{
			name: "成功更新新闻",
			news: &domain.News{
				ID:       1,
				Title:    "更新后的标题",
				Content:  "更新后的内容",
				AuthorID: 1,
			},
		},
		{
			name: "标题为空",
			news: &domain.News{
				ID:      1,
				Title:   "",
				Content: "更新后的内容",
			},
			expectedError: "新闻标题不能为空",
		},
		{
			name: "内容为空",
			news: &domain.News{
				ID:      1,
				Title:   "更新后的标题",
				Content: "",
			},
			expectedError: "新闻内容不能为空",
		},
		{
			name: "更新失败",
			news: &domain.News{
				ID:      1,
				Title:   "更新后的标题",
				Content: "更新后的内容",
			},
			updateError:   errors.New("update failed"),
			expectedError: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockNewsRepo := new(MockNewsRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			if tt.news.Title != "" && tt.news.Content != "" {
				mockNewsRepo.On("Update", tt.news).Return(tt.updateError)
			}

			// 创建服务
			service := services.NewNewsService(mockNewsRepo, mockUserRepo)

			// 执行测试
			err := service.UpdateNews(tt.news)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockNewsRepo.AssertExpectations(t)
		})
	}
}

// TestNewsService_DeleteNews 测试删除新闻
func TestNewsService_DeleteNews(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockNews      *domain.News
		getError      error
		deleteError   error
		expectedError string
	}{
		{
			name: "成功删除新闻",
			id:   1,
			mockNews: &domain.News{
				ID:      1,
				Title:   "要删除的新闻",
				Content: "新闻内容",
			},
		},
		{
			name:          "新闻不存在",
			id:            999,
			mockNews:      nil,
			expectedError: "新闻不存在",
		},
		{
			name:          "获取新闻失败",
			id:            1,
			getError:      errors.New("get news failed"),
			expectedError: "get news failed",
		},
		{
			name: "删除失败",
			id:   1,
			mockNews: &domain.News{
				ID:      1,
				Title:   "要删除的新闻",
				Content: "新闻内容",
			},
			deleteError:   errors.New("delete failed"),
			expectedError: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockNewsRepo := new(MockNewsRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockNewsRepo.On("GetByID", tt.id).Return(tt.mockNews, tt.getError)
			if tt.getError == nil && tt.mockNews != nil {
				mockNewsRepo.On("Delete", tt.id).Return(tt.deleteError)
			}

			// 创建服务
			service := services.NewNewsService(mockNewsRepo, mockUserRepo)

			// 执行测试
			err := service.DeleteNews(tt.id)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockNewsRepo.AssertExpectations(t)
		})
	}
}