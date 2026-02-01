#!/bin/sh
set -e

# 等待数据库就绪
echo "Waiting for database..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if nc -z db 3306 2>/dev/null; then
        echo "Database is ready!"
        break
    fi
    attempt=$((attempt + 1))
    echo "等待数据库启动... ($attempt/$max_attempts)"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "Error: Database connection timeout!"
    exit 1
fi

# 额外等待 MySQL 完全初始化
sleep 5

# 启动 Go 后端 (后台运行)
echo "Starting Backend..."
./server &
BACKEND_PID=$!

# 等待后端启动
sleep 3

# 检查后端是否运行
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "Error: Backend failed to start!"
    exit 1
fi

echo "Backend started successfully (PID: $BACKEND_PID)"

# 启动 Nginx (前台运行)
echo "Starting Nginx..."
nginx -g "daemon off;"
