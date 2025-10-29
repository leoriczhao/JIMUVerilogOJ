#!/usr/bin/env python3
"""
论坛管理测试模块（重构版）
展示使用新的 RBAC 测试框架进行测试的最佳实践
"""

from base_test import BaseAPITester
from colorama import Back


class ForumTester(BaseAPITester):
    """论坛管理测试类（基于 RBAC）"""

    def __init__(self):
        super().__init__()
        # 初始化用户池
        self.setup_user_pool()
        # 存储创建的帖子ID
        self.student_post_id = None
        self.teacher_post_id = None
        self.reply_id = None

    def run_tests(self):
        """主测试流程"""
        test_results = []

        try:
            # 1. 公开接口测试（无需登录）
            self.print_section_header("公开接口测试", Back.CYAN)
            test_results.append(("获取帖子列表(公开)", self.test_list_posts_public()))

            # 2. 学生角色测试（创建、查看、回复）
            self.print_section_header("学生角色测试", Back.BLUE)
            self.login_as('student')
            test_results.append(("学生-创建帖子", self.test_create_post_as_student()))
            test_results.append(("学生-查看帖子详情", self.test_get_post_detail()))
            test_results.append(("学生-创建回复", self.test_create_reply_as_student()))
            test_results.append(("学生-查看回复列表", self.test_list_replies()))
            test_results.append(("学生-删除其他人帖子(应拒绝)", self.test_student_delete_others_post()))

            # 3. 教师角色测试（创建、更新、删除自己的帖子）
            self.print_section_header("教师角色测试", Back.GREEN)
            self.login_as('teacher')
            test_results.append(("教师-创建帖子", self.test_create_post_as_teacher()))
            test_results.append(("教师-更新自己的帖子", self.test_update_own_post()))
            test_results.append(("教师-创建回复", self.test_create_reply_as_teacher()))

            # 4. 管理员角色测试（管理所有帖子、删除、锁定、置顶）
            self.print_section_header("管理员角色测试", Back.MAGENTA)
            self.login_as('admin')
            test_results.append(("管理员-更新任意帖子", self.test_admin_update_any_post()))
            test_results.append(("管理员-删除帖子", self.test_admin_delete_post()))

            # 5. 权限边界测试
            self.print_section_header("权限边界测试", Back.RED)
            test_results.append(("未登录-创建帖子(应拒绝)", self.test_unauthorized_create_post()))
            test_results.append(("未登录-创建回复(应拒绝)", self.test_unauthorized_create_reply()))

        finally:
            # 6. 清理测试数据
            self.cleanup()

        # 7. 打印测试报告
        return self.print_test_summary(test_results)

    # ========== 公开接口测试 ==========

    def test_list_posts_public(self):
        """测试公开获取帖子列表（无需登录）"""
        self.clear_token()  # 确保无认证状态

        response = self.make_request(
            "GET", "/forum/posts",
            expect_status=200,
            module="forum",
        )

        if response and "posts" in response:
            total = response.get("total", 0)
            posts = response.get("posts") or []  # 处理 posts 为 None 的情况
            self.log_success(f"获取到 {total} 个帖子，当前页 {len(posts)} 个")
            return True
        elif response:
            self.log_success("成功获取帖子列表")
            return True
        return False

    # ========== 学生角色测试 ==========

    def test_create_post_as_student(self):
        """学生创建帖子（应该成功）"""
        post_data = {
            "title": f"学生测试帖子_{self.generate_unique_name()}",
            "content": "这是学生创建的测试帖子内容。",
            "category": "技术讨论",
            "tags": ["Verilog", "测试"]
        }

        response = self.make_request(
            "POST", "/forum/posts",
            data=post_data,
            expect_status=201,
            module="forum",
        )

        if response and "post" in response:
            post_id = response["post"].get("id")
            if post_id:
                self.student_post_id = post_id
                self.mark_for_cleanup('post', post_id)
                self.log_success(f"学生成功创建帖子，ID: {post_id}")
                return True
        return False

    def test_get_post_detail(self):
        """学生查看帖子详情（应该成功）"""
        if not self.student_post_id:
            self.log_warning("没有可查看的帖子，跳过测试")
            return True

        response = self.make_request(
            "GET", f"/forum/posts/{self.student_post_id}",
            expect_status=200,
            module="forum",
        )

        if response and "post" in response:
            post = response["post"]
            self.log_success(f"成功获取帖子详情: {post.get('title', 'N/A')}")
            return True
        elif response:
            self.log_success("成功获取帖子详情")
            return True
        return False

    def test_create_reply_as_student(self):
        """学生创建回复（应该成功）"""
        if not self.student_post_id:
            self.log_warning("没有可回复的帖子，跳过测试")
            return True

        reply_data = {
            "content": "这是学生的测试回复。"
        }

        response = self.make_request(
            "POST", f"/forum/posts/{self.student_post_id}/replies",
            data=reply_data,
            expect_status=201,
            module="forum",
        )

        if response and "reply" in response:
            reply_id = response["reply"].get("id")
            if reply_id:
                self.reply_id = reply_id
                self.log_success(f"学生成功创建回复，ID: {reply_id}")
                return True
        return False

    def test_list_replies(self):
        """学生查看回复列表（应该成功）"""
        if not self.student_post_id:
            self.log_warning("没有可查看回复的帖子，跳过测试")
            return True

        response = self.make_request(
            "GET", f"/forum/posts/{self.student_post_id}/replies",
            expect_status=200,
            module="forum",
        )

        if response and "replies" in response:
            replies = response.get("replies", [])
            self.log_success(f"成功获取 {len(replies)} 条回复")
            return True
        elif response:
            self.log_success("成功获取回复列表")
            return True
        return False

    def test_student_delete_others_post(self):
        """学生删除其他人的帖子（应该被拒绝）"""
        # 使用 teacher 创建的帖子ID（如果存在）
        # 如果还没有 teacher 帖子，先切换到 teacher 创建一个
        if not self.teacher_post_id:
            self.login_as('teacher')
            post_data = {
                "title": f"教师临时帖子_{self.generate_unique_name()}",
                "content": "这是一个用于测试学生权限的临时帖子",  # 至少10个字符
                "category": "测试"
            }
            response = self.make_request(
                "POST", "/forum/posts",
                data=post_data,
                expect_status=201,
                module="forum",
            )
            if response and "post" in response:
                self.teacher_post_id = response["post"].get("id")
                self.mark_for_cleanup('post', self.teacher_post_id)

            # 切换回 student
            self.login_as('student')

        if not self.teacher_post_id:
            self.log_warning("无法创建测试帖子，跳过测试")
            return True

        # 学生尝试删除教师的帖子
        response = self.make_request(
            "DELETE", f"/forum/posts/{self.teacher_post_id}",
            expect_status=403,
            module="forum",
        )

        return self.assert_forbidden(response)

    # ========== 教师角色测试 ==========

    def test_create_post_as_teacher(self):
        """教师创建帖子（应该成功）"""
        post_data = {
            "title": f"教师测试帖子_{self.generate_unique_name()}",
            "content": "这是教师创建的测试帖子内容。",
            "category": "经验分享",
            "tags": ["FPGA", "教学"]
        }

        response = self.make_request(
            "POST", "/forum/posts",
            data=post_data,
            expect_status=201,
            module="forum",
        )

        if response and "post" in response:
            post_id = response["post"].get("id")
            if post_id:
                self.teacher_post_id = post_id
                self.mark_for_cleanup('post', post_id)
                self.log_success(f"教师成功创建帖子，ID: {post_id}")
                return True
        return False

    def test_update_own_post(self):
        """教师更新自己的帖子（应该成功）"""
        if not self.teacher_post_id:
            self.log_warning("没有可更新的帖子，跳过测试")
            return True

        update_data = {
            "title": f"教师更新后的帖子_{self.generate_unique_name()}",
            "content": "这是更新后的内容。教师可以更新自己创建的帖子，包括标题、内容和分类等信息。这段内容足够长以满足后端验证要求。",
            "category": "技术讨论"
        }

        response = self.make_request(
            "PUT", f"/forum/posts/{self.teacher_post_id}",
            data=update_data,
            expect_status=200,
            module="forum",
        )

        return response is not None

    def test_create_reply_as_teacher(self):
        """教师创建回复（应该成功）"""
        # 教师可以回复自己或他人的帖子，这里回复学生的帖子
        if not self.student_post_id:
            self.log_warning("没有可回复的帖子，跳过测试")
            return True

        reply_data = {
            "content": "这是教师的测试回复。"
        }

        response = self.make_request(
            "POST", f"/forum/posts/{self.student_post_id}/replies",
            data=reply_data,
            expect_status=201,
            module="forum",
        )

        if response and "reply" in response:
            self.log_success("教师成功创建回复")
            return True
        return False

    # ========== 管理员角色测试 ==========

    def test_admin_update_any_post(self):
        """管理员更新任意帖子（应该成功）"""
        # 更新学生的帖子
        if not self.student_post_id:
            self.log_warning("没有可更新的帖子，跳过测试")
            return True

        update_data = {
            "title": "管理员更新的帖子标题",
            "content": "管理员可以更新任意用户的帖子。"
        }

        response = self.make_request(
            "PUT", f"/forum/posts/{self.student_post_id}",
            data=update_data,
            expect_status=200,
            module="forum",
        )

        return response is not None

    def test_admin_delete_post(self):
        """管理员删除帖子（应该成功）"""
        # 删除教师的帖子
        if not self.teacher_post_id:
            self.log_warning("没有可删除的帖子，跳过测试")
            return True

        response = self.make_request(
            "DELETE", f"/forum/posts/{self.teacher_post_id}",
            expect_status=200,
            module="forum",
        )

        if response:
            # 已删除，从清理列表中移除
            self.cleanup_items = [
                item for item in self.cleanup_items
                if not (item['type'] == 'post' and item['id'] == self.teacher_post_id)
            ]
            self.log_success(f"管理员成功删除帖子 {self.teacher_post_id}")
            return True
        return False

    # ========== 权限边界测试 ==========

    def test_unauthorized_create_post(self):
        """未登录创建帖子（应该被拒绝 401）"""
        self.clear_token()

        post_data = {
            "title": "未授权测试帖子",
            "content": "这应该失败",
            "category": "测试"
        }

        response = self.make_request(
            "POST", "/forum/posts",
            data=post_data,
            expect_status=401,
            module="forum",
        )

        return self.assert_unauthorized(response)

    def test_unauthorized_create_reply(self):
        """未登录创建回复（应该被拒绝 401）"""
        if not self.student_post_id:
            self.log_warning("没有可回复的帖子，跳过测试")
            return True

        self.clear_token()

        reply_data = {
            "content": "未授权测试回复"
        }

        response = self.make_request(
            "POST", f"/forum/posts/{self.student_post_id}/replies",
            data=reply_data,
            expect_status=401,
            module="forum",
        )

        return self.assert_unauthorized(response)


def main():
    """主函数"""
    print("\n" + "=" * 60)
    print(" 论坛管理测试模块（重构版 - 基于 RBAC）")
    print("=" * 60 + "\n")

    tester = ForumTester()
    success = tester.run_tests()

    return 0 if success else 1


if __name__ == "__main__":
    import sys
    sys.exit(main())
