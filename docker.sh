#!/bin/bash
# docker.sh
# 自动化 Docker 部署脚本 (优化版 - 解决挂载冲突 & 动态端口注入)
# 功能同步自 docker.ps1

set -e # 遇到错误立即退出

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

print_color() {
    echo -e "${2}${1}${NC}"
}

# 配置文件路径
FRONTEND_CONFIG="web/frontend/public/config.yaml"
FRONTEND_CONFIG_BAK="web/frontend/public/config.yaml.bak"

# 还原配置文件的函数
cleanup() {
    EXIT_CODE=$?
    if [ -f "$FRONTEND_CONFIG_BAK" ]; then
        cp -f "$FRONTEND_CONFIG_BAK" "$FRONTEND_CONFIG"
        rm -f "$FRONTEND_CONFIG_BAK"
        print_color "   [Build Cleanup] 已还原前端配置文件至初始状态" "$GRAY"
    fi
    exit $EXIT_CODE
}

# 设置 trap，脚本退出时(无论正常还是错误)都执行 cleanup
trap cleanup EXIT INT TERM

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

# 2. 从 docker-compose.yaml 读取配置 (提取端口)
print_color "2. 读取 Docker 配置..." "$CYAN"

# 默认值
BACKEND_PORT="8080"
ADMIN_PORT="3001"
PUBLIC_PORT="3002"
DB_VOLUME="yanblog_db_data"

# 简单的逐行解析逻辑
CURRENT_SERVICE=""
while IFS= read -r line; do
    trim_line=$(echo "$line" | xargs) # 去除首尾空白
    
    if [[ "$trim_line" == "backend:" ]]; then CURRENT_SERVICE="backend"; fi
    if [[ "$trim_line" == "frontend-admin:" ]]; then CURRENT_SERVICE="frontend-admin"; fi
    if [[ "$trim_line" == "frontend-public:" ]]; then CURRENT_SERVICE="frontend-public"; fi
    if [[ "$trim_line" == "db:" ]]; then CURRENT_SERVICE="db"; fi

    # 提取端口 - "HOST:CONTAINER" (兼容引号)
    if [[ "$trim_line" =~ -[[:space:]]*[\"\']?([0-9]+):8080[\"\']? ]]; then
        if [ "$CURRENT_SERVICE" == "backend" ]; then BACKEND_PORT="${BASH_REMATCH[1]}"; fi
    fi
    if [[ "$trim_line" =~ -[[:space:]]*[\"\']?([0-9]+):80[\"\']? ]]; then
        if [ "$CURRENT_SERVICE" == "frontend-admin" ]; then ADMIN_PORT="${BASH_REMATCH[1]}"; fi
        if [ "$CURRENT_SERVICE" == "frontend-public" ]; then PUBLIC_PORT="${BASH_REMATCH[1]}"; fi
    fi

    # 提取 Volume - NAME:/var/lib/mysql
    # 注意：这里需要匹配映射到 mysql 目录的行
    if [[ "$trim_line" =~ -[[:space:]]*([a-zA-Z0-9_\-]+):/var/lib/mysql ]]; then
        if [ "$CURRENT_SERVICE" == "db" ]; then DB_VOLUME="${BASH_REMATCH[1]}"; fi
    fi
done < docker-compose.yaml

print_color "   检测到后台管理端口: $ADMIN_PORT" "$GRAY"

# 3. 准备构建环境
print_color "3. 准备构建环境..." "$CYAN"

# --- 3.1 同步静态资源 ---
print_color "   正在同步静态资源..." "$GRAY"
BACKEND_STATIC="web/backend/public/static"
mkdir -p "$BACKEND_STATIC"

if [ -d "web/frontend/public/static" ]; then
    cp -rf web/frontend/public/static/* "$BACKEND_STATIC/"
fi
if [ -f "web/frontend/public/favicon.svg" ]; then
    cp -f web/frontend/public/favicon.svg web/backend/public/
fi

# --- 3.2 注入端口配置 ---
if [ -f "$FRONTEND_CONFIG" ]; then
    # 备份
    cp "$FRONTEND_CONFIG" "$FRONTEND_CONFIG_BAK"
    
    # 替换端口 (使用临时文件以兼容 Linux/Mac)
    # 正则解释：^ 头部, ( *) 捕获缩进, dev_admin_port: 匹配key, *[0-9]* 匹配可能的旧值, \(.*\) 匹配注释及剩余
    if grep -q "dev_admin_port:" "$FRONTEND_CONFIG"; then
        sed "s/^\( *\)dev_admin_port: *[0-9]*\(.*\)/\1dev_admin_port: $ADMIN_PORT\2/" "$FRONTEND_CONFIG" > "${FRONTEND_CONFIG}.tmp" && mv "${FRONTEND_CONFIG}.tmp" "$FRONTEND_CONFIG"
        print_color "   [Build Prep] 已注入端口配置: $(grep 'dev_admin_port' $FRONTEND_CONFIG)" "$CYAN"
    else
        print_color "   [警告] 未找到 dev_admin_port 配置项" "$YELLOW"
    fi
else
    print_color "   [警告] 找不到配置文件: $FRONTEND_CONFIG" "$YELLOW"
fi

# 4. 执行部署
print_color "4. 构建并启动 Docker 服务..." "$CYAN"
print_color "   (使用 Named Volumes 存储数据)" "$GRAY"

$COMPOSE_CMD up -d --build

# 5. 检查结果
print_color "\n>>> 部署成功！" "$GREEN"
print_color "\n访问地址:"
print_color "-------------------------------------------"
print_color "   DB 数据:     存储在 Docker Volume '$DB_VOLUME' 中"
print_color "   API 后端:    http://localhost:$BACKEND_PORT"
print_color "   后台管理:    http://localhost:$ADMIN_PORT"
print_color "   前台页面:    http://localhost:$PUBLIC_PORT"
print_color "-------------------------------------------"

# trap cleanup 会在退出时自动运行