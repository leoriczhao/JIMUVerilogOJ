package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"verilog-oj/backend/internal/services"
)

// RequirePermission 权限检查中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		userIDValue, userExists := c.Get("user_id")
		roleValue, roleExists := c.Get("role")

		if !userExists || !roleExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "需要登录",
			})
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的用户信息",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的角色信息",
			})
			return
		}

		// 检查权限
		if !DefaultRBAC.HasPermission(userID, role, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "权限不足",
				"details": gin.H{
					"required_permission": permission,
					"user_role":           role,
				},
			})
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 需要任意一个权限
func RequireAnyPermission(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, userExists := c.Get("user_id")
		roleValue, roleExists := c.Get("role")

		if !userExists || !roleExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "需要登录",
			})
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的用户信息",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的角色信息",
			})
			return
		}

		// 检查是否有任意一个权限
		if !DefaultRBAC.HasAnyPermission(userID, role, permissions) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "权限不足",
				"details": gin.H{
					"required_permissions": permissions,
					"user_role":            role,
				},
			})
			return
		}

		c.Next()
	}
}

// RequireAllPermissions 需要所有权限
func RequireAllPermissions(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, userExists := c.Get("user_id")
		roleValue, roleExists := c.Get("role")

		if !userExists || !roleExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "需要登录",
			})
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的用户信息",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的角色信息",
			})
			return
		}

		// 检查是否有所有权限
		if !DefaultRBAC.HasAllPermissions(userID, role, permissions) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "权限不足",
				"details": gin.H{
					"required_permissions": permissions,
					"user_role":            role,
				},
			})
			return
		}

		c.Next()
	}
}

// ResourceOwnershipFunc 资源所有权检查函数类型
type ResourceOwnershipFunc func(c *gin.Context) (uint, bool)

// RequireOwnershipOrPermission 要求资源所有权或特定权限
func RequireOwnershipOrPermission(permission string, getOwnerID ResourceOwnershipFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, userExists := c.Get("user_id")
		roleValue, roleExists := c.Get("role")

		if !userExists || !roleExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "需要登录",
			})
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的用户信息",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的角色信息",
			})
			return
		}

		// 首先检查权限
		if DefaultRBAC.HasPermission(userID, role, permission) {
			c.Next()
			return
		}

		// 检查资源所有权
		if getOwnerID != nil {
			if ownerID, hasOwner := getOwnerID(c); hasOwner && ownerID == userID {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "权限不足",
			"details": gin.H{
				"required_permission": permission,
				"user_role":           role,
				"is_owner":            false,
			},
		})
	}
}

// RequireRoleOrOwnership 要求特定角色或资源所有权
func RequireRoleOrOwnership(allowedRoles []string, getOwnerID ResourceOwnershipFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, userExists := c.Get("user_id")
		roleValue, roleExists := c.Get("role")

		if !userExists || !roleExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "需要登录",
			})
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的用户信息",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "无效的角色信息",
			})
			return
		}

		// 检查角色权限
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		// 检查资源所有权
		if getOwnerID != nil {
			if ownerID, hasOwner := getOwnerID(c); hasOwner && ownerID == userID {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "权限不足",
			"details": gin.H{
				"allowed_roles": allowedRoles,
				"user_role":     role,
				"is_owner":      false,
			},
		})
	}
}

// OptionalAuthPermission 可选认证权限检查
func OptionalAuthPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, userExists := c.Get("user_id")
		roleValue, roleExists := c.Get("role")

		// 如果用户未登录，继续执行（但可能需要在业务层处理）
		if !userExists || !roleExists {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.Set("authenticated", false)
			c.Next()
			return
		}

		// 设置认证状态
		c.Set("authenticated", true)
		c.Set("has_permission", DefaultRBAC.HasPermission(userID, role, permission))

		c.Next()
	}
}

// GetResourceOwnerFromParam 从路径参数获取资源所有者ID
func GetResourceOwnerFromParam(param string) ResourceOwnershipFunc {
	return func(c *gin.Context) (uint, bool) {
		paramValue := c.Param(param)
		if paramValue == "" {
			return 0, false
		}

		ownerID, err := strconv.ParseUint(paramValue, 10, 32)
		if err != nil {
			return 0, false
		}

		return uint(ownerID), true
	}
}

// GetProblemOwner 从题目获取所有者ID（直接使用现有ProblemService）
func GetProblemOwner(problemIDParam string) ResourceOwnershipFunc {
	return func(c *gin.Context) (uint, bool) {
		problemIDStr := c.Param(problemIDParam)
		if problemIDStr == "" {
			return 0, false
		}

		problemID, err := strconv.ParseUint(problemIDStr, 10, 32)
		if err != nil {
			return 0, false
		}

		// 从gin上下文中获取ProblemService
		service, exists := c.Get("problem_service")
		if !exists {
			_ = c.Error(fmt.Errorf("problem service not available"))
			return 0, false
		}

		problemService, ok := service.(*services.ProblemService)
		if !ok {
			_ = c.Error(fmt.Errorf("invalid problem service type"))
			return 0, false
		}

		// 通过现有Service查询题目所有者
		problem, err := problemService.GetProblem(uint(problemID))
		if err != nil {
			// 记录错误但不阻止请求
			_ = c.Error(fmt.Errorf("failed to query problem: %w", err))
			return 0, false
		}

		return problem.AuthorID, true
	}
}

// GetForumPostOwner 从论坛帖子获取所有者ID（直接使用现有ForumService）
func GetForumPostOwner(postIDParam string) ResourceOwnershipFunc {
	return func(c *gin.Context) (uint, bool) {
		postIDStr := c.Param(postIDParam)
		if postIDStr == "" {
			return 0, false
		}

		postID, err := strconv.ParseUint(postIDStr, 10, 32)
		if err != nil {
			return 0, false
		}

		// 从gin上下文中获取ForumService
		service, exists := c.Get("forum_service")
		if !exists {
			_ = c.Error(fmt.Errorf("forum service not available"))
			return 0, false
		}

		forumService, ok := service.(*services.ForumService)
		if !ok {
			_ = c.Error(fmt.Errorf("invalid forum service type"))
			return 0, false
		}

		// 通过现有Service查询帖子所有者
		post, err := forumService.GetPost(uint(postID))
		if err != nil {
			// 记录错误但不阻止请求
			_ = c.Error(fmt.Errorf("failed to query forum post: %w", err))
			return 0, false
		}

		return post.AuthorID, true
	}
}
