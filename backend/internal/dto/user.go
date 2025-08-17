package dto

import "time"

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Nickname  string `json:"nickname"`
	School    string `json:"school"`
	StudentID string `json:"student_id"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateProfileRequest 更新用户信息请求
type UserUpdateProfileRequest struct {
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	School    string `json:"school"`
	StudentID string `json:"student_id"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	School    string    `json:"school"`
	StudentID string    `json:"student_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresIn int64        `json:"expires_in"`
}

// UserRegisterResponse 用户注册响应
type UserRegisterResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}

// UserProfileResponse 用户信息响应
type UserProfileResponse struct {
	User UserResponse `json:"user"`
}

// UserUpdateResponse 用户更新响应
type UserUpdateResponse struct {
	Message string       `json:"message"`
	User    UserResponse `json:"user"`
}
