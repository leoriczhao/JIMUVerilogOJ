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

// MockProblemService is a mock for ProblemService
type MockProblemService struct {
	mock.Mock
}

func (m *MockProblemService) CreateProblem(problem *domain.Problem) error {
	args := m.Called(problem)
	return args.Error(0)
}

func (m *MockProblemService) GetProblem(id uint) (*domain.Problem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Problem), args.Error(1)
}

func (m *MockProblemService) ListProblems(page, limit int, filters map[string]interface{}) ([]domain.Problem, int64, error) {
	args := m.Called(page, limit, filters)
	return args.Get(0).([]domain.Problem), args.Get(1).(int64), args.Error(2)
}

func (m *MockProblemService) UpdateProblem(problem *domain.Problem) error {
	args := m.Called(problem)
	return args.Error(0)
}

func (m *MockProblemService) DeleteProblem(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProblemService) GetTestCases(problemID uint) ([]domain.TestCase, error) {
	args := m.Called(problemID)
	return args.Get(0).([]domain.TestCase), args.Error(1)
}

func (m *MockProblemService) AddTestCase(testCase *domain.TestCase) error {
	args := m.Called(testCase)
	return args.Error(0)
}

func TestProblemHandler_ListProblems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/problems?page=1&limit=10", nil)

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		problems := []domain.Problem{
			{ID: 1, Title: "Problem 1", Description: "Description 1"},
			{ID: 2, Title: "Problem 2", Description: "Description 2"},
		}
		mockService.On("ListProblems", 1, 10, mock.AnythingOfType("map[string]interface {}")).Return(problems, int64(2), nil)

		handler.ListProblems(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ProblemListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Problems))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/problems", nil)

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		mockService.On("ListProblems", 1, 20, mock.AnythingOfType("map[string]interface {}")).Return([]domain.Problem{}, int64(0), errors.New("service error"))

		handler.ListProblems(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProblemHandler_GetProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		problem := &domain.Problem{ID: 1, Title: "Test Problem", Description: "Test Description"}
		mockService.On("GetProblem", uint(1)).Return(problem, nil)

		handler.GetProblem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "problem")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		handler.GetProblem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Problem Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		mockService.On("GetProblem", uint(999)).Return(nil, errors.New("not found"))

		handler.GetProblem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProblemHandler_CreateProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		req := dto.ProblemCreateRequest{
			Title:       "New Problem",
			Description: "New Description",
			Difficulty:  "Easy",
			TimeLimit:   1000,
			MemoryLimit: 128,
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/problems", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("CreateProblem", mock.AnythingOfType("*domain.Problem")).Return(nil)

		handler.CreateProblem(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置user_id

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		req := dto.ProblemCreateRequest{
			Title:       "New Problem",
			Description: "New Description",
			Difficulty:  "Easy",
			TimeLimit:   1000,
			MemoryLimit: 128,
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/problems", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProblem(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		c.Request, _ = http.NewRequest(http.MethodPost, "/problems", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProblem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProblemHandler_UpdateProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		req := dto.ProblemUpdateRequest{
			Title:       "Updated Problem",
			Description: "Updated Description",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPut, "/problems/1", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// 模拟用户是问题的作者
		existingProblem := &domain.Problem{ID: 1, Title: "Old Title", Description: "Old Description", AuthorID: 1}
		mockService.On("GetProblem", uint(1)).Return(existingProblem, nil)
		mockService.On("UpdateProblem", mock.AnythingOfType("*domain.Problem")).Return(nil)

		handler.UpdateProblem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置user_id
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		req := dto.ProblemUpdateRequest{
			Title: "Updated Problem",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPut, "/problems/1", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateProblem(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Permission Denied", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(2)) // 不同的用户ID
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		req := dto.ProblemUpdateRequest{
			Title: "Updated Problem",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPut, "/problems/1", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// 模拟问题属于其他用户
		existingProblem := &domain.Problem{ID: 1, Title: "Old Title", AuthorID: 1}
		mockService.On("GetProblem", uint(1)).Return(existingProblem, nil)

		handler.UpdateProblem(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		handler.UpdateProblem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProblemHandler_DeleteProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin") // 设置为管理员角色
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		// DeleteProblem需要先GetProblem检查权限
		problem := &domain.Problem{ID: 1, AuthorID: 1}
		mockService.On("GetProblem", uint(1)).Return(problem, nil)
		mockService.On("DeleteProblem", uint(1)).Return(nil)

		handler.DeleteProblem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置user_id
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		handler.DeleteProblem(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin")
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		handler.DeleteProblem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProblemHandler_GetTestCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		// Set user context for permission check
		c.Set("user_id", uint(1))
		c.Set("role", "admin")

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		problem := &domain.Problem{
			ID:       1,
			AuthorID: 1,
			Title:    "Test Problem",
		}
		testCases := []domain.TestCase{
			{ID: 1, ProblemID: 1, Input: "input1", Output: "output1", IsSample: false},
			{ID: 2, ProblemID: 1, Input: "input2", Output: "output2", IsSample: true},
		}
		mockService.On("GetProblem", uint(1)).Return(problem, nil)
		mockService.On("GetTestCases", uint(1)).Return(testCases, nil)

		handler.GetTestCases(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.TestCaseListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.TestCases))
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		handler.GetTestCases(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Student Only Sees Sample Test Cases", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		// Set student context (not author, not privileged)
		c.Set("user_id", uint(2))
		c.Set("role", "student")

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		problem := &domain.Problem{
			ID:       1,
			AuthorID: 1, // Different from user_id
			Title:    "Test Problem",
		}
		testCases := []domain.TestCase{
			{ID: 1, ProblemID: 1, Input: "input1", Output: "output1", IsSample: false},
			{ID: 2, ProblemID: 1, Input: "input2", Output: "output2", IsSample: true},
			{ID: 3, ProblemID: 1, Input: "input3", Output: "output3", IsSample: false},
		}
		mockService.On("GetProblem", uint(1)).Return(problem, nil)
		mockService.On("GetTestCases", uint(1)).Return(testCases, nil)

		handler.GetTestCases(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.TestCaseListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Student should only see 1 sample test case
		assert.Equal(t, 1, len(response.TestCases))
		assert.True(t, response.TestCases[0].IsSample)
		mockService.AssertExpectations(t)
	})

	t.Run("Author Sees All Test Cases", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		// Set author context
		c.Set("user_id", uint(1))
		c.Set("role", "student")

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		problem := &domain.Problem{
			ID:       1,
			AuthorID: 1, // Same as user_id
			Title:    "Test Problem",
		}
		testCases := []domain.TestCase{
			{ID: 1, ProblemID: 1, Input: "input1", Output: "output1", IsSample: false},
			{ID: 2, ProblemID: 1, Input: "input2", Output: "output2", IsSample: true},
		}
		mockService.On("GetProblem", uint(1)).Return(problem, nil)
		mockService.On("GetTestCases", uint(1)).Return(testCases, nil)

		handler.GetTestCases(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.TestCaseListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Author should see all test cases
		assert.Equal(t, 2, len(response.TestCases))
		mockService.AssertExpectations(t)
	})

	t.Run("Non-Author Teacher Only Sees Sample Test Cases", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		// Set teacher context (but not the author)
		c.Set("user_id", uint(2))
		c.Set("role", "teacher")

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		problem := &domain.Problem{
			ID:       1,
			AuthorID: 1, // Different from user_id
			Title:    "Test Problem",
		}
		testCases := []domain.TestCase{
			{ID: 1, ProblemID: 1, Input: "input1", Output: "output1", IsSample: false},
			{ID: 2, ProblemID: 1, Input: "input2", Output: "output2", IsSample: true},
		}
		mockService.On("GetProblem", uint(1)).Return(problem, nil)
		mockService.On("GetTestCases", uint(1)).Return(testCases, nil)

		handler.GetTestCases(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.TestCaseListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Non-author teacher should only see sample test cases (like regular users)
		assert.Equal(t, 1, len(response.TestCases))
		assert.True(t, response.TestCases[0].IsSample)
		mockService.AssertExpectations(t)
	})
}

func TestProblemHandler_AddTestCase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin") // 设置为管理员
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		req := dto.TestCaseAddRequest{
			Input:  "test input",
			Output: "test output",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/problems/1/testcases", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// AddTestCase需要先GetProblem检查权限
		problem := &domain.Problem{ID: 1, AuthorID: 1}
		mockService.On("GetProblem", uint(1)).Return(problem, nil)
		mockService.On("AddTestCase", mock.AnythingOfType("*domain.TestCase")).Return(nil)

		handler.AddTestCase(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置user_id
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		req := dto.TestCaseAddRequest{
			Input:  "sample input",
			Output: "sample output",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/problems/1/testcases", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		handler.AddTestCase(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin")
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		handler.AddTestCase(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin")
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockProblemService)
		handler := NewProblemHandler(mockService)

		c.Request, _ = http.NewRequest(http.MethodPost, "/problems/1/testcases", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddTestCase(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
