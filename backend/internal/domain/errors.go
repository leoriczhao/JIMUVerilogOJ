package domain

import "errors"

// 用户相关错误
var (
	ErrInvalidRole     = errors.New("invalid user role")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserInactive    = errors.New("user is inactive")
)

// 题目相关错误
var (
	ErrProblemNotFound    = errors.New("problem not found")
	ErrInvalidDifficulty  = errors.New("invalid problem difficulty")
	ErrProblemNotPublic   = errors.New("problem is not public")
	ErrUnauthorizedAccess = errors.New("unauthorized access to problem")
)

// 提交相关错误
var (
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrInvalidLanguage    = errors.New("invalid programming language")
	ErrInvalidStatus      = errors.New("invalid submission status")
	ErrCodeEmpty          = errors.New("code cannot be empty")
)

// 论坛相关错误
var (
	ErrPostNotFound     = errors.New("forum post not found")
	ErrReplyNotFound    = errors.New("forum reply not found")
	ErrPostLocked       = errors.New("forum post is locked")
	ErrInvalidCategory  = errors.New("invalid forum category")
	ErrContentEmpty     = errors.New("content cannot be empty")
)

// 新闻相关错误
var (
	ErrNewsNotFound      = errors.New("news not found")
	ErrNewsNotPublished  = errors.New("news is not published")
	ErrInvalidNewsStatus = errors.New("invalid news status")
	ErrTitleEmpty        = errors.New("title cannot be empty")
)

// 通用错误
var (
	ErrInvalidID       = errors.New("invalid ID")
	ErrInvalidInput    = errors.New("invalid input")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInternalError   = errors.New("internal server error")
)