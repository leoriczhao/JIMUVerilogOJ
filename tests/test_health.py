#!/usr/bin/env python3
"""
健康检查测试模块
测试API服务的基本健康状态
"""
       
from base_test import BaseAPITester, BASE_URL
from colorama import Back
import requests


class HealthTester(BaseAPITester):
    """健康检查测试类"""
    
    def test_health_check(self):
        """测试健康检查接口"""
        self.print_section_header("测试健康检查接口", Back.BLUE)
        
        try:
            # 测试健康检查接口
            health_url = f"{BASE_URL}/health"
            self.log_info(f"GET {health_url}")
            response = self.session.get(health_url, timeout=10)
            
            if response.status_code == 200:
                data = response.json()
                self.log_success("健康检查接口正常")
                self.log_info(f"响应: {data}")
                return True
            else:
                self.log_error(f"健康检查失败: {response.status_code}")
                return False
        except Exception as e:
            self.log_error(f"健康检查异常: {str(e)}")
            return False

    def test_api_root(self):
        """测试API根路径"""
        self.log_info("测试API根路径")
        response = self.make_request("GET", "", expect_status=404)
        # 根路径返回404是正常的
        return response is not None

    def run_tests(self):
        """运行健康检查测试"""
        print("🏥 开始健康检查测试")
        print("=" * 50)
        
        test_results = [
            ("健康检查接口", self.test_health_check()),
            ("API根路径", self.test_api_root()),
        ]
        
        return self.print_test_summary(test_results)


def main():
    """主函数"""
    tester = HealthTester()
    success = tester.run_tests()
    exit(0 if success else 1)


if __name__ == "__main__":
    main() 