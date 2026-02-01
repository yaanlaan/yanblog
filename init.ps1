# init.ps1
# 用于初始化配置文件，将 config_template.yaml 复制为 config.yaml

# 后端配置模板复制
if (Test-Path "config/config_template.yaml") {
    Copy-Item -Path "config/config_template.yaml" -Destination "config/config.yaml" -Force
    Write-Host "config/config_template.yaml 已复制为 config/config.yaml"
} else {
    Write-Host "未找到 config/config_template.yaml，未进行复制。"
}

# 前端配置模板复制
if (Test-Path "web/frontend/public/config_template.yaml") {
    Copy-Item -Path "web/frontend/public/config_template.yaml" -Destination "web/frontend/public/config.yaml" -Force
    Write-Host "web/frontend/public/config_template.yaml 已复制为 web/frontend/public/config.yaml"
} else {
    Write-Host "未找到 web/frontend/public/config_template.yaml，未进行复制。"
}