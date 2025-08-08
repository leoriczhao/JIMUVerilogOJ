package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminHandler 处理管理端接口
type AdminHandler struct{}

func NewAdminHandler() *AdminHandler { return &AdminHandler{} }

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
	c.JSON(http.StatusOK, gin.H{
		"users":         0,
		"problems":      0,
		"submissions":   0,
		"judges_online": 1,
	})
}

