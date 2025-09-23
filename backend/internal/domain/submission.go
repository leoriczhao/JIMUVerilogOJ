package domain

import "time"

// Submission 提交记录领域实体
type Submission struct {
	ID        uint
	ProblemID uint
	UserID    uint

	// 代码信息
	Code     string
	Language string

	// 判题结果
	Status       string
	Score        int
	ErrorMessage string

	// 测试统计
	TotalTests  int
	PassedTests int

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}
