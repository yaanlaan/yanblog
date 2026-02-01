#!/bin/sh

# 启动 Go 后端 (后台运行)
# 确保后端已经编译为 server 可执行文件
echo "Starting Backend..."
./server &

# 启动 Nginx (前台运行)
echo "Starting Nginx..."
nginx -g "daemon off;"
