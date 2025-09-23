package domain

import "time"

// User 用户领域实体
type User struct {
	ID       uint
	Username string
	Email    string
	Password string

	// 用户信息
	Nickname  string
	Avatar    string
	School    string
	StudentID string
	Role      string

	// 统计信息
	Solved    int
	Submitted int
	Rating    int

	// 状态
	IsActive    bool
	LastLoginAt time.Time

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}
