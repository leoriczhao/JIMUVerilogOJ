package domain

import "time"

// ForumPost 论坛帖子领域实体
type ForumPost struct {
	ID       uint
	Title    string
	Content  string
	AuthorID uint

	// 分类和标签
	Category string
	Tags     []string

	// 统计信息
	ViewCount  int
	ReplyCount int
	LikeCount  int

	// 状态标志
	IsSticky bool
	IsLocked bool
	IsPublic bool

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ForumReply 论坛回复领域实体
type ForumReply struct {
	ID       uint
	PostID   uint
	AuthorID uint
	Content  string

	// 回复关系
	ParentID *uint // 父回复ID，用于嵌套回复

	// 统计信息
	LikeCount int

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ForumLike 论坛点赞领域实体
type ForumLike struct {
	ID     uint
	UserID uint

	// 多态关联
	TargetType string // "post" 或 "reply"
	TargetID   uint   // 帖子ID或回复ID

	// 时间戳
	CreatedAt time.Time
}
