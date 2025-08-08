package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly 确保用户具备管理员角色
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "需要管理员权限",
			})
			return
		}

		role, _ := roleValue.(string)
		if role != "admin" && role != "super_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "需要管理员权限",
			})
			return
		}

		c.Next()
	}
}

