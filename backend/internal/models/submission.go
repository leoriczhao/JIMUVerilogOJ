package models

import (
	"time"

	"gorm.io/gorm"
)

// Submission 提交记录模型
type Submission struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联信息
	UserID    uint    `json:"user_id" gorm:"not null"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
	ProblemID uint    `json:"problem_id" gorm:"not null"`
	Problem   Problem `json:"problem" gorm:"foreignKey:ProblemID"`

	// 代码信息
	Code     string `json:"code" gorm:"type:text"`
	Language string `json:"language" gorm:"default:verilog"`

	// 判题结果
	Status       string `json:"status" gorm:"default:pending"` // pending, judging, accepted, wrong_answer, time_limit_exceeded, etc.
	Score        int    `json:"score" gorm:"default:0"`
	RunTime      int    `json:"run_time" gorm:"default:0"` // 毫秒
	Memory       int    `json:"memory" gorm:"default:0"`   // KB
	ErrorMessage string `json:"error_message" gorm:"type:text"`

	// 判题详情
	PassedTests int    `json:"passed_tests" gorm:"default:0"`
	TotalTests  int    `json:"total_tests" gorm:"default:0"`
	JudgeID     string `json:"judge_id"` // 用于与判题服务通信的ID
}

// JudgeStatus 判题状态常量
const (
	StatusPending             = "pending"
	StatusJudging             = "judging"
	StatusAccepted            = "accepted"
	StatusWrongAnswer         = "wrong_answer"
	StatusTimeLimitExceeded   = "time_limit_exceeded"
	StatusMemoryLimitExceeded = "memory_limit_exceeded"
	StatusRuntimeError        = "runtime_error"
	StatusCompileError        = "compile_error"
	StatusSystemError         = "system_error"
)
