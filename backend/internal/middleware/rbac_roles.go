package middleware

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"
)

// RBAC 角色权限管理
type RBAC struct {
	rolePermissions map[string][]string
	userCache       map[string][]string
	userCacheKeys   map[uint][]string  // 记录用户ID对应的缓存键
	cacheMutex      sync.RWMutex
}

// 全局RBAC实例
var DefaultRBAC = &RBAC{
	rolePermissions: make(map[string][]string),
	userCache:       make(map[string][]string),
	userCacheKeys:   make(map[uint][]string),
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
	// 使用基于哈希的缓存键，避免键冲突
	cacheKey := generateCacheKey(userID, role)

	// 检查缓存
	rbac.cacheMutex.RLock()
	if cachedPerms, exists := rbac.userCache[cacheKey]; exists {
		rbac.cacheMutex.RUnlock()
		return cachedPerms
	}
	rbac.cacheMutex.RUnlock()

	// 获取角色权限
	permissions := rbac.GetRolePermissions(role)

	// 缓存结果并记录键
	rbac.cacheMutex.Lock()
	rbac.userCache[cacheKey] = permissions

	// 记录用户ID对应的缓存键，便于后续清理
	if rbac.userCacheKeys[userID] == nil {
		rbac.userCacheKeys[userID] = []string{}
	}
	rbac.userCacheKeys[userID] = append(rbac.userCacheKeys[userID], cacheKey)

	rbac.cacheMutex.Unlock()

	return permissions
}

// generateCacheKey 生成基于哈希的缓存键，避免键冲突
func generateCacheKey(userID uint, role string) string {
	data := fmt.Sprintf("%d|%s|%s", userID, role, "rbac_v1")
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:16]) // 使用前16字节，既保证唯一性又控制键长度
}

// RoleHash 计算角色的简单哈希值（保留用于向后兼容）
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

	// 获取用户的所有缓存键
	cacheKeys, exists := rbac.userCacheKeys[userID]
	if !exists {
		return // 没有缓存，无需清理
	}

	// 删除所有相关的缓存条目
	for _, cacheKey := range cacheKeys {
		delete(rbac.userCache, cacheKey)
	}

	// 清除键记录
	delete(rbac.userCacheKeys, userID)
}

// ClearAllCache 清除所有权限缓存
func (rbac *RBAC) ClearAllCache() {
	rbac.cacheMutex.Lock()
	rbac.userCache = make(map[string][]string)
	rbac.userCacheKeys = make(map[uint][]string)
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
