package handlers

import (
	"net/http"
	"strconv"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/dto"
	"verilog-oj/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SubmissionService 接口定义
type SubmissionService interface {
	CreateSubmission(problemID uint, code, language string, userID uint) (*domain.Submission, error)
	GetSubmission(id uint) (*domain.Submission, error)
	ListSubmissions(page, limit int, userID, problemID uint, status string) (*services.SubmissionListResult, error)
	UpdateSubmissionStatus(id uint, status string, score int, runTime, memory int, errorMessage string, passedTests, totalTests int) error
	GetUserSubmissions(userID uint, page, limit int) (*services.SubmissionListResult, error)
	GetProblemSubmissions(problemID uint, page, limit int) (*services.SubmissionListResult, error)
	GetSubmissionStats(userID uint) (map[string]interface{}, error)
	DeleteSubmission(id uint, userID uint, userRole string) error
}

// SubmissionHandler 提交处理器
type SubmissionHandler struct {
	submissionService SubmissionService
}

// NewSubmissionHandler 创建提交处理器
func NewSubmissionHandler(submissionService interface{}) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService.(SubmissionService),
	}
}

// ListSubmissions 获取提交列表
func (h *SubmissionHandler) ListSubmissions(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	userIDStr := c.Query("user_id")
	problemIDStr := c.Query("problem_id")
	status := c.Query("status")

	// 解析分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// 解析过滤参数
	var userID, problemID uint
	if userIDStr != "" {
		if id, parseErr := strconv.ParseUint(userIDStr, 10, 32); parseErr == nil {
			userID = uint(id)
		}
	}
	if problemIDStr != "" {
		if id, parseErr := strconv.ParseUint(problemIDStr, 10, 32); parseErr == nil {
			problemID = uint(id)
		}
	}

	// 获取提交列表
	response, err := h.submissionService.ListSubmissions(page, limit, userID, problemID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "获取提交列表失败：" + err.Error(),
		})
		return
	}

	// 转换为DTO响应
	var submissionResponses []dto.SubmissionResponse
	for _, submission := range response.Submissions {
		submissionResponses = append(submissionResponses, dto.SubmissionDomainToResponse(&submission))
	}

	c.JSON(http.StatusOK, dto.SubmissionListResponse{
		Submissions: submissionResponses,
		Total:       response.Total,
		Page:        response.Page,
		Limit:       response.Limit,
	})
}

// GetSubmission 获取提交详情
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	// 获取提交ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的提交ID",
		})
		return
	}

	// 获取提交详情
	submission, err := h.submissionService.GetSubmission(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "submission_not_found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SubmissionDetailsResponse{
		Submission: dto.SubmissionDomainToResponse(submission),
	})
}

// CreateSubmission 创建提交
func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	// 解析请求体
	var req dto.SubmissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	// 创建提交
	submission, err := h.submissionService.CreateSubmission(req.ProblemID, req.Code, req.Language, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "creation_failed",
			"message": "创建提交失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SubmissionCreateResponse{
		Message:    "提交创建成功",
		Submission: dto.SubmissionDomainToResponse(submission),
	})
}

// GetUserSubmissions 获取用户的提交记录
func (h *SubmissionHandler) GetUserSubmissions(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// 获取用户提交记录
	response, err := h.submissionService.GetUserSubmissions(userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "获取用户提交记录失败：" + err.Error(),
		})
		return
	}

	// 转换为DTO响应
	var submissionResponses []dto.SubmissionResponse
	for _, submission := range response.Submissions {
		submissionResponses = append(submissionResponses, dto.SubmissionDomainToResponse(&submission))
	}

	c.JSON(http.StatusOK, dto.SubmissionListResponse{
		Submissions: submissionResponses,
		Total:       response.Total,
		Page:        response.Page,
		Limit:       response.Limit,
	})
}

// GetProblemSubmissions 获取题目的提交记录
func (h *SubmissionHandler) GetProblemSubmissions(c *gin.Context) {
	// 获取题目ID
	problemIDStr := c.Param("id")
	problemID, err := strconv.ParseUint(problemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的题目ID",
		})
		return
	}

	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	// 获取题目提交记录
	response, err := h.submissionService.GetProblemSubmissions(uint(problemID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "获取题目提交记录失败：" + err.Error(),
		})
		return
	}

	// 转换为DTO响应
	var submissionResponses []dto.SubmissionResponse
	for _, submission := range response.Submissions {
		submissionResponses = append(submissionResponses, dto.SubmissionDomainToResponse(&submission))
	}

	c.JSON(http.StatusOK, dto.SubmissionListResponse{
		Submissions: submissionResponses,
		Total:       response.Total,
		Page:        response.Page,
		Limit:       response.Limit,
	})
}

// GetSubmissionStats 获取提交统计信息
func (h *SubmissionHandler) GetSubmissionStats(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	// 获取统计信息
	stats, err := h.submissionService.GetSubmissionStats(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "获取统计信息失败：" + err.Error(),
		})
		return
	}

	// Safely extract stats with nil checks
	totalSubmissions := 0
	if val, ok := stats["total_submissions"]; ok && val != nil {
		totalSubmissions = int(val.(int64))
	}

	acceptedSubmissions := 0
	if val, ok := stats["accepted_submissions"]; ok && val != nil {
		acceptedSubmissions = int(val.(int64))
	}

	solvedProblems := 0
	if val, ok := stats["solved_problems"]; ok && val != nil {
		solvedProblems = int(val.(int64))
	}

	acceptanceRate := 0.0
	if val, ok := stats["acceptance_rate"]; ok && val != nil {
		acceptanceRate = val.(float64)
	}

	c.JSON(http.StatusOK, dto.SubmissionStatsResponse{
		Stats: dto.SubmissionStats{
			TotalSubmissions:    totalSubmissions,
			AcceptedSubmissions: acceptedSubmissions,
			SolvedProblems:      solvedProblems,
			AcceptanceRate:      acceptanceRate,
		},
	})
}

// DeleteSubmission 删除提交记录
func (h *SubmissionHandler) DeleteSubmission(c *gin.Context) {
	// 获取当前用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		userRole = "student" // 默认角色
	}

	// 获取提交ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的提交ID",
		})
		return
	}

	// 删除提交记录
	err = h.submissionService.DeleteSubmission(uint(id), userID.(uint), userRole.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "delete_failed",
			"message": "删除提交记录失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SubmissionDeleteResponse{
		Message: "提交记录删除成功",
	})
}
