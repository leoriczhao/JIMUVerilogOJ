#!/usr/bin/env python3
"""
题目管理测试模块（重构版）
展示使用新的 RBAC 测试框架进行测试的最佳实践
"""

from base_test import BaseAPITester
from colorama import Back


class ProblemTester(BaseAPITester):
    """题目管理测试类（基于 RBAC）"""

    def __init__(self):
        super().__init__()
        # 初始化用户池
        self.setup_user_pool()

    def run_tests(self):
        """主测试流程"""
        test_results = []

        try:
            # 1. 公开接口测试
            test_results.append(("获取题目列表(公开)", self.test_list_problems_public()))

            # 2. 学生角色测试
            self.print_section_header("学生角色测试", Back.BLUE)
            self.login_as('student')
            test_results.append(("学生-查看题目列表", self.test_list_problems()))
            test_results.append(("学生-创建题目(应拒绝)", self.test_create_problem_as_student()))

            # 3. 教师角色测试
            self.print_section_header("教师角色测试", Back.GREEN)
            self.login_as('teacher')
            test_results.append(("教师-创建题目", self.test_create_problem_as_teacher()))
            test_results.append(("教师-更新自己的题目", self.test_update_own_problem()))
            test_results.append(("教师-添加测试用例", self.test_add_testcase()))

            # 4. 管理员角色测试
            self.print_section_header("管理员角色测试", Back.MAGENTA)
            self.login_as('admin')
            test_results.append(("管理员-更新任意题目", self.test_update_any_problem()))
            test_results.append(("管理员-删除题目", self.test_delete_problem()))

            # 5. 权限边界测试
            self.print_section_header("权限边界测试", Back.RED)
            test_results.append(("未登录-创建题目(应拒绝)", self.test_create_problem_unauthorized()))

        finally:
            # 6. 清理测试数据
            self.cleanup()

        # 7. 打印测试报告
        return self.print_test_summary(test_results)

    # ========== 公开接口测试 ==========

    def test_list_problems_public(self):
        """测试公开获取题目列表（无需登录）"""
        self.print_section_header("测试公开获取题目列表", Back.CYAN)

        # 清除 token，模拟未登录用户
        self.clear_token()

        response = self.make_request(
            "GET", "/problems",
            expect_status=200,
            module="problem"
        )

        return response is not None

    # ========== 学生角色测试 ==========

    def test_list_problems(self):
        """学生查看题目列表（应该成功）"""
        response = self.make_request(
            "GET", "/problems",
            expect_status=200,
            module="problem"
        )

        if response:
            self.log_success(f"获取到 {response.get('total', 0)} 个题目")
            return True
        return False

    def test_create_problem_as_student(self):
        """学生创建题目（应该被拒绝 403）"""
        problem_data = {
            "title": "学生尝试创建的题目",
            "description": "这个请求应该被拒绝",
            "difficulty": "Easy",
            "time_limit": 1000,
            "memory_limit": 256,
            "tags": ["测试"]
        }

        response = self.make_request(
            "POST", "/problems",
            data=problem_data,
            expect_status=403,  # 期望权限不足
            module="problem"
        )

        # 验证是否正确返回 403
        return self.assert_forbidden(response)

    # ========== 教师角色测试 ==========

    def test_create_problem_as_teacher(self):
        """教师创建题目（应该成功）"""
        problem_data = {
            "title": f"测试题目_{self.generate_unique_name()}",
            "description": "这是一个教师创建的测试题目。",
            "difficulty": "Medium",
            "time_limit": 1500,
            "memory_limit": 512,
            "tags": ["测试", "教师创建"]
        }

        response = self.make_request(
            "POST", "/problems",
            data=problem_data,
            expect_status=201,
            module="problem"
        )

        if response:
            # 从响应中获取题目 ID（可能在 response['problem']['id'] 或 response['id']）
            problem_id = None
            if 'problem' in response and 'id' in response['problem']:
                problem_id = response['problem']['id']
            elif 'id' in response:
                problem_id = response['id']

            if problem_id:
                self.log_success(f"教师成功创建题目，ID: {problem_id}")

                # 标记用于清理
                self.mark_for_cleanup('problem', problem_id)

                # 保存用于后续测试
                self.problem_id = problem_id
                return True

        return False

    def test_update_own_problem(self):
        """教师更新自己创建的题目（应该成功）"""
        if not hasattr(self, 'problem_id'):
            self.log_warning("没有可更新的题目，跳过测试")
            return True

        update_data = {
            "title": f"更新后的题目_{self.generate_unique_name()}",
            "description": "题目已被更新",
            "difficulty": "Hard"
        }

        response = self.make_request(
            "PUT", f"/problems/{self.problem_id}",
            data=update_data,
            expect_status=200,
            module="problem"
        )

        return response is not None

    def test_add_testcase(self):
        """教师添加测试用例（应该成功）"""
        if not hasattr(self, 'problem_id'):
            self.log_warning("没有可添加测试用例的题目，跳过测试")
            return True

        testcase_data = {
            "input": "test input",
            "expected_output": "test output",
            "is_sample": True,
            "score": 10
        }

        response = self.make_request(
            "POST", f"/problems/{self.problem_id}/testcases",
            data=testcase_data,
            expect_status=201,
            module="problem"
        )

        return response is not None

    # ========== 管理员角色测试 ==========

    def test_update_any_problem(self):
        """管理员更新任意题目（应该成功）"""
        if not hasattr(self, 'problem_id'):
            self.log_warning("没有可更新的题目，跳过测试")
            return True

        update_data = {
            "title": "管理员更新的题目",
            "description": "管理员可以更新任意题目"
        }

        response = self.make_request(
            "PUT", f"/problems/{self.problem_id}",
            data=update_data,
            expect_status=200,
            module="problem"
        )

        return response is not None

    def test_delete_problem(self):
        """管理员删除题目（应该成功）"""
        if not hasattr(self, 'problem_id'):
            self.log_warning("没有可删除的题目，跳过测试")
            return True

        response = self.make_request(
            "DELETE", f"/problems/{self.problem_id}",
            expect_status=200,
            module="problem"
        )

        if response:
            # 已删除，无需清理
            self.cleanup_items = [
                item for item in self.cleanup_items
                if not (item['type'] == 'problem' and item['id'] == self.problem_id)
            ]
            return True

        return False

    # ========== 权限边界测试 ==========

    def test_create_problem_unauthorized(self):
        """未登录创建题目（应该被拒绝 401）"""
        self.clear_token()

        problem_data = {
            "title": "未授权创建"
        }

        response = self.make_request(
            "POST", "/problems",
            data=problem_data,
            expect_status=401,  # 期望未认证
            module="problem"
        )

        return self.assert_unauthorized(response)


def main():
    """主函数"""
    print("\n" + "=" * 60)
    print(" 题目管理测试模块（重构版 - 基于 RBAC）")
    print("=" * 60 + "\n")

    tester = ProblemTester()
    success = tester.run_tests()

    return 0 if success else 1


if __name__ == "__main__":
    import sys
    sys.exit(main())
