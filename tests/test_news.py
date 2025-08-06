#!/usr/bin/env python3
"""
新闻管理测试模块
测试新闻的增删改查功能
"""

from base_test import BaseAPITester
from colorama import Back


class NewsTester(BaseAPITester):
    """新闻管理测试类"""
    
    def __init__(self):
        super().__init__()
        self.news_id = None
        self.news_ids = []  # 存储多个新闻ID

    def setup_auth(self):
        """设置认证（需要管理员权限）"""
        admin_user = {
            "username": self.generate_unique_name("news_admin"),
            "email": f"{self.generate_unique_name('news_admin')}@example.com",
            "password": "password123",
            "nickname": "News Admin",
            "role": "admin"
        }
        reg_res = self.make_request("POST", "/users/register", data=admin_user, expect_status=201)
        if not reg_res:
            return False
        login_res = self.make_request("POST", "/users/login", data={"username": admin_user["username"], "password": admin_user["password"]})
        if not login_res or "token" not in login_res:
            return False
        self.set_token(login_res["token"])
        return True

    def test_create_news(self):
        """测试创建新闻"""
        self.print_section_header("测试创建新闻", Back.LIGHTGREEN_EX)
        news_data = {
            "title": f"测试新闻_{self.generate_unique_name()}",
            "content": "这是一条测试新闻的内容。",
            "is_published": True
        }
        response = self.make_request("POST", "/news", data=news_data, expect_status=201)
        if response and "news" in response and "id" in response["news"]:
            self.news_id = response["news"]["id"]
            self.news_ids.append(self.news_id)
            self.log_success(f"创建新闻成功，ID: {self.news_id}")
            return True
        self.log_error("创建新闻失败")
        return False

    def test_create_multiple_news(self):
        """测试创建多个新闻"""
        self.print_section_header("测试创建多个新闻", Back.LIGHTGREEN_EX)
        categories = ["技术", "公告", "活动"]
        success_count = 0
        
        for i, category in enumerate(categories, 1):
            news_data = {
                "title": f"批量测试新闻_{i}_{category}_{self.generate_unique_name()}",
                "content": f"这是第{i}条{category}类型的测试新闻内容。包含了丰富的信息和详细的描述。",
                "summary": f"第{i}条{category}新闻摘要",
                "category": category,
                "is_published": True,
                "is_featured": i == 1  # 第一条设为推荐
            }
            response = self.make_request("POST", "/news", data=news_data, expect_status=201)
            if response and "news" in response and "id" in response["news"]:
                news_id = response["news"]["id"]
                self.news_ids.append(news_id)
                self.log_success(f"创建新闻{i}成功，ID: {news_id}")
                success_count += 1
            else:
                self.log_error(f"创建新闻{i}失败")
        
        self.log_info(f"批量创建完成，成功创建{success_count}条新闻，总共{len(self.news_ids)}条新闻")
        return success_count > 0

    def test_get_news_detail(self):
        """测试获取新闻详情"""
        self.print_section_header("测试获取新闻详情", Back.LIGHTGREEN_EX)
        if not self.news_id:
            self.log_warning("没有新闻ID，跳过测试")
            return True
        response = self.make_request("GET", f"/news/{self.news_id}", expect_status=200)
        return response is not None

    def test_list_news(self):
        """测试获取新闻列表"""
        self.print_section_header("测试获取新闻列表", Back.LIGHTGREEN_EX)
        response = self.make_request("GET", "/news", expect_status=200)
        return response is not None

    def test_update_news(self):
        """测试更新新闻"""
        self.print_section_header("测试更新新闻", Back.LIGHTGREEN_EX)
        if not self.news_id:
            self.log_warning("没有新闻ID，跳过测试")
            return True
        update_data = {"title": f"更新后的新闻_{self.generate_unique_name()}"}
        response = self.make_request("PUT", f"/news/{self.news_id}", data=update_data, expect_status=200)
        return response is not None

    def test_delete_news(self):
        """测试删除新闻"""
        self.print_section_header("测试删除新闻", Back.LIGHTGREEN_EX)
        if not self.news_ids:
            self.log_warning("没有剩余的新闻ID，跳过测试")
            return True
        
        # 使用剩余ID列表中的第一个ID进行删除
        news_id = self.news_ids[0]
        response = self.make_request("DELETE", f"/news/{news_id}", expect_status=200)
        if response:
            self.log_success(f"删除新闻 {news_id} 成功")
            self.news_ids.remove(news_id)
            if self.news_id == news_id:
                self.news_id = None
            return True
        return False

    def test_delete_partial_news(self):
        """测试删除部分新闻（保留一些用于查看列表效果）"""
        self.print_section_header("测试删除部分新闻", Back.LIGHTGREEN_EX)
        if not self.news_ids:
            self.log_warning("没有新闻ID列表，跳过测试")
            return True
        
        # 只删除一半的新闻，保留其他的
        delete_count = len(self.news_ids) // 2
        if delete_count == 0:
            delete_count = 1  # 至少删除一个
        
        deleted_count = 0
        for i in range(delete_count):
            if i < len(self.news_ids):
                news_id = self.news_ids[i]
                response = self.make_request("DELETE", f"/news/{news_id}", expect_status=200)
                if response:
                    self.log_success(f"删除新闻 {news_id} 成功")
                    deleted_count += 1
                else:
                    self.log_error(f"删除新闻 {news_id} 失败")
        
        # 更新新闻ID列表，移除已删除的
        self.news_ids = self.news_ids[delete_count:]
        remaining_count = len(self.news_ids)
        
        self.log_info(f"删除了{deleted_count}条新闻，还剩余{remaining_count}条新闻用于列表展示")
        return deleted_count > 0

    def test_unauthorized_news_operations(self):
        """测试未授权新闻操作"""
        self.print_section_header("测试未授权新闻操作", Back.RED)
        
        old_token = self.token
        self.clear_token()
        
        news_data = {
            "title": "未授权测试新闻",
            "content": "这应该失败"
        }
        
        response = self.make_request("POST", "/news", data=news_data, expect_status=401)
        
        self.set_token(old_token)
        
        if response is None:
            self.log_error("未授权新闻操作检查失败")
            return False
        else:
            self.log_success("未授权新闻操作正确被拒绝")
            return True

    def test_news_pagination(self):
        """测试新闻分页功能"""
        self.print_section_header("测试新闻分页功能", Back.LIGHTGREEN_EX)
        
        response = self.make_request("GET", "/news?page=1&limit=5", expect_status=200)
        
        if response is None:
            self.log_error("新闻分页测试失败")
            return False
        else:
            self.log_success("新闻分页功能正常")
            return True

    def test_create_draft_news(self):
        """测试创建草稿新闻"""
        self.print_section_header("测试创建草稿新闻", Back.LIGHTGREEN_EX)
        
        draft_data = {
            "title": f"草稿新闻_{self.generate_unique_name()}",
            "content": "这是一条草稿新闻的内容。",
            "is_published": False
        }
        
        response = self.make_request("POST", "/news", data=draft_data, expect_status=201)
        
        if response and "news" in response and "id" in response["news"]:
            draft_id = response["news"]["id"]
            self.log_success(f"创建草稿新闻成功，ID: {draft_id}")
            
            # 清理：删除草稿新闻
            self.make_request("DELETE", f"/news/{draft_id}", expect_status=200)
            return True
        else:
            self.log_error("创建草稿新闻失败")
            return False

    def test_update_news_publish_status(self):
        """测试更新新闻发布状态"""
        self.print_section_header("测试更新新闻发布状态", Back.LIGHTGREEN_EX)
        
        if not self.news_id:
            self.log_warning("没有新闻ID，跳过发布状态测试")
            return True
        
        # 先设置为未发布
        update_data = {"is_published": False}
        response = self.make_request("PUT", f"/news/{self.news_id}", data=update_data, expect_status=200)
        
        if response is None:
            self.log_error("更新新闻发布状态失败")
            return False
        
        # 再设置为已发布
        update_data = {"is_published": True}
        response = self.make_request("PUT", f"/news/{self.news_id}", data=update_data, expect_status=200)
        
        if response is None:
            self.log_error("更新新闻发布状态失败")
            return False
        else:
            self.log_success("更新新闻发布状态成功")
            return True

    def run_tests(self):
        """运行所有新闻模块测试"""
        self.print_section_header("新闻管理测试模块", Back.BLUE)

        if not self.setup_auth():
            self.log_error("管理员认证失败，测试终止")
            return False

        test_flow = [
            ("获取新闻列表 (初始)", self.test_list_news),
            ("测试新闻分页", self.test_news_pagination),
            ("测试未授权操作", self.test_unauthorized_news_operations),
            ("创建新闻", self.test_create_news),
            ("批量创建多个新闻", self.test_create_multiple_news),
            ("获取新闻列表 (创建后)", self.test_list_news),
            ("获取新闻详情", self.test_get_news_detail),
            ("更新新闻", self.test_update_news),
            ("测试更新发布状态", self.test_update_news_publish_status),
            ("测试创建草稿", self.test_create_draft_news),
            ("删除部分新闻", self.test_delete_partial_news),
            ("获取新闻列表 (部分删除后)", self.test_list_news),
            ("删除单个新闻", self.test_delete_news),
            ("获取新闻列表 (最终)", self.test_list_news),
        ]

        test_results = []
        for name, test_func in test_flow:
            test_results.append((name, test_func()))

        return self.print_test_summary(test_results)

def main():
    """主函数"""
    tester = NewsTester()
    success = tester.run_tests()
    exit(0 if success else 1)


if __name__ == "__main__":
    main()
