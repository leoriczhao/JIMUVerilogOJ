#!/usr/bin/env python3
"""
提交管理测试模块
测试代码提交、状态查询等功能
"""

from base_test import BaseAPITester
from colorama import Back


class SubmissionTester(BaseAPITester):
    """提交管理测试类"""
    
    def __init__(self):
        super().__init__()
        self.submission_id = None
        self.problem_id = None
        self.admin_token = None
        self.student_token = None

    def setup_users(self):
        """创建并登录管理员和学生用户"""
        # 创建管理员
        admin_user = {
            "username": self.generate_unique_name("sub_admin"),
            "email": f"{self.generate_unique_name('admin')}@example.com",
            "password": "password123",
            "nickname": "Admin Submitter",
            "role": "admin"
        }
        reg_admin_res = self.make_request("POST", "/users/register", data=admin_user, expect_status=201)
        if not reg_admin_res:
            return False
        login_admin_res = self.make_request("POST", "/users/login", data={"username": admin_user["username"], "password": admin_user["password"]})
        if not login_admin_res or "token" not in login_admin_res:
            return False
        self.admin_token = login_admin_res["token"]

        # 创建学生
        student_user = {
            "username": self.generate_unique_name("sub_student"),
            "email": f"{self.generate_unique_name('student')}@example.com",
            "password": "password123",
            "nickname": "Student Submitter",
            "role": "student"
        }
        reg_student_res = self.make_request("POST", "/users/register", data=student_user, expect_status=201)
        if not reg_student_res:
            return False
        login_student_res = self.make_request("POST", "/users/login", data={"username": student_user["username"], "password": student_user["password"]})
        if not login_student_res or "token" not in login_student_res:
            return False
        self.student_token = login_student_res["token"]
        return True

    def create_test_problem(self):
        """使用管理员账户创建测试题目"""
        self.set_token(self.admin_token)
        problem_data = {
            "title": f"提交测试题目_{self.generate_unique_name()}",
            "description": "这是一个用于提交测试的题目",
            "difficulty": "Easy",
            "time_limit": 1000,
            "memory_limit": 128
        }
        response = self.make_request("POST", "/problems", data=problem_data, expect_status=201)
        if response and "problem" in response and "id" in response["problem"]:
            self.problem_id = response["problem"]["id"]
            self.log_success(f"创建测试题目成功，ID: {self.problem_id}")
            # Make the problem public
            update_data = {"is_public": True}
            update_response = self.make_request("PUT", f"/problems/{self.problem_id}", data=update_data, expect_status=200)
            if update_response:
                self.log_success(f"题目 {self.problem_id} 已设为公开")
                return True
            else:
                self.log_error(f"设置题目 {self.problem_id} 为公开失败")
                return False
        self.log_error("创建测试题目失败")
        return False

    def test_student_create_submission(self):
        """测试学生创建提交"""
        self.print_section_header("测试学生创建提交", Back.YELLOW)
        self.set_token(self.student_token)
        verilog_code = "module test; endmodule"
        submission_data = {"problem_id": self.problem_id, "code": verilog_code, "language": "verilog"}
        response = self.make_request("POST", "/submissions", data=submission_data, expect_status=201)
        if response and "submission" in response and "id" in response["submission"]:
            self.submission_id = response["submission"]["id"]
            self.log_success(f"学生创建提交成功，ID: {self.submission_id}")
            return True
        self.log_error("学生创建提交失败")
        return False

    def test_get_submission_detail(self):
        """测试获取提交详情"""
        self.print_section_header("测试获取提交详情", Back.YELLOW)
        response = self.make_request("GET", f"/submissions/{self.submission_id}", expect_status=200)
        return response is not None

    def test_get_user_submissions(self):
        """测试获取用户提交记录"""
        self.print_section_header("测试获取用户提交记录", Back.YELLOW)
        response = self.make_request("GET", "/submissions/user", expect_status=200)
        return response is not None

    def test_get_submission_stats(self):
        """测试获取提交统计"""
        self.print_section_header("测试获取提交统计", Back.YELLOW)
        response = self.make_request("GET", "/submissions/stats", expect_status=200)
        return response is not None

    def test_delete_submission(self):
        """测试删除提交"""
        self.print_section_header("测试删除提交", Back.YELLOW)
        if not self.submission_id:
            self.log_warning("没有提交ID，跳过测试")
            return True
        response = self.make_request("DELETE", f"/submissions/{self.submission_id}", expect_status=200)
        if response:
            self.log_success(f"删除提交 {self.submission_id} 成功")
            self.submission_id = None  # 清除ID
            return True
        return False

    def test_unauthorized_delete_submission(self):
        """测试未授权删除提交"""
        self.print_section_header("测试未授权删除提交", Back.RED)
        # 先创建一个新的提交用于测试
        self.set_token(self.student_token)
        verilog_code = "module test2; endmodule"
        submission_data = {"problem_id": self.problem_id, "code": verilog_code, "language": "verilog"}
        create_response = self.make_request("POST", "/submissions", data=submission_data, expect_status=201)
        if not create_response or "submission" not in create_response:
            self.log_warning("创建测试提交失败，跳过测试")
            return True
        
        test_submission_id = create_response["submission"]["id"]
        
        # 清除token模拟未授权
        old_token = self.token
        self.clear_token()
        response = self.make_request("DELETE", f"/submissions/{test_submission_id}", expect_status=401)
        self.set_token(old_token)  # 恢复token
        return response is not None

    def test_list_submissions(self):
        """测试获取提交列表"""
        self.print_section_header("测试获取提交列表", Back.YELLOW)
        response = self.make_request("GET", "/submissions", expect_status=200)
        return response is not None

    def run_tests(self):
        """运行所有测试"""
        self.print_section_header("提交管理测试模块", Back.BLUE)

        if not self.setup_users():
            self.log_error("用户设置失败，测试终止")
            return False

        if not self.create_test_problem():
            self.log_error("创建测试题目失败，测试终止")
            return False

        # Now run the submission tests with the created problem
        test_flow = [
            ("获取提交列表", self.test_list_submissions),
            ("学生创建提交", self.test_student_create_submission),
            ("获取提交详情", self.test_get_submission_detail),
            ("获取用户提交记录", self.test_get_user_submissions),
            ("获取提交统计", self.test_get_submission_stats),
            ("删除提交", self.test_delete_submission),
            ("未授权删除提交", self.test_unauthorized_delete_submission),
        ]

        test_results = []
        for name, test_func in test_flow:
            test_results.append((name, test_func()))

        return self.print_test_summary(test_results)


def main():
    """主函数"""
    tester = SubmissionTester()
    success = tester.run_tests()
    exit(0 if success else 1)


if __name__ == "__main__":
    main()
