package dto

import "time"

// ForumPostCreateRequest 创建帖子请求
type ForumPostCreateRequest struct {
	Title    string   `json:"title" binding:"required,min=5,max=100"`
	Content  string   `json:"content" binding:"required,min=10"`
	Category string   `json:"category" binding:"required"`
	Tags     []string `json:"tags"`
}

// ForumPostUpdateRequest 更新帖子请求
type ForumPostUpdateRequest struct {
	Title    string   `json:"title" binding:"omitempty,min=5,max=100"`
	Content  string   `json:"content" binding:"omitempty,min=10"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	IsLocked *bool    `json:"is_locked"`
}

// ForumReplyCreateRequest 创建回复请求
type ForumReplyCreateRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	ParentID *uint  `json:"parent_id"`
}

// ForumPostResponse 帖子响应
type ForumPostResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	UserID     uint      `json:"user_id"`
	Category   string    `json:"category"`
	Tags       []string  `json:"tags"`
	ViewCount  int       `json:"view_count"`
	ReplyCount int       `json:"reply_count"`
	IsLocked   bool      `json:"is_locked"`
	IsSticky   bool      `json:"is_sticky"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ForumPostListResponse 帖子列表响应
type ForumPostListResponse struct {
	Posts []ForumPostResponse `json:"posts"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

// ForumPostCreateResponse 创建帖子响应
type ForumPostCreateResponse struct {
	Message string            `json:"message"`
	Post    ForumPostResponse `json:"post"`
}

// ForumPostUpdateResponse 更新帖子响应
type ForumPostUpdateResponse struct {
	Message string            `json:"message"`
	Post    ForumPostResponse `json:"post"`
}

// ForumPostDeleteResponse 删除帖子响应
type ForumPostDeleteResponse struct {
	Message string `json:"message"`
}

// ForumReplyResponse 回复响应
type ForumReplyResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ForumReplyListResponse 回复列表响应
type ForumReplyListResponse struct {
	Replies []ForumReplyResponse `json:"replies"`
	Total   int64                `json:"total"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
}

// ForumReplyCreateResponse 创建回复响应
type ForumReplyCreateResponse struct {
	Message string             `json:"message"`
	Reply   ForumReplyResponse `json:"reply"`
}
