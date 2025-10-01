package models

import (
	"time"

	"gorm.io/gorm"
)

// ForumPost 论坛帖子模型
type ForumPost struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 基本信息
	Title   string `json:"title" gorm:"not null" validate:"required"`
	Content string `json:"content" gorm:"type:text" validate:"required"`

	// 关联信息
	UserID uint `json:"user_id" gorm:"not null"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	// 分类和标签
	Category string `json:"category" gorm:"size:50"` // general, question, discussion, announcement
	Tags     string `json:"tags" gorm:"type:text"`   // JSON数组字符串

	// 统计信息
	ViewCount  int `json:"view_count" gorm:"default:0"`
	ReplyCount int `json:"reply_count" gorm:"default:0"`
	LikeCount  int `json:"like_count" gorm:"default:0"`

	// 状态
	IsSticky bool `json:"is_sticky" gorm:"default:false"` // 置顶
	IsLocked bool `json:"is_locked" gorm:"default:false"` // 锁定
	IsPublic bool `json:"is_public" gorm:"default:true"`
}

// ForumReply 论坛回复模型
type ForumReply struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 基本信息
	Content string `json:"content" gorm:"type:text" validate:"required"`

	// 关联信息
	PostID uint      `json:"post_id" gorm:"not null"`
	Post   ForumPost `json:"post" gorm:"foreignKey:PostID"`
	UserID uint      `json:"user_id" gorm:"not null"`
	User   User      `json:"user" gorm:"foreignKey:UserID"`

	// 回复关系
	ParentID *uint       `json:"parent_id"` // 回复的回复ID，为null表示直接回复帖子
	Parent   *ForumReply `json:"parent" gorm:"foreignKey:ParentID"`

	// 统计信息
	LikeCount int `json:"like_count" gorm:"default:0"`
}

// ForumLike 点赞记录模型
type ForumLike struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	UserID uint `json:"user_id" gorm:"not null"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	// 点赞目标（帖子或回复）
	PostID  *uint       `json:"post_id"`
	Post    *ForumPost  `json:"post" gorm:"foreignKey:PostID"`
	ReplyID *uint       `json:"reply_id"`
	Reply   *ForumReply `json:"reply" gorm:"foreignKey:ReplyID"`
}
