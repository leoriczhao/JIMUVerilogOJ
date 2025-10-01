package handlers

import (
	"net/http"
	"time"
	"verilog-oj/backend/internal/config"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/dto"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService 接口定义
type UserService interface {
	CreateUser(username, email, password, role string) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	UpdateUser(user *domain.User) error
}

// UserHandler 用户处理器
type UserHandler struct {
	userService UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// 使用DTO包中的结构体
// RegisterRequest = dto.UserRegisterRequest
// LoginRequest = dto.UserLoginRequest
// LoginResponse = dto.UserLoginResponse
// UpdateProfileRequest = dto.UserUpdateProfileRequest

// 从配置读取JWT密钥
var jwtSecret = []byte(config.LoadConfig().JWT.Secret)

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	if existingUser, _ := h.userService.GetUserByUsername(req.Username); existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "username_exists",
			"message": "用户名已存在",
		})
		return
	}

	// 创建用户
	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "registration_failed",
			"message": "注册失败：" + err.Error(),
		})
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.School != "" {
		user.School = req.School
	}
	if req.StudentID != "" {
		user.StudentID = req.StudentID
	}

	if err := h.userService.UpdateUser(user); err != nil {
		// 如果更新失败，至少返回基本用户信息
		_ = c.Error(err)
	}

	c.JSON(http.StatusCreated, dto.UserRegisterResponse{
		Message: "注册成功",
		User:    dto.UserToResponse(user),
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	// 获取用户
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_credentials",
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); compareErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_credentials",
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成JWT Token
	token, expiresIn, err := h.generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_failed",
			"message": "Token生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserLoginResponse{
		Token:     token,
		User:      dto.UserToResponse(user),
		ExpiresIn: expiresIn,
	})
}

// GetProfile 获取用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "未授权访问",
		})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "user_not_found",
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserProfileResponse{
		User: dto.UserToResponse(user),
	})
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "未授权访问",
		})
		return
	}

	var req dto.UserUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "user_not_found",
			"message": "用户不存在",
		})
		return
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.School != "" {
		user.School = req.School
	}
	if req.StudentID != "" {
		user.StudentID = req.StudentID
	}

	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "update_failed",
			"message": "更新失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserUpdateResponse{
		Message: "更新成功",
		User:    dto.UserToResponse(user),
	})
}

// generateToken 生成JWT Token
func (h *UserHandler) generateToken(user *domain.User) (string, int64, error) {
	expiresIn := time.Hour * 24 * 7 // 7天过期
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(expiresIn).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, int64(expiresIn.Seconds()), nil
}
