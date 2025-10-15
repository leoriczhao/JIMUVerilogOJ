package middleware

import (
	"testing"
)

// 测试权限匹配逻辑
func TestRBAC_MatchPermission(t *testing.T) {
	rbac := &RBAC{}

	tests := []struct {
		name     string
		userPerm string
		required string
		expected bool
	}{
		{
			name:     "完全匹配",
			userPerm: "user.profile.read",
			required: "user.profile.read",
			expected: true,
		},
		{
			name:     "全权限通配符",
			userPerm: "*",
			required: "any.permission",
			expected: true,
		},
		{
			name:     "资源级通配符匹配",
			userPerm: "user.*",
			required: "user.profile.read",
			expected: true,
		},
		{
			name:     "资源级通配符不匹配",
			userPerm: "problem.*",
			required: "user.profile.read",
			expected: false,
		},
		{
			name:     "无权限",
			userPerm: "forum.post.create",
			required: "user.profile.read",
			expected: false,
		},
		{
			name:     "操作级通配符匹配",
			userPerm: "profile.*",
			required: "profile.read",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rbac.matchPermission(tt.userPerm, tt.required)
			if result != tt.expected {
				t.Errorf("matchPermission(%s, %s) = %v, want %v",
					tt.userPerm, tt.required, result, tt.expected)
			}
		})
	}
}

// 测试角色权限获取
func TestRBAC_GetRolePermissions(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected []string
	}{
		{
			name:     "学生权限",
			role:     "student",
			expected: StudentPermissions,
		},
		{
			name:     "教师权限",
			role:     "teacher",
			expected: TeacherPermissions,
		},
		{
			name:     "管理员权限",
			role:     "admin",
			expected: AdminPermissions,
		},
		{
			name:     "超级管理员权限",
			role:     "super_admin",
			expected: SuperAdminPermissions,
		},
		{
			name:     "无效角色",
			role:     "invalid_role",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultRBAC.GetRolePermissions(tt.role)

			if len(result) != len(tt.expected) {
				t.Errorf("GetRolePermissions(%s) length = %d, want %d",
					tt.role, len(result), len(tt.expected))
				return
			}

			// 创建权限集合进行比较
			expectedSet := make(map[string]bool)
			for _, perm := range tt.expected {
				expectedSet[perm] = true
			}

			for _, perm := range result {
				if !expectedSet[perm] {
					t.Errorf("GetRolePermissions(%s) contains unexpected permission: %s",
						tt.role, perm)
				}
			}
		})
	}
}

// 测试用户权限检查
func TestRBAC_HasPermission(t *testing.T) {
	// 清除所有缓存以避免测试间的干扰
	DefaultRBAC.ClearAllCache()

	tests := []struct {
		name       string
		userID     uint
		role       string
		permission string
		expected   bool
	}{
		{
			name:       "学生基础权限",
			userID:     1,
			role:       "student",
			permission: PermUserProfileRead,
			expected:   true,
		},
		{
			name:       "学生管理员权限",
			userID:     2,
			role:       "student",
			permission: PermManageUsers,
			expected:   false,
		},
		{
			name:       "教师题目权限",
			userID:     3,
			role:       "teacher",
			permission: PermProblemCreate,
			expected:   true,
		},
		{
			name:       "教师用户管理权限",
			userID:     4,
			role:       "teacher",
			permission: PermManageUsers,
			expected:   false,
		},
		{
			name:       "管理员所有权限",
			userID:     5,
			role:       "admin",
			permission: PermManageUsers,
			expected:   true,
		},
		{
			name:       "超级管理员全权限",
			userID:     6,
			role:       "super_admin",
			permission: "any.permission",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultRBAC.HasPermission(tt.userID, tt.role, tt.permission)
			if result != tt.expected {
				t.Errorf("HasPermission(%d, %s, %s) = %v, want %v",
					tt.userID, tt.role, tt.permission, result, tt.expected)
			}
		})
	}
}

