#!/bin/bash

echo "========================================="
echo "  博客系统初始化向导"
echo "========================================="
echo ""

# 检查配置文件是否存在
if [ ! -f "config/backend/config.yaml" ]; then
    echo "✓ 检测到 config/backend/config.yaml 不存在，正在创建..."
    cp config/config_template.yaml config/backend/config.yaml
    echo "✓ 已创建 config/backend/config.yaml"
    echo ""
else
    echo "ℹ️  config/backend/config.yaml 已存在"
    echo ""
fi

# 生成 JWT 密钥
echo "🔐 正在生成 JWT 密钥..."
JWT_KEY=$(openssl rand -hex 32 2>/dev/null)

if [ -z "$JWT_KEY" ]; then
    echo "⚠️  无法生成 JWT 密钥，请手动运行：openssl rand -hex 32"
    echo "   然后将结果填入 config/backend/config.yaml 的 JwtKey 字段"
    echo ""
else
    echo "✓ JWT 密钥已生成"
    # 替换配置文件中的 JwtKey
    sed -i.bak "s/^JwtKey:.*/JwtKey: $JWT_KEY/" config/backend/config.yaml
    rm -f config/backend/config.yaml.bak
    echo "✓ JWT 密钥已写入配置文件"
    echo ""
fi

# 提示用户修改数据库密码
echo "========================================="
echo "  配置完成！"
echo "========================================="
echo ""
echo "请检查以下配置项："
echo "  1. database.DbPassWord - 数据库密码（强烈建议修改）"
echo "  2. database.DbHost - Docker 部署为 db，本地部署为 localhost"
echo "  3. server.SiteUrl - 网站域名（可选）"
echo ""
echo "下一步："
echo "  - Docker 部署：运行 ./docker.sh"
echo "  - 本地部署：运行 go run main.go"
echo ""
echo "详细使用说明请查看：QUICK_START.md"
echo "========================================="
