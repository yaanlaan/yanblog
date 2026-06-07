#!/bin/sh
set -e

# 检测数据库类型，决定是否需要等待 MySQL
DB_TYPE=$(grep -i 'Db:' /app/config/config.yaml 2>/dev/null | awk '{print $2}' | tr '[:lower:]' '[:upper:]')

if [ "$DB_TYPE" = "MYSQL" ] || [ "$DB_TYPE" = "MARIADB" ]; then
    echo "检测到数据库类型为 MySQL，等待数据库就绪..."
    max_attempts=30
    attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if nc -z db 3306 2>/dev/null; then
            echo "数据库已就绪！"
            break
        fi
        attempt=$((attempt + 1))
        echo "等待数据库启动... ($attempt/$max_attempts)"
        sleep 2
    done

    if [ $attempt -eq $max_attempts ]; then
        echo "错误：数据库连接超时！"
        exit 1
    fi

    # 额外等待 MySQL 完全初始化
    sleep 5
else
    echo "数据库类型为 $DB_TYPE，跳过 MySQL 等待。"
fi

# 启动 Go 后端 (后台运行)
echo "启动后端服务..."
./server &
BACKEND_PID=$!

# 等待后端启动
sleep 3

# 检查后端是否运行
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "错误：后端启动失败！"
    exit 1
fi

echo "后端启动成功 (PID: $BACKEND_PID)"

# 启动 Nginx (前台运行)
echo "启动 Nginx..."
nginx -g "daemon off;"
