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

# 2. 从 docker-compose.yaml 读取配置 (提取端口)
Print-Color "2. 读取 Docker 配置..."
$ComposeLines = Get-Content "docker-compose.yaml"
$CurrentService = ""
$BackendPort = "8080" # 默认值
$AdminPort = "3001"   # 默认值
$PublicPort = "3002"  # 默认值
$DbVolume = "yanblog_db_data" # 默认值

foreach ($line in $ComposeLines) {
    $trimLine = $line.Trim()
    
    # 识别服务块
    if ($trimLine -match "^backend:$") { $CurrentService = "backend" }
    elseif ($trimLine -match "^frontend-admin:$") { $CurrentService = "frontend-admin" }
    elseif ($trimLine -match "^frontend-public:$") { $CurrentService = "frontend-public" }
    elseif ($trimLine -match "^db:$") { $CurrentService = "db" }
    
    # 提取端口映射 - 优化正则以支持更多格式 (如无引号、带空格等)
    # 匹配格式: - "3001:80" 或 - 3001:80 或 - '3001:80'
    if ($trimLine -match '-\s*["'']?(\d+):8080["'']?') {
        if ($CurrentService -eq "backend") { $BackendPort = $matches[1] }
    }
    if ($trimLine -match '-\s*["'']?(\d+):80["'']?') {
        if ($CurrentService -eq "frontend-admin") { $AdminPort = $matches[1] }
        if ($CurrentService -eq "frontend-public") { $PublicPort = $matches[1] }
    }
    
    # 提取 Volume 映射
    if ($trimLine -match '- ([a-zA-Z0-9_\-]+):/var/lib/mysql') {
        if ($CurrentService -eq "db") { $DbVolume = $matches[1] }
    }
}
Write-Host "   检测到后台管理端口: $AdminPort" -ForegroundColor DarkGray

# 3. 准备构建环境 (资源同步 & 配置注入)
Print-Color "3. 准备构建环境..."

# 定义需要修改的文件和备份路径
$FrontendConfigFile = "web/frontend/public/config.yaml"
$FrontendConfigBak = "web/frontend/public/config.yaml.bak"

# 标记是否需要还原
$RestoreFrontendConfig = $false

try {
    # --- 3.1 同步静态资源 (Frontend -> Backend) ---
    # 解决后台管理系统引用前台图片资源 404 的问题
    Print-Color "   正在同步静态资源..." "Gray"
    $BackendPublicStatic = "web/backend/public/static"
    
    if (!(Test-Path $BackendPublicStatic)) { 
        New-Item -ItemType Directory -Path $BackendPublicStatic -Force | Out-Null 
    }
    
    # 复制 static 目录
    if (Test-Path "web/frontend/public/static") {
        Copy-Item -Path "web/frontend/public/static/*" -Destination $BackendPublicStatic -Recurse -Force
    }
    
    # 复制 favicon
    if (Test-Path "web/frontend/public/favicon.svg") { 
        Copy-Item -Path "web/frontend/public/favicon.svg" -Destination "web/backend/public/" -Force 
    }

    # --- 3.2 注入端口配置 (核心逻辑) ---
    # 临时修改配置文件，将 docker-compose 中的端口写入，确保构建出的镜像包含正确端口
    if (Test-Path $FrontendConfigFile) {
        # A. 备份原始文件
        Copy-Item -Path $FrontendConfigFile -Destination $FrontendConfigBak -Force
        $RestoreFrontendConfig = $true
        
        # B. 读取并替换端口
        $ConfContent = Get-Content $FrontendConfigFile -Raw -Encoding UTF8
        
        # 使用正则替换，同时支持 key: value 和 key:value
        if ($ConfContent -match "(?m)^dev_admin_port:\s*\d+") {
            $ConfContent = $ConfContent -replace "(?m)^dev_admin_port:\s*\d+", "dev_admin_port: $AdminPort"
            $ConfContent | Set-Content $FrontendConfigFile -Encoding UTF8
            
            # 验证修改
            $CheckLine = Get-Content $FrontendConfigFile | Select-String "dev_admin_port"
            Write-Host "   [Build Prep] 已注入端口配置: $CheckLine" -ForegroundColor Cyan
        } else {
            Write-Warning "   未在配置文件中找到 'dev_admin_port' 字段，跳过注入。"
            Write-Host "   DEBUG: 文件内容预览 (前500字符):" -ForegroundColor Gray
            Write-Host $ConfContent.Substring(0, [Math]::Min($ConfContent.Length, 500))
        }
    } else {
        Write-Warning "   找不到前端配置文件: $FrontendConfigFile"
    }

    # 4. 执行部署
    Print-Color "4. 构建并启动 Docker 服务..."
    Print-Color "   (使用 Named Volumes 存储数据，解决 Windows 文件占用问题)" "Gray"
    
    # 重新构建 (-build) 确保修改后的 config 被打入镜像
    docker-compose up -d --build
}
catch {
    Print-Color "`n[严重错误] Docker 启动失败。" "Red"
    Write-Error $_
    # 不需要 exit，以便 finally 块执行清理
}
finally {
    # 5. 清理和还原
    if ($RestoreFrontendConfig -and (Test-Path $FrontendConfigBak)) {
        # 还原配置文件，防止 git 将临时修改识别为变更
        Copy-Item -Path $FrontendConfigBak -Destination $FrontendConfigFile -Force
        Remove-Item $FrontendConfigBak -Force
        Write-Host "   [Build Cleanup] 已还原前端配置文件至初始状态" -ForegroundColor DarkGray
    }
}

# 6. 检查结果
if ($LASTEXITCODE -eq 0) {
    Print-Color "`n>>> 部署成功！" "Green"
    
    Print-Color "`n访问地址:" "White"
    Print-Color "-------------------------------------------"
    Print-Color "   DB 数据:     存储在 Docker Volume '$DbVolume' 中"
    Print-Color "   API 后端:    http://localhost:$BackendPort"
    Print-Color "   后台管理:    http://localhost:$AdminPort"
    Print-Color "   前台页面:    http://localhost:$PublicPort"
    Print-Color "-------------------------------------------"
}
