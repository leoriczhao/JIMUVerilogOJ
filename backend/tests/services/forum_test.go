package services

import (
	"errors"
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockForumRepository Mock 论坛仓储
type MockForumRepository struct {
	mock.Mock
}

func (m *MockForumRepository) ToggleLike(userID uint, postID, replyID *uint) error {
	args := m.Called(userID, postID, replyID)
	return args.Error(0)
}

func (m *MockForumRepository) CreatePost(post *domain.ForumPost) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockForumRepository) GetPostByID(id uint) (*domain.ForumPost, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.ForumPost), args.Error(1)
}

func (m *MockForumRepository) ListPosts(page, limit int, filters map[string]interface{}) ([]domain.ForumPost, int64, error) {
	args := m.Called(page, limit, filters)
	return args.Get(0).([]domain.ForumPost), args.Get(1).(int64), args.Error(2)
}

func (m *MockForumRepository) UpdatePost(post *domain.ForumPost) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockForumRepository) DeletePost(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForumRepository) IncrementPostViewCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForumRepository) IncrementPostReplyCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForumRepository) DecrementPostReplyCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForumRepository) CreateReply(reply *domain.ForumReply) error {
	args := m.Called(reply)
	return args.Error(0)
}

func (m *MockForumRepository) GetRepliesByPostID(postID uint, page, limit int) ([]domain.ForumReply, int64, error) {
	args := m.Called(postID, page, limit)
	return args.Get(0).([]domain.ForumReply), args.Get(1).(int64), args.Error(2)
}

func (m *MockForumRepository) UpdateReply(reply *domain.ForumReply) error {
	args := m.Called(reply)
	return args.Error(0)
}

func (m *MockForumRepository) DeleteReply(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForumRepository) DeleteRepliesByPostID(postID uint) error {
	args := m.Called(postID)
	return args.Error(0)
}

func (m *MockForumRepository) CreateLike(like *domain.ForumLike) error {
	args := m.Called(like)
	return args.Error(0)
}

func (m *MockForumRepository) DeleteLike(userID uint, postID, replyID *uint) error {
	args := m.Called(userID, postID, replyID)
	return args.Error(0)
}

func (m *MockForumRepository) CheckLikeExists(userID uint, postID, replyID *uint) (bool, error) {
	args := m.Called(userID, postID, replyID)
	return args.Bool(0), args.Error(1)
}

func (m *MockForumRepository) GetLikeCount(postID, replyID *uint) (int64, error) {
	args := m.Called(postID, replyID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockForumRepository) GetUserLikes(userID uint, page, limit int) ([]domain.ForumLike, int64, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]domain.ForumLike), args.Get(1).(int64), args.Error(2)
}

