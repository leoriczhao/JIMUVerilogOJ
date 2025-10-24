#!/usr/bin/env python3
"""
权限管理模块
与 Go 项目的 RBAC 权限定义对齐
同步自: backend/internal/middleware/rbac_permissions.go
"""

# ========== 角色权限映射 ==========
# 直接从 Go 的 StudentPermissions, TeacherPermissions, AdminPermissions 同步

ROLE_PERMISSIONS = {
    'student': [
        # 用户资料权限
        'user.profile.read',
        'user.profile.update',
        'user.password.change',
        'user.avatar.upload',
        # 题目权限
        'problem.read',
        'problem.list',
        # 测试用例权限
        'testcase.read.sample',
        # 提交权限
        'submission.create',
        'submission.read',
        'submission.list',
        # 论坛权限
        'forum.post.create',
        'forum.post.read',
        'forum.reply.create',
        'forum.reply.read',
        # 新闻权限
        'news.read',
        'news.list',
        # 统计权限
        'stats.read',
    ],

    'teacher': [
        # 继承 student 的所有权限
        'user.profile.read',
        'user.profile.update',
        'user.password.change',
        'user.avatar.upload',
        'problem.read',
        'problem.list',
        'testcase.read.sample',
        'submission.create',
        'submission.read',
        'submission.list',
        'forum.post.create',
        'forum.post.read',
        'forum.reply.create',
        'forum.reply.read',
        'news.read',
        'news.list',
        'stats.read',
        # teacher 特有权限
        'problem.create',
        'problem.update.own',
        'problem.delete.own',
        'testcase.read.own',
        'testcase.create',
        'testcase.update',
        'testcase.delete',
        'news.create',
        'news.update',
        'news.delete',
    ],

    'admin': [
        # 继承 teacher 的所有权限
        'user.profile.read',
        'user.profile.update',
        'user.password.change',
        'user.avatar.upload',
        'problem.read',
        'problem.list',
        'testcase.read.sample',
        'submission.create',
        'submission.read',
        'submission.list',
        'forum.post.create',
        'forum.post.read',
        'forum.reply.create',
        'forum.reply.read',
        'news.read',
        'news.list',
        'stats.read',
        'problem.create',
        'problem.update.own',
        'problem.delete.own',
        'testcase.read.own',
        'testcase.create',
        'testcase.update',
        'testcase.delete',
        'news.create',
        'news.update',
        'news.delete',
        # admin 特有权限
        'user.create',
        'user.read',
        'user.update',
        'user.delete',
        'user.ban',
        'user.unban',
        'problem.update.all',
        'problem.delete.all',
        'problem.publish',
        'problem.archive',
        'testcase.read.all',
        'submission.manage',
        'submission.delete',
        'submission.rejudge',
        'forum.edit.all',
        'forum.moderate',
        'forum.post.lock',
        'forum.post.sticky',
        'forum.delete',
        'news.publish',
        'news.archive',
        'manage.users',
        'manage.system',
        'manage.config',
        'manage.content',
        'stats.admin',
    ]
}


def has_permission(role: str, permission: str) -> bool:
    """
    检查角色是否有指定权限

    Args:
        role: 角色名称 ('student', 'teacher', 'admin')
        permission: 权限字符串，如 'problem.create'

    Returns:
        True 如果角色拥有该权限，否则 False
    """
    if role not in ROLE_PERMISSIONS:
        return False

    role_perms = ROLE_PERMISSIONS[role]

    # 完全匹配
    if permission in role_perms:
        return True

    # 检查通配符权限 (如 'problem.*' 匹配 'problem.create')
    for perm in role_perms:
        if perm.endswith('.*'):
            prefix = perm[:-2]  # 移除 '.*'
            if permission.startswith(prefix + '.'):
                return True

    # 检查全权限通配符
    if '*' in role_perms:
        return True

    return False


def get_minimum_role(permission: str) -> str:
    """
    获取拥有某权限的最低角色

    Args:
        permission: 权限字符串

    Returns:
        最低角色名称，如果没有角色拥有该权限返回 None
    """
    # 按权限从低到高检查
    for role in ['student', 'teacher', 'admin']:
        if has_permission(role, permission):
            return role
    return None


def get_role_permissions(role: str) -> list:
    """
    获取角色的所有权限列表

    Args:
        role: 角色名称

    Returns:
        权限列表
    """
    return ROLE_PERMISSIONS.get(role, [])


def suggest_role_for_permissions(permissions: list) -> str:
    """
    根据权限列表推荐合适的角色

    Args:
        permissions: 权限列表

    Returns:
        推荐的角色名称
    """
    if not permissions:
        return None

    # 找到能满足所有权限的最低角色
    for role in ['student', 'teacher', 'admin']:
        if all(has_permission(role, perm) for perm in permissions):
            return role

    return 'admin'  # 默认返回 admin


# ========== 公开接口列表 ==========
# 这些接口不需要认证
PUBLIC_ENDPOINTS = [
    ('POST', '/users/register'),
    ('POST', '/users/login'),
    ('GET', '/health'),
]


def is_public_endpoint(method: str, endpoint: str) -> bool:
    """
    判断端点是否为公开接口（无需认证）

    Args:
        method: HTTP 方法
        endpoint: API 端点路径

    Returns:
        True 如果是公开接口
    """
    # 规范化端点路径（移除参数）
    normalized_endpoint = endpoint.split('?')[0]

    return (method.upper(), normalized_endpoint) in PUBLIC_ENDPOINTS
