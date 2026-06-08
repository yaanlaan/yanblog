@echo off
chcp 65001 >nul
echo =========================================
echo   YanBlog 初始化向导
echo =========================================
echo.

REM 确保配置目录存在
if not exist "config\backend" mkdir config\backend
if not exist "config\frontend" mkdir config\frontend

REM --- 后端配置 ---
if not exist "config\backend\config.yaml" (
    echo ✓ 创建后端配置...
    if exist "config\config_template.yaml" (
        copy "config\config_template.yaml" "config\backend\config.yaml" >nul
    ) else (
        echo server:> "config\backend\config.yaml"
        echo   AppMode: debug>> "config\backend\config.yaml"
        echo   HttpPort: :8080>> "config\backend\config.yaml"
        echo   SiteUrl: "">> "config\backend\config.yaml"
        echo.>> "config\backend\config.yaml"
        echo database:>> "config\backend\config.yaml"
        echo   Db: SQLite>> "config\backend\config.yaml"
        echo   DbHost: localhost>> "config\backend\config.yaml"
        echo   DbPort: 3306>> "config\backend\config.yaml"
        echo   DbUser: root>> "config\backend\config.yaml"
        echo   DbPassWord: rootpassword>> "config\backend\config.yaml"
        echo   DbName: yanblog.db>> "config\backend\config.yaml"
        echo.>> "config\backend\config.yaml"
        echo JwtKey: "">> "config\backend\config.yaml"
        echo.>> "config\backend\config.yaml"
        echo weather:>> "config\backend\config.yaml"
        echo   Provider: openweathermap>> "config\backend\config.yaml"
        echo   ApiKey: "">> "config\backend\config.yaml"
        echo   DefaultCity: Hefei>> "config\backend\config.yaml"
        echo.>> "config\backend\config.yaml"
        echo FrontEndConfigPath: config/frontend/config.yaml>> "config\backend\config.yaml"
    )
    echo ✓ 已创建 config\backend\config.yaml
) else (
    echo ℹ️  config\backend\config.yaml 已存在，跳过
)

REM 生成 JWT 密钥
echo.
echo 🔐 正在生成 JWT 密钥...
where openssl >nul 2>nul
if %errorlevel% equ 0 (
    for /f "delims=" %%i in ('openssl rand -hex 32') do set JWT_KEY=%%i
    echo ✓ JWT 密钥已生成
    powershell -Command "(Get-Content 'config\backend\config.yaml') -replace '^JwtKey:.*', 'JwtKey: %JWT_KEY%' | Set-Content 'config\backend\config.yaml'"
    echo ✓ JWT 密钥已写入配置文件
) else (
    echo ⚠️  未检测到 OpenSSL，请手动设置 JWT 密钥
)

REM --- 前端配置 ---
if not exist "config\frontend\config.yaml" (
    echo.
    echo ✓ 创建前端配置...
    if exist "web\frontend\public\config_template.yaml" (
        copy "web\frontend\public\config_template.yaml" "config\frontend\config.yaml" >nul
    ) else if exist "web\frontend\public\config.yaml" (
        copy "web\frontend\public\config.yaml" "config\frontend\config.yaml" >nul
    ) else (
        echo ⚠️  未找到前端配置模板，请手动创建 config\frontend\config.yaml
    )
    if exist "config\frontend\config.yaml" echo ✓ 已创建 config\frontend\config.yaml
) else (
    echo ℹ️  config\frontend\config.yaml 已存在，跳过
)

echo.
echo =========================================
echo   配置完成！
echo =========================================
echo.
echo 请检查以下配置项：
echo   1. database.DbPassWord — 数据库密码
echo   2. database.DbHost — Docker 部署为 db，本地为 localhost
echo   3. server.SiteUrl — 网站域名
echo.
echo 下一步：
echo   Docker 部署：.\docker.ps1
echo   本地部署：go run main.go
echo.
echo 详细说明：QUICK_START.md
echo =========================================
