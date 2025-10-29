package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"verilog-oj/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminService is a mock for handlers.AdminService
type MockAdminService struct {
	mock.Mock
}

// GetSystemStats provides a mock function for the service
func (m *MockAdminService) GetSystemStats() (*services.SystemStats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SystemStats), args.Error(1)
}

func TestAdminHandler_WhoAmI(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Set values in the context as if they came from a middleware
	c.Set("user_id", uint(1))
	c.Set("username", "admin")
	c.Set("role", "admin")

	// Create mock services (the WhoAmI method doesn't use them, but we need to pass them)
	mockAdminService := &MockAdminService{}
	mockUserService := &MockUserService{}
	handler := NewAdminHandler(mockAdminService, mockUserService)

	// Execute
	handler.WhoAmI(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["user_id"]) // JSON numbers are float64 by default
	assert.Equal(t, "admin", response["username"])
	assert.Equal(t, "admin", response["role"])
}

func TestAdminHandler_Stats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockService := new(MockAdminService)
		mockUserService := new(MockUserService)
		expectedStats := &services.SystemStats{Users: 10, Problems: 20, Submissions: 30, JudgesOnline: 1}
		mockService.On("GetSystemStats").Return(expectedStats, nil)

		handler := NewAdminHandler(mockService, mockUserService)

		// Execute
		handler.Stats(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response services.SystemStats
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, &response)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Setup
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockService := new(MockAdminService)
		mockUserService := new(MockUserService)
		mockService.On("GetSystemStats").Return(nil, errors.New("db error"))

		handler := NewAdminHandler(mockService, mockUserService)

		// Execute
		handler.Stats(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
