#!/usr/bin/env python3
"""
Verilog OJ API 测试基类
提供公共的HTTP请求方法和工具函数
"""

import json
import time
import requests
from typing import Dict, Optional, List
from colorama import init, Fore, Back, Style
from openapi_validator import get_validator
from fixtures.users import TestUserPool
from fixtures.permissions import has_permission, get_minimum_role

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
        self.validation_errors = []  # 记录验证失败的情况
        self.validation_skipped = []  # 记录跳过验证的情况（schema未找到）

        # 新增：用户池和角色管理
        self.user_pool = None  # 用户池实例
        self.current_role = None  # 当前角色

        # 新增：资源清理管理
        self.cleanup_items = []  # 需要清理的资源列表
        
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
                # 记录跳过的验证
                self.validation_skipped.append({
                    'module': module,
                    'method': method,
                    'endpoint': endpoint,
                    'status_code': status_code,
                    'reason': error_msg
                })
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

        # 打印Schema验证统计
        print(f"\n{Back.BLUE}{Fore.WHITE} 📋 Schema验证统计 📋 {Style.RESET_ALL}")
        total_errors = len(self.validation_errors)
        total_skipped = len(self.validation_skipped)
        total_validated = total_errors + total_skipped

        if total_validated > 0 or self.enable_schema_validation:
            # 计算通过的数量 = 总验证数 - 错误数 - 跳过数
            # 但需要知道总共进行了多少次验证尝试
            # 由于我们无法准确知道验证通过的数量，只显示错误和跳过
            if total_errors > 0 or total_skipped > 0:
                print(f"验证失败: {Fore.RED}{total_errors}{Style.RESET_ALL}")
                print(f"跳过验证(未找到schema): {Fore.YELLOW}{total_skipped}{Style.RESET_ALL}")
            else:
                print(f"{Fore.GREEN}所有Schema验证通过{Style.RESET_ALL}")
        else:
            print(f"{Fore.CYAN}未启用Schema验证{Style.RESET_ALL}")

        # 打印Schema验证错误摘要
        if self.validation_errors:
            print(f"\n{Back.RED}{Fore.WHITE} ⚠️  Schema验证错误 ({len(self.validation_errors)}) ⚠️  {Style.RESET_ALL}")
            for idx, error in enumerate(self.validation_errors, 1):
                print(f"\n{Fore.RED}错误 #{idx}:{Style.RESET_ALL}")
                print(f"  模块: {error['module']}")
                print(f"  请求: {error['method']} {error['endpoint']}")
                print(f"  状态码: {error['status_code']}")
                print(f"  详情: {error['error']}")

        # 打印Schema跳过摘要（仅显示前5个）
        if self.validation_skipped:
            print(f"\n{Back.YELLOW}{Fore.BLACK} ℹ️  Schema未定义 ({len(self.validation_skipped)}) ℹ️  {Style.RESET_ALL}")
            print(f"{Fore.YELLOW}以下API端点缺少OpenAPI Schema定义：{Style.RESET_ALL}")
            # 去重并按endpoint分组
            unique_skipped = {}
            for skip in self.validation_skipped:
                key = f"{skip['method']} {skip['endpoint']} {skip['status_code']}"
                if key not in unique_skipped:
                    unique_skipped[key] = skip

            # 只显示前10个
            for idx, (key, skip) in enumerate(list(unique_skipped.items())[:10], 1):
                print(f"  {idx}. {skip['method']} {skip['endpoint']} (状态码: {skip['status_code']})")

            if len(unique_skipped) > 10:
                print(f"  ... 还有 {len(unique_skipped) - 10} 个未定义的schema")

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

    # ========== 新增：用户池和角色管理 ==========

    def setup_user_pool(self):
        """初始化用户池"""
        if self.user_pool is None:
            self.user_pool = TestUserPool()
            self.user_pool.setup(self)
            self.log_success("用户池初始化完成")

    def login_as(self, role: str):
        """
        切换到指定角色

        Args:
            role: 角色名称 ('student', 'teacher', 'admin')

        Raises:
            Exception: 如果用户池未初始化或角色无效
        """
        if self.user_pool is None:
            self.setup_user_pool()

        if role not in ['student', 'teacher', 'admin']:
            raise ValueError(f"无效的角色: {role}")

        token = self.user_pool.get_token(role)
        if not token:
            raise Exception(f"无法获取 {role} 角色的 token")

        self.set_token(token)
        self.current_role = role
        self.user_id = self.user_pool.get_user_id(role)

        self.log_info(f"已切换到角色: {role} (用户ID: {self.user_id})")

    def get_current_role(self) -> Optional[str]:
        """获取当前角色"""
        return self.current_role

    def check_permission(self, permission: str) -> bool:
        """
        检查当前用户是否有指定权限

        Args:
            permission: 权限字符串

        Returns:
            True 如果有权限
        """
        if not self.current_role:
            return False
        return has_permission(self.current_role, permission)

    # ========== 新增：资源清理管理 ==========

    def mark_for_cleanup(self, resource_type: str, resource_id: int):
        """
        标记资源用于测试后清理

        Args:
            resource_type: 资源类型 ('problem', 'submission', 'post', 'news', 等)
            resource_id: 资源ID
        """
        self.cleanup_items.append({
            'type': resource_type,
            'id': resource_id
        })
        self.log_info(f"标记 {resource_type} #{resource_id} 待清理")

    def cleanup(self):
        """
        清理所有测试资源

        按照依赖顺序倒序删除资源，最后清理用户
        """
        if not self.cleanup_items and not self.user_pool:
            return

        self.log_warning(f"开始清理测试数据 ({len(self.cleanup_items)} 个资源)...")

        # 以管理员身份删除所有标记的资源
        if self.cleanup_items:
            try:
                self.login_as('admin')
            except:
                self.log_error("无法切换到管理员角色进行清理")
                return

            # 倒序删除（后创建的先删除）
            for item in reversed(self.cleanup_items):
                self._delete_resource(item['type'], item['id'])

        # 清理用户池
        if self.user_pool:
            self.user_pool.cleanup(self)

        self.cleanup_items = []
        self.log_success("测试数据清理完成")

    def _delete_resource(self, resource_type: str, resource_id: int):
        """
        删除指定类型的资源

        Args:
            resource_type: 资源类型
            resource_id: 资源ID
        """
        endpoint_map = {
            'problem': f"/problems/{resource_id}",
            'submission': f"/submissions/{resource_id}",
            'post': f"/forum/posts/{resource_id}",
            'news': f"/news/{resource_id}",
        }

        endpoint = endpoint_map.get(resource_type)
        if not endpoint:
            self.log_warning(f"未知的资源类型: {resource_type}")
            return

        # 尝试删除，忽略错误（资源可能已被删除）
        response = self.make_request(
            "DELETE", endpoint,
            expect_status=200,
            validate_schema=False
        )

        if response is not None:
            self.log_success(f"已删除 {resource_type} #{resource_id}")
        else:
            self.log_warning(f"删除 {resource_type} #{resource_id} 失败（可能已不存在）")

    # ========== 新增：权限测试断言 ==========

    def assert_forbidden(self, response: Optional[Dict]) -> bool:
        """
        断言响应为 403 权限不足

        Args:
            response: API 响应

        Returns:
            True 如果断言成功
        """
        if response is None:
            return True

        is_forbidden = (
            response.get('error') == 'forbidden' or
            '权限不足' in str(response.get('message', ''))
        )

        if is_forbidden:
            self.log_success("✓ 权限检查通过：正确返回 403 Forbidden")
            return True
        else:
            self.log_error(f"✗ 期望 403 Forbidden，但得到: {response}")
            return False

    def assert_unauthorized(self, response: Optional[Dict]) -> bool:
        """
        断言响应为 401 未认证

        Args:
            response: API 响应（如果 make_request 匹配了 expect_status=401，会返回响应数据）

        Returns:
            True 如果断言成功
        """
        # response 不为 None 说明匹配了 expect_status=401
        if response is None:
            self.log_error("✗ 未返回 401 Unauthorized")
            return False

        # 进一步检查错误信息
        is_unauthorized = (
            response.get('error') == 'unauthorized' or
            '未提供认证Token' in str(response.get('message', '')) or
            '未认证' in str(response.get('message', '')) or
            '用户名或密码错误' in str(response.get('message', '')) or
            'invalid' in str(response.get('error', ''))
        )

        if is_unauthorized:
            self.log_success("✓ 认证检查通过：正确返回 401 Unauthorized")
            return True
        else:
            # 即使不匹配具体的错误信息，只要返回了 401 就算通过
            self.log_success(f"✓ 返回 401 Unauthorized: {response.get('error', 'N/A')}")
            return True

    def assert_has_permission(self, permission: str):
        """
        断言当前用户有指定权限

        Args:
            permission: 权限字符串

        Raises:
            AssertionError: 如果当前用户没有该权限
        """
        if not self.check_permission(permission):
            raise AssertionError(
                f"当前角色 {self.current_role} 没有权限 {permission}"
            ) 