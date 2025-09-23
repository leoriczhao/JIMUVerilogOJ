#!/usr/bin/env python3
"""
题目管理测试模块
测试题目的增删改查功能
"""

from base_test import BaseAPITester
from colorama import Back


class ProblemTester(BaseAPITester):
    """题目管理测试类"""
    
    def __init__(self):
        super().__init__()
        self.problem_id = None
        self.problem_ids = []  # 存储多个问题ID

    def setup_auth(self):
        """设置认证（需要先登录获取token）"""
        test_user = {
            "username": self.generate_unique_name("problemtester"),
            "email": f"{self.generate_unique_name('test')}@example.com",
            "password": "password123",
            "nickname": "题目测试员",
            "role": "admin"
        }
        
        reg_response = self.make_request("POST", "/users/register", data=test_user, expect_status=201)
        if not reg_response:
            return False
            
        login_data = {
            "username": test_user["username"],
            "password": test_user["password"]
        }
        login_response = self.make_request("POST", "/users/login", data=login_data, expect_status=200)
        
        if login_response and "token" in login_response:
            self.set_token(login_response["token"])
            return True
        return False

    def test_list_problems(self):
        """测试获取题目列表"""
        self.print_section_header("测试获取题目列表", Back.MAGENTA)
        response = self.make_request("GET", "/problems", expect_status=200)
        return response is not None

    def test_create_problem(self):
        """测试创建题目"""
        self.print_section_header("测试创建题目", Back.MAGENTA)
        problem_data = {
            "title": f"测试题目_{self.generate_unique_name()}",
            "description": "这是一个测试题目的描述。",
            "difficulty": "Easy",
            "time_limit": 1000,
            "memory_limit": 256,
            "tags": ["基础", "逻辑"]
        }
        response = self.make_request("POST", "/problems", data=problem_data, expect_status=201)
        if response and "problem" in response and "id" in response["problem"]:
            self.problem_id = response["problem"]["id"]
            self.problem_ids.append(self.problem_id)
            self.log_success(f"创建题目成功，ID: {self.problem_id}")
            return True
        self.log_error("创建题目失败")
        return False

    def test_create_multiple_problems(self):
        """测试创建多个题目"""
        self.print_section_header("测试创建多个题目", Back.MAGENTA)
        difficulties = ["Easy", "Medium", "Hard"]
        success_count = 0
        
        for i, difficulty in enumerate(difficulties, 1):
            problem_data = {
                "title": f"批量测试题目_{i}_{self.generate_unique_name()}",
                "description": f"这是第{i}个测试题目，难度为{difficulty}。",
                "difficulty": difficulty,
                "time_limit": 1000 + i * 500,
                "memory_limit": 256 + i * 128,
                "tags": ["批量测试", difficulty]
            }
            response = self.make_request("POST", "/problems", data=problem_data, expect_status=201)
            if response and "problem" in response and "id" in response["problem"]:
                problem_id = response["problem"]["id"]
                self.problem_ids.append(problem_id)
                self.log_success(f"创建题目{i}成功，ID: {problem_id}")
                success_count += 1
            else:
                self.log_error(f"创建题目{i}失败")
        
        self.log_info(f"批量创建完成，成功创建{success_count}个题目，总共{len(self.problem_ids)}个题目")
        return success_count > 0

    def test_get_problem_detail(self):
        """测试获取题目详情"""
        self.print_section_header("测试获取题目详情", Back.MAGENTA)
        if not self.problem_id:
            self.log_warning("没有题目ID，跳过测试")
            return True
        response = self.make_request("GET", f"/problems/{self.problem_id}", expect_status=200)
        return response is not None

    def test_update_problem(self):
        """测试更新题目"""
        self.print_section_header("测试更新题目", Back.MAGENTA)
        if not self.problem_id:
            self.log_warning("没有题目ID，跳过测试")
            return True
        update_data = {
            "title": f"更新后的题目_{self.generate_unique_name()}",
            "difficulty": "Medium"
        }
        response = self.make_request("PUT", f"/problems/{self.problem_id}", data=update_data, expect_status=200)
        return response is not None

    def test_delete_problem(self):
        """测试删除题目"""
        self.print_section_header("测试删除题目", Back.MAGENTA)
        if not self.problem_ids:
            self.log_warning("没有剩余的题目ID，跳过测试")
            return True
        
        # 使用剩余ID列表中的第一个ID进行删除
        problem_id = self.problem_ids[0]
        response = self.make_request("DELETE", f"/problems/{problem_id}", expect_status=200)
        if response:
            self.log_success(f"删除题目 {problem_id} 成功")
            self.problem_ids.remove(problem_id)
            if self.problem_id == problem_id:
                self.problem_id = None  # 清除ID
            return True
        return False

    def test_delete_partial_problems(self):
        """测试删除部分题目（保留一些用于查看列表效果）"""
        self.print_section_header("测试删除部分题目", Back.MAGENTA)
        if not self.problem_ids:
            self.log_warning("没有题目ID列表，跳过测试")
            return True
        
        # 只删除一半的题目，保留其他的
        delete_count = len(self.problem_ids) // 2
        if delete_count == 0:
            delete_count = 1  # 至少删除一个
        
        deleted_count = 0
        for i in range(delete_count):
            if i < len(self.problem_ids):
                problem_id = self.problem_ids[i]
                response = self.make_request("DELETE", f"/problems/{problem_id}", expect_status=200)
                if response:
                    self.log_success(f"删除题目 {problem_id} 成功")
                    deleted_count += 1
                else:
                    self.log_error(f"删除题目 {problem_id} 失败")
        
        # 更新问题ID列表，移除已删除的
        self.problem_ids = self.problem_ids[delete_count:]
        remaining_count = len(self.problem_ids)
        
        self.log_info(f"删除了{deleted_count}个题目，还剩余{remaining_count}个题目用于列表展示")
        return deleted_count > 0

    def test_unauthorized_create_problem(self):
        """测试未授权创建题目"""
        self.print_section_header("测试未授权创建题目", Back.RED)
        old_token = self.token
        self.clear_token() # 模拟未授权
        problem_data = {"title": "未授权"}
        response = self.make_request("POST", "/problems", data=problem_data, expect_status=401)
        self.set_token(old_token) # 恢复
        return response is not None

    def test_get_testcases(self):
        """测试获取题目测试用例"""
        self.print_section_header("测试获取题目测试用例", Back.MAGENTA)
        if not self.problem_id:
            self.log_warning("没有题目ID，跳过测试")
            return True
        response = self.make_request("GET", f"/problems/{self.problem_id}/testcases", expect_status=200)
        return response is not None

    def test_add_testcase(self):
        """测试添加测试用例"""
        self.print_section_header("测试添加测试用例", Back.MAGENTA)
        if not self.problem_id:
            self.log_warning("没有题目ID，跳过测试")
            return True
        testcase_data = {
            "input": "1 2 3",
            "output": "6",
            "is_sample": True
        }
        response = self.make_request("POST", f"/problems/{self.problem_id}/testcases", data=testcase_data, expect_status=201)
        return response is not None

    def test_get_problem_submissions(self):
        """测试获取题目提交记录"""
        self.print_section_header("测试获取题目提交记录", Back.MAGENTA)
        if not self.problem_id:
            self.log_warning("没有题目ID，跳过测试")
            return True
        response = self.make_request("GET", f"/problems/{self.problem_id}/submissions", expect_status=200)
        return response is not None

    def test_unauthorized_add_testcase(self):
        """测试未授权添加测试用例"""
        self.print_section_header("测试未授权添加测试用例", Back.RED)
        if not self.problem_id:
            self.log_warning("没有题目ID，跳过测试")
            return True
        old_token = self.token
        self.clear_token() # 模拟未授权
        testcase_data = {"input": "test", "output": "test"}
        response = self.make_request("POST", f"/problems/{self.problem_id}/testcases", data=testcase_data, expect_status=401)
        self.set_token(old_token) # 恢复
        return response is not None

    def run_tests(self):
        """运行所有测试"""
        self.print_section_header("题目管理测试模块", Back.BLUE)
        
        if not self.setup_auth():
            self.log_error("认证设置失败，跳过需要认证的测试")
            test_flow = [
                ("获取题目列表", self.test_list_problems),
                ("未授权创建题目", self.test_unauthorized_create_problem),
            ]
        else:
            test_flow = [
                ("获取题目列表 (初始)", self.test_list_problems),
                ("创建题目", self.test_create_problem),
                ("批量创建多个题目", self.test_create_multiple_problems),
                ("获取题目列表 (创建后)", self.test_list_problems),
                ("获取题目详情", self.test_get_problem_detail),
                ("获取题目测试用例", self.test_get_testcases),
                ("添加测试用例", self.test_add_testcase),
                ("获取题目提交记录", self.test_get_problem_submissions),
                ("更新题目", self.test_update_problem),
                ("删除部分题目", self.test_delete_partial_problems),
                ("获取题目列表 (部分删除后)", self.test_list_problems),
                ("删除单个题目", self.test_delete_problem),
                ("获取题目列表 (最终)", self.test_list_problems),
                ("未授权创建题目", self.test_unauthorized_create_problem),
                ("未授权添加测试用例", self.test_unauthorized_add_testcase),
            ]

        test_results = []
        for name, test_func in test_flow:
            test_results.append((name, test_func()))

        return self.print_test_summary(test_results)


def main():
    """主函数"""
    tester = ProblemTester()
    success = tester.run_tests()
    exit(0 if success else 1)


if __name__ == "__main__":
    main()
