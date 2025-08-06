package handlers

import (
	"net/http"
	"strconv"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/dto"

	"github.com/gin-gonic/gin"
)

// ProblemService 接口定义
type ProblemService interface {
	CreateProblem(problem *domain.Problem) error
	GetProblem(id uint) (*domain.Problem, error)
	ListProblems(page, limit int, filters map[string]interface{}) ([]domain.Problem, int64, error)
	UpdateProblem(problem *domain.Problem) error
	DeleteProblem(id uint) error
	GetTestCases(problemID uint) ([]domain.TestCase, error)
	AddTestCase(testCase *domain.TestCase) error
}

// ProblemHandler 题目处理器
type ProblemHandler struct {
	problemService ProblemService
}

// NewProblemHandler 创建题目处理器
func NewProblemHandler(problemService interface{}) *ProblemHandler {
	return &ProblemHandler{
		problemService: problemService.(ProblemService),
	}
}

// 使用DTO包中的结构体
// CreateProblemRequest = dto.ProblemCreateRequest
// UpdateProblemRequest = dto.ProblemUpdateRequest
// TestCaseRequest = dto.TestCaseCreateRequest
// AddTestCaseRequest = dto.TestCaseCreateRequest

// ListProblems 获取题目列表
func (h *ProblemHandler) ListProblems(c *gin.Context) {
	// 获取分页参数
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
	if difficulty := c.Query("difficulty"); difficulty != "" {
		filters["difficulty"] = difficulty
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	// 默认只显示公开题目
	filters["is_public"] = true

	// 获取题目列表
	problems, total, err := h.problemService.ListProblems(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "获取题目列表失败：" + err.Error(),
		})
		return
	}

	// 转换为DTO响应
	var problemResponses []dto.ProblemResponse
	for _, problem := range problems {
		problemResponses = append(problemResponses, dto.ProblemDomainToResponse(&problem))
	}

	c.JSON(http.StatusOK, dto.ProblemListResponse{
		Problems: problemResponses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

// GetProblem 获取题目详情
func (h *ProblemHandler) GetProblem(c *gin.Context) {
	// 获取题目ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的题目ID",
		})
		return
	}

	// 获取题目详情
	problem, err := h.problemService.GetProblem(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "problem_not_found",
			"message": "题目不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problem": dto.ProblemDomainToResponse(problem),
	})
}

// CreateProblem 创建题目
func (h *ProblemHandler) CreateProblem(c *gin.Context) {
	var req dto.ProblemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	// 设置默认值
	if req.TimeLimit == 0 {
		req.TimeLimit = 1000
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 128
	}

	// 创建题目
	problem := &domain.Problem{
		Title:       req.Title,
		Description: req.Description,
		InputDesc:   req.InputDesc,
		OutputDesc:  req.OutputDesc,
		Difficulty:  req.Difficulty,
		Category:    req.Category,
		Tags:        req.Tags,
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
		IsPublic:    false, // 默认私有
		AuthorID:    userID.(uint),
	}

	if err := h.problemService.CreateProblem(problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "creation_failed",
			"message": "创建题目失败：" + err.Error(),
		})
		return
	}

	// 创建测试用例
	if len(req.TestCases) > 0 {
		for _, tc := range req.TestCases {
			testCase := &domain.TestCase{
				ProblemID: problem.ID,
				Input:     tc.Input,
				Output:    tc.Output,
				IsSample:  tc.IsSample,
			}
			if err := h.problemService.AddTestCase(testCase); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "creation_failed",
					"message": "创建测试用例失败：" + err.Error(),
				})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, dto.ProblemCreateResponse{
		Message: "题目创建成功",
		Problem: dto.ProblemDomainToResponse(problem),
	})
}

