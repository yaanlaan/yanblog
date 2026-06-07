$CONFIG_DIR = "config"
$BACKEND_CONFIG = "$CONFIG_DIR/backend/config.yaml"
$FRONTEND_CONFIG = "$CONFIG_DIR/frontend/config.yaml"
$BACKEND_TEMPLATE = "$CONFIG_DIR/config_template.yaml"
$FRONTEND_TEMPLATE = "$CONFIG_DIR/frontend/config_template.yaml"
$OLD_CONFIG = "config/config.yaml"
$OLD_FRONTEND_CONFIG = "web/frontend/public/config.yaml"

Write-Host "=== YanBlog Configuration Initialization ==="

if (-not (Test-Path "$CONFIG_DIR/backend")) {
    New-Item -ItemType Directory -Path "$CONFIG_DIR/backend" | Out-Null
}
if (-not (Test-Path "$CONFIG_DIR/frontend")) {
    New-Item -ItemType Directory -Path "$CONFIG_DIR/frontend" | Out-Null
}

if (Test-Path $OLD_CONFIG -and -not (Test-Path $BACKEND_CONFIG)) {
    Write-Host "Migrating old backend config..."
    Copy-Item $OLD_CONFIG $BACKEND_CONFIG
    Write-Host "Old config migrated to $BACKEND_CONFIG"
} elseif (Test-Path $BACKEND_TEMPLATE -and -not (Test-Path $BACKEND_CONFIG)) {
    Copy-Item $BACKEND_TEMPLATE $BACKEND_CONFIG
    Write-Host "$BACKEND_TEMPLATE 已复制为 $BACKEND_CONFIG"
} elseif (-not (Test-Path $BACKEND_CONFIG)) {
    $defaultConfig = @"
server:
  AppMode: debug
  HttpPort: :8080
  SiteUrl: ""

database:
  Db: MYSQL
  DbHost: db
  DbPort: 3306
  DbUser: root
  DbPassWord: rootpassword
  DbName: yanblog

JwtKey: ""

weather:
  Provider: openweathermap
  ApiKey: ""
  DefaultCity: Hefei

FrontEndConfigPath: config/frontend/config.yaml
"@
    $defaultConfig | Out-File -FilePath $BACKEND_CONFIG -Encoding utf8
    Write-Host "Created default backend config at $BACKEND_CONFIG"
}

if (Test-Path $OLD_FRONTEND_CONFIG -and -not (Test-Path $FRONTEND_CONFIG)) {
    Write-Host "Migrating old frontend config..."
    Copy-Item $OLD_FRONTEND_CONFIG $FRONTEND_CONFIG
    Write-Host "Old frontend config migrated to $FRONTEND_CONFIG"
} elseif (Test-Path $FRONTEND_TEMPLATE -and -not (Test-Path $FRONTEND_CONFIG)) {
    Copy-Item $FRONTEND_TEMPLATE $FRONTEND_CONFIG
    Write-Host "$FRONTEND_TEMPLATE 已复制为 $FRONTEND_CONFIG"
} elseif (-not (Test-Path $FRONTEND_CONFIG)) {
    Write-Host "Error: No frontend config template found!" -ForegroundColor Red
    exit 1
}

Write-Host "=== Configuration Initialization Complete ==="
Write-Host ""
Write-Host "Please edit the following files to configure your blog:"
Write-Host "- Backend: $BACKEND_CONFIG"
Write-Host "- Frontend: $FRONTEND_CONFIG"
Write-Host ""
Write-Host "Important:"
Write-Host "- Set JwtKey in backend config (run: openssl rand -hex 32)"
Write-Host "- Configure database connection if not using Docker"
Write-Host "- Customize frontend settings (blog name, author info, etc.)"