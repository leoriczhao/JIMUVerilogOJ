#!/usr/bin/env python3
"""
Verilog OJ 完整API测试套件
整合所有模块的测试，提供完整的测试报告
"""

import sys
import time
from typing import List, Tuple
from colorama import init, Fore, Back, Style

# 导入各个测试模块
from test_health import HealthTester
from test_user import UserTester
from test_problem import ProblemTester
from test_submission import SubmissionTester
from test_forum import ForumTester
from test_news import NewsTester

# 初始化colorama
init(autoreset=True)


class ComprehensiveAPITester:
    """综合API测试器"""

    def __init__(self, run_individual_summaries=False):
        self.run_individual_summaries = run_individual_summaries
        self.start_time = None
        self.end_time = None
        # 收集所有模块的schema验证问题
        self.all_validation_errors = []
        self.all_validation_skipped = []
        
    def log_success(self, message: str):
        """成功日志"""
        print(f"{Fore.GREEN}✅ {message}{Style.RESET_ALL}")
        
    def log_error(self, message: str):
        """错误日志"""
        print(f"{Fore.RED}❌ {message}{Style.RESET_ALL}")
        
    def log_info(self, message: str):
        """信息日志"""
        print(f"{Fore.CYAN}ℹ️  {message}{Style.RESET_ALL}")

    def print_main_header(self):
        """打印主标题"""
        print(f"\n{Back.CYAN}{Fore.BLACK} 🚀 Verilog OJ 完整API测试套件 🚀 {Style.RESET_ALL}")
        print(f"{Fore.CYAN}API 基础地址: http://localhost:8080/api/v1{Style.RESET_ALL}")
        print(f"{Fore.CYAN}测试时间: {time.strftime('%Y-%m-%d %H:%M:%S')}{Style.RESET_ALL}")
        print("=" * 80)

    def run_health_tests(self) -> Tuple[str, bool]:
        """运行健康检查测试"""
        print(f"\n{Back.BLUE}{Fore.WHITE} 🏥 健康检查测试模块 {Style.RESET_ALL}")
        tester = HealthTester()
        # 重写run_tests方法以避免重复的总结
        test_results = [
            ("健康检查接口", tester.test_health_check()),
            ("API根路径", tester.test_api_root()),
        ]
        success = all(result for _, result in test_results)
        return "健康检查模块", success

    def run_user_tests(self) -> Tuple[str, bool]:
        """运行用户管理测试"""
        print(f"\n{Back.GREEN}{Fore.WHITE} 👤 用户管理测试模块 {Style.RESET_ALL}")
        tester = UserTester()
        # 使用 UserTester 的 run_tests 方法，它会正确调用所有测试
        success = tester.run_tests()
        # 收集schema验证问题
        self.all_validation_errors.extend(tester.validation_errors)
        self.all_validation_skipped.extend(tester.validation_skipped)
        return "user模块", success

    def run_problem_tests(self) -> Tuple[str, bool]:
        """运行题目管理测试"""
        print(f"\n{Back.MAGENTA}{Fore.WHITE} 📚 题目管理测试模块 {Style.RESET_ALL}")
        tester = ProblemTester()

        # The ProblemTester's run_tests method has the correct sequential logic.
        success = tester.run_tests()
        # 收集schema验证问题
        self.all_validation_errors.extend(tester.validation_errors)
        self.all_validation_skipped.extend(tester.validation_skipped)
        return "题目管理模块", success

    def run_submission_tests(self) -> Tuple[str, bool]:
        """运行提交管理测试"""
        print(f"\n{Back.YELLOW}{Fore.BLACK} 📝 提交管理测试模块 {Style.RESET_ALL}")
        tester = SubmissionTester()

        # The SubmissionTester's run_tests method has the correct sequential logic.
        success = tester.run_tests()
        # 收集schema验证问题
        self.all_validation_errors.extend(tester.validation_errors)
        self.all_validation_skipped.extend(tester.validation_skipped)
        return "提交管理模块", success

    def run_forum_tests(self) -> Tuple[str, bool]:
        """运行论坛管理测试"""
        print(f"\n{Back.CYAN}{Fore.WHITE} 💬 论坛管理测试模块 {Style.RESET_ALL}")
        tester = ForumTester()

        # The ForumTester's run_tests method has the correct sequential logic.
        success = tester.run_tests()
        # 收集schema验证问题
        self.all_validation_errors.extend(tester.validation_errors)
        self.all_validation_skipped.extend(tester.validation_skipped)
        return "论坛管理模块", success

    def run_news_tests(self) -> Tuple[str, bool]:
        """运行新闻管理测试"""
        print(f"\n{Back.LIGHTGREEN_EX}{Fore.BLACK} 📰 新闻管理测试模块 {Style.RESET_ALL}")
        tester = NewsTester()

        # The NewsTester's run_tests method has the correct sequential logic.
        success = tester.run_tests()
        # 收集schema验证问题
        self.all_validation_errors.extend(tester.validation_errors)
        self.all_validation_skipped.extend(tester.validation_skipped)
        return "新闻管理模块", success

    def print_comprehensive_summary(self, module_results: List[Tuple[str, bool]]):
        """打印综合测试总结"""
        print(f"\n{Back.CYAN}{Fore.BLACK} 📊 综合测试结果总结 📊 {Style.RESET_ALL}")
        print("=" * 80)
        
        passed_modules = 0
        total_modules = len(module_results)
        
        for module_name, result in module_results:
            if result:
                passed_modules += 1
                status = f"{Fore.GREEN}✅ 通过{Style.RESET_ALL}"
            else:
                status = f"{Fore.RED}❌ 失败{Style.RESET_ALL}"
            print(f"{module_name:<20} {status}")
        
        print("=" * 80)
        print(f"总模块数: {total_modules}")
        print(f"通过模块: {Fore.GREEN}{passed_modules}{Style.RESET_ALL}")
        print(f"失败模块: {Fore.RED}{total_modules - passed_modules}{Style.RESET_ALL}")
        print(f"通过率: {Fore.CYAN}{passed_modules/total_modules*100:.1f}%{Style.RESET_ALL}")
        
        # 测试时间统计
        if self.start_time and self.end_time:
            duration = self.end_time - self.start_time
            print(f"测试耗时: {Fore.YELLOW}{duration:.2f}秒{Style.RESET_ALL}")

        # Schema验证统计
        self._print_schema_summary()

        if passed_modules == total_modules:
            print(f"\n{Fore.GREEN}🎉 所有模块测试通过！Verilog OJ API 完全正常！{Style.RESET_ALL}")
            print(f"{Fore.GREEN}🚀 系统已准备好投入使用！{Style.RESET_ALL}")
            return True
        else:
            print(f"\n{Fore.YELLOW}⚠️  部分模块测试失败，请检查相应的API实现{Style.RESET_ALL}")
            failed_modules = [name for name, result in module_results if not result]
            print(f"{Fore.RED}失败的模块: {', '.join(failed_modules)}{Style.RESET_ALL}")
            return False

    def _print_schema_summary(self):
        """打印Schema验证汇总"""
        if not self.all_validation_errors and not self.all_validation_skipped:
            return

        print(f"\n{Back.BLUE}{Fore.WHITE} 📋 OpenAPI Schema验证汇总 📋 {Style.RESET_ALL}")
        print("=" * 80)

        # 统计
        total_errors = len(self.all_validation_errors)
        total_skipped = len(self.all_validation_skipped)

        print(f"Schema验证错误: {Fore.RED}{total_errors}{Style.RESET_ALL}")
        print(f"Schema未定义: {Fore.YELLOW}{total_skipped}{Style.RESET_ALL}")

        # 显示验证错误
        if self.all_validation_errors:
            print(f"\n{Fore.RED}⚠️  Schema验证失败的API：{Style.RESET_ALL}")
            for idx, error in enumerate(self.all_validation_errors[:5], 1):
                print(f"  {idx}. [{error['module']}] {error['method']} {error['endpoint']} (状态码: {error['status_code']})")
                print(f"     {Fore.RED}错误: {error['error'][:100]}...{Style.RESET_ALL}")
            if total_errors > 5:
                print(f"  ... 还有 {total_errors - 5} 个验证错误")

        # 显示未定义schema的API
        if self.all_validation_skipped:
            print(f"\n{Fore.YELLOW}ℹ️  缺少Schema定义的API端点：{Style.RESET_ALL}")
            # 去重
            unique_skipped = {}
            for skip in self.all_validation_skipped:
                key = f"{skip['method']} {skip['endpoint']} {skip['status_code']}"
                if key not in unique_skipped:
                    unique_skipped[key] = skip

            for idx, (key, skip) in enumerate(list(unique_skipped.items())[:10], 1):
                print(f"  {idx}. [{skip['module']}] {skip['method']} {skip['endpoint']} (状态码: {skip['status_code']})")

            if len(unique_skipped) > 10:
                print(f"  ... 还有 {len(unique_skipped) - 10} 个未定义的schema")

            print(f"\n{Fore.YELLOW}💡 建议: 在 docs/openapi/ 目录下补充这些API的Schema定义{Style.RESET_ALL}")

    def run_all_tests(self, modules_to_run=None):
        """运行所有测试"""
        self.start_time = time.time()
        self.print_main_header()
        
        # 定义所有可用的测试模块
        available_modules = {
            'health': self.run_health_tests,
            'user': self.run_user_tests,
            'problem': self.run_problem_tests,
            'submission': self.run_submission_tests,
            'forum': self.run_forum_tests,
            'news': self.run_news_tests,
        }
        
        # 如果没有指定模块，运行所有模块
        if modules_to_run is None:
            modules_to_run = list(available_modules.keys())
        
        module_results = []
        
        for module_name in modules_to_run:
            if module_name in available_modules:
                try:
                    result = available_modules[module_name]()
                    module_results.append(result)
                    # 模块间添加延迟以避免速率限制
                    time.sleep(1)
                except Exception as e:
                    self.log_error(f"{module_name}模块测试异常: {str(e)}")
                    module_results.append((f"{module_name}模块", False))
            else:
                self.log_error(f"未知的测试模块: {module_name}")
        
        self.end_time = time.time()
        
        # 打印综合结果
        return self.print_comprehensive_summary(module_results)


