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

	news, total, err := h.newsService.GetNewsList(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve news"})
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

	status := "draft"
	if req.IsPublished {
		status = "published"
	}

	news := &domain.News{
		Title:    req.Title,
		Content:  req.Content,
		Summary:  req.Summary,
		Category: req.Category,
		Tags:     req.Tags,
		Status:   status,
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
	if req.IsPublished != nil {
		if *req.IsPublished {
			news.Status = "published"
		} else {
			news.Status = "draft"
		}
	}
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
