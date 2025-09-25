package handlers

import (
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
}

func NewAdminHandler(adminService AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
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
