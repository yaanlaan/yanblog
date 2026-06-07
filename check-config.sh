#!/bin/bash

echo "========================================="
echo "  配置检查工具"
echo "========================================="
echo ""

# 检查配置文件
if [ ! -f "config/backend/config.yaml" ]; then
    echo "❌ 配置文件不存在：config/backend/config.yaml"
    echo "   解决方法：运行 ./setup.sh 或复制 config/config_template.yaml"
    exit 1
else
    echo "✓ 配置文件存在"
fi

# 检查数据库密码
DB_PASS=$(grep "DbPassWord:" config/backend/config.yaml | awk '{print $2}')
if [ "$DB_PASS" == "rootpassword" ]; then
    echo "⚠️  警告：数据库密码仍为默认值，建议修改"
else
    echo "✓ 数据库密码已修改"
fi

# 检查 JWT 密钥
JWT_KEY=$(grep "JwtKey:" config/backend/config.yaml | awk '{print $2}')
if [ -z "$JWT_KEY" ] || [ "$JWT_KEY" == "" ]; then
    echo "❌ JWT 密钥未设置"
    echo "   解决方法：运行 ./setup.sh 自动生成"
    echo "   或手动运行：openssl rand -hex 32"
    exit 1
elif [ ${#JWT_KEY} -lt 32 ]; then
    echo "⚠️  警告：JWT 密钥长度不足（当前：${#JWT_KEY} 位，建议：32 位）"
else
    echo "✓ JWT 密钥已设置（长度：${#JWT_KEY} 位）"
fi

# 检查数据库主机
DB_HOST=$(grep "DbHost:" config/backend/config.yaml | awk '{print $2}')
echo "ℹ️  数据库主机：$DB_HOST"

# 检查端口
HTTP_PORT=$(grep "HttpPort:" config/backend/config.yaml | awk '{print $2}')
echo "ℹ️  服务端口：$HTTP_PORT"

echo ""
echo "========================================="
echo "  检查完成！"
echo "========================================="
echo ""

if [ -f "go.mod" ]; then
    echo "下一步："
    echo "  - Docker 部署：./docker.sh"
    echo "  - 本地部署：go run main.go"
else
    echo "提示：请确保在项目根目录运行此脚本"
fi

echo ""
