#!/usr/bin/env python3
"""
测试用户池管理模块
提供统一的测试用户创建、管理和清理功能
"""

import time
from typing import Dict, Optional


class TestUserPool:
    """
    测试用户池
    管理不同角色的测试用户，提供快速角色切换功能
    """

    def __init__(self):
        self.users = {
            'student': None,
            'teacher': None,
            'admin': None
        }
        self.tokens = {}
        self.user_ids = {}
        self.initialized = False

    def setup(self, api_tester):
        """
        初始化所有测试用户

        Args:
            api_tester: BaseAPITester 实例
        """
        if self.initialized:
            return

        # 1. 设置管理员（使用预置的 admin 账户）
        self._setup_admin(api_tester)

        # 2. 创建 teacher 用户
        self._setup_teacher(api_tester)

        # 3. 创建 student 用户
        self._setup_student(api_tester)

        self.initialized = True
        api_tester.log_success("用户池初始化完成")

    def _setup_admin(self, api_tester):
        """设置管理员账户"""
        # 使用预置的 admin/admin123 账户
        admin_login = {
            "username": "admin",
            "password": "admin123"
        }

        response = api_tester.make_request(
            "POST", "/users/login",
            data=admin_login,
            expect_status=200,
            validate_schema=False  # 避免循环依赖
        )

        if response and "token" in response:
            self.tokens['admin'] = response["token"]
            self.users['admin'] = response.get("user", {})
            self.user_ids['admin'] = response.get("user", {}).get("id")
            api_tester.log_success("管理员账户设置成功")
        else:
            raise Exception("无法登录管理员账户，请确保数据库中存在 admin 用户")

    def _setup_teacher(self, api_tester):
        """创建教师用户"""
        timestamp = int(time.time()) % 10000
        teacher_data = {
            "username": f"test_teacher_{timestamp}",
            "email": f"teacher_{timestamp}@test.com",
            "password": "test123456",
            "nickname": "测试教师"
        }

        # 1. 注册用户（默认为 student 角色）
        reg_response = api_tester.make_request(
            "POST", "/users/register",
            data=teacher_data,
            expect_status=201,
            validate_schema=False
        )

        if not reg_response or "user" not in reg_response:
            raise Exception("无法创建 teacher 测试用户")

        user_id = reg_response["user"]["id"]

        # 2. 使用管理员权限提升为 teacher
        api_tester.set_token(self.tokens['admin'])
        update_response = api_tester.make_request(
            "PUT", f"/admin/users/{user_id}/role",
            data={"role": "teacher"},
            expect_status=200,
            validate_schema=False
        )

        if not update_response:
            raise Exception(f"无法将用户 {user_id} 提升为 teacher 角色")

        # 3. 登录获取 teacher token
        login_response = api_tester.make_request(
            "POST", "/users/login",
            data={
                "username": teacher_data["username"],
                "password": teacher_data["password"]
            },
            expect_status=200,
            validate_schema=False
        )

        if login_response and "token" in login_response:
            self.tokens['teacher'] = login_response["token"]
            self.users['teacher'] = login_response.get("user", {})
            self.user_ids['teacher'] = user_id
            api_tester.log_success(f"教师用户创建成功 (ID: {user_id})")
        else:
            raise Exception("无法登录 teacher 用户")

    def _setup_student(self, api_tester):
        """创建学生用户"""
        timestamp = int(time.time()) % 10000
        student_data = {
            "username": f"test_student_{timestamp}",
            "email": f"student_{timestamp}@test.com",
            "password": "test123456",
            "nickname": "测试学生"
        }

        # 注册用户（默认就是 student 角色）
        reg_response = api_tester.make_request(
            "POST", "/users/register",
            data=student_data,
            expect_status=201,
            validate_schema=False
        )

        if not reg_response or "user" not in reg_response:
            raise Exception("无法创建 student 测试用户")

        user_id = reg_response["user"]["id"]

        # 登录获取 token
        login_response = api_tester.make_request(
            "POST", "/users/login",
            data={
                "username": student_data["username"],
                "password": student_data["password"]
            },
            expect_status=200,
            validate_schema=False
        )

        if login_response and "token" in login_response:
            self.tokens['student'] = login_response["token"]
            self.users['student'] = login_response.get("user", {})
            self.user_ids['student'] = user_id
            api_tester.log_success(f"学生用户创建成功 (ID: {user_id})")
        else:
            raise Exception("无法登录 student 用户")

    def get_user(self, role: str) -> Optional[Dict]:
        """
        获取指定角色的用户信息

        Args:
            role: 角色名称

        Returns:
            用户信息字典
        """
        return self.users.get(role)

    def get_token(self, role: str) -> Optional[str]:
        """
        获取指定角色的 token

        Args:
            role: 角色名称

        Returns:
            JWT token
        """
        return self.tokens.get(role)

    def get_user_id(self, role: str) -> Optional[int]:
        """
        获取指定角色的用户 ID

        Args:
            role: 角色名称

        Returns:
            用户 ID
        """
        return self.user_ids.get(role)

    def cleanup(self, api_tester):
        """
        清理测试用户

        Args:
            api_tester: BaseAPITester 实例
        """
        if not self.initialized:
            return

        # 使用管理员权限删除测试用户
        api_tester.set_token(self.tokens.get('admin'))

        # 删除 student 和 teacher（不删除 admin）
        for role in ['student', 'teacher']:
            user_id = self.user_ids.get(role)
            if user_id:
                # 注意：需要后端实现用户删除接口
                # 暂时通过数据库清理或保留测试用户
                api_tester.log_info(f"标记用户 {role} (ID: {user_id}) 待清理")

        api_tester.log_success("用户池清理完成")
        self.initialized = False
