#!/usr/bin/env python3
"""
新闻管理测试模块（重构版）
展示使用新的 RBAC 测试框架进行测试的最佳实践
"""

from base_test import BaseAPITester
from colorama import Back


class NewsTester(BaseAPITester):
    """新闻管理测试类（基于 RBAC）"""

    def __init__(self):
        super().__init__()
        # 初始化用户池
        self.setup_user_pool()
        # 存储创建的新闻ID
        self.teacher_news_id = None
        self.admin_news_id = None

    def run_tests(self):
        """主测试流程"""
        test_results = []

        try:
            # 1. 公开接口测试（查看新闻列表）
            self.print_section_header("公开接口测试", Back.CYAN)
            test_results.append(("获取新闻列表(公开)", self.test_list_news_public()))

            # 2. 学生角色测试（只读权限）
            self.print_section_header("学生角色测试", Back.BLUE)
            self.login_as('student')
            test_results.append(("学生-查看新闻列表", self.test_list_news()))
            test_results.append(("学生-创建新闻(应拒绝)", self.test_student_create_news()))

            # 3. 教师角色测试（创建、更新、删除新闻）
            self.print_section_header("教师角色测试", Back.GREEN)
            self.login_as('teacher')
            test_results.append(("教师-创建新闻", self.test_create_news_as_teacher()))
            test_results.append(("教师-查看新闻详情", self.test_get_news_detail()))
            test_results.append(("教师-更新自己的新闻", self.test_update_own_news()))
            test_results.append(("教师-删除自己的新闻", self.test_delete_own_news()))

            # 4. 管理员角色测试（管理所有新闻、发布、归档）
            self.print_section_header("管理员角色测试", Back.MAGENTA)
            self.login_as('admin')
            test_results.append(("管理员-创建新闻", self.test_create_news_as_admin()))
            test_results.append(("管理员-更新任意新闻", self.test_admin_update_any_news()))
            test_results.append(("管理员-删除新闻", self.test_admin_delete_news()))

            # 5. 权限边界测试
            self.print_section_header("权限边界测试", Back.RED)
            test_results.append(("未登录-创建新闻(应拒绝)", self.test_unauthorized_create_news()))

        finally:
            # 6. 清理测试数据
            self.cleanup()

        # 7. 打印测试报告
        return self.print_test_summary(test_results)

    # ========== 公开接口测试 ==========

    def test_list_news_public(self):
        """测试公开获取新闻列表（无需登录）"""
        self.clear_token()

        response = self.make_request(
            "GET", "/news",
            expect_status=200,
            module="news",
            validate_schema=False
        )

        if response and "news" in response:
            news_list = response.get("news") or []
            total = response.get("total", 0)
            self.log_success(f"获取到 {total} 条新闻，当前页 {len(news_list)} 条")
            return True
        elif response:
            self.log_success("成功获取新闻列表")
            return True
        return False

    # ========== 学生角色测试 ==========

    def test_list_news(self):
        """学生查看新闻列表（应该成功）"""
        response = self.make_request(
            "GET", "/news",
            expect_status=200,
            module="news",
            validate_schema=False
        )

        if response and "news" in response:
            news_list = response.get("news") or []
            self.log_success(f"成功获取 {len(news_list)} 条新闻")
            return True
        elif response:
            self.log_success("成功获取新闻列表")
            return True
        return False

    def test_student_create_news(self):
        """学生创建新闻（应该被拒绝 403）"""
        news_data = {
            "title": "学生尝试创建的新闻",
            "content": "这个请求应该被拒绝，因为学生没有创建新闻的权限。",
            "is_published": True
        }

        response = self.make_request(
            "POST", "/news",
            data=news_data,
            expect_status=403,
            module="news",
            validate_schema=False
        )

        return self.assert_forbidden(response)

    # ========== 教师角色测试 ==========

    def test_create_news_as_teacher(self):
        """教师创建新闻（应该成功）"""
        news_data = {
            "title": f"教师测试新闻_{self.generate_unique_name()}",
            "content": "这是教师创建的测试新闻内容。教师可以创建新闻来发布公告和技术文章。",
            "summary": "教师测试新闻摘要",
            "category": "技术",
            "is_published": True
        }

        response = self.make_request(
            "POST", "/news",
            data=news_data,
            expect_status=201,
            module="news",
            validate_schema=False
        )

        if response and "news" in response:
            news_id = response["news"].get("id")
            if news_id:
                self.teacher_news_id = news_id
                self.mark_for_cleanup('news', news_id)
                self.log_success(f"教师成功创建新闻，ID: {news_id}")
                return True
        return False

    def test_get_news_detail(self):
        """教师查看新闻详情（应该成功）"""
        if not self.teacher_news_id:
            self.log_warning("没有可查看的新闻，跳过测试")
            return True

        response = self.make_request(
            "GET", f"/news/{self.teacher_news_id}",
            expect_status=200,
            module="news",
            validate_schema=False
        )

        if response and "news" in response:
            news = response["news"]
            self.log_success(f"成功获取新闻详情: {news.get('title', 'N/A')}")
            return True
        elif response:
            self.log_success("成功获取新闻详情")
            return True
        return False

    def test_update_own_news(self):
        """教师更新自己的新闻（应该成功）"""
        if not self.teacher_news_id:
            self.log_warning("没有可更新的新闻，跳过测试")
            return True

        update_data = {
            "title": f"教师更新后的新闻_{self.generate_unique_name()}",
            "content": "这是更新后的新闻内容。教师可以更新自己创建的新闻，包括标题、内容、分类等信息。",
            "category": "公告"
        }

        response = self.make_request(
            "PUT", f"/news/{self.teacher_news_id}",
            data=update_data,
            expect_status=200,
            module="news",
            validate_schema=False
        )

        return response is not None

    def test_delete_own_news(self):
        """教师删除自己的新闻（应该成功）"""
        if not self.teacher_news_id:
            self.log_warning("没有可删除的新闻，跳过测试")
            return True

        response = self.make_request(
            "DELETE", f"/news/{self.teacher_news_id}",
            expect_status=200,
            module="news",
            validate_schema=False
        )

        if response:
            # 已删除，从清理列表中移除
            self.cleanup_items = [
                item for item in self.cleanup_items
                if not (item['type'] == 'news' and item['id'] == self.teacher_news_id)
            ]
            self.log_success(f"教师成功删除自己的新闻 {self.teacher_news_id}")
            self.teacher_news_id = None  # 清除ID
            return True
        return False

    # ========== 管理员角色测试 ==========

    def test_create_news_as_admin(self):
        """管理员创建新闻（应该成功）"""
        news_data = {
            "title": f"管理员测试新闻_{self.generate_unique_name()}",
            "content": "这是管理员创建的测试新闻内容。管理员可以创建、发布和管理所有新闻。",
            "summary": "管理员测试新闻摘要",
            "category": "公告",
            "is_published": True,
            "is_featured": True  # 设为推荐新闻
        }

        response = self.make_request(
            "POST", "/news",
            data=news_data,
            expect_status=201,
            module="news",
            validate_schema=False
        )

        if response and "news" in response:
            news_id = response["news"].get("id")
            if news_id:
                self.admin_news_id = news_id
                self.mark_for_cleanup('news', news_id)
                self.log_success(f"管理员成功创建新闻，ID: {news_id}")
                return True
        return False

    def test_admin_update_any_news(self):
        """管理员更新任意新闻（应该成功）"""
        # 先创建一个教师的新闻（如果还不存在）
        if not self.teacher_news_id:
            self.login_as('teacher')
            news_data = {
                "title": f"教师临时新闻_{self.generate_unique_name()}",
                "content": "用于测试管理员权限的新闻",
                "is_published": True
            }
            response = self.make_request(
                "POST", "/news",
                data=news_data,
                expect_status=201,
                module="news",
                validate_schema=False
            )
            if response and "news" in response:
                self.teacher_news_id = response["news"].get("id")
                self.mark_for_cleanup('news', self.teacher_news_id)

            # 切换回 admin
            self.login_as('admin')

        if not self.teacher_news_id:
            self.log_warning("无法创建测试新闻，跳过测试")
            return True

        # 管理员更新教师的新闻
        update_data = {
            "title": "管理员更新的新闻标题",
            "content": "管理员可以更新任意用户创建的新闻。"
        }

        response = self.make_request(
            "PUT", f"/news/{self.teacher_news_id}",
            data=update_data,
            expect_status=200,
            module="news",
            validate_schema=False
        )

        return response is not None

    def test_admin_delete_news(self):
        """管理员删除新闻（应该成功）"""
        # 删除管理员自己创建的新闻
        if not self.admin_news_id:
            self.log_warning("没有可删除的新闻，跳过测试")
            return True

        response = self.make_request(
            "DELETE", f"/news/{self.admin_news_id}",
            expect_status=200,
            module="news",
            validate_schema=False
        )

        if response:
            # 已删除，从清理列表中移除
            self.cleanup_items = [
                item for item in self.cleanup_items
                if not (item['type'] == 'news' and item['id'] == self.admin_news_id)
            ]
            self.log_success(f"管理员成功删除新闻 {self.admin_news_id}")
            return True
        return False

    # ========== 权限边界测试 ==========

    def test_unauthorized_create_news(self):
        """未登录创建新闻（应该被拒绝 401）"""
        self.clear_token()

        news_data = {
            "title": "未授权测试新闻",
            "content": "这应该失败",
            "is_published": True
        }

        response = self.make_request(
            "POST", "/news",
            data=news_data,
            expect_status=401,
            module="news",
            validate_schema=False
        )

        return self.assert_unauthorized(response)


def main():
    """主函数"""
    print("\n" + "=" * 60)
    print(" 新闻管理测试模块（重构版 - 基于 RBAC）")
    print("=" * 60 + "\n")

    tester = NewsTester()
    success = tester.run_tests()

    return 0 if success else 1


if __name__ == "__main__":
    import sys
    sys.exit(main())
