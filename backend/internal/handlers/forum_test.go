package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockForumService is a mock for ForumService
type MockForumService struct {
	mock.Mock
}

func (m *MockForumService) CreatePost(post *domain.ForumPost) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockForumService) GetPost(id uint) (*domain.ForumPost, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ForumPost), args.Error(1)
}

func (m *MockForumService) ListPosts(page, limit int, filters map[string]interface{}) ([]domain.ForumPost, int64, error) {
	args := m.Called(page, limit, filters)
	return args.Get(0).([]domain.ForumPost), args.Get(1).(int64), args.Error(2)
}

func (m *MockForumService) UpdatePost(post *domain.ForumPost) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockForumService) DeletePost(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockForumService) CreateReply(reply *domain.ForumReply) error {
	args := m.Called(reply)
	return args.Error(0)
}

func (m *MockForumService) ListReplies(postID uint, page, limit int) ([]domain.ForumReply, int64, error) {
	args := m.Called(postID, page, limit)
	return args.Get(0).([]domain.ForumReply), args.Get(1).(int64), args.Error(2)
}

func TestForumHandler_ListPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/forum/posts?page=1&limit=10", nil)

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		posts := []domain.ForumPost{
			{ID: 1, Title: "Test Post 1", Content: "Content 1", AuthorID: 1},
			{ID: 2, Title: "Test Post 2", Content: "Content 2", AuthorID: 2},
		}
		mockService.On("ListPosts", 1, 10, mock.AnythingOfType("map[string]interface {}")).Return(posts, int64(2), nil)

		handler.ListPosts(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ForumPostListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Posts))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/forum/posts", nil)

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		mockService.On("ListPosts", 1, 15, mock.AnythingOfType("map[string]interface {}")).Return([]domain.ForumPost{}, int64(0), errors.New("service error"))

		handler.ListPosts(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestForumHandler_GetPost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		post := &domain.ForumPost{ID: 1, Title: "Test Post", Content: "Test Content", AuthorID: 1}
		mockService.On("GetPost", uint(1)).Return(post, nil)

		handler.GetPost(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "post")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		handler.GetPost(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Post Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		mockService.On("GetPost", uint(999)).Return(nil, errors.New("not found"))

		handler.GetPost(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestForumHandler_CreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		req := dto.ForumPostCreateRequest{
			Title:    "New Post",
			Content:  "New Content",
			Category: "General",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/forum/posts", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("CreatePost", mock.AnythingOfType("*domain.ForumPost")).Return(nil)

		handler.CreatePost(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		c.Request, _ = http.NewRequest(http.MethodPost, "/forum/posts", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreatePost(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestForumHandler_CreateReply(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		req := dto.ForumReplyCreateRequest{
			Content: "Reply content",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/forum/posts/1/replies", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		post := &domain.ForumPost{ID: 1, IsLocked: false}
		mockService.On("GetPost", uint(1)).Return(post, nil)
		mockService.On("CreateReply", mock.AnythingOfType("*domain.ForumReply")).Return(nil)

		handler.CreateReply(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Post ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		handler.CreateReply(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestForumHandler_ListReplies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/forum/posts/1/replies?page=1&limit=10", nil)

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		replies := []domain.ForumReply{
			{ID: 1, Content: "Reply 1", PostID: 1, AuthorID: 1},
			{ID: 2, Content: "Reply 2", PostID: 1, AuthorID: 2},
		}
		mockService.On("ListReplies", uint(1), 1, 10).Return(replies, int64(2), nil)

		handler.ListReplies(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ForumReplyListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Replies))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Post ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockForumService)
		handler := NewForumHandler(mockService)

		handler.ListReplies(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}