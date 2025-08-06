package dto

import "time"

// NewsCreateRequest 创建新闻请求
type NewsCreateRequest struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Summary     string   `json:"summary"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"is_published"`
}

// NewsUpdateRequest 更新新闻请求
type NewsUpdateRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Summary     string   `json:"summary"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	IsPublished *bool    `json:"is_published"`
}

// NewsResponse 新闻响应
type NewsResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	IsPublished bool      `json:"is_published"`
	AuthorID    uint      `json:"author_id"`
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewsListResponse 新闻列表响应
type NewsListResponse struct {
	News  []NewsResponse `json:"news"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// NewsCreateResponse 创建新闻响应
type NewsCreateResponse struct {
	Message string       `json:"message"`
	News    NewsResponse `json:"news"`
}

// NewsUpdateResponse 更新新闻响应
type NewsUpdateResponse struct {
	Message string       `json:"message"`
	News    NewsResponse `json:"news"`
}

// NewsDeleteResponse 删除新闻响应
type NewsDeleteResponse struct {
	Message string `json:"message"`
}