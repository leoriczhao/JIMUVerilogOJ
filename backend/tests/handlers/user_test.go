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
	"golang.org/x/crypto/bcrypt"
)

// MockUserService is a mock for handlers.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(username, email, password, role string) (*domain.User, error) {
	args := m.Called(username, email, password, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockService := new(MockUserService)
		handler := handlers.NewUserHandler(mockService)

		req := dto.UserRegisterRequest{Username: "newuser", Email: "new@test.com", Password: "password"}
		// The user object that CreateUser is expected to return
		createdUser := &domain.User{ID: 1, Username: req.Username, Email: req.Email}

		// Setup mock expectations
		mockService.On("GetUserByUsername", req.Username).Return(nil, nil)
		mockService.On("CreateUser", req.Username, req.Email, req.Password, "").Return(createdUser, nil)
		mockService.On("UpdateUser", createdUser).Return(nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Execute
		handler.Register(c)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Username Exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockService := new(MockUserService)
		handler := handlers.NewUserHandler(mockService)

		req := dto.UserRegisterRequest{Username: "existinguser", Email: "new@test.com", Password: "password"}
		existingUser := &domain.User{ID: 1, Username: req.Username}

		mockService.On("GetUserByUsername", req.Username).Return(existingUser, nil)

		reqBody, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1)) // Mock user ID from auth middleware

	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	expectedUser := &domain.User{ID: 1, Username: "testuser", Email: "test@test.com"}
	mockService.On("GetUserByID", uint(1)).Return(expectedUser, nil)

	handler.GetProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.UserProfileResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedUser.Username, resp.User.Username)
	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1))

	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	req := dto.UserUpdateProfileRequest{Nickname: "Updated Nickname"}
	userToUpdate := &domain.User{ID: 1, Username: "testuser"}

	mockService.On("GetUserByID", uint(1)).Return(userToUpdate, nil)
	// We check if the user object passed to UpdateUser has the new nickname
	mockService.On("UpdateUser", mock.MatchedBy(func(u *domain.User) bool {
		return u.Nickname == req.Nickname
	})).Return(nil)

	reqBody, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateProfile(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	req := dto.UserLoginRequest{Username: "testuser", Password: password}
	user := &domain.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockService.On("GetUserByUsername", "testuser").Return(user, nil)

	reqBody, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.UserLoginResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, user.Username, resp.User.Username)
	mockService.AssertExpectations(t)
}