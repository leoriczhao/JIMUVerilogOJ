package middleware

import (
	"fmt"
	"strings"
	"sync"
)

// RBAC 角色权限管理
type RBAC struct {
	rolePermissions map[string][]string
	userCache       map[string][]string
	cacheMutex      sync.RWMutex
}

// 全局RBAC实例
var DefaultRBAC = &RBAC{
	rolePermissions: make(map[string][]string),
	userCache:       make(map[string][]string),
}

// 初始化角色权限映射
func init() {
	// 初始化角色权限
	DefaultRBAC.initRolePermissions()
}

// initRolePermissions 初始化角色权限映射
func (rbac *RBAC) initRolePermissions() {
	rbac.rolePermissions = map[string][]string{
		"student":     StudentPermissions,
		"teacher":     TeacherPermissions,
		"admin":       AdminPermissions,
		"super_admin": SuperAdminPermissions,
	}
}

// GetRolePermissions 获取角色的所有权限
func (rbac *RBAC) GetRolePermissions(role string) []string {
	permissions, exists := rbac.rolePermissions[role]
	if !exists {
		return []string{}
	}
	return permissions
}

// GetUserPermissions 获取用户的所有权限
func (rbac *RBAC) GetUserPermissions(userID uint, role string) []string {
	// 使用用户ID和角色名组成唯一缓存键，确保不同角色不会相互污染
	cacheKey := fmt.Sprintf("%d|%s", userID, role)

	// 检查缓存
	rbac.cacheMutex.RLock()
	if cachedPerms, exists := rbac.userCache[cacheKey]; exists {
		rbac.cacheMutex.RUnlock()
		return cachedPerms
	}
	rbac.cacheMutex.RUnlock()

	// 获取角色权限
	permissions := rbac.GetRolePermissions(role)

	// 缓存结果
	rbac.cacheMutex.Lock()
	rbac.userCache[cacheKey] = permissions
	rbac.cacheMutex.Unlock()

	return permissions
}

// RoleHash 计算角色的简单哈希值
func RoleHash(role string) int {
	hash := 0
	for _, c := range role {
		hash = hash*31 + int(c)
	}
	return hash
}

// HasPermission 检查用户是否有指定权限
func (rbac *RBAC) HasPermission(userID uint, role, permission string) bool {
	userPermissions := rbac.GetUserPermissions(userID, role)

	for _, userPerm := range userPermissions {
		if rbac.matchPermission(userPerm, permission) {
			return true
		}
	}

	return false
}

// matchPermission 权限匹配（支持通配符）
func (rbac *RBAC) matchPermission(userPerm, requiredPerm string) bool {
	// 完全匹配
	if userPerm == requiredPerm {
		return true
	}

	// 全权限通配符
	if userPerm == PermAll {
		return true
	}

	// 资源级通配符匹配
	if strings.HasSuffix(userPerm, ".*") {
		prefix := strings.TrimSuffix(userPerm, ".*")
		if strings.HasPrefix(requiredPerm, prefix) {
			return true
		}
	}

	// 操作级通配符匹配
	if strings.HasSuffix(userPerm, ".*") && strings.Contains(requiredPerm, ".") {
		parts := strings.Split(requiredPerm, ".")
		if len(parts) >= 2 {
			resourceWildcard := parts[0] + ".*"
			if userPerm == resourceWildcard {
				return true
			}
		}
	}

	return false
}

// HasAnyPermission 检查用户是否有任意一个指定权限
func (rbac *RBAC) HasAnyPermission(userID uint, role string, permissions []string) bool {
	for _, permission := range permissions {
		if rbac.HasPermission(userID, role, permission) {
			return true
		}
	}
	return false
}

// HasAllPermissions 检查用户是否有所有指定权限
func (rbac *RBAC) HasAllPermissions(userID uint, role string, permissions []string) bool {
	for _, permission := range permissions {
		if !rbac.HasPermission(userID, role, permission) {
			return false
		}
	}
	return true
}

// ClearUserCache 清除特定用户的所有角色权限缓存
func (rbac *RBAC) ClearUserCache(userID uint) {
	rbac.cacheMutex.Lock()
	defer rbac.cacheMutex.Unlock()

	prefix := fmt.Sprintf("%d|", userID)
	keysToDelete := make([]string, 0)
	for key := range rbac.userCache {
		if strings.HasPrefix(key, prefix) {
			keysToDelete = append(keysToDelete, key)
		}
	}
	for _, key := range keysToDelete {
		delete(rbac.userCache, key)
	}
}

// ClearAllCache 清除所有权限缓存
func (rbac *RBAC) ClearAllCache() {
	rbac.cacheMutex.Lock()
	rbac.userCache = make(map[string][]string)
	rbac.cacheMutex.Unlock()
}

// AddRolePermission 添加角色权限（运行时）
func (rbac *RBAC) AddRolePermission(role, permission string) {
	rbac.cacheMutex.Lock()
	defer rbac.cacheMutex.Unlock()

	permissions, exists := rbac.rolePermissions[role]
	if !exists {
		permissions = []string{}
	}

	// 避免重复添加
	for _, perm := range permissions {
		if perm == permission {
			return
		}
	}

	rbac.rolePermissions[role] = append(permissions, permission)
}

// RemoveRolePermission 移除角色权限（运行时）
func (rbac *RBAC) RemoveRolePermission(role, permission string) {
	rbac.cacheMutex.Lock()
	defer rbac.cacheMutex.Unlock()

	permissions, exists := rbac.rolePermissions[role]
	if !exists {
		return
	}

	// 过滤掉要移除的权限
	newPermissions := []string{}
	for _, perm := range permissions {
		if perm != permission {
			newPermissions = append(newPermissions, perm)
		}
	}

	rbac.rolePermissions[role] = newPermissions
}

// GetRoleList 获取所有角色列表
func (rbac *RBAC) GetRoleList() []string {
	roles := make([]string, 0, len(rbac.rolePermissions))
	for role := range rbac.rolePermissions {
		roles = append(roles, role)
	}
	return roles
}

// IsRoleValid 检查角色是否有效
func (rbac *RBAC) IsRoleValid(role string) bool {
	_, exists := rbac.rolePermissions[role]
	return exists
}

// GetPermissionStats 获取权限统计信息
func (rbac *RBAC) GetPermissionStats() map[string]int {
	stats := make(map[string]int)
	for role, permissions := range rbac.rolePermissions {
		stats[role] = len(permissions)
	}
	return stats
}
