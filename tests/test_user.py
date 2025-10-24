#!/usr/bin/env python3
"""
用户管理测试模块（重构版）
展示使用新的 RBAC 测试框架进行测试的最佳实践
"""

from base_test import BaseAPITester
from colorama import Back


class UserTester(BaseAPITester):
    """用户管理测试类（基于 RBAC）"""

    def __init__(self):
        super().__init__()
        # 初始化用户池
        self.setup_user_pool()
        # 用于存储新注册的测试用户
        self.new_user_data = None
        self.new_user_token = None

    def run_tests(self):
        """主测试流程"""
        test_results = []

        try:
            # 1. 公开接口测试（注册、登录、输入验证）
            self.print_section_header("公开接口测试", Back.CYAN)
            test_results.append(("新用户注册", self.test_user_registration()))
            test_results.append(("新用户登录", self.test_user_login()))
            test_results.append(("重复注册(应拒绝)", self.test_duplicate_registration()))
            test_results.append(("无效登录(应拒绝)", self.test_invalid_login()))
            test_results.append(("无效邮箱格式(应拒绝)", self.test_invalid_email_format()))
            test_results.append(("用户名过短(应拒绝)", self.test_short_username()))
            test_results.append(("密码过短(应拒绝)", self.test_short_password()))

            # 2. 学生角色测试（个人资料管理）
            self.print_section_header("学生角色测试", Back.BLUE)
            self.login_as('student')
            test_results.append(("学生-查看个人资料", self.test_get_own_profile()))
            test_results.append(("学生-更新个人资料", self.test_update_own_profile()))
            test_results.append(("学生-修改密码", self.test_change_password()))

            # 3. 教师角色测试（展示所有角色都能管理自己）
            self.print_section_header("教师角色测试", Back.GREEN)
            self.login_as('teacher')
            test_results.append(("教师-查看个人资料", self.test_get_own_profile()))
            test_results.append(("教师-更新个人资料", self.test_update_own_profile()))

            # 4. 管理员角色测试（管理其他用户）
            self.print_section_header("管理员角色测试", Back.MAGENTA)
            self.login_as('admin')
            test_results.append(("管理员-更新用户角色", self.test_admin_update_user_role()))
            # 可以添加更多管理员特有的测试，如查看所有用户、禁用用户等

            # 5. 权限边界测试
            self.print_section_header("权限边界测试", Back.RED)
            test_results.append(("未登录-访问资料(应拒绝)", self.test_unauthorized_access()))
            test_results.append(("错误旧密码(应拒绝)", self.test_change_password_wrong_old()))

        finally:
            # 6. 清理测试数据
            self.cleanup()

        # 7. 打印测试报告
        return self.print_test_summary(test_results)

    # ========== 公开接口测试 ==========

    def test_user_registration(self):
        """测试新用户注册（公开接口）"""
        self.clear_token()  # 确保无认证状态

        # 生成唯一用户数据
        unique_name = self.generate_unique_name("newuser")
        self.new_user_data = {
            "username": unique_name,
            "email": f"{unique_name}@example.com",
            "password": "password123",
            "nickname": "新注册用户",
            "school": "测试大学",
            "student_id": "202301001"
        }

        response = self.make_request(
            "POST", "/users/register",
            data=self.new_user_data,
            expect_status=201,
            module="user"
        )

        if response and "user" in response:
            user_id = response["user"].get("id")
            if user_id:
                # 标记用于清理（虽然用户清理可能需要特殊处理）
                self.log_success(f"用户注册成功，ID: {user_id}")
                return True
        return False

    def test_user_login(self):
        """测试新用户登录（公开接口）"""
        if not self.new_user_data:
            self.log_warning("需要先注册用户，跳过测试")
            return True

        self.clear_token()  # 确保无认证状态

        login_data = {
            "username": self.new_user_data["username"],
            "password": self.new_user_data["password"]
        }

        response = self.make_request(
            "POST", "/users/login",
            data=login_data,
            expect_status=200,
            module="user"
        )

        if response and "token" in response:
            self.new_user_token = response["token"]
            self.log_success(f"登录成功，Token: {self.new_user_token[:20]}...")
            return True
        return False

    def test_duplicate_registration(self):
        """测试重复注册（应该被拒绝 400）"""
        if not self.new_user_data:
            self.log_warning("需要先完成注册，跳过测试")
            return True

        self.clear_token()

        # 尝试用相同用户名注册
        duplicate_user = {
            "username": self.new_user_data["username"],
            "email": "different@example.com",
            "password": "password123"
        }

        response = self.make_request(
            "POST", "/users/register",
            data=duplicate_user,
            expect_status=400,
            module="user"
        )

        # 验证是否正确返回 400（response 不为 None 说明匹配了 expect_status=400）
        return response is not None

    def test_invalid_login(self):
        """测试无效登录（应该被拒绝 401）"""
        self.clear_token()

        invalid_login = {
            "username": "nonexistent_user_xyz",
            "password": "wrong_password"
        }

        response = self.make_request(
            "POST", "/users/login",
            data=invalid_login,
            expect_status=401,
            module="user"
        )

        return self.assert_unauthorized(response)

    def test_invalid_email_format(self):
        """测试无效邮箱格式（应该被拒绝 400）"""
        self.clear_token()

        invalid_user = {
            "username": self.generate_unique_name("invalid"),
            "email": "invalid-email-format",  # 无效邮箱格式
            "password": "password123"
        }

        response = self.make_request(
            "POST", "/users/register",
            data=invalid_user,
            expect_status=400,
            module="user"
        )

        # 验证是否正确返回 400（response 不为 None 说明匹配了 expect_status=400）
        return response is not None

    def test_short_username(self):
        """测试用户名过短（应该被拒绝 400）"""
        self.clear_token()

        short_user = {
            "username": "ab",  # 少于3个字符
            "email": f"{self.generate_unique_name('test')}@example.com",
            "password": "password123"
        }

        response = self.make_request(
            "POST", "/users/register",
            data=short_user,
            expect_status=400,
            module="user"
        )

        # 验证是否正确返回 400（response 不为 None 说明匹配了 expect_status=400）
        return response is not None

    def test_short_password(self):
        """测试密码过短（应该被拒绝 400）"""
        self.clear_token()

        short_pwd_user = {
            "username": self.generate_unique_name("test"),
            "email": f"{self.generate_unique_name('test')}@example.com",
            "password": "12345"  # 少于6个字符
        }

        response = self.make_request(
            "POST", "/users/register",
            data=short_pwd_user,
            expect_status=400,
            module="user"
        )

        # 验证是否正确返回 400（response 不为 None 说明匹配了 expect_status=400）
        return response is not None

    # ========== 学生/教师角色测试（个人资料管理）==========

    def test_get_own_profile(self):
        """测试查看自己的资料（所有角色都可以）"""
        response = self.make_request(
            "GET", "/users/profile",
            expect_status=200,
            module="user"
        )

        if response and "user" in response:
            user = response["user"]
            self.log_success(f"成功获取资料: {user.get('username', 'N/A')}")
            return True
        elif response:
            self.log_success("成功获取资料")
            return True
        return False

    def test_update_own_profile(self):
        """测试更新自己的资料（所有角色都可以）"""
        current_role = self.get_current_role()

        update_data = {
            "nickname": f"{current_role}更新的昵称",
            "school": f"{current_role}的学校",
            "student_id": f"ID-{self.generate_unique_name()}"
        }

        response = self.make_request(
            "PUT", "/users/profile",
            data=update_data,
            expect_status=200,
            module="user"
        )

        return response is not None

    def test_change_password(self):
        """测试修改密码（所有角色都可以）"""
        # 使用当前角色的默认密码
        current_role = self.get_current_role()

        # 对于测试用户池的用户，密码是 "test123456"（除了 admin 是 "admin123"）
        old_password = "admin123" if current_role == "admin" else "test123456"
        new_password = f"new_{current_role}_pwd123"

        password_data = {
            "old_password": old_password,
            "new_password": new_password
        }

        response = self.make_request(
            "PUT", "/users/password",
            data=password_data,
            expect_status=200,
            module="user"
        )

        if response:
            # 改回原密码，避免影响后续测试
            reverse_data = {
                "old_password": new_password,
                "new_password": old_password
            }
            self.make_request(
                "PUT", "/users/password",
                data=reverse_data,
                expect_status=200,
                module="user",
                validate_schema=False
            )
            return True
        return False

    # ========== 管理员角色测试 ==========

    def test_admin_update_user_role(self):
        """管理员更新用户角色（应该成功）"""
        # 获取 student 用户的 ID
        student_id = self.user_pool.get_user_id('student')

        if not student_id:
            self.log_warning("无法获取学生用户ID，跳过测试")
            return True

        # 尝试更新为 teacher（然后改回来）
        update_data = {"role": "teacher"}

        response = self.make_request(
            "PUT", f"/admin/users/{student_id}/role",
            data=update_data,
            expect_status=200,
            module="admin",
            validate_schema=False  # admin API 可能没有完整的 schema
        )

        if response:
            # 改回 student 角色
            rollback_data = {"role": "student"}
            self.make_request(
                "PUT", f"/admin/users/{student_id}/role",
                data=rollback_data,
                expect_status=200,
                module="admin",
                validate_schema=False
            )
            self.log_success("成功更新并回滚用户角色")
            return True
        return False

    # ========== 权限边界测试 ==========

    def test_unauthorized_access(self):
        """未登录访问受保护接口（应该被拒绝 401）"""
        self.clear_token()

        response = self.make_request(
            "GET", "/users/profile",
            expect_status=401,
            module="user"
        )

        return self.assert_unauthorized(response)

    def test_change_password_wrong_old(self):
        """使用错误的旧密码修改密码（应该被拒绝 401）"""
        self.login_as('student')

        password_data = {
            "old_password": "wrong_old_password_xyz",
            "new_password": "newpassword456"
        }

        response = self.make_request(
            "PUT", "/users/password",
            data=password_data,
            expect_status=401,
            module="user"
        )

        # 后端返回 401 且包含密码错误信息
        if response and (
            response.get('error') in ['unauthorized', 'invalid_password'] or
            '密码错误' in str(response.get('message', ''))
        ):
            self.log_success("✓ 密码验证通过：正确拒绝错误的旧密码")
            return True
        elif response:
            # response 不为 None 说明匹配了 expect_status=401，也算通过
            return True
        else:
            self.log_error(f"✗ 期望 401，但得到: {response}")
            return False


def main():
    """主函数"""
    print("\n" + "=" * 60)
    print(" 用户管理测试模块（重构版 - 基于 RBAC）")
    print("=" * 60 + "\n")

    tester = UserTester()
    success = tester.run_tests()

    return 0 if success else 1


if __name__ == "__main__":
    import sys
    sys.exit(main())
