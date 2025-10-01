package handlers

import (
	"net/http"
	"strconv"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

// ForumService 接口定义
type ForumService interface {
	CreatePost(post *domain.ForumPost) error
	GetPost(id uint) (*domain.ForumPost, error)
	ListPosts(page, limit int, filters map[string]interface{}) ([]domain.ForumPost, int64, error)
	UpdatePost(post *domain.ForumPost) error
	DeletePost(id uint) error
	CreateReply(reply *domain.ForumReply) error
	ListReplies(postID uint, page, limit int) ([]domain.ForumReply, int64, error)
}

// ForumHandler 论坛处理器
type ForumHandler struct {
	forumService ForumService
}

// NewForumHandler 创建论坛处理器
func NewForumHandler(forumService interface{}) *ForumHandler {
	return &ForumHandler{
		forumService: forumService.(ForumService),
	}
}

// 使用DTO包中的结构体
// CreatePostRequest = dto.ForumPostCreateRequest
// UpdatePostRequest = dto.ForumPostUpdateRequest
// CreateReplyRequest = dto.ForumReplyCreateRequest

// ListPosts 获取帖子列表
func (h *ForumHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))
	category := c.Query("category")

	// 构建过滤条件，只有非空值才添加过滤
	filters := make(map[string]interface{})
	if category != "" {
		filters["category"] = category
	}

	posts, total, err := h.forumService.ListPosts(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list posts"})
		return
	}

	// 转换为DTO响应
	var postResponses []dto.ForumPostResponse
	for _, post := range posts {
		postResponses = append(postResponses, dto.ForumPostDomainToResponse(&post))
	}

	c.JSON(http.StatusOK, dto.ForumPostListResponse{
		Posts: postResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// GetPost 获取帖子详情
func (h *ForumHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	post, err := h.forumService.GetPost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": dto.ForumPostDomainToResponse(post)})
}

// CreatePost 创建帖子
func (h *ForumHandler) CreatePost(c *gin.Context) {
	var req dto.ForumPostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")

	post := &domain.ForumPost{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userID.(uint),
		Category: req.Category,
		Tags:     req.Tags,
	}

	if err := h.forumService.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, dto.ForumPostCreateResponse{
		Message: "Post created successfully",
		Post:    dto.ForumPostDomainToResponse(post),
	})
}

// UpdatePost 更新帖子
func (h *ForumHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var req dto.ForumPostUpdateRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	post, err := h.forumService.GetPost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.AuthorID != userID.(uint) && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Category != "" {
		post.Category = req.Category
	}
	// 只有管理员可以修改锁定状态
	if req.IsLocked != nil && (userRole == "admin" || userRole == "super_admin") {
		post.IsLocked = *req.IsLocked
	}
	if len(req.Tags) > 0 {
		post.Tags = req.Tags
	}

	if err := h.forumService.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, dto.ForumPostUpdateResponse{
		Message: "Post updated successfully",
		Post:    dto.ForumPostDomainToResponse(post),
	})
}

// DeletePost 删除帖子
func (h *ForumHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("role")

	post, err := h.forumService.GetPost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.AuthorID != userID.(uint) && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	if err := h.forumService.DeletePost(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, dto.ForumPostDeleteResponse{
		Message: "Post deleted successfully",
	})
}

// ListReplies 获取回复列表
func (h *ForumHandler) ListReplies(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	replies, total, err := h.forumService.ListReplies(uint(postID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list replies"})
		return
	}
	// 转换为DTO响应
	var replyResponses []dto.ForumReplyResponse
	for _, reply := range replies {
		replyResponses = append(replyResponses, dto.ForumReplyDomainToResponse(&reply))
	}

	c.JSON(http.StatusOK, dto.ForumReplyListResponse{
		Replies: replyResponses,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// CreateReply 创建回复
func (h *ForumHandler) CreateReply(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var req dto.ForumReplyCreateRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}
	userID, _ := c.Get("user_id")

	// 检查帖子是否存在且未锁定
	post, err := h.forumService.GetPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if post.IsLocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Post is locked"})
		return
	}

	reply := &domain.ForumReply{
		Content:  req.Content,
		PostID:   uint(postID),
		AuthorID: userID.(uint),
		ParentID: req.ParentID,
	}

	if err := h.forumService.CreateReply(reply); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
		return
	}
	c.JSON(http.StatusCreated, dto.ForumReplyCreateResponse{
		Message: "Reply created successfully",
		Reply:   dto.ForumReplyDomainToResponse(reply),
	})
}