def main():
    """主函数"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Verilog OJ API 测试套件')
    parser.add_argument('--modules', '-m', nargs='+', 
                       choices=['health', 'user', 'problem', 'submission', 'forum', 'news'],
                       help='指定要运行的测试模块')
    parser.add_argument('--list', '-l', action='store_true', 
                       help='列出所有可用的测试模块')
    
    args = parser.parse_args()
    
    if args.list:
        print("可用的测试模块:")
        print("  health     - 健康检查测试")
        print("  user       - 用户管理测试")
        print("  problem    - 题目管理测试")
        print("  submission - 提交管理测试")
        print("  forum      - 论坛管理测试")
        print("  news       - 新闻管理测试")
        return
    
    print("Verilog OJ API 完整测试套件")
    print("确保后端服务已启动在 http://localhost:8080")
    print("")
    
    try:
        tester = ComprehensiveAPITester()
        success = tester.run_all_tests(args.modules)
        sys.exit(0 if success else 1)
    except KeyboardInterrupt:
        print(f"\n{Fore.YELLOW}测试被用户中断{Style.RESET_ALL}")
        sys.exit(1)
    except Exception as e:
        print(f"\n{Fore.RED}测试异常: {str(e)}{Style.RESET_ALL}")
        sys.exit(1)


if __name__ == "__main__":
    main() 