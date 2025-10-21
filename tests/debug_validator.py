#!/usr/bin/env python3
"""
Debug OpenAPI validator schema resolution
"""

from openapi_validator import get_validator
import json

validator = get_validator()

# 检查user模块的schemas
if 'user' in validator.schemas:
    print("User module schemas:")
    for key in validator.schemas['user'].keys():
        print(f"  - {key}")

    # 检查一个具体的schema
    schema_key = "POST /users/register 201"
    if schema_key in validator.schemas['user']:
        schema = validator.schemas['user'][schema_key]
        print(f"\nSchema for {schema_key}:")
        print(json.dumps(schema, indent=2, ensure_ascii=False))
else:
    print("User module not loaded!")

print("\nAll loaded modules:", list(validator.schemas.keys()))
