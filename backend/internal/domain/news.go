package domain

import "time"

// News 新闻领域实体
type News struct {
	ID       uint
	Title    string
	Content  string
	Summary  string
	AuthorID uint
	
	// 状态和分类
	Status   string
	Category string
	Tags     []string
	
	// 特殊标志
	IsFeatured bool
	
	// 统计信息
	ViewCount int
	
	// 时间戳
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
}