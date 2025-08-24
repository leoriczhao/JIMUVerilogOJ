package handlers

import (
	"net/http"
	"strconv"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

// NewsService defines the interface for news-related services.
type NewsService interface {
	CreateNews(news *domain.News) error
	GetNewsByID(id uint) (*domain.News, error)
	GetNewsList(page, limit int) ([]domain.News, int64, error)
	GetNewsListWithFilters(page, limit int, filters map[string]interface{}) ([]domain.News, int64, error)
	UpdateNews(news *domain.News) error
	DeleteNews(id uint) error
}

// NewsHandler handles API requests for news.
type NewsHandler struct {
	newsService NewsService
}

// NewNewsHandler creates a new NewsHandler.
func NewNewsHandler(newsService interface{}) *NewsHandler {
	return &NewsHandler{
		newsService: newsService.(NewsService),
	}
}

// ListNews handles the request to get a paginated list of news articles.
func (h *NewsHandler) ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 限制每页数量
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}
	if page < 1 {
		page = 1
	}

	// 构建过滤条件
	filters := make(map[string]interface{})
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	// 权限控制：普通用户只能看已发布的新闻，管理员可以看所有
	userRole, _ := c.Get("role")
	if userRole != "admin" && userRole != "super_admin" {
		filters["is_published"] = true
	}

	news, total, err := h.newsService.GetNewsListWithFilters(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "获取新闻列表失败：" + err.Error(),
		})
		return
	}

	// 转换为DTO响应
	var newsResponses []dto.NewsResponse
	for _, newsItem := range news {
		newsResponses = append(newsResponses, dto.NewsDomainToResponse(&newsItem))
	}

	c.JSON(http.StatusOK, dto.NewsListResponse{
		News:  newsResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// GetNews handles the request to get a single news article by its ID.
func (h *NewsHandler) GetNews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	news, err := h.newsService.GetNewsByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"news": dto.NewsDomainToResponse(news)})
}

// CreateNews handles the request to create a new news article.
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var req dto.NewsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	news := &domain.News{
		Title:    req.Title,
		Content:  req.Content,
		Summary:  req.Summary,
		Category: req.Category,
		Tags:     req.Tags,
		AuthorID: userID.(uint),
	}

	if err := h.newsService.CreateNews(news); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
		return
	}
	c.JSON(http.StatusCreated, dto.NewsCreateResponse{
		Message: "News created successfully",
		News:    dto.NewsDomainToResponse(news),
	})
}

// UpdateNews handles the request to update a news article.
func (h *NewsHandler) UpdateNews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.NewsUpdateRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	news, err := h.newsService.GetNewsByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	if req.Title != "" {
		news.Title = req.Title
	}
	if req.Content != "" {
		news.Content = req.Content
	}
	if req.Summary != "" {
		news.Summary = req.Summary
	}
	if req.Category != "" {
		news.Category = req.Category
	}
	// IsPublished 字段在 domain 转换时处理
	if len(req.Tags) > 0 {
		news.Tags = req.Tags
	}

	if err := h.newsService.UpdateNews(news); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}
	c.JSON(http.StatusOK, dto.NewsUpdateResponse{
		Message: "News updated successfully",
		News:    dto.NewsDomainToResponse(news),
	})
}

// DeleteNews handles the request to delete a news article.
func (h *NewsHandler) DeleteNews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.newsService.DeleteNews(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
		return
	}
	c.JSON(http.StatusOK, dto.NewsDeleteResponse{
		Message: "News deleted successfully",
	})
}
