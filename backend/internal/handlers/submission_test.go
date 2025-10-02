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
	"verilog-oj/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSubmissionService is a mock for SubmissionService
type MockSubmissionService struct {
	mock.Mock
}

func (m *MockSubmissionService) CreateSubmission(problemID uint, code, language string, userID uint) (*domain.Submission, error) {
	args := m.Called(problemID, code, language, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Submission), args.Error(1)
}

func (m *MockSubmissionService) GetSubmission(id uint) (*domain.Submission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Submission), args.Error(1)
}

func (m *MockSubmissionService) ListSubmissions(page, limit int, userID, problemID uint, status string) (*services.SubmissionListResult, error) {
	args := m.Called(page, limit, userID, problemID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SubmissionListResult), args.Error(1)
}

func (m *MockSubmissionService) UpdateSubmissionStatus(id uint, status string, score int, runTime, memory int, errorMessage string, passedTests, totalTests int) error {
	args := m.Called(id, status, score, runTime, memory, errorMessage, passedTests, totalTests)
	return args.Error(0)
}

func (m *MockSubmissionService) GetUserSubmissions(userID uint, page, limit int) (*services.SubmissionListResult, error) {
	args := m.Called(userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SubmissionListResult), args.Error(1)
}

func (m *MockSubmissionService) GetProblemSubmissions(problemID uint, page, limit int) (*services.SubmissionListResult, error) {
	args := m.Called(problemID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SubmissionListResult), args.Error(1)
}

func (m *MockSubmissionService) GetSubmissionStats(userID uint) (map[string]interface{}, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockSubmissionService) DeleteSubmission(id uint, userID uint, userRole string) error {
	args := m.Called(id, userID, userRole)
	return args.Error(0)
}

func TestSubmissionHandler_ListSubmissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/submissions?page=1&limit=10", nil)

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		result := &services.SubmissionListResult{
			Submissions: []domain.Submission{
				{ID: 1, UserID: 1, ProblemID: 1, Code: "code1", Status: "Accepted"},
				{ID: 2, UserID: 2, ProblemID: 2, Code: "code2", Status: "Wrong Answer"},
			},
			Total: 2,
			Page:  1,
			Limit: 10,
		}
		mockService.On("ListSubmissions", 1, 10, uint(0), uint(0), "").Return(result, nil)

		handler.ListSubmissions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.SubmissionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Submissions))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/submissions", nil)

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		mockService.On("ListSubmissions", 1, 20, uint(0), uint(0), "").Return(nil, errors.New("service error"))

		handler.ListSubmissions(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSubmissionHandler_GetSubmission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		submission := &domain.Submission{ID: 1, UserID: 1, ProblemID: 1, Code: "test code", Status: "Accepted"}
		mockService.On("GetSubmission", uint(1)).Return(submission, nil)

		handler.GetSubmission(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "submission")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		handler.GetSubmission(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Submission Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		mockService.On("GetSubmission", uint(999)).Return(nil, errors.New("not found"))

		handler.GetSubmission(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSubmissionHandler_CreateSubmission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		req := dto.SubmissionCreateRequest{
			ProblemID: 1,
			Code:      "test code",
			Language:  "verilog",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/submissions", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		submission := &domain.Submission{ID: 1, UserID: 1, ProblemID: 1, Code: "test code", Language: "verilog"}
		mockService.On("CreateSubmission", uint(1), "test code", "verilog", uint(1)).Return(submission, nil)

		handler.CreateSubmission(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置user_id

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		req := dto.SubmissionCreateRequest{
			ProblemID: 1,
			Code:      "test code",
			Language:  "verilog",
		}
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/submissions", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateSubmission(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		c.Request, _ = http.NewRequest(http.MethodPost, "/submissions", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateSubmission(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSubmissionHandler_GetUserSubmissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Request, _ = http.NewRequest(http.MethodGet, "/users/1/submissions?page=1&limit=10", nil)

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		result := &services.SubmissionListResult{
			Submissions: []domain.Submission{
				{ID: 1, UserID: 1, ProblemID: 1, Code: "code1", Status: "Accepted"},
				{ID: 2, UserID: 1, ProblemID: 2, Code: "code2", Status: "Wrong Answer"},
			},
			Total: 2,
			Page:  1,
			Limit: 10,
		}
		mockService.On("GetUserSubmissions", uint(1), 1, 10).Return(result, nil)

		handler.GetUserSubmissions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.SubmissionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Submissions))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置user_id

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		handler.GetUserSubmissions(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestSubmissionHandler_GetProblemSubmissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/problems/1/submissions?page=1&limit=10", nil)

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		result := &services.SubmissionListResult{
			Submissions: []domain.Submission{
				{ID: 1, UserID: 1, ProblemID: 1, Code: "code1", Status: "Accepted"},
				{ID: 2, UserID: 2, ProblemID: 1, Code: "code2", Status: "Wrong Answer"},
			},
			Total: 2,
			Page:  1,
			Limit: 10,
		}
		mockService.On("GetProblemSubmissions", uint(1), 1, 10).Return(result, nil)

		handler.GetProblemSubmissions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.SubmissionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Submissions))
		assert.Equal(t, int64(2), response.Total)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Problem ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		handler.GetProblemSubmissions(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSubmissionHandler_GetSubmissionStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Params = []gin.Param{{Key: "userId", Value: "1"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		stats := map[string]interface{}{
			"total_submissions":    int64(10),
			"accepted_submissions": int64(7),
			"solved_problems":      int64(5),
			"acceptance_rate":      0.7,
		}
		mockService.On("GetSubmissionStats", uint(1)).Return(stats, nil)

		handler.GetSubmissionStats(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.SubmissionStatsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 10, response.Stats.TotalSubmissions)
		assert.Equal(t, 7, response.Stats.AcceptedSubmissions)
		assert.Equal(t, 5, response.Stats.SolvedProblems)
		assert.Equal(t, 0.7, response.Stats.AcceptanceRate)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置user_id
		c.Params = []gin.Param{{Key: "userId", Value: "1"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		handler.GetSubmissionStats(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		mockService.On("GetSubmissionStats", uint(1)).Return(nil, errors.New("service error"))

		handler.GetSubmissionStats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSubmissionHandler_DeleteSubmission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin")
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		mockService.On("DeleteSubmission", uint(1), uint(1), "admin").Return(nil)

		handler.DeleteSubmission(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin")
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		handler.DeleteSubmission(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", uint(1))
		c.Set("role", "admin")
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := new(MockSubmissionService)
		handler := NewSubmissionHandler(mockService)

		mockService.On("DeleteSubmission", uint(1), uint(1), "admin").Return(errors.New("service error"))

		handler.DeleteSubmission(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
