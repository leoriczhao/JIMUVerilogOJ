package dto

import "time"

// SubmissionCreateRequest 创建提交请求
type SubmissionCreateRequest struct {
	ProblemID uint   `json:"problem_id" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Language  string `json:"language"`
}

// SubmissionResponse 提交响应
type SubmissionResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ProblemID    uint      `json:"problem_id"`
	Code         string    `json:"code"`
	Language     string    `json:"language"`
	Status       string    `json:"status"`
	Score        int       `json:"score"`
	RunTime      int       `json:"run_time"`
	Memory       int       `json:"memory"`
	ErrorMessage string    `json:"error_message"`
	PassedTests  int       `json:"passed_tests"`
	TotalTests   int       `json:"total_tests"`
	JudgeID      string    `json:"judge_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SubmissionListResponse 提交列表响应
type SubmissionListResponse struct {
	Submissions []SubmissionResponse `json:"submissions"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

// SubmissionCreateResponse 创建提交响应
type SubmissionCreateResponse struct {
	Message    string             `json:"message"`
	Submission SubmissionResponse `json:"submission"`
}

// SubmissionStats 提交统计
type SubmissionStats struct {
	TotalSubmissions    int     `json:"total_submissions"`
	AcceptedSubmissions int     `json:"accepted_submissions"`
	SolvedProblems      int     `json:"solved_problems"`
	AcceptanceRate      float64 `json:"acceptance_rate"`
}

// SubmissionStatsResponse 提交统计响应
type SubmissionStatsResponse struct {
	Stats SubmissionStats `json:"stats"`
}

// SubmissionDeleteResponse 删除提交响应
type SubmissionDeleteResponse struct {
	Message string `json:"message"`
}

// SubmissionDetailsResponse 提交详情响应
type SubmissionDetailsResponse struct {
	Submission SubmissionResponse `json:"submission"`
}
