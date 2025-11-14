#!/bin/bash

# API测试脚本

BASE_URL="http://localhost:8080"

echo "=== TiDB Gin Demo API测试 ==="
echo

# 1. 健康检查
echo "1. 健康检查:"
curl -s "$BASE_URL/health" | jq .
echo

# 2. 创建用户
echo "2. 创建用户:"
curl -s -X POST "$BASE_URL/api/users/" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "age": 25
  }' | jq .
echo

# 3. 创建另一个用户
echo "3. 创建另一个用户:"
curl -s -X POST "$BASE_URL/api/users/" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "李四",
    "email": "lisi@example.com",
    "age": 30
  }' | jq .
echo

# 4. 获取所有用户
echo "4. 获取所有用户:"
curl -s "$BASE_URL/api/users/" | jq .
echo

# 5. 获取单个用户（ID=1）
echo "5. 获取单个用户(ID=1):"
curl -s "$BASE_URL/api/users/1" | jq .
echo

# 6. 更新用户
echo "6. 更新用户(ID=1):"
curl -s -X PUT "$BASE_URL/api/users/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三更新",
    "age": 26
  }' | jq .
echo

# 7. 再次获取用户确认更新
echo "7. 确认更新后的用户:"
curl -s "$BASE_URL/api/users/1" | jq .
echo

# 8. 删除用户
echo "8. 删除用户(ID=2):"
curl -s -X DELETE "$BASE_URL/api/users/2" | jq .
echo

# 9. 最终用户列表
echo "9. 最终用户列表:"
curl -s "$BASE_URL/api/users/" | jq .
echo

echo "=== 测试完成 ==="