#!/usr/bin/env python3
"""
论坛管理测试模块
测试论坛帖子和回复功能
"""

from base_test import BaseAPITester
from colorama import Back


class ForumTester(BaseAPITester):
    """论坛管理测试类"""
    
    def __init__(self):
        super().__init__()
        self.post_id = None
        self.post_ids = []  # 存储多个帖子ID
        self.user_id = None
    
    def setup_auth(self):
        """设置认证（需要先登录获取token）"""
        # 创建测试用户并登录
        test_user = {
            "username": self.generate_unique_name("forumtester"),
            "email": f"{self.generate_unique_name('test')}@example.com",
            "password": "password123",
            "nickname": "论坛测试员"
        }
        
        # 注册用户
        reg_response = self.make_request("POST", "/users/register", data=test_user, expect_status=201)
        if not reg_response:
            return False
        
        self.user_id = reg_response.get("user", {}).get("id")
            
        # 登录获取token
        login_data = {
            "username": test_user["username"],
            "password": test_user["password"]
        }
        login_response = self.make_request("POST", "/users/login", data=login_data, expect_status=200)
        
        if login_response and "token" in login_response:
            self.set_token(login_response["token"])
            return True
        return False

    def test_list_forum_posts(self):
        """测试获取论坛帖子列表"""
        self.print_section_header("测试获取论坛帖子列表", Back.CYAN)
        
        response = self.make_request("GET", "/forum/posts", expect_status=200)
        
        if response is None:
            self.log_error("获取论坛帖子列表失败")
            return False
        else:
            self.log_success("获取论坛帖子列表成功")
            return True

    def test_create_forum_post(self):
        """测试创建论坛帖子"""
        self.print_section_header("测试创建论坛帖子", Back.CYAN)
        
        if not self.token:
            self.log_error("需要认证才能创建帖子")
            return False
        
        post_data = {
            "title": f"测试帖子_{self.generate_unique_name()}",
            "content": "这是一个测试帖子的内容，用于验证论坛功能。",
            "category": "技术讨论",
            "tags": ["Verilog", "FPGA", "测试"]
        }
        
        response = self.make_request("POST", "/forum/posts", data=post_data, expect_status=201)
        
        if response is None:
            self.log_error("创建论坛帖子失败")
            return False
        else:
            self.log_success("创建论坛帖子成功")
            if "post" in response and "id" in response["post"]:
                self.post_id = response["post"]["id"]
                self.post_ids.append(self.post_id)
                self.log_info(f"创建的帖子ID: {self.post_id}")
            return True

    def test_create_multiple_forum_posts(self):
        """测试创建多个论坛帖子"""
        self.print_section_header("测试创建多个论坛帖子", Back.CYAN)
        
        if not self.token:
            self.log_error("需要认证才能创建帖子")
            return False
        
        categories = ["技术讨论", "问题求助", "经验分享"]
        success_count = 0
        
        for i, category in enumerate(categories, 1):
            post_data = {
                "title": f"批量测试帖子_{i}_{category}_{self.generate_unique_name()}",
                "content": f"这是第{i}个{category}类型的测试帖子内容。包含了详细的技术讨论和问题描述。",
                "category": category,
                "tags": ["批量测试", category, "Verilog"]
            }
            
            response = self.make_request("POST", "/forum/posts", data=post_data, expect_status=201)
            
            if response and "post" in response and "id" in response["post"]:
                post_id = response["post"]["id"]
                self.post_ids.append(post_id)
                self.log_success(f"创建帖子{i}成功，ID: {post_id}")
                success_count += 1
            else:
                self.log_error(f"创建帖子{i}失败")
        
        self.log_info(f"批量创建完成，成功创建{success_count}个帖子，总共{len(self.post_ids)}个帖子")
        return success_count > 0
            
    def test_get_forum_post_detail(self):
        """测试获取论坛帖子详情"""
        self.print_section_header("测试获取论坛帖子详情", Back.CYAN)
        
        if not self.post_id:
            self.log_warning("没有可用的帖子ID，跳过获取详情测试")
            return True

        response = self.make_request("GET", f"/forum/posts/{self.post_id}", expect_status=200)
        
        if response is None:
            self.log_error(f"获取帖子 {self.post_id} 详情失败")
            return False
        else:
            self.log_success(f"获取帖子 {self.post_id} 详情成功")
            return True

    def test_update_forum_post(self):
        """测试更新论坛帖子"""
        self.print_section_header("测试更新论坛帖子", Back.CYAN)
        
        if not self.token or not self.post_id:
            self.log_warning("需要认证或没有帖子ID，跳过更新测试")
            return True
        
        update_data = {
            "title": f"更新后的帖子_{self.generate_unique_name()}",
            "content": "这是更新后的帖子内容。",
            "category": "经验分享"
        }
        
        response = self.make_request("PUT", f"/forum/posts/{self.post_id}", data=update_data, expect_status=200)
        
        if response is None:
            self.log_error(f"更新帖子 {self.post_id} 失败")
            return False
        else:
            self.log_success(f"更新帖子 {self.post_id} 成功")
            return True

    def test_list_forum_replies(self):
        """测试获取帖子回复列表"""
        self.print_section_header("测试获取帖子回复列表", Back.CYAN)
        
        if not self.post_id:
            self.log_warning("没有可用的帖子ID，跳过获取回复列表测试")
            return True
        
        response = self.make_request("GET", f"/forum/posts/{self.post_id}/replies", expect_status=200)
        
        if response is None:
            self.log_error(f"获取帖子 {self.post_id} 的回复列表失败")
            return False
        else:
            self.log_success(f"获取帖子 {self.post_id} 的回复列表成功")
            return True

    def test_create_forum_reply(self):
        """测试创建帖子回复"""
        self.print_section_header("测试创建帖子回复", Back.CYAN)
        
        if not self.token or not self.post_id:
            self.log_warning("需要认证或没有帖子ID，跳过创建回复测试")
            return True
        
        reply_data = {
            "content": "这是一个测试回复，用于验证论坛回复功能。",
            "reply_to": None
        }
        
        response = self.make_request("POST", f"/forum/posts/{self.post_id}/replies", data=reply_data, expect_status=201)
        
        if response is None:
            self.log_error(f"为帖子 {self.post_id} 创建回复失败")
            return False
        else:
            self.log_success(f"为帖子 {self.post_id} 创建回复成功")
            return True

    def test_delete_forum_post(self):
        """测试删除论坛帖子"""
        self.print_section_header("测试删除论坛帖子", Back.CYAN)
        
        if not self.token:
            self.log_warning("需要认证，跳过删除测试")
            return True
            
        # 使用剩余的帖子ID之一进行删除测试
        if not self.post_ids:
            self.log_warning("没有剩余的帖子ID，跳过删除测试")
            return True
            
        # 选择第一个剩余的帖子ID进行删除
        delete_post_id = self.post_ids[0]
        
        response = self.make_request("DELETE", f"/forum/posts/{delete_post_id}", expect_status=200)
        
        if response is None:
            self.log_error(f"删除帖子 {delete_post_id} 失败")
            return False
        else:
            self.log_success(f"删除帖子 {delete_post_id} 成功")
            self.post_ids.remove(delete_post_id)
            if self.post_id == delete_post_id:
                self.post_id = None
            return True

    def test_delete_partial_forum_posts(self):
        """测试删除部分论坛帖子（保留一些用于查看列表效果）"""
        self.print_section_header("测试删除部分论坛帖子", Back.CYAN)
        
        if not self.token:
            self.log_warning("需要认证才能删除帖子")
            return True
            
        if not self.post_ids:
            self.log_warning("没有帖子ID列表，跳过测试")
            return True
        
        # 只删除一半的帖子，保留其他的
        delete_count = len(self.post_ids) // 2
        if delete_count == 0:
            delete_count = 1  # 至少删除一个
        
        deleted_count = 0
        for i in range(delete_count):
            if i < len(self.post_ids):
                post_id = self.post_ids[i]
                response = self.make_request("DELETE", f"/forum/posts/{post_id}", expect_status=200)
                if response:
                    self.log_success(f"删除帖子 {post_id} 成功")
                    deleted_count += 1
                else:
                    self.log_error(f"删除帖子 {post_id} 失败")
        
        # 更新帖子ID列表，移除已删除的
        self.post_ids = self.post_ids[delete_count:]
        remaining_count = len(self.post_ids)
        
        self.log_info(f"删除了{deleted_count}个帖子，还剩余{remaining_count}个帖子用于列表展示")
        return deleted_count > 0
            
    def test_delete_non_existent_post(self):
        """测试删除不存在的帖子"""
        self.print_section_header("测试删除不存在的帖子", Back.CYAN)
        
        if not self.token:
            self.log_error("需要认证才能进行此项测试")
            return False
            
        response = self.make_request("DELETE", "/forum/posts/999999", expect_status=404)
        
        if response is None:
            self.log_error("删除不存在的帖子测试失败")
            return False
        else:
            self.log_success("删除不存在的帖子按预期失败 (404)")
            return True

    def test_unauthorized_forum_operations(self):
        """测试未授权论坛操作"""
        self.print_section_header("测试未授权论坛操作", Back.RED)
        
        old_token = self.token
        self.clear_token()
        
        post_data = {
            "title": "未授权测试帖子",
            "content": "这应该失败",
            "category": "测试"
        }
        
        response = self.make_request("POST", "/forum/posts", data=post_data, expect_status=401)
        
        self.set_token(old_token)
        
        if response is None:
            self.log_error("未授权论坛操作检查失败")
            return False
        else:
            self.log_success("未授权论坛操作正确被拒绝")
            return True

    def test_unauthorized_reply_operations(self):
        """测试未授权回复操作"""
        self.print_section_header("测试未授权回复操作", Back.RED)
        
        if not self.post_id:
            self.log_warning("没有帖子ID，跳过未授权回复测试")
            return True
        
        old_token = self.token
        self.clear_token()
        
        reply_data = {
            "content": "未授权回复测试"
        }
        
        response = self.make_request("POST", f"/forum/posts/{self.post_id}/replies", data=reply_data, expect_status=401)
        
        self.set_token(old_token)
        
        if response is None:
            self.log_error("未授权回复操作检查失败")
            return False
        else:
            self.log_success("未授权回复操作正确被拒绝")
            return True

    def test_reply_with_parent_id(self):
        """测试创建带父回复的回复"""
        self.print_section_header("测试创建带父回复的回复", Back.CYAN)
        
        if not self.post_id:
            self.log_warning("没有帖子ID，跳过父回复测试")
            return True
        
        # 先创建一个回复
        reply_data = {
            "content": "这是第一个回复"
        }
        first_reply = self.make_request("POST", f"/forum/posts/{self.post_id}/replies", data=reply_data, expect_status=201)
        
        if not first_reply or "reply" not in first_reply:
            self.log_error("创建第一个回复失败")
            return False
        
        first_reply_id = first_reply["reply"]["id"]
        
        # 创建对第一个回复的回复
        nested_reply_data = {
            "content": "这是对第一个回复的回复",
            "parent_id": first_reply_id
        }
        
        response = self.make_request("POST", f"/forum/posts/{self.post_id}/replies", data=nested_reply_data, expect_status=201)
        
        if response is None:
            self.log_error("创建嵌套回复失败")
            return False
        else:
            self.log_success("创建嵌套回复成功")
            return True

    def test_reply_pagination(self):
        """测试回复分页功能"""
        self.print_section_header("测试回复分页功能", Back.CYAN)
        
        if not self.post_id:
            self.log_warning("没有帖子ID，跳过分页测试")
            return True
        
        # 测试分页参数
        response = self.make_request("GET", f"/forum/posts/{self.post_id}/replies?page=1&limit=5", expect_status=200)
        
        if response is None:
            self.log_error("回复分页测试失败")
            return False
        else:
            self.log_success("回复分页功能正常")
            return True

    def run_tests(self):
        """运行论坛管理测试"""
        self.print_section_header("论坛管理测试模块", Back.BLUE)
        
        auth_setup = self.setup_auth()
        if not auth_setup:
            self.log_error("认证设置失败，跳过所有需要认证的测试")
        
        # 定义测试流程
        test_flow = [
            ("获取论坛帖子列表", self.test_list_forum_posts),
            ("未授权操作检查", self.test_unauthorized_forum_operations),
        ]
        
        if auth_setup:
            authed_tests = [
                ("创建论坛帖子", self.test_create_forum_post),
                ("批量创建多个论坛帖子", self.test_create_multiple_forum_posts),
                ("获取论坛帖子列表 (创建后)", self.test_list_forum_posts),
                ("获取帖子详情", self.test_get_forum_post_detail),
                ("更新论坛帖子", self.test_update_forum_post),
                ("创建帖子回复", self.test_create_forum_reply),
                ("获取帖子回复", self.test_list_forum_replies),
                ("测试回复分页", self.test_reply_pagination),
                ("测试嵌套回复", self.test_reply_with_parent_id),
                ("测试未授权回复", self.test_unauthorized_reply_operations),
                ("删除部分论坛帖子", self.test_delete_partial_forum_posts),
                ("获取论坛帖子列表 (部分删除后)", self.test_list_forum_posts),
                ("删除单个论坛帖子", self.test_delete_forum_post),
                ("获取论坛帖子列表 (最终)", self.test_list_forum_posts),
                ("删除不存在的帖子", self.test_delete_non_existent_post),
            ]
            test_flow.extend(authed_tests)

        test_results = []
        for name, test_func in test_flow:
            # 如果之前的测试创建了帖子，但后续测试需要它却不存在，则跳过
            if '帖子' in name and name not in ["获取论坛帖子列表", "创建论坛帖子", "未授权操作检查", "删除不存在的帖子"]:
                 if not self.post_id:
                    self.log_warning(f"因为没有帖子ID，跳过测试: {name}")
                    test_results.append((name, True)) # Mark as skipped but not failed
                    continue
            
            test_results.append((name, test_func()))

        return self.print_test_summary(test_results)


def main():
    """主函数"""
    tester = ForumTester()
    success = tester.run_tests()
    exit(0 if success else 1)


if __name__ == "__main__":
    main()