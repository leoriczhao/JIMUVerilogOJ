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
from openapi_validator import get_validator

# 初始化colorama
init(autoreset=True)

# API基础URL
BASE_URL = "http://localhost:8080"
API_BASE = f"{BASE_URL}/api/v1"

class BaseAPITester:
    """API测试基类"""

    def __init__(self, enable_schema_validation: bool = True):
        self.session = requests.Session()
        self.token = None
        self.user_id = None
        self.enable_schema_validation = enable_schema_validation
        self.validator = get_validator() if enable_schema_validation else None
        self.validation_errors = []  # 记录所有验证错误
        
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
                    headers: Optional[Dict] = None, expect_status: int = 200,
                    module: Optional[str] = None, validate_schema: Optional[bool] = None) -> Optional[Dict]:
        """
        发送HTTP请求

        Args:
            method: HTTP方法
            endpoint: API端点
            data: 请求数据
            headers: 请求头
            expect_status: 期望的状态码
            module: OpenAPI模块名（如'user', 'problem'），用于schema验证
            validate_schema: 是否验证schema，默认使用实例设置

        Returns:
            响应数据，如果请求失败返回None
        """
        # 添加请求间延迟以避免速率限制
        time.sleep(0.1)
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

            # 执行schema验证
            should_validate = validate_schema if validate_schema is not None else self.enable_schema_validation
            if should_validate and self.validator and module and isinstance(response_data, dict):
                self._validate_response_schema(module, method, endpoint, response.status_code, response_data)

            if response.status_code == expect_status:
                self.log_success(f"请求成功 - {method.upper()} {endpoint}")
                return response_data
            else:
                self.log_error(f"请求失败 - 期望状态码 {expect_status}，实际 {response.status_code}")
                return None

        except requests.exceptions.RequestException as e:
            self.log_error(f"请求异常: {str(e)}")
            return None

    def _validate_response_schema(self, module: str, method: str, endpoint: str,
                                  status_code: int, response_data: Dict):
        """
        验证响应schema

        Args:
            module: OpenAPI模块名
            method: HTTP方法
            endpoint: API端点
            status_code: 状态码
            response_data: 响应数据
        """
        is_valid, error_msg = self.validator.validate_response(
            module, method, endpoint, status_code, response_data
        )

        if is_valid:
            if error_msg:  # 如果有提示信息（如schema未找到）
                self.log_warning(f"Schema验证跳过: {error_msg}")
            else:
                self.log_success("✓ Schema验证通过")
        else:
            self.log_error(f"✗ Schema验证失败:\n{error_msg}")
            self.validation_errors.append({
                'module': module,
                'method': method,
                'endpoint': endpoint,
                'status_code': status_code,
                'error': error_msg
            })

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

        # 打印Schema验证错误摘要
        if self.validation_errors:
            print(f"\n{Back.RED}{Fore.WHITE} ⚠️  Schema验证错误 ({len(self.validation_errors)}) ⚠️  {Style.RESET_ALL}")
            for idx, error in enumerate(self.validation_errors, 1):
                print(f"\n{Fore.RED}错误 #{idx}:{Style.RESET_ALL}")
                print(f"  模块: {error['module']}")
                print(f"  请求: {error['method']} {error['endpoint']}")
                print(f"  状态码: {error['status_code']}")
                print(f"  详情: {error['error']}")

        if passed == total:
            if self.validation_errors:
                print(f"\n{Fore.YELLOW}⚠️  功能测试通过，但存在Schema验证错误{Style.RESET_ALL}")
                return False
            else:
                print(f"\n{Fore.GREEN}🎉 所有测试通过！{Style.RESET_ALL}")
                return True
        else:
            print(f"\n{Fore.YELLOW}⚠️  部分测试失败，请检查实现{Style.RESET_ALL}")
            return False

    def get_validation_errors(self) -> list:
        """获取所有schema验证错误"""
        return self.validation_errors

    def clear_validation_errors(self):
        """清除schema验证错误记录"""
        self.validation_errors = [] 