// 测试任意权限检查
func TestRBAC_HasAnyPermission(t *testing.T) {
	// 清除所有缓存以避免测试间的干扰
	DefaultRBAC.ClearAllCache()

	tests := []struct {
		name        string
		userID      uint
		role        string
		permissions []string
		expected    bool
	}{
		{
			name:        "学生有部分权限",
			userID:      1,
			role:        "student",
			permissions: []string{PermManageUsers, PermUserProfileRead},
			expected:    true,
		},
		{
			name:        "学生无权限",
			userID:      2,
			role:        "student",
			permissions: []string{PermManageUsers, PermManageSystem},
			expected:    false,
		},
		{
			name:        "教师有权限",
			userID:      3,
			role:        "teacher",
			permissions: []string{PermProblemCreate, PermManageUsers},
			expected:    true,
		},
		{
			name:        "管理员有权限",
			userID:      4,
			role:        "admin",
			permissions: []string{PermManageUsers, PermManageSystem},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultRBAC.HasAnyPermission(tt.userID, tt.role, tt.permissions)
			if result != tt.expected {
				t.Errorf("HasAnyPermission(%d, %s, %v) = %v, want %v",
					tt.userID, tt.role, tt.permissions, result, tt.expected)
			}
		})
	}
}

// 测试所有权限检查
func TestRBAC_HasAllPermissions(t *testing.T) {
	// 清除所有缓存以避免测试间的干扰
	DefaultRBAC.ClearAllCache()

	tests := []struct {
		name        string
		userID      uint
		role        string
		permissions []string
		expected    bool
	}{
		{
			name:        "学生有所有基础权限",
			userID:      1,
			role:        "student",
			permissions: []string{PermUserProfileRead, PermUserPasswordChange},
			expected:    true,
		},
		{
			name:        "学生没有管理员权限",
			userID:      2,
			role:        "student",
			permissions: []string{PermUserProfileRead, PermManageUsers},
			expected:    false,
		},
		{
			name:        "管理员有所有权限",
			userID:      3,
			role:        "admin",
			permissions: []string{PermManageUsers, PermManageSystem, PermProblemCreate},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultRBAC.HasAllPermissions(tt.userID, tt.role, tt.permissions)
			if result != tt.expected {
				t.Errorf("HasAllPermissions(%d, %s, %v) = %v, want %v",
					tt.userID, tt.role, tt.permissions, result, tt.expected)
			}
		})
	}
}

