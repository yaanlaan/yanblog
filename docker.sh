#!/bin/bash
# docker.sh
# 自动化 Docker 部署脚本 (优化版 - 兼容 docker-compose V1 和 V2)
# 对应 Windows 的 docker.ps1 逻辑

set -e # 遇到错误立即退出

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_color() {
    echo -e "${2}${1}${NC}"
}

# 检测可用命令
if command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    print_color "[错误] 未找到 docker-compose 或 docker compose 命令" "$RED"
    exit 1
fi

print_color ">>> 开始 Docker 部署流程 (使用: $COMPOSE_CMD)..." "$CYAN"

# 1. 环境清理
print_color "1. 清理旧环境 ($COMPOSE_CMD down)..." "$CYAN"
# 忽略清理过程中的错误
$COMPOSE_CMD down --remove-orphans > /dev/null 2>&1 || true
$COMPOSE_CMD rm -f > /dev/null 2>&1 || true

# 2. 准备构建上下文 (配置文件)
print_color "2. 准备构建上下文 (配置文件)..." "$CYAN"

# 确保目标目录存在
mkdir -p docker_field/frontend
mkdir -p docker_field/backend

# 复制配置文件
CONFIG_FRONT="web/frontend/public/config.yaml"
CONFIG_BACK="config/config.yaml"
DEST_FRONT="docker_field/frontend/config.yaml"
DEST_BACK="docker_field/backend/config.yaml"

# --- 新增: 同步前端静态资源到后台 (解决后台图片 404 问题) ---
print_color "2.1 同步静态资源 (Frontend -> Backend)..." "$CYAN"
# 确保后台 public 目录存在
mkdir -p web/backend/public/static
# 复制 static 目录下的图片资源
cp -r web/frontend/public/static/* web/backend/public/static/ 2>/dev/null || true
# 复制可能用到的根目录图标
cp web/frontend/public/favicon.svg web/backend/public/ 2>/dev/null || true
echo "   已同步静态资源到后台构建目录"
# -------------------------------------------------------

if [ -f "$CONFIG_FRONT" ]; then
    cp -f "$CONFIG_FRONT" "$DEST_FRONT"
    echo "   已同步: $CONFIG_FRONT -> $DEST_FRONT"
else
    print_color "   [警告] 前端配置文件未找到: $CONFIG_FRONT" "$YELLOW"
fi

if [ -f "$CONFIG_BACK" ]; then
    cp -f "$CONFIG_BACK" "$DEST_BACK"
    echo "   已同步: $CONFIG_BACK -> $DEST_BACK"
else
    print_color "   [警告] 后端配置文件未找到: $CONFIG_BACK" "$YELLOW"
fi

# 3. 执行部署
print_color "3. 构建并启动 Docker 服务..." "$CYAN"
print_color "   (使用 Named Volumes 存储数据)" "$CYAN"

if $COMPOSE_CMD up -d --build; then
    # 动态获取端口函数
    get_port() {
        local service=$1
        local default=$2
        # 尝试从 docker-compose.yaml 提取端口 (匹配 - "8080:8080" 或 - 8080:8080)
        # 1. 找服务名 2. 找ports块 3. 提取主机端口
        local port=$(grep -A 15 "^[[:space:]]*${service}:" docker-compose.yaml | grep -m 1 "ports:" -A 5 | grep -m 1 -E "\-\s+[\"']?[0-9]+:[0-9]+" | sed -E 's/.*-\s+["'"'"']?([0-9]+):.*/\1/')
        
        if [ -z "$port" ]; then
            echo "$default"
        else
            echo "$port"
        fi
    }

    PORT_BACKEND=$(get_port "backend" "8080")
    PORT_ADMIN=$(get_port "frontend-admin" "3001")
    PORT_PUBLIC=$(get_port "frontend-public" "3002")

    print_color "\n>>> 部署成功！" "$GREEN"
    echo "-------------------------------------------"
    echo -e "   API 后端:    http://localhost:$PORT_BACKEND"
    echo -e "   后台管理:    http://localhost:$PORT_ADMIN"
    echo -e "   前台页面:    http://localhost:$PORT_PUBLIC"
    echo "-------------------------------------------"
else
    print_color "\n[严重错误] Docker 启动失败。" "$RED"
    exit 1
fi
