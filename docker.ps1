# docker.ps1
# 用于在 Docker 构建前准备配置文件，并自动运行 docker-compose

# 创建目标目录（如果不存在）
if (!(Test-Path -Path "docker_field/frontend")) {
    New-Item -ItemType Directory -Path "docker_field/frontend" | Out-Null
}
if (!(Test-Path -Path "docker_field/backend")) {
    New-Item -ItemType Directory -Path "docker_field/backend" | Out-Null
}

# 复制前端配置
Copy-Item -Path "web/frontend/public/config.yaml" -Destination "docker_field/frontend/config.yaml" -Force

# 复制后端配置
Copy-Item -Path "config/config.yaml" -Destination "docker_field/backend/config.yaml" -Force

Write-Host "配置文件已复制到 docker_field 目录下。"

# 自动运行 docker-compose 部署
Write-Host "即将执行 Docker 部署..."
docker-compose up -d --build
Write-Host "Docker 服务已启动。"
