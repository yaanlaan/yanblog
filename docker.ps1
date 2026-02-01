# Ensure LF line endings
Write-Host ">>> Starting YanBlog Deployment (Single Container Mode)..." -ForegroundColor Green

# --- 0. Auto-Configure Ports ---
$ComposeFile = "docker-compose.yaml"
$ConfigFile = "web/frontend/public/config.yaml"

if (Test-Path $ComposeFile) {
    # Extract Admin Port (mapped to internal 81)
    # Regex looks for: - "3011:81" or - 3011:81
    $Content = Get-Content $ComposeFile -Raw
    if ($Content -match '-\s*"?(\d+):81"?') {
        $AdminPort = $matches[1]
        Write-Host ">>> Detected Admin Port from docker-compose: $AdminPort" -ForegroundColor Cyan
        
        if (Test-Path $ConfigFile) {
            $ConfigContent = Get-Content $ConfigFile -Raw
            # Update dev_admin_port
            if ($ConfigContent -match 'dev_admin_port:\s*\d+') {
                $NewConfig = $ConfigContent -replace 'dev_admin_port:\s*\d+', "dev_admin_port: $AdminPort"
                Set-Content -Path $ConfigFile -Value $NewConfig -NoNewline
                Write-Host ">>> Updated $ConfigFile with port $AdminPort" -ForegroundColor Green
            }
        }
    }
}
# -----------------------------

# 1. Cleaning up old containers
Write-Host "1. Removing old containers..." -ForegroundColor Cyan
docker compose down
if ($LASTEXITCODE -ne 0) {
    Write-Host ">>> Failed to remove old containers." -ForegroundColor Red
    exit $LASTEXITCODE
}

# 2. Build and Start
Write-Host "2. Building and Starting..." -ForegroundColor Cyan
# Rebuild forced to update Go binary with new paths
docker compose up -d --build
if ($LASTEXITCODE -ne 0) {
    Write-Host ">>> Deployment Failed! See error details above." -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host ">>> Deployment Complete!" -ForegroundColor Green
Write-Host "    Frontend: http://localhost:3002" -ForegroundColor Green
Write-Host "    Admin:    http://localhost:3011" -ForegroundColor Green

