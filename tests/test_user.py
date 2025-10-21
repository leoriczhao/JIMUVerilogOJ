#!/usr/bin/env python3
"""
用户管理测试模块
测试用户注册、登录、信息管理等功能
"""

from base_test import BaseAPITester
from colorama import Back


class UserTester(BaseAPITester):
    """用户管理测试类"""
    
    def __init__(self):
        super().__init__()
        self.test_user_data = None
    
    def test_user_registration(self):
        """测试用户注册"""
        self.print_section_header("测试用户注册", Back.GREEN)

        # 生成唯一用户名
        timestamp = self.generate_unique_name("testuser")
        test_user = {
            "username": timestamp,
            "email": f"{timestamp}@example.com",
            "password": "password123",
            "nickname": "测试用户",
            "school": "测试大学",
            "student_id": "202301001"
        }

        response = self.make_request("POST", "/users/register", data=test_user,
                                     expect_status=201, module="user")

        if response is None:
            self.log_error("用户注册失败")
            return False
        else:
            self.log_success("用户注册成功")
            self.test_user_data = test_user
            return True

    def test_user_login(self):
        """测试用户登录"""
        self.print_section_header("测试用户登录", Back.GREEN)

        if not self.test_user_data:
            self.log_error("需要先注册用户")
            return False

        login_data = {
            "username": self.test_user_data["username"],
            "password": self.test_user_data["password"]
        }

        response = self.make_request("POST", "/users/login", data=login_data,
                                     expect_status=200, module="user")

        if response is None:
            self.log_error("用户登录失败")
            return False
        else:
            self.set_token(response["token"])
            if "user" in response:
                self.user_id = response["user"].get("id")
            self.log_success("用户登录成功")
            self.log_info(f"获取到Token: {self.token[:20]}...")
            return True

    def test_get_profile(self):
        """测试获取用户信息"""
        self.print_section_header("测试获取用户信息", Back.GREEN)

        if not self.token:
            self.log_error("需要先登录获取Token")
            return False

        response = self.make_request("GET", "/users/profile",
                                     expect_status=200, module="user")

        if response is None:
            self.log_error("获取用户信息失败")
            return False
        else:
            self.log_success("获取用户信息成功")
            return True

    def test_update_profile(self):
        """测试更新用户信息"""
        self.print_section_header("测试更新用户信息", Back.GREEN)

        if not self.token:
            self.log_error("需要先登录获取Token")
            return False

        update_data = {
            "nickname": "更新后的昵称",
            "school": "新的学校",
            "student_id": "202301002"
        }

        response = self.make_request("PUT", "/users/profile", data=update_data,
                                     expect_status=200, module="user")

        if response is None:
            self.log_error("更新用户信息失败")
            return False
        else:
            self.log_success("更新用户信息成功")
            return True

    def test_duplicate_registration(self):
        """测试重复注册"""
        self.print_section_header("测试重复注册", Back.GREEN)

        if not self.test_user_data:
            self.log_warning("需要先完成正常注册")
            return True  # 跳过测试

        # 尝试用相同用户名注册
        duplicate_user = {
            "username": self.test_user_data["username"],
            "email": "different@example.com",
            "password": "password123"
        }

        response = self.make_request("POST", "/users/register", data=duplicate_user,
                                     expect_status=400, module="user")

        if response is None:
            self.log_error("重复注册检查失败")
            return False
        else:
            self.log_success("重复注册正确被拒绝")
            return True

    def test_invalid_login(self):
        """测试无效登录"""
        self.print_section_header("测试无效登录", Back.GREEN)

        invalid_login = {
            "username": "nonexistent_user",
            "password": "wrong_password"
        }

        response = self.make_request("POST", "/users/login", data=invalid_login,
                                     expect_status=401, module="user")

        if response is None:
            self.log_error("无效登录检查失败")
            return False
        else:
            self.log_success("无效登录正确被拒绝")
            return True

    def test_change_password(self):
        """测试修改密码"""
        self.print_section_header("测试修改密码", Back.GREEN)

        if not self.token:
            self.log_error("需要先登录获取Token")
            return False

        # 修改密码
        password_data = {
            "old_password": self.test_user_data["password"],
            "new_password": "newpassword123"
        }

        response = self.make_request("PUT", "/users/password", data=password_data,
                                     expect_status=200, module="user")

        if response is None:
            self.log_error("修改密码失败")
            return False
        else:
            self.log_success("修改密码成功")
            # 更新测试数据中的密码
            self.test_user_data["password"] = password_data["new_password"]
            return True

    def test_change_password_wrong_old(self):
        """测试使用错误的旧密码修改密码"""
        self.print_section_header("测试错误的旧密码", Back.GREEN)

        if not self.token:
            self.log_error("需要先登录获取Token")
            return False

        # 使用错误的旧密码
        password_data = {
            "old_password": "wrongoldpassword",
            "new_password": "newpassword456"
        }

        response = self.make_request("PUT", "/users/password", data=password_data,
                                     expect_status=401, module="user")

        if response is None:
            self.log_error("错误旧密码检查失败")
            return False
        else:
            self.log_success("错误旧密码正确被拒绝")
            return True

    def test_login_with_new_password(self):
        """测试使用新密码登录"""
        self.print_section_header("测试新密码登录", Back.GREEN)

        if not self.test_user_data:
            self.log_warning("需要先完成密码修改")
            return True

        # 清除旧token
        self.clear_token()

        # 使用新密码登录
        login_data = {
            "username": self.test_user_data["username"],
            "password": self.test_user_data["password"]
        }

        response = self.make_request("POST", "/users/login", data=login_data,
                                     expect_status=200, module="user")

        if response is None:
            self.log_error("新密码登录失败")
            return False
        else:
            self.set_token(response["token"])
            self.log_success("新密码登录成功")
            return True

    def test_unauthorized_access(self):
        """测试未授权访问"""
        self.print_section_header("测试未授权访问", Back.RED)

        # 临时清除token
        old_token = self.token
        self.clear_token()

        # 测试需要认证的接口
        self.log_info("测试未授权访问用户信息")
        response = self.make_request("GET", "/users/profile",
                                     expect_status=401, module="user")

        # 恢复token
        self.set_token(old_token)

        if response is None:
            self.log_error("未授权访问检查失败")
            return False
        else:
            self.log_success("未授权访问正确被拒绝")
            return True

    def test_invalid_email_format(self):
        """测试无效的邮箱格式"""
        self.print_section_header("测试无效邮箱格式", Back.YELLOW)

        invalid_user = {
            "username": self.generate_unique_name("invalid"),
            "email": "invalid-email-format",  # 无效邮箱格式
            "password": "password123"
        }

        response = self.make_request("POST", "/users/register", data=invalid_user,
                                     expect_status=400, module="user")

        if response is None:
            self.log_error("无效邮箱格式检查失败")
            return False
        else:
            self.log_success("无效邮箱格式正确被拒绝")
            return True

    def test_short_username(self):
        """测试用户名过短"""
        self.print_section_header("测试用户名过短", Back.YELLOW)

        short_user = {
            "username": "ab",  # 少于3个字符
            "email": f"{self.generate_unique_name('test')}@example.com",
            "password": "password123"
        }

        response = self.make_request("POST", "/users/register", data=short_user,
                                     expect_status=400, module="user")

        if response is None:
            self.log_error("用户名过短检查失败")
            return False
        else:
            self.log_success("用户名过短正确被拒绝")
            return True

    def test_short_password(self):
        """测试密码过短"""
        self.print_section_header("测试密码过短", Back.YELLOW)

        short_pwd_user = {
            "username": self.generate_unique_name("test"),
            "email": f"{self.generate_unique_name('test')}@example.com",
            "password": "12345"  # 少于6个字符
        }

        response = self.make_request("POST", "/users/register", data=short_pwd_user,
                                     expect_status=400, module="user")

        if response is None:
            self.log_error("密码过短检查失败")
            return False
        else:
            self.log_success("密码过短正确被拒绝")
            return True

    def run_tests(self):
        """运行用户管理测试"""
        print("👤 开始用户管理测试")
        print("=" * 50)

        test_results = [
            # 基本功能测试
            ("用户注册", self.test_user_registration()),
            ("用户登录", self.test_user_login()),
            ("获取用户信息", self.test_get_profile()),
            ("更新用户信息", self.test_update_profile()),

            # 密码管理测试
            ("修改密码", self.test_change_password()),
            ("错误旧密码", self.test_change_password_wrong_old()),
            ("新密码登录", self.test_login_with_new_password()),

            # 错误处理测试
            ("重复注册检查", self.test_duplicate_registration()),
            ("无效登录检查", self.test_invalid_login()),
            ("未授权访问检查", self.test_unauthorized_access()),

            # 输入验证测试
            ("无效邮箱格式", self.test_invalid_email_format()),
            ("用户名过短", self.test_short_username()),
            ("密码过短", self.test_short_password()),
        ]

        return self.print_test_summary(test_results)


def main():
    """主函数"""
    tester = UserTester()
    success = tester.run_tests()
    exit(0 if success else 1)


if __name__ == "__main__":
    main() 