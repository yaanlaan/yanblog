$COMPOSE_FILE = "docker-compose.yaml"
$CONFIG_DIR = "config"
$BACKEND_CONFIG = "$CONFIG_DIR/backend/config.yaml"
$FRONTEND_CONFIG = "$CONFIG_DIR/frontend/config.yaml"

Write-Host ">>> Starting YanBlog Deployment (Single Container Mode)..."

if (-not (Test-Path "$CONFIG_DIR/backend") -or -not (Test-Path $BACKEND_CONFIG)) {
    Write-Host ">>> Initializing configuration..."
    .\setup.bat
}

if (Test-Path $COMPOSE_FILE) {
    $content = Get-Content $COMPOSE_FILE -Raw
    $match = [regex]::Match($content, '- "\?([0-9]*):81"')
    if ($match.Success) {
        $ADMIN_PORT = $match.Groups[1].Value
        Write-Host ">>> Detected Admin Port from docker-compose: $ADMIN_PORT"
        
        if (Test-Path $FRONTEND_CONFIG) {
            $configContent = Get-Content $FRONTEND_CONFIG -Raw
            $configContent = $configContent -replace 'dev_admin_port: [0-9]*', "dev_admin_port: $ADMIN_PORT"
            $configContent | Out-File -FilePath $FRONTEND_CONFIG -Encoding utf8 -NoNewline
            Write-Host ">>> Updated $FRONTEND_CONFIG with port $ADMIN_PORT"
        }
    }
}

Write-Host ">>> Removing old containers..."
docker compose down

Write-Host ">>> Building and Starting..."
docker compose up -d --build

Write-Host ">>> Deployment Complete!"
Write-Host "    Frontend: http://localhost:3002"
Write-Host "    Admin:    http://localhost:3011"
Write-Host "    API:      http://localhost:8080"
Write-Host ""
Write-Host ">>> Configuration Files:"
Write-Host "    Backend: $BACKEND_CONFIG"
Write-Host "    Frontend: $FRONTEND_CONFIG"
Write-Host ""
Write-Host ">>> Tips:"
Write-Host "    - Modify config files and restart container to apply changes"
Write-Host "    - Use environment variables to override config values"
Write-Host "    - Run 'docker compose logs -f' to view logs"