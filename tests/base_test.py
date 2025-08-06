#!/usr/bin/env python3
"""
Verilog OJ API 测试基类
提供公共的HTTP请求方法和工具函数
"""

import json
import time
import requests
from typing import Dict, Optional
from colorama import init, Fore, Back, Style

# 初始化colorama
init(autoreset=True)

# API基础URL
BASE_URL = "http://localhost:8080"
API_BASE = f"{BASE_URL}/api/v1"

class BaseAPITester:
    """API测试基类"""
    
    def __init__(self):
        self.session = requests.Session()
        self.token = None
        self.user_id = None
        
    def log_success(self, message: str):
        """成功日志"""
        print(f"{Fore.GREEN}✅ {message}{Style.RESET_ALL}")
        
    def log_error(self, message: str):
        """错误日志"""
        print(f"{Fore.RED}❌ {message}{Style.RESET_ALL}")
        
    def log_info(self, message: str):
        """信息日志"""
        print(f"{Fore.CYAN}ℹ️  {message}{Style.RESET_ALL}")
        
    def log_warning(self, message: str):
        """警告日志"""
        print(f"{Fore.YELLOW}⚠️  {message}{Style.RESET_ALL}")

    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None, 
                    headers: Optional[Dict] = None, expect_status: int = 200) -> Optional[Dict]:
        """发送HTTP请求"""
        url = f"{API_BASE}{endpoint}"
        
        # 设置默认headers
        default_headers = {"Content-Type": "application/json"}
        if self.token:
            default_headers["Authorization"] = f"Bearer {self.token}"
        
        if headers:
            default_headers.update(headers)
            
        try:
            self.log_info(f"{method.upper()} {url}")
            if data:
                self.log_info(f"请求数据: {json.dumps(data, ensure_ascii=False, indent=2)}")
                
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                headers=default_headers,
                timeout=10
            )
            
            self.log_info(f"响应状态: {response.status_code}")
            
            try:
                response_data = response.json()
                self.log_info(f"响应数据: {json.dumps(response_data, ensure_ascii=False, indent=2)}")
            except:
                response_data = {"raw_response": response.text}
                self.log_info(f"响应内容: {response.text}")
            
            if response.status_code == expect_status:
                self.log_success(f"请求成功 - {method.upper()} {endpoint}")
                return response_data
            else:
                self.log_error(f"请求失败 - 期望状态码 {expect_status}，实际 {response.status_code}")
                return None
                
        except requests.exceptions.RequestException as e:
            self.log_error(f"请求异常: {str(e)}")
            return None

    def set_token(self, token: str):
        """设置认证token"""
        self.token = token
        
    def clear_token(self):
        """清除认证token"""
        self.token = None
        
    def generate_unique_name(self, prefix: str = "test") -> str:
        """生成唯一名称"""
        timestamp = int(time.time())
        # 使用更短的时间戳，确保用户名不超过20个字符
        short_timestamp = timestamp % 1000  # 只取后3位数字
        return f"{prefix}{short_timestamp}"
        
    def print_section_header(self, title: str, color=Back.BLUE):
        """打印测试章节标题"""
        print(f"\n{color}{Fore.WHITE} {title} {Style.RESET_ALL}\n")
        
    def print_test_summary(self, test_results: list) -> bool:
        """打印测试结果总结"""
        print(f"\n{Back.CYAN}{Fore.BLACK} 📊 测试结果总结 📊 {Style.RESET_ALL}")
        print("=" * 60)
        
        passed = 0
        total = 0
        
        for test_name, result in test_results:
            total += 1
            if result:
                passed += 1
                status = f"{Fore.GREEN}✅ 通过{Style.RESET_ALL}"
            else:
                status = f"{Fore.RED}❌ 失败{Style.RESET_ALL}"
            print(f"{test_name:<20} {status}")
        
        print("=" * 60)
        print(f"总测试数: {total}")
        print(f"通过数: {Fore.GREEN}{passed}{Style.RESET_ALL}")
        print(f"失败数: {Fore.RED}{total - passed}{Style.RESET_ALL}")
        print(f"通过率: {Fore.CYAN}{passed/total*100:.1f}%{Style.RESET_ALL}")
        
        if passed == total:
            print(f"\n{Fore.GREEN}🎉 所有测试通过！{Style.RESET_ALL}")
            return True
        else:
            print(f"\n{Fore.YELLOW}⚠️  部分测试失败，请检查实现{Style.RESET_ALL}")
            return False 