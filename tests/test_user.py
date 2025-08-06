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
        
        response = self.make_request("POST", "/users/register", data=test_user, expect_status=201)
        
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
        
        response = self.make_request("POST", "/users/login", data=login_data, expect_status=200)
        
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
            
        response = self.make_request("GET", "/users/profile", expect_status=200)
        
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
        
        response = self.make_request("PUT", "/users/profile", data=update_data, expect_status=200)
        
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
        
        response = self.make_request("POST", "/users/register", data=duplicate_user, expect_status=400)
        
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
        
        response = self.make_request("POST", "/users/login", data=invalid_login, expect_status=401)
        
        if response is None:
            self.log_error("无效登录检查失败")
            return False
        else:
            self.log_success("无效登录正确被拒绝")
            return True

    def test_unauthorized_access(self):
        """测试未授权访问"""
        self.print_section_header("测试未授权访问", Back.RED)
        
        # 临时清除token
        old_token = self.token
        self.clear_token()
        
        # 测试需要认证的接口
        self.log_info("测试未授权访问用户信息")
        response = self.make_request("GET", "/users/profile", expect_status=401)
        
        # 恢复token
        self.set_token(old_token)
        
        if response is None:
            self.log_error("未授权访问检查失败")
            return False
        else:
            self.log_success("未授权访问正确被拒绝")
            return True

    def run_tests(self):
        """运行用户管理测试"""
        print("👤 开始用户管理测试")
        print("=" * 50)
        
        test_results = [
            ("用户注册", self.test_user_registration()),
            ("用户登录", self.test_user_login()),
            ("获取用户信息", self.test_get_profile()),
            ("更新用户信息", self.test_update_profile()),
            ("重复注册检查", self.test_duplicate_registration()),
            ("无效登录检查", self.test_invalid_login()),
            ("未授权访问检查", self.test_unauthorized_access()),
        ]
        
        return self.print_test_summary(test_results)


def main():
    """主函数"""
    tester = UserTester()
    success = tester.run_tests()
    exit(0 if success else 1)


if __name__ == "__main__":
    main() 