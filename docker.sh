#!/bin/bash
# docker.sh
# 自动化 Docker 部署脚本 (优化版 - 解决挂载冲突)
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

print_color ">>> 开始 Docker 部署流程..." "$CYAN"

# 1. 环境清理
print_color "1. 清理旧环境 (docker-compose down)..." "$CYAN"
# 忽略清理过程中的错误
docker-compose down --remove-orphans > /dev/null 2>&1 || true
docker-compose rm -f > /dev/null 2>&1 || true

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

if docker-compose up -d --build; then
    print_color "\n>>> 部署成功！" "$GREEN"
    echo "-------------------------------------------"
    echo -e "   API 后端:    http://localhost:8080"
    echo -e "   后台管理:    http://localhost:3001"
    echo -e "   前台页面:    http://localhost:3002"
    echo "-------------------------------------------"
else
    print_color "\n[严重错误] Docker 启动失败。" "$RED"
    exit 1
fi