// UpdateProblem 更新题目
func (h *ProblemHandler) UpdateProblem(c *gin.Context) {
	// 获取题目ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的题目ID",
		})
		return
	}

	var req dto.ProblemUpdateRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "请求参数错误：" + bindErr.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	// 获取原题目
	problem, err := h.problemService.GetProblem(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "problem_not_found",
			"message": "题目不存在",
		})
		return
	}

	// 检查权限（只有作者或管理员可以更新）
	if problem.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "permission_denied",
			"message": "没有权限更新此题目",
		})
		return
	}

	// 更新字段
	if req.Title != "" {
		problem.Title = req.Title
	}
	if req.Description != "" {
		problem.Description = req.Description
	}
	if req.InputDesc != "" {
		problem.InputDesc = req.InputDesc
	}
	if req.OutputDesc != "" {
		problem.OutputDesc = req.OutputDesc
	}
	if req.Difficulty != "" {
		problem.Difficulty = req.Difficulty
	}
	if req.Category != "" {
		problem.Category = req.Category
	}
	if req.TimeLimit > 0 {
		problem.TimeLimit = req.TimeLimit
	}
	if req.MemoryLimit > 0 {
		problem.MemoryLimit = req.MemoryLimit
	}
	if req.IsPublic != nil {
		problem.IsPublic = *req.IsPublic
	}

	// 处理标签
	if len(req.Tags) > 0 {
		problem.Tags = req.Tags
	}

	// 更新题目
	if err := h.problemService.UpdateProblem(problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "update_failed",
			"message": "更新题目失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ProblemUpdateResponse{
		Message: "题目更新成功",
		Problem: dto.ProblemDomainToResponse(problem),
	})
}

// DeleteProblem 删除题目
func (h *ProblemHandler) DeleteProblem(c *gin.Context) {
	// 获取题目ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的题目ID",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	// 获取题目检查权限
	problem, err := h.problemService.GetProblem(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "problem_not_found",
			"message": "题目不存在",
		})
		return
	}

	// 检查权限（只有作者或管理员可以删除）
	if problem.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "permission_denied",
			"message": "没有权限删除此题目",
		})
		return
	}

	// 删除题目
	if err := h.problemService.DeleteProblem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "delete_failed",
			"message": "删除题目失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ProblemDeleteResponse{
		Message: "题目删除成功",
	})
}

// GetTestCases 获取题目的测试用例列表
func (h *ProblemHandler) GetTestCases(c *gin.Context) {
	// 获取题目ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的题目ID",
		})
		return
	}

	// 获取测试用例
	testCases, err := h.problemService.GetTestCases(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "problem_not_found",
			"message": "题目不存在",
		})
		return
	}

	// 转换为DTO响应
	var testCaseResponses []dto.TestCaseResponse
	for _, tc := range testCases {
		testCaseResponses = append(testCaseResponses, dto.TestCaseDomainToResponse(&tc))
	}

	c.JSON(http.StatusOK, dto.TestCaseListResponse{
		TestCases: testCaseResponses,
	})
}

// AddTestCase 为题目添加测试用例
func (h *ProblemHandler) AddTestCase(c *gin.Context) {
	// 获取题目ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "无效的题目ID",
		})
		return
	}

	var req dto.TestCaseAddRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "请求参数错误：" + bindErr.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "用户未认证",
		})
		return
	}

	// 获取题目检查权限
	problem, err := h.problemService.GetProblem(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "problem_not_found",
			"message": "题目不存在",
		})
		return
	}

	// 检查权限（只有作者或管理员可以添加测试用例）
	if problem.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "permission_denied",
			"message": "没有权限为此题目添加测试用例",
		})
		return
	}

	// 创建测试用例
	testCase := &domain.TestCase{
		ProblemID: uint(id),
		Input:     req.Input,
		Output:    req.Output,
		IsSample:  req.IsSample,
	}

	if err := h.problemService.AddTestCase(testCase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "creation_failed",
			"message": "添加测试用例失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.TestCaseAddResponse{
		Message:  "测试用例添加成功",
		TestCase: dto.TestCaseDomainToResponse(testCase),
	})
}
