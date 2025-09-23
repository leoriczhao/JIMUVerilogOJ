package domain

import "time"

// Problem 题目领域实体
type Problem struct {
	ID          uint
	Title       string
	Description string
	InputDesc   string
	OutputDesc  string

	// 难度和分类
	Difficulty string
	Category   string
	Tags       []string

	// 限制条件
	TimeLimit   int // 毫秒
	MemoryLimit int // MB

	// 统计信息
	SubmitCount   int
	AcceptedCount int

	// 状态和作者
	IsPublic bool
	AuthorID uint

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TestCase 测试用例领域实体
type TestCase struct {
	ID        uint
	ProblemID uint
	Input     string
	Output    string
	IsSample  bool

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}
