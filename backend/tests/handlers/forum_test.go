package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/dto"
	"verilog-oj/backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
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
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.ForumReply), args.Get(1).(int64), args.Error(2)
}

func TestForumHandler_ListPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/posts?page=1&limit=10", nil)

	mockService := new(MockForumService)
	handler := handlers.NewForumHandler(mockService)

	expectedPosts := []domain.ForumPost{{ID: 1, Title: "Test Post"}}
	mockService.On("ListPosts", 1, 10, mock.AnythingOfType("map[string]interface {}")).Return(expectedPosts, int64(1), nil)

	handler.ListPosts(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.ForumPostListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Posts, 1)
	mockService.AssertExpectations(t)
}

func TestForumHandler_GetPost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	mockService := new(MockForumService)
	handler := handlers.NewForumHandler(mockService)

	expectedPost := &domain.ForumPost{ID: 1, Title: "Test Post"}
	mockService.On("GetPost", uint(1)).Return(expectedPost, nil)

	handler.GetPost(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestForumHandler_CreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1))

	mockService := new(MockForumService)
	handler := handlers.NewForumHandler(mockService)

	req := dto.ForumPostCreateRequest{Title: "New Post", Content: "This is a valid content string.", Category: "general"}
	mockService.On("CreatePost", mock.AnythingOfType("*domain.ForumPost")).Return(nil)

	reqBody, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreatePost(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestForumHandler_UpdatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1))
	c.Set("role", "user")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	mockService := new(MockForumService)
	handler := handlers.NewForumHandler(mockService)

	// User is author, so update is allowed
	mockService.On("GetPost", uint(1)).Return(&domain.ForumPost{ID: 1, AuthorID: 1}, nil)
	mockService.On("UpdatePost", mock.AnythingOfType("*domain.ForumPost")).Return(nil)

	req := dto.ForumPostUpdateRequest{Title: "Updated"}
	reqBody, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest(http.MethodPut, "/posts/1", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdatePost(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestForumHandler_DeletePost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1))
	c.Set("role", "admin") // Admin can delete any post
	c.Params = gin.Params{gin.Param{Key: "id", Value: "2"}}

	mockService := new(MockForumService)
	handler := handlers.NewForumHandler(mockService)

	mockService.On("GetPost", uint(2)).Return(&domain.ForumPost{ID: 2, AuthorID: 99}, nil) // Different author
	mockService.On("DeletePost", uint(2)).Return(nil)

	c.Request, _ = http.NewRequest(http.MethodDelete, "/posts/2", nil)

	handler.DeletePost(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestForumHandler_CreateReply(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1))
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}} // Post ID

	mockService := new(MockForumService)
	handler := handlers.NewForumHandler(mockService)

	// Mock that the post exists and is not locked
	mockService.On("GetPost", uint(1)).Return(&domain.ForumPost{ID: 1, IsLocked: false}, nil)
	mockService.On("CreateReply", mock.AnythingOfType("*domain.ForumReply")).Return(nil)

	req := dto.ForumReplyCreateRequest{Content: "A new reply"}
	reqBody, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest(http.MethodPost, "/posts/1/replies", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateReply(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}