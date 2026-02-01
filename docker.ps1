# docker.ps1
# 自动化 Docker 部署脚本 (优化版 - 解决挂载冲突)
# 策略：
# 1. 使用 Docker Named Volumes 替代 Bind Mounts 解决 Windows 权限/路径问题
# 2. 配置文件直接构建到镜像中，不再依赖运行时挂载

$ErrorActionPreference = "Stop"

function Print-Color {
    param([string]$Message, [ConsoleColor]$Color = "Cyan")
    Write-Host $Message -ForegroundColor $Color
}

Print-Color ">>> 开始 Docker 部署流程..."

# 1. 环境清理
Print-Color "1. 清理旧环境 (docker-compose down)..."
try {
    # 移除容器、网络和所有相关资源
    docker-compose down --remove-orphans 2>&1 | Out-Null
    docker-compose rm -f 2>&1 | Out-Null
} catch {
    Write-Warning "清理过程遇到轻微错误（可忽略）"
}

# 3. 准备配置文件（供 Build 使用）
# 虽然我们不再挂载文件，但在 Build 之前准备好文件确保它们被正确 COPY 进镜像
Print-Color "2. 准备构建上下文 (配置文件)..."

# 确保目标目录存在（用于暂存配置）
$DockerFieldPath = Join-Path (Get-Item .).FullName "docker_field"
$DirsToCreate = @(
    "$DockerFieldPath\frontend",
    "$DockerFieldPath\backend"
)
foreach ($Dir in $DirsToCreate) {
    if (!(Test-Path $Dir)) {
        New-Item -ItemType Directory -Path $Dir -Force | Out-Null
    }
}

# 复制最新的配置到 docker_field (用于 Backend 构建)
$ConfigMap = @{
    "web/frontend/public/config.yaml" = "$DockerFieldPath/frontend/config.yaml"
    "config/config.yaml"              = "$DockerFieldPath/backend/config.yaml"
}

foreach ($Src in $ConfigMap.Keys) {
    $Dest = $ConfigMap[$Src]
    if (Test-Path $Src) {
        Copy-Item -Path $Src -Destination $Dest -Force
        Write-Host "   配置已就绪: $Src" -ForegroundColor DarkGray
    } else {
        Print-Color "   [警告] 缺配置文件: $Src (请先运行 init.ps1)" "Yellow"
    }
}

# 4. 执行部署
Print-Color "3. 构建并启动 Docker 服务..."
Print-Color "   (使用 Named Volumes 存储数据，解决 Windows 文件占用问题)" "Gray"

try {
    # 重新构建 (-build) 确保最新的 config 被打入镜像
    docker-compose up -d --build
} catch {
    Print-Color "`n[严重错误] Docker 启动失败。" "Red"
    Write-Error $_
    exit 1
}

# 5. 检查结果
if ($LASTEXITCODE -eq 0) {
    Print-Color "`n>>> 部署成功！" "Green"
    
    Print-Color "`n访问地址:" "White"
    Print-Color "-------------------------------------------"
    Print-Color "   DB 数据:     存储在 Docker Volume 'yanblog_db_data' 中"
    Print-Color "   API 后端:    http://localhost:8080"
    Print-Color "   后台管理:    http://localhost:3001"
    Print-Color "   前台页面:    http://localhost:3002"
    Print-Color "-------------------------------------------"
}
