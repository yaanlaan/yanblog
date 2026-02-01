Write-Host ">>> Starting YanBlog Deployment (Single Container Mode)..." -ForegroundColor Green

# 1. Cleaning up old containers
Write-Host "1. Removing old containers..." -ForegroundColor Cyan
docker-compose down

# 2. Build and Start
Write-Host "2. Building and Starting..." -ForegroundColor Cyan
# Rebuild forced to update Go binary with new paths
docker-compose up -d --build

Write-Host ">>> Deployment Complete!" -ForegroundColor Green
Write-Host "    Frontend: http://localhost:3002" -ForegroundColor Green
Write-Host "    Admin:    http://localhost:3011" -ForegroundColor Green

