#!/bin/bash
set -e

echo "========================================="
echo "  YanBlog 初始化向导"
echo "========================================="
echo ""

# 确保配置目录存在
mkdir -p config/backend config/frontend

# --- 后端配置 ---
if [ ! -f "config/backend/config.yaml" ]; then
    echo "✓ 创建后端配置..."
    if [ -f "config/config_template.yaml" ]; then
        cp config/config_template.yaml config/backend/config.yaml
    else
        # 内联默认配置
        cat > config/backend/config.yaml << 'EOF'
server:
  AppMode: debug
  HttpPort: :8080
  SiteUrl: ""

database:
  Db: SQLite
  DbHost: localhost
  DbPort: 3306
  DbUser: root
  DbPassWord: rootpassword
  DbName: yanblog.db

JwtKey: ""

weather:
  Provider: openweathermap
  ApiKey: ""
  DefaultCity: Hefei

FrontEndConfigPath: config/frontend/config.yaml
EOF
    fi
    echo "✓ 已创建 config/backend/config.yaml"
else
    echo "ℹ️  config/backend/config.yaml 已存在，跳过"
fi

# 生成 JWT 密钥
echo ""
echo "🔐 正在生成 JWT 密钥..."
JWT_KEY=$(openssl rand -hex 32 2>/dev/null)

if [ -z "$JWT_KEY" ]; then
    echo "⚠️  无法生成 JWT 密钥（openssl 未安装），请手动设置"
else
    sed -i.bak "s/^JwtKey:.*/JwtKey: $JWT_KEY/" config/backend/config.yaml
    rm -f config/backend/config.yaml.bak
    echo "✓ JWT 密钥已写入配置文件"
fi

# --- 前端配置 ---
if [ ! -f "config/frontend/config.yaml" ]; then
    echo ""
    echo "✓ 创建前端配置..."
    if [ -f "web/frontend/public/config_template.yaml" ]; then
        cp web/frontend/public/config_template.yaml config/frontend/config.yaml
    elif [ -f "web/frontend/public/config.yaml" ]; then
        cp web/frontend/public/config.yaml config/frontend/config.yaml
    else
        echo "⚠️  未找到前端配置模板，请手动创建 config/frontend/config.yaml"
    fi
    [ -f "config/frontend/config.yaml" ] && echo "✓ 已创建 config/frontend/config.yaml"
else
    echo "ℹ️  config/frontend/config.yaml 已存在，跳过"
fi

echo ""
echo "========================================="
echo "  配置完成！"
echo "========================================="
echo ""
echo "请检查以下配置项："
echo "  1. database.DbPassWord — 数据库密码"
echo "  2. database.DbHost — Docker 部署为 db，本地为 localhost"
echo "  3. server.SiteUrl — 网站域名"
echo ""
echo "下一步："
echo "  Docker 部署：./docker.sh"
echo "  本地部署：go run main.go"
echo ""
echo "详细说明：QUICK_START.md"
echo "========================================="
