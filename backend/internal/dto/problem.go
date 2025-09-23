package dto

import "time"

// ProblemCreateRequest 创建题目请求
type ProblemCreateRequest struct {
	Title       string            `json:"title" binding:"required"`
	Description string            `json:"description" binding:"required"`
	InputDesc   string            `json:"input_desc"`
	OutputDesc  string            `json:"output_desc"`
	Difficulty  string            `json:"difficulty" binding:"required,oneof=Easy Medium Hard"`
	Category    string            `json:"category"`
	Tags        []string          `json:"tags"`
	TimeLimit   int               `json:"time_limit" binding:"min=100,max=30000"`
	MemoryLimit int               `json:"memory_limit" binding:"min=16,max=1024"`
	TestCases   []TestCaseRequest `json:"test_cases"`
}

// ProblemUpdateRequest 更新题目请求
type ProblemUpdateRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	InputDesc   string   `json:"input_desc"`
	OutputDesc  string   `json:"output_desc"`
	Difficulty  string   `json:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	TimeLimit   int      `json:"time_limit" binding:"omitempty,min=100,max=30000"`
	MemoryLimit int      `json:"memory_limit" binding:"omitempty,min=16,max=1024"`
	IsPublic    *bool    `json:"is_public"`
}

// TestCaseRequest 测试用例请求
type TestCaseRequest struct {
	Input    string `json:"input" binding:"required"`
	Output   string `json:"output" binding:"required"`
	IsSample bool   `json:"is_sample"`
}

// TestCaseAddRequest 添加测试用例请求
type TestCaseAddRequest struct {
	Input    string `json:"input" binding:"required"`
	Output   string `json:"output" binding:"required"`
	IsSample bool   `json:"is_sample"`
}

// ProblemResponse 题目响应
type ProblemResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	InputDesc   string             `json:"input_desc"`
	OutputDesc  string             `json:"output_desc"`
	Difficulty  string             `json:"difficulty"`
	Category    string             `json:"category"`
	Tags        []string           `json:"tags"`
	TimeLimit   int                `json:"time_limit"`
	MemoryLimit int                `json:"memory_limit"`
	IsPublic    bool               `json:"is_public"`
	AuthorID    uint               `json:"author_id"`
	SubmitCount int                `json:"submit_count"`
	AcceptCount int                `json:"accept_count"`
	TestCases   []TestCaseResponse `json:"test_cases,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ProblemListResponse 题目列表响应
type ProblemListResponse struct {
	Problems []ProblemResponse `json:"problems"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// ProblemCreateResponse 创建题目响应
type ProblemCreateResponse struct {
	Message string          `json:"message"`
	Problem ProblemResponse `json:"problem"`
}

// ProblemUpdateResponse 更新题目响应
type ProblemUpdateResponse struct {
	Message string          `json:"message"`
	Problem ProblemResponse `json:"problem"`
}

// ProblemDeleteResponse 删除题目响应
type ProblemDeleteResponse struct {
	Message string `json:"message"`
}

// TestCaseResponse 测试用例响应
type TestCaseResponse struct {
	ID        uint      `json:"id"`
	ProblemID uint      `json:"problem_id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	IsSample  bool      `json:"is_sample"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TestCaseListResponse 测试用例列表响应
type TestCaseListResponse struct {
	TestCases []TestCaseResponse `json:"test_cases"`
}

// TestCaseAddResponse 添加测试用例响应
type TestCaseAddResponse struct {
	Message  string           `json:"message"`
	TestCase TestCaseResponse `json:"test_case"`
}
