#!/bin/bash
# docker.sh
# 用于在 Docker 构建前准备配置文件（Linux 版），并自动运行 docker-compose

# 创建目标目录（如果不存在）
mkdir -p docker_field/frontend
mkdir -p docker_field/backend

# 复制前端配置
cp -f web/frontend/public/config.yaml docker_field/frontend/config.yaml

# 复制后端配置
cp -f config/config.yaml docker_field/backend/config.yaml

echo "配置文件已复制到 docker_field 目录下。"

echo "即将执行 Docker 部署..."
docker-compose up -d --build
echo "Docker 服务已启动。"
