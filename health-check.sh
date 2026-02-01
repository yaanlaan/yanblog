#!/bin/bash

echo "==================================="
echo "YanBlog 部署健康检查"
echo "==================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

check_pass() {
    echo -e "${GREEN}✓${NC} $1"
}

check_fail() {
    echo -e "${RED}✗${NC} $1"
}

check_warn() {
    echo -e "${YELLOW}!${NC} $1"
}

# 1. 检查容器状态
echo "1. 检查容器状态..."
if docker compose ps | grep -q "yanblog.*Up"; then
    check_pass "应用容器运行中"
else
    check_fail "应用容器未运行"
    exit 1
fi

if docker compose ps | grep -q "db.*Up"; then
    check_pass "数据库容器运行中"
else
    check_fail "数据库容器未运行"
    exit 1
fi
echo ""

# 2. 检查数据库连接
echo "2. 检查数据库连接..."
if docker compose exec -T db mysqladmin ping -h localhost -u root -prootpassword 2>/dev/null | grep -q "alive"; then
    check_pass "数据库连接正常"
else
    check_fail "数据库连接失败"
fi
echo ""

# 3. 检查后端进程
echo "3. 检查后端进程..."
if docker compose exec -T yanblog ps aux 2>/dev/null | grep -q "[s]erver"; then
    check_pass "Go 后端进程运行中"
else
    check_fail "Go 后端进程未运行"
fi
echo ""

# 4. 检查 Nginx
echo "4. 检查 Nginx..."
if docker compose exec -T yanblog ps aux 2>/dev/null | grep -q "[n]ginx"; then
    check_pass "Nginx 进程运行中"
else
    check_fail "Nginx 进程未运行"
fi
echo ""

# 5. 检查配置文件
echo "5. 检查配置文件..."
if docker compose exec -T yanblog test -f /app/config/config.yaml 2>/dev/null; then
    check_pass "后端配置文件存在"
else
    check_fail "后端配置文件缺失"
fi

if docker compose exec -T yanblog test -f /app/config/frontend_config.yaml 2>/dev/null; then
    check_pass "前端配置文件存在"
else
    check_fail "前端配置文件缺失"
fi
echo ""

# 6. 检查静态文件
echo "6. 检查静态文件..."
if docker compose exec -T yanblog test -f /usr/share/nginx/html/web/static/about.md 2>/dev/null; then
    check_pass "about.md 文件存在"
else
    check_warn "about.md 文件缺失（可选）"
fi

if docker compose exec -T yanblog test -d /usr/share/nginx/html/web/static 2>/dev/null; then
    check_pass "静态资源目录存在"
else
    check_fail "静态资源目录缺失"
fi
echo ""

# 7. 检查端口监听
echo "7. 检查端口监听..."
if docker compose exec -T yanblog netstat -tuln 2>/dev/null | grep -q ":80"; then
    check_pass "端口 80 (前端) 监听中"
else
    check_fail "端口 80 未监听"
fi

if docker compose exec -T yanblog netstat -tuln 2>/dev/null | grep -q ":81"; then
    check_pass "端口 81 (后台) 监听中"
else
    check_fail "端口 81 未监听"
fi

if docker compose exec -T yanblog netstat -tuln 2>/dev/null | grep -q ":8080"; then
    check_pass "端口 8080 (API) 监听中"
else
    check_fail "端口 8080 未监听"
fi
echo ""

# 8. 测试 HTTP 访问
echo "8. 测试 HTTP 访问..."

# 获取容器 IP
CONTAINER_IP=$(docker compose ps yanblog -q | xargs docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}')

# 测试前端
if docker compose exec -T yanblog wget -q -O- http://127.0.0.1:80 2>/dev/null | grep -q "<!DOCTYPE html>"; then
    check_pass "前端页面可访问 (http://localhost:3002)"
else
    check_fail "前端页面访问失败"
fi

# 测试后台
if docker compose exec -T yanblog wget -q -O- http://127.0.0.1:81 2>/dev/null | grep -q "<!DOCTYPE html>"; then
    check_pass "后台页面可访问 (http://localhost:3011)"
else
    check_fail "后台页面访问失败"
fi

# 测试 API
if docker compose exec -T yanblog wget -q -O- http://127.0.0.1:8080/api/v1/category 2>/dev/null | grep -q "data"; then
    check_pass "API 接口可访问 (http://localhost:8080)"
else
    check_warn "API 接口访问失败（可能数据库未初始化）"
fi
echo ""

# 9. 检查数据卷
echo "9. 检查数据卷..."
if docker volume ls | grep -q "yanblog_mysql_data"; then
    check_pass "MySQL 数据卷已创建"
else
    check_fail "MySQL 数据卷缺失"
fi

if docker volume ls | grep -q "yanblog_uploads_data"; then
    check_pass "上传文件数据卷已创建"
else
    check_fail "上传文件数据卷缺失"
fi
echo ""

echo "==================================="
echo "健康检查完成！"
echo "==================================="
echo ""
echo "访问地址："
echo "  前端: http://localhost:3002"
echo "  后台: http://localhost:3011"
echo ""
echo "如有问题，请查看："
echo "  docker compose logs -f"
echo ""
