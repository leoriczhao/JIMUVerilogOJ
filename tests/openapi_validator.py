#!/usr/bin/env python3
"""
OpenAPI Schema 验证器
用于验证API响应是否符合OpenAPI规范定义的schema
"""

import os
import yaml
import json
from pathlib import Path
from typing import Dict, Optional, Any, Tuple
from jsonschema import validate, ValidationError, Draft7Validator, RefResolver
from colorama import Fore, Style


class OpenAPIValidator:
    """OpenAPI Schema验证器"""

    def __init__(self, openapi_dir: str = "../docs/openapi"):
        """
        初始化验证器

        Args:
            openapi_dir: OpenAPI文档所在目录
        """
        self.openapi_dir = Path(__file__).parent / openapi_dir
        self.specs = {}
        self.schemas = {}
        self._load_specs()

    def _load_yaml_file(self, file_path: Path) -> Dict:
        """加载YAML文件"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                return yaml.safe_load(f)
        except Exception as e:
            print(f"{Fore.RED}加载OpenAPI文件失败 {file_path}: {str(e)}{Style.RESET_ALL}")
            return {}

    def _resolve_ref(self, ref: str, base_path: Path, base_document: Optional[Dict] = None) -> Optional[Dict]:
        """
        解析$ref引用

        Args:
            ref: $ref引用字符串，如 './models/user.yaml#/components/schemas/User' 或 '#/components/schemas/User'
            base_path: 当前文件的路径
            base_document: 当前文档（用于解析同文件内的引用）

        Returns:
            解析后的schema字典
        """
        # 处理同文件内的引用（以#开头）
        if ref.startswith('#'):
            if base_document is None:
                # 如果没有传入基础文档，尝试加载当前文件
                base_document = self._load_yaml_file(base_path)

            if not base_document:
                return None

            # 解析JSON路径
            json_path = ref.lstrip('#').strip('/')
            if not json_path:
                return base_document

            path_parts = json_path.split('/')
            result = base_document
            for part in path_parts:
                if isinstance(result, dict) and part in result:
                    result = result[part]
                else:
                    return None

            # 递归解析嵌套的$ref
            if isinstance(result, dict) and '$ref' in result:
                return self._resolve_ref(result['$ref'], base_path, base_document)

            return result

        # 处理跨文件引用（以.开头）
        if not ref.startswith('.'):
            return None

        # 分离文件路径和JSON路径
        parts = ref.split('#')
        if len(parts) != 2:
            return None

        file_ref, json_path = parts

        # 解析文件路径
        ref_file = (base_path.parent / file_ref).resolve()
        if not ref_file.exists():
            return None

        # 加载引用的文件
        ref_data = self._load_yaml_file(ref_file)
        if not ref_data:
            return None

        # 解析JSON路径
        path_parts = json_path.strip('/').split('/')
        result = ref_data
        for part in path_parts:
            if isinstance(result, dict) and part in result:
                result = result[part]
            else:
                return None

        # 递归解析嵌套的$ref（传入新的文档作为base_document）
        if isinstance(result, dict) and '$ref' in result:
            return self._resolve_ref(result['$ref'], ref_file, ref_data)

        return result

    def _resolve_all_refs(self, schema: Any, base_path: Path, base_document: Optional[Dict] = None) -> Any:
        """
        递归解析schema中的所有$ref引用

        Args:
            schema: 需要解析的schema
            base_path: 当前文件的路径
            base_document: 当前文档（用于解析同文件内的引用）

        Returns:
            解析后的schema
        """
        if isinstance(schema, dict):
            if '$ref' in schema:
                resolved = self._resolve_ref(schema['$ref'], base_path, base_document)
                if resolved:
                    # 递归解析返回的schema
                    return self._resolve_all_refs(resolved, base_path, base_document)
                return schema

            # 递归处理字典中的每个值
            return {k: self._resolve_all_refs(v, base_path, base_document) for k, v in schema.items()}

        elif isinstance(schema, list):
            # 递归处理列表中的每个元素
            return [self._resolve_all_refs(item, base_path, base_document) for item in schema]

        return schema

    def _load_specs(self):
        """加载所有OpenAPI规范文件"""
        if not self.openapi_dir.exists():
            print(f"{Fore.YELLOW}⚠️  OpenAPI目录不存在: {self.openapi_dir}{Style.RESET_ALL}")
            return

        # 加载主API规范文件
        spec_files = {
            'user': 'user.yaml',
            'problem': 'problem.yaml',
            'submission': 'submission.yaml',
            'forum': 'forum.yaml',
            'news': 'news.yaml',
        }

        for name, filename in spec_files.items():
            spec_path = self.openapi_dir / filename
            if spec_path.exists():
                spec = self._load_yaml_file(spec_path)
                if spec:
                    self.specs[name] = spec
                    # 解析并缓存所有schema
                    self._extract_schemas(name, spec, spec_path)

    def _extract_schemas(self, module_name: str, spec: Dict, spec_path: Path):
        """
        从OpenAPI规范中提取并解析schema定义

        Args:
            module_name: 模块名称
            spec: OpenAPI规范字典
            spec_path: 规范文件路径
        """
        if module_name not in self.schemas:
            self.schemas[module_name] = {}

        # 标准化路径，避免相对路径导致的缓存键不一致
        spec_path = spec_path.resolve()

        # 加载referenced schema文件，建立完整的文档上下文
        # 如果spec中有components引用，需要加载这些文件
        models_dir = spec_path.parent / 'models'
        schema_documents = {}

        # 预加载所有模型文件
        if models_dir.exists():
            for model_file in models_dir.glob('*.yaml'):
                resolved_model_file = model_file.resolve()
                model_data = self._load_yaml_file(resolved_model_file)
                if model_data:
                    schema_documents[resolved_model_file] = model_data

        # 从paths中提取每个endpoint的响应schema
        paths = spec.get('paths', {})
        for path, path_item in paths.items():
            for method, operation in path_item.items():
                if method.upper() not in ['GET', 'POST', 'PUT', 'DELETE', 'PATCH']:
                    continue

                operation_id = operation.get('operationId', f"{method}_{path}")

                # 提取请求body schema
                request_body = operation.get('requestBody', {})
                if request_body:
                    content = request_body.get('content', {}).get('application/json', {})
                    if 'schema' in content:
                        # 解析schema时，需要知道被引用文件的位置
                        request_schema = self._resolve_schema_with_context(
                            content['schema'], spec_path, schema_documents
                        )
                        self.schemas[module_name][f"{operation_id}_request"] = request_schema

                # 提取响应schema
                responses = operation.get('responses', {})
                for status_code, response in responses.items():
                    content = response.get('content', {}).get('application/json', {})
                    if 'schema' in content:
                        response_schema = self._resolve_schema_with_context(
                            content['schema'], spec_path, schema_documents
                        )
                        key = f"{method.upper()} {path} {status_code}"
                        self.schemas[module_name][key] = response_schema

    def _resolve_schema_with_context(self, schema: Any, base_path: Path,
                                     schema_documents: Dict[Path, Dict]) -> Any:
        """
        使用上下文信息解析schema

        Args:
            schema: 需要解析的schema
            base_path: 当前文件路径
            schema_documents: 已加载的schema文档字典

        Returns:
            解析后的schema
        """
        if not isinstance(schema, dict):
            return schema

        base_path = base_path.resolve()

        if '$ref' in schema:
            ref = schema['$ref']

            # 处理跨文件引用
            if ref.startswith('./'):
                parts = ref.split('#')
                if len(parts) == 2:
                    file_ref, json_path = parts
                    ref_file = (base_path.parent / file_ref).resolve()

                    # 从预加载的文档中获取
                    ref_document = schema_documents.get(ref_file)
                    if ref_document:
                        # 解析JSON路径
                        path_parts = json_path.strip('/').split('/')
                        result = ref_document
                        for part in path_parts:
                            if isinstance(result, dict) and part in result:
                                result = result[part]
                            else:
                                # 无法找到引用，返回原schema
                                return schema

                        # 递归解析嵌套的引用（现在有了文档上下文）
                        # 继续使用ref_file作为基础路径，因为我们现在在referenced文件的上下文中
                        return self._resolve_schema_with_context(result, ref_file, schema_documents)
                    else:
                        # 文档未加载，返回原schema
                        return schema

            # 处理同文件内的引用（以#开头，没有文件路径）
            elif ref.startswith('#'):
                # 这种情况通常出现在模型文件内部
                # 我们需要找到包含此schema的文档
                # 由于base_path可能是模型文件，尝试从该文件加载
                current_doc = self._load_yaml_file(base_path)
                if current_doc:
                    json_path = ref.lstrip('#').strip('/')
                    path_parts = json_path.split('/')
                    result = current_doc
                    for part in path_parts:
                        if isinstance(result, dict) and part in result:
                            result = result[part]
                        else:
                            return schema

                    # 递归解析
                    return self._resolve_schema_with_context(result, base_path, schema_documents)

            return schema

        # 递归处理字典和列表中的所有值
        if isinstance(schema, dict):
            return {k: self._resolve_schema_with_context(v, base_path, schema_documents)
                   for k, v in schema.items()}
        elif isinstance(schema, list):
            return [self._resolve_schema_with_context(item, base_path, schema_documents)
                   for item in schema]

        return schema

    def validate_response(self, module: str, method: str, endpoint: str,
                         status_code: int, response_data: Any) -> Tuple[bool, Optional[str]]:
        """
        验证响应数据是否符合OpenAPI schema

        Args:
            module: 模块名称（如 'user', 'problem'）
            method: HTTP方法（如 'GET', 'POST'）
            endpoint: API端点（如 '/users/profile'）
            status_code: HTTP状态码
            response_data: 响应数据

        Returns:
            (是否验证通过, 错误信息)
        """
        if module not in self.schemas:
            return True, f"模块 '{module}' 没有加载OpenAPI规范"

        # 构建schema key
        # 移除/api/v1前缀（如果存在）
        clean_endpoint = endpoint.replace('/api/v1', '')
        schema_key = f"{method.upper()} {clean_endpoint} {status_code}"

        if schema_key not in self.schemas[module]:
            return True, f"未找到schema: {schema_key}"

        schema = self.schemas[module][schema_key]

        try:
            # 创建一个自定义的RefResolver来处理文件引用
            # 但由于我们已经解析了所有引用，直接验证即可
            validator = Draft7Validator(schema)
            errors = list(validator.iter_errors(response_data))

            if errors:
                error_messages = []
                for error in errors:
                    path = '.'.join(str(p) for p in error.path) if error.path else 'root'
                    error_messages.append(f"  - {path}: {error.message}")
                return False, "Schema验证失败:\n" + "\n".join(error_messages)

            return True, None

        except Exception as e:
            return False, f"验证过程出错: {str(e)}"

    def validate_request(self, module: str, operation_id: str, request_data: Any) -> Tuple[bool, Optional[str]]:
        """
        验证请求数据是否符合OpenAPI schema

        Args:
            module: 模块名称
            operation_id: 操作ID
            request_data: 请求数据

        Returns:
            (是否验证通过, 错误信息)
        """
        if module not in self.schemas:
            return True, f"模块 '{module}' 没有加载OpenAPI规范"

        schema_key = f"{operation_id}_request"

        if schema_key not in self.schemas[module]:
            return True, f"未找到请求schema: {schema_key}"

        schema = self.schemas[module][schema_key]

        try:
            validator = Draft7Validator(schema)
            errors = list(validator.iter_errors(request_data))

            if errors:
                error_messages = []
                for error in errors:
                    path = '.'.join(str(p) for p in error.path) if error.path else 'root'
                    error_messages.append(f"  - {path}: {error.message}")
                return False, "请求Schema验证失败:\n" + "\n".join(error_messages)

            return True, None

        except Exception as e:
            return False, f"验证过程出错: {str(e)}"

    def get_available_schemas(self, module: str = None) -> Dict[str, list]:
        """
        获取可用的schema列表

        Args:
            module: 模块名称，如果为None则返回所有模块

        Returns:
            {module_name: [schema_keys]}
        """
        if module:
            return {module: list(self.schemas.get(module, {}).keys())}
        return {m: list(schemas.keys()) for m, schemas in self.schemas.items()}


# 全局验证器实例
_validator = None


def get_validator() -> OpenAPIValidator:
    """获取全局验证器实例（单例模式）"""
    global _validator
    if _validator is None:
        _validator = OpenAPIValidator()
    return _validator
