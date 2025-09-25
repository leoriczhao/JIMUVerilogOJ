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

// MockNewsService is a mock for NewsService
type MockNewsService struct {
	mock.Mock
}

func (m *MockNewsService) CreateNews(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsService) GetNewsByID(id uint) (*domain.News, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.News), args.Error(1)
}

func (m *MockNewsService) GetNewsList(page, limit int) ([]domain.News, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]domain.News), args.Get(1).(int64), args.Error(2)
}

func (m *MockNewsService) GetNewsListWithFilters(page, limit int, filters map[string]interface{}) ([]domain.News, int64, error) {
	args := m.Called(page, limit, filters)
	return args.Get(0).([]domain.News), args.Get(1).(int64), args.Error(2)
}

func (m *MockNewsService) UpdateNews(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsService) DeleteNews(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNewsHandler_ListNews(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/news?page=1&limit=10", nil)

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		news := []domain.News{
			{ID: 1, Title: "News 1", Content: "Content 1", AuthorID: 1},
			{ID: 2, Title: "News 2", Content: "Content 2", AuthorID: 2},
		}
		mockService.On("GetNewsListWithFilters", 1, 10, mock.AnythingOfType("map[string]interface {}")).Return(news, int64(2), nil)

		handler.ListNews(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.NewsListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.News))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/news", nil)

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		mockService.On("GetNewsListWithFilters", 1, 20, mock.AnythingOfType("map[string]interface {}")).Return([]domain.News{}, int64(0), errors.New("service error"))

		handler.ListNews(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestNewsHandler_GetNews(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		news := &domain.News{ID: 1, Title: "Test News", Content: "Test Content", AuthorID: 1}
		mockService.On("GetNewsByID", uint(1)).Return(news, nil)

		handler.GetNews(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "news")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		handler.GetNews(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("News Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		mockService.On("GetNewsByID", uint(999)).Return(nil, errors.New("not found"))

		handler.GetNews(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestNewsHandler_CreateNews(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		req := dto.NewsCreateRequest{
			Title:   "New News",
			Content: "New Content",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/news", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateNews", mock.AnythingOfType("*domain.News")).Return(nil)

		handler.CreateNews(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		c.Request, _ = http.NewRequest(http.MethodPost, "/news", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNews(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestNewsHandler_UpdateNews(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		req := dto.NewsUpdateRequest{
			Title:   "Updated News",
			Content: "Updated Content",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPut, "/news/1", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		existingNews := &domain.News{ID: 1, Title: "Old Title", Content: "Old Content", AuthorID: 1}
		mockService.On("GetNewsByID", uint(1)).Return(existingNews, nil)
		mockService.On("UpdateNews", mock.AnythingOfType("*domain.News")).Return(nil)

		handler.UpdateNews(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		handler.UpdateNews(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestNewsHandler_DeleteNews(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		mockService.On("DeleteNews", uint(1)).Return(nil)

		handler.DeleteNews(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		handler.DeleteNews(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("News Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		mockService := new(MockNewsService)
		handler := NewNewsHandler(mockService)

		mockService.On("DeleteNews", uint(999)).Return(errors.New("not found"))

		handler.DeleteNews(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}