// TestForumService_CreatePost 测试创建帖子
func TestForumService_CreatePost(t *testing.T) {
	tests := []struct {
		name          string
		post          *domain.ForumPost
		mockUser      *domain.User
		userRepoError error
		createError   error
		expectedError string
	}{
		{
			name: "成功创建帖子",
			post: &domain.ForumPost{
				Title:    "测试帖子标题",
				Content:  "这是一个测试帖子内容",
				AuthorID: 1,
				Category: "技术讨论",
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
		},
		{
			name: "标题为空",
			post: &domain.ForumPost{
				Title:    "",
				Content:  "这是一个测试帖子内容",
				AuthorID: 1,
			},
			expectedError: "帖子标题不能为空",
		},
		{
			name: "内容为空",
			post: &domain.ForumPost{
				Title:    "测试帖子标题",
				Content:  "",
				AuthorID: 1,
			},
			expectedError: "帖子内容不能为空",
		},
		{
			name: "作者ID为空",
			post: &domain.ForumPost{
				Title:    "测试帖子标题",
				Content:  "这是一个测试帖子内容",
				AuthorID: 0,
			},
			expectedError: "用户ID不能为空",
		},
		{
			name: "作者不存在",
			post: &domain.ForumPost{
				Title:    "测试帖子标题",
				Content:  "这是一个测试帖子内容",
				AuthorID: 1,
			},
			mockUser:      nil,
			expectedError: "用户不存在",
		},
		{
			name: "用户仓储错误",
			post: &domain.ForumPost{
				Title:    "测试帖子标题",
				Content:  "这是一个测试帖子内容",
				AuthorID: 1,
			},
			userRepoError: errors.New("user repo error"),
			expectedError: "user repo error",
		},
		{
			name: "创建帖子失败",
			post: &domain.ForumPost{
				Title:    "测试帖子标题",
				Content:  "这是一个测试帖子内容",
				AuthorID: 1,
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			createError:   errors.New("create failed"),
			expectedError: "create failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockForumRepo := new(MockForumRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			// 只有在标题和内容都不为空且 AuthorID 不为 0 时才会调用 GetByID
			if tt.post.AuthorID != 0 && tt.post.Title != "" && tt.post.Content != "" {
				mockUserRepo.On("GetByID", tt.post.AuthorID).Return(tt.mockUser, tt.userRepoError)
			}

			if tt.userRepoError == nil && tt.mockUser != nil && tt.post.Title != "" && tt.post.Content != "" {
				mockForumRepo.On("CreatePost", tt.post).Return(tt.createError)
			}

			// 创建服务
			service := services.NewForumService(mockForumRepo, mockUserRepo)

			// 执行测试
			err := service.CreatePost(tt.post)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockForumRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestForumService_GetPost 测试获取帖子详情
func TestForumService_GetPost(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockPost      *domain.ForumPost
		repoError     error
		viewError     error
		expectedError string
	}{
		{
			name: "成功获取帖子",
			id:   1,
			mockPost: &domain.ForumPost{
				ID:        1,
				Title:     "测试帖子",
				Content:   "测试内容",
				AuthorID:  1,
				ViewCount: 10,
			},
		},
		{
			name:          "帖子不存在",
			id:            999,
			mockPost:      nil,
			expectedError: "帖子不存在",
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
			mockForumRepo := new(MockForumRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			mockForumRepo.On("GetPostByID", tt.id).Return(tt.mockPost, tt.repoError)
			if tt.mockPost != nil && tt.repoError == nil {
				mockForumRepo.On("IncrementPostViewCount", tt.id).Return(tt.viewError)
			}

			// 创建服务
			service := services.NewForumService(mockForumRepo, mockUserRepo)

			// 执行测试
			result, err := service.GetPost(tt.id)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockPost.ID, result.ID)
				assert.Equal(t, tt.mockPost.Title, result.Title)
				assert.Equal(t, tt.mockPost.Content, result.Content)
			}

			// 验证 Mock 调用
			mockForumRepo.AssertExpectations(t)
		})
	}
}

// TestForumService_ListPosts 测试获取帖子列表
func TestForumService_ListPosts(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		category      string
		mockPosts     []domain.ForumPost
		mockTotal     int64
		repoError     error
		expectedPage  int
		expectedLimit int
		expectedError string
	}{
		{
			name:     "成功获取帖子列表",
			page:     1,
			limit:    10,
			category: "技术讨论",
			mockPosts: []domain.ForumPost{
				{
					ID:       1,
					Title:    "帖子1",
					Content:  "内容1",
					Category: "技术讨论",
				},
				{
					ID:       2,
					Title:    "帖子2",
					Content:  "内容2",
					Category: "技术讨论",
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
			mockPosts:     []domain.ForumPost{},
			mockTotal:     0,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "限制修正",
			page:          1,
			limit:         0,
			mockPosts:     []domain.ForumPost{},
			mockTotal:     0,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "限制超过最大值",
			page:          1,
			limit:         150,
			mockPosts:     []domain.ForumPost{},
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
			mockForumRepo := new(MockForumRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望 - 模拟服务层的参数处理逻辑
			expectedPage := tt.page
			expectedLimit := tt.limit
			if expectedPage <= 0 {
				expectedPage = 1
			}
			if expectedLimit <= 0 || expectedLimit > 100 {
				expectedLimit = 20
			}

			// 设置 Mock 期望 - 使用 AnythingOfType 来匹配 map 类型
			mockForumRepo.On("ListPosts", expectedPage, expectedLimit, mock.AnythingOfType("map[string]interface {}")).Return(tt.mockPosts, tt.mockTotal, tt.repoError)

			// 创建服务
			service := services.NewForumService(mockForumRepo, mockUserRepo)

			// 执行测试
			filters := map[string]interface{}{}
			if tt.category != "" {
				filters["category"] = tt.category
			}
			result, total, err := service.ListPosts(tt.page, tt.limit, filters)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockPosts, result)
				assert.Equal(t, tt.mockTotal, total)
			}

			// 验证 Mock 调用
			mockForumRepo.AssertExpectations(t)
		})
	}
}

// TestForumService_CreateReply 测试创建回复
func TestForumService_CreateReply(t *testing.T) {
	tests := []struct {
		name          string
		reply         *domain.ForumReply
		mockPost      *domain.ForumPost
		mockUser      *domain.User
		postRepoError error
		userRepoError error
		createError   error
		expectedError string
	}{
		{
			name: "成功创建回复",
			reply: &domain.ForumReply{
				Content:  "这是一个测试回复",
				PostID:   1,
				AuthorID: 1,
			},
			mockPost: &domain.ForumPost{
				ID:      1,
				Title:   "测试帖子",
				Content: "测试内容",
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
		},
		{
			name: "内容为空",
			reply: &domain.ForumReply{
				Content:  "",
				PostID:   1,
				AuthorID: 1,
			},
			expectedError: "回复内容不能为空",
		},
		{
			name: "帖子ID为空",
			reply: &domain.ForumReply{
				Content:  "这是一个测试回复",
				PostID:   0,
				AuthorID: 1,
			},
			expectedError: "帖子ID不能为空",
		},
		{
			name: "作者ID为空",
			reply: &domain.ForumReply{
				Content:  "这是一个测试回复",
				PostID:   1,
				AuthorID: 0,
			},
			expectedError: "用户ID不能为空",
		},
		{
			name: "帖子不存在",
			reply: &domain.ForumReply{
				Content:  "这是一个测试回复",
				PostID:   1,
				AuthorID: 1,
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			mockPost:      nil,
			expectedError: "帖子不存在",
		},
		{
			name: "作者不存在",
			reply: &domain.ForumReply{
				Content:  "这是一个测试回复",
				PostID:   1,
				AuthorID: 1,
			},
			mockPost: &domain.ForumPost{
				ID:      1,
				Title:   "测试帖子",
				Content: "测试内容",
			},
			mockUser:      nil,
			expectedError: "用户不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockForumRepo := new(MockForumRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			// 只有在内容、用户ID和帖子ID都不为空时才会调用验证方法
			if tt.reply.Content != "" && tt.reply.AuthorID != 0 && tt.reply.PostID != 0 {
				mockUserRepo.On("GetByID", tt.reply.AuthorID).Return(tt.mockUser, tt.userRepoError)
				// 如果用户验证通过，再设置帖子验证的 Mock
				if tt.userRepoError == nil && tt.mockUser != nil {
					mockForumRepo.On("GetPostByID", tt.reply.PostID).Return(tt.mockPost, tt.postRepoError)
				}
			}
			// 如果所有验证都通过且内容不为空，设置创建回复的 Mock
			if tt.postRepoError == nil && tt.userRepoError == nil && tt.mockPost != nil && tt.mockUser != nil && tt.reply.Content != "" {
				mockForumRepo.On("CreateReply", tt.reply).Return(tt.createError)
				// 如果创建成功，还需要增加帖子回复数
				if tt.createError == nil {
					mockForumRepo.On("IncrementPostReplyCount", tt.reply.PostID).Return(nil)
				}
			}

			// 创建服务
			service := services.NewForumService(mockForumRepo, mockUserRepo)

			// 执行测试
			err := service.CreateReply(tt.reply)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockForumRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestForumService_GetReplies 测试获取回复列表
func TestForumService_GetReplies(t *testing.T) {
	tests := []struct {
		name          string
		postID        uint
		page          int
		limit         int
		mockReplies   []domain.ForumReply
		mockTotal     int64
		repoError     error
		expectedPage  int
		expectedLimit int
		expectedError string
	}{
		{
			name:   "成功获取回复列表",
			postID: 1,
			page:   1,
			limit:  10,
			mockReplies: []domain.ForumReply{
				{
					ID:       1,
					Content:  "回复1",
					PostID:   1,
					AuthorID: 1,
				},
				{
					ID:       2,
					Content:  "回复2",
					PostID:   1,
					AuthorID: 2,
				},
			},
			mockTotal:     2,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "页码修正",
			postID:        1,
			page:          0,
			limit:         10,
			mockReplies:   []domain.ForumReply{},
			mockTotal:     0,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "限制修正",
			postID:        1,
			page:          1,
			limit:         0,
			mockReplies:   []domain.ForumReply{},
			mockTotal:     0,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "仓储错误",
			postID:        1,
			page:          1,
			limit:         10,
			repoError:     errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockForumRepo := new(MockForumRepository)
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

			// 首先验证帖子是否存在
			if tt.postID != 0 {
				mockPost := &domain.ForumPost{ID: tt.postID}
				mockForumRepo.On("GetPostByID", tt.postID).Return(mockPost, nil)
				// 如果帖子验证通过，再设置获取回复列表的 Mock
				mockForumRepo.On("GetRepliesByPostID", tt.postID, expectedPage, expectedLimit).Return(tt.mockReplies, tt.mockTotal, tt.repoError)
			}

			// 创建服务
			service := services.NewForumService(mockForumRepo, mockUserRepo)

			// 执行测试
			result, total, err := service.ListReplies(tt.postID, tt.page, tt.limit)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReplies, result)
				assert.Equal(t, tt.mockTotal, total)
			}

			// 验证 Mock 调用
			mockForumRepo.AssertExpectations(t)
		})
	}
}

// TestForumService_LikePost 测试点赞帖子
func TestForumService_LikePost(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		postID        uint
		mockPost      *domain.ForumPost
		mockUser      *domain.User
		isLiked       bool
		postRepoError error
		userRepoError error
		likeError     error
		likeRepoError error
		expectedError string
	}{
		{
			name:   "成功点赞帖子",
			userID: 1,
			postID: 1,
			mockPost: &domain.ForumPost{
				ID:      1,
				Title:   "测试帖子",
				Content: "测试内容",
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			isLiked: false,
		},
		{
			name:   "仓储错误",
			userID: 1,
			postID: 1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			isLiked:       false,
			likeRepoError: errors.New("database error"),
			expectedError: "database error",
		},
		{
			name:   "用户不存在",
			userID: 999,
			postID: 1,
			mockPost: &domain.ForumPost{
				ID:      1,
				Title:   "测试帖子",
				Content: "测试内容",
			},
			mockUser:      nil,
			expectedError: "用户不存在",
		},
		{
			name:   "取消点赞",
			userID: 1,
			postID: 1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			isLiked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockForumRepo := new(MockForumRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			postIDPtr := &tt.postID

			// 首先验证用户是否存在
			if tt.userID != 0 {
				mockUserRepo.On("GetByID", tt.userID).Return(tt.mockUser, nil)
			}

			// 如果用户存在，检查点赞状态
			if tt.mockUser != nil {
				mockForumRepo.On("CheckLikeExists", tt.userID, postIDPtr, (*uint)(nil)).Return(tt.isLiked, nil)

				// 根据点赞状态执行相应操作
				if tt.isLiked {
					// 取消点赞
					mockForumRepo.On("DeleteLike", tt.userID, postIDPtr, (*uint)(nil)).Return(tt.likeRepoError)
				} else {
					// 添加点赞
					mockForumRepo.On("CreateLike", mock.AnythingOfType("*domain.ForumLike")).Return(tt.likeRepoError)
				}
			}

			// 创建服务
			service := services.NewForumService(mockForumRepo, mockUserRepo)

			// 执行测试
			err := service.ToggleLike(tt.userID, postIDPtr, nil)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockForumRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestForumService_UnlikePost 测试取消点赞帖子
func TestForumService_UnlikePost(t *testing.T) {
	tests := []struct {
		name            string
		userID          uint
		postID          uint
		mockPost        *domain.ForumPost
		mockUser        *domain.User
		isLiked         bool
		postRepoError   error
		userRepoError   error
		likeError       error
		unlikeRepoError error
		expectedError   string
	}{
		{
			name:   "成功取消点赞",
			userID: 1,
			postID: 1,
			mockPost: &domain.ForumPost{
				ID:      1,
				Title:   "测试帖子",
				Content: "测试内容",
			},
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			isLiked: true,
		},
		{
			name:   "仓储错误",
			userID: 1,
			postID: 1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			likeError:     errors.New("repository error"),
			expectedError: "repository error",
		},
		{
			name:   "用户不存在",
			userID: 999,
			postID: 1,
			mockPost: &domain.ForumPost{
				ID:      1,
				Title:   "测试帖子",
				Content: "测试内容",
			},
			mockUser:      nil,
			expectedError: "用户不存在",
		},
		{
			name:   "添加点赞",
			userID: 1,
			postID: 1,
			mockUser: &domain.User{
				ID:       1,
				Username: "testuser",
			},
			isLiked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockForumRepo := new(MockForumRepository)
			mockUserRepo := new(MockUserRepository)

			// 设置 Mock 期望
			postIDPtr := &tt.postID

			// 首先验证用户是否存在
			if tt.userID != 0 {
				mockUserRepo.On("GetByID", tt.userID).Return(tt.mockUser, nil)
			}

			// 如果用户存在，检查点赞状态
			if tt.mockUser != nil {
				mockForumRepo.On("CheckLikeExists", tt.userID, postIDPtr, (*uint)(nil)).Return(tt.isLiked, tt.likeError)

				// 如果CheckLikeExists没有错误，根据点赞状态执行相应操作
				if tt.likeError == nil {
					if tt.isLiked {
						// 取消点赞
						mockForumRepo.On("DeleteLike", tt.userID, postIDPtr, (*uint)(nil)).Return(tt.unlikeRepoError)
					} else {
						// 添加点赞
						mockForumRepo.On("CreateLike", mock.AnythingOfType("*domain.ForumLike")).Return(tt.unlikeRepoError)
					}
				}
			}

			// 创建服务
			service := services.NewForumService(mockForumRepo, mockUserRepo)

			// 执行测试
			err := service.ToggleLike(tt.userID, postIDPtr, nil)

			// 验证结果
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// 验证 Mock 调用
			mockForumRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

// TestForumService_GetLikeCount 测试获取点赞数
func TestForumService_GetLikeCount(t *testing.T) {
	mockForumRepo := new(MockForumRepository)
	mockUserRepo := new(MockUserRepository)

	postID := uint(1)
	postIDPtr := &postID
	expectedCount := int64(10)

	mockForumRepo.On("GetLikeCount", postIDPtr, (*uint)(nil)).Return(expectedCount, nil)

	service := services.NewForumService(mockForumRepo, mockUserRepo)
	result, err := service.GetLikeCount(postIDPtr, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
	mockForumRepo.AssertExpectations(t)
}

// TestForumService_CheckUserLike 测试检查用户是否点赞
func TestForumService_CheckUserLike(t *testing.T) {
	mockForumRepo := new(MockForumRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uint(1)
	postID := uint(1)
	postIDPtr := &postID
	expectedLiked := true

	mockForumRepo.On("CheckLikeExists", userID, postIDPtr, (*uint)(nil)).Return(expectedLiked, nil)

	service := services.NewForumService(mockForumRepo, mockUserRepo)
	result, err := service.CheckUserLike(userID, postIDPtr, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedLiked, result)
	mockForumRepo.AssertExpectations(t)
}
