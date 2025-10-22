package handlers

import (
	"fmt"
	"net/http"
	"verilog-oj/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AdminService defines the interface for the admin service.
type AdminService interface {
	GetSystemStats() (*services.SystemStats, error)
}

// AdminHandler 处理管理端接口
type AdminHandler struct {
	adminService AdminService
	userService  UserService
}

func NewAdminHandler(adminService AdminService, userService UserService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		userService:  userService,
	}
}

// WhoAmI 返回当前管理员的基础信息
func (h *AdminHandler) WhoAmI(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
		"role":     role,
	})
}

// Stats 系统基础统计（占位）
func (h *AdminHandler) Stats(c *gin.Context) {
	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// UpdateUserRoleRequest 更新用户角色请求
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=student teacher admin"`
}

// UpdateUserRole 更新用户角色（仅管理员可用）
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	// 获取用户ID参数
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "用户ID不能为空",
		})
		return
	}

	// 解析请求体
	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	// 转换用户ID
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "无效的用户ID",
		})
		return
	}

	// 获取用户
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "user_not_found",
			"message": "用户不存在",
		})
		return
	}

	// 更新角色
	user.Role = req.Role
	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "update_failed",
			"message": "更新用户角色失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户角色更新成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
