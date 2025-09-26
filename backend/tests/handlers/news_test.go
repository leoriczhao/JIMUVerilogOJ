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
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.News), args.Get(1).(int64), args.Error(2)
}

func (m *MockNewsService) GetNewsListWithFilters(page, limit int, filters map[string]interface{}) ([]domain.News, int64, error) {
	args := m.Called(page, limit, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
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
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/news", nil)

	mockService := new(MockNewsService)
	handler := handlers.NewNewsHandler(mockService)

	expectedNews := []domain.News{{ID: 1, Title: "Test News"}}
	mockService.On("GetNewsListWithFilters", 1, 20, mock.AnythingOfType("map[string]interface {}")).Return(expectedNews, int64(1), nil)

	handler.ListNews(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.NewsListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp.News, 1)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_GetNews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	mockService := new(MockNewsService)
	handler := handlers.NewNewsHandler(mockService)

	expectedNews := &domain.News{ID: 1, Title: "Test News"}
	mockService.On("GetNewsByID", uint(1)).Return(expectedNews, nil)

	handler.GetNews(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_CreateNews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1))

	mockService := new(MockNewsService)
	handler := handlers.NewNewsHandler(mockService)

	req := dto.NewsCreateRequest{Title: "New News", Content: "Content"}
	mockService.On("CreateNews", mock.AnythingOfType("*domain.News")).Return(nil)

	reqBody, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest(http.MethodPost, "/news", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateNews(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_UpdateNews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	mockService := new(MockNewsService)
	handler := handlers.NewNewsHandler(mockService)

	mockService.On("GetNewsByID", uint(1)).Return(&domain.News{ID: 1}, nil)
	mockService.On("UpdateNews", mock.AnythingOfType("*domain.News")).Return(nil)

	req := dto.NewsUpdateRequest{Title: "Updated"}
	reqBody, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest(http.MethodPut, "/news/1", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateNews(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_DeleteNews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	mockService := new(MockNewsService)
	handler := handlers.NewNewsHandler(mockService)

	mockService.On("DeleteNews", uint(1)).Return(nil)

	c.Request, _ = http.NewRequest(http.MethodDelete, "/news/1", nil)

	handler.DeleteNews(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}