// 测试角色权限管理
func TestRBAC_RolePermissionManagement(t *testing.T) {
	// 创建临时RBAC实例
	rbac := &RBAC{
		rolePermissions: make(map[string][]string),
		userCache:       make(map[string][]string),
	}

	// 初始化基础角色
	rbac.initRolePermissions()

	testRole := "test_role"
	testPermission := "test.permission"

	// 测试添加权限
	rbac.AddRolePermission(testRole, testPermission)
	permissions := rbac.GetRolePermissions(testRole)

	found := false
	for _, perm := range permissions {
		if perm == testPermission {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("AddRolePermission did not add permission %s to role %s", testPermission, testRole)
	}

	// 测试重复添加
	rbac.AddRolePermission(testRole, testPermission)
	permissions = rbac.GetRolePermissions(testRole)
	if len(permissions) != 1 {
		t.Errorf("AddRolePermission should not add duplicate permissions, got %d", len(permissions))
	}

	// 测试移除权限
	rbac.RemoveRolePermission(testRole, testPermission)
	permissions = rbac.GetRolePermissions(testRole)

	found = false
	for _, perm := range permissions {
		if perm == testPermission {
			found = true
			break
		}
	}

	if found {
		t.Errorf("RemoveRolePermission did not remove permission %s from role %s", testPermission, testRole)
	}
}

// 测试权限缓存
func TestRBAC_PermissionCache(t *testing.T) {
	userID := uint(1)
	role := "student"

	// 清除所有缓存
	DefaultRBAC.ClearAllCache()

	// 第一次获取权限（应该从角色权限计算）
	permissions1 := DefaultRBAC.GetUserPermissions(userID, role)

	// 验证缓存中有数据
	cacheCount := 0
	DefaultRBAC.cacheMutex.RLock()
	for range DefaultRBAC.userCache {
		cacheCount++
	}
	DefaultRBAC.cacheMutex.RUnlock()

	if cacheCount == 0 {
		t.Error("Expected cache to contain data after permission lookup")
	}

	// 第二次获取权限（应该从缓存获取）
	permissions2 := DefaultRBAC.GetUserPermissions(userID, role)

	// 验证结果一致
	if len(permissions1) != len(permissions2) {
		t.Errorf("Cached permissions do not match original permissions")
	}

	// 清除所有缓存验证功能
	DefaultRBAC.ClearAllCache()

	// 验证缓存被完全清除
	afterClearCount := 0
	DefaultRBAC.cacheMutex.RLock()
	for range DefaultRBAC.userCache {
		afterClearCount++
	}
	DefaultRBAC.cacheMutex.RUnlock()

	if afterClearCount != 0 {
		t.Errorf("Expected all cache to be cleared, but found %d items", afterClearCount)
	}
}

// 验证不同角色的缓存隔离（防回归）
func TestRBAC_CacheIsolationBetweenRoles(t *testing.T) {
	userID := uint(42)

	DefaultRBAC.ClearAllCache()

	// 先以 admin 身份缓存权限
	if !DefaultRBAC.HasPermission(userID, "admin", PermManageUsers) {
		t.Fatalf("admin should have %s", PermManageUsers)
	}

	// 同一用户以 student 身份不应拥有管理员权限
	if DefaultRBAC.HasPermission(userID, "student", PermManageUsers) {
		t.Fatalf("student role should not inherit cached admin permissions")
	}

	// 反向：缓存 student 后，admin 仍应有权限
	DefaultRBAC.ClearUserCache(userID)
	if DefaultRBAC.HasPermission(userID, "student", PermUserProfileRead) == false {
		t.Fatalf("student should have %s", PermUserProfileRead)
	}
	if !DefaultRBAC.HasPermission(userID, "admin", PermManageUsers) {
		t.Fatalf("admin should still have %s after student cache", PermManageUsers)
	}
}

// 测试角色变更时缓存隔离
func TestRBAC_RoleCacheIsolation(t *testing.T) {
	userID := uint(1)
	studentRole := "student"
	teacherRole := "teacher"

	// 清除所有缓存
	DefaultRBAC.ClearAllCache()

	// 获取学生权限
	studentPerms := DefaultRBAC.GetUserPermissions(userID, studentRole)

	// 获取教师权限（不同角色）
	teacherPerms := DefaultRBAC.GetUserPermissions(userID, teacherRole)

	// 验证不同角色的权限不同
	if len(studentPerms) == len(teacherPerms) {
		// 学生和教师权限应该不同
		isDifferent := false
		for _, sp := range studentPerms {
			found := false
			for _, tp := range teacherPerms {
				if sp == tp {
					found = true
					break
				}
			}
			if !found {
				isDifferent = true
				break
			}
		}
		if !isDifferent {
			t.Error("Student and teacher permissions should be different")
		}
	}

	// 验证缓存中有2个条目（不同角色）
	cacheCount := 0
	DefaultRBAC.cacheMutex.RLock()
	for _ = range DefaultRBAC.userCache {
		cacheCount++
	}
	DefaultRBAC.cacheMutex.RUnlock()

	if cacheCount != 2 {
		t.Errorf("Expected 2 cache entries for different roles, got %d", cacheCount)
	}

	// 清除特定用户的缓存
	DefaultRBAC.ClearUserCache(userID)

	// 验证缓存被清除
	afterClearCount := 0
	DefaultRBAC.cacheMutex.RLock()
	for _ = range DefaultRBAC.userCache {
		afterClearCount++
	}
	DefaultRBAC.cacheMutex.RUnlock()

	if afterClearCount != 0 {
		t.Errorf("Expected user cache to be cleared, but found %d items", afterClearCount)
	}
}

// 基准测试
func BenchmarkRBAC_HasPermission(b *testing.B) {
	userID := uint(1)
	role := "student"
	permission := PermUserProfileRead

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DefaultRBAC.HasPermission(userID, role, permission)
	}
}

func BenchmarkRBAC_MatchPermission(b *testing.B) {
	rbac := &RBAC{}
	userPerm := "user.*"
	required := "user.profile.read"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbac.matchPermission(userPerm, required)
	}
}
