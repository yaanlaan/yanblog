@echo off
chcp 65001 >nul
echo =========================================
echo   博客系统初始化向导
echo =========================================
echo.

REM 检查配置文件是否存在
if not exist "config\backend\config.yaml" (
    echo ✓ 检测到 config\backend\config.yaml 不存在，正在创建...
    copy "config\config_template.yaml" "config\backend\config.yaml"
    echo ✓ 已创建 config\backend\config.yaml
    echo.
) else (
    echo ℹ️  config\backend\config.yaml 已存在
    echo.
)

REM 尝试生成 JWT 密钥
echo  正在生成 JWT 密钥...
where openssl >nul 2>nul
if %errorlevel% equ 0 (
    for /f "delims=" %%i in ('openssl rand -hex 32') do set JWT_KEY=%%i
    echo ✓ JWT 密钥已生成
    REM 使用 PowerShell 替换 JwtKey
    powershell -Command "(Get-Content 'config\backend\config.yaml') -replace '^JwtKey:.*', 'JwtKey: %JWT_KEY%' | Set-Content 'config\backend\config.yaml'"
    echo ✓ JWT 密钥已写入配置文件
    echo.
) else (
    echo ⚠️  未检测到 OpenSSL，无法自动生成 JWT 密钥
    echo.
    echo 请手动执行以下步骤：
    echo   1. 打开 config\backend\config.yaml 文件
    echo   2. 找到 JwtKey: 这一行
    echo   3. 访问 https://randomkeygen.com/ 生成一个随机字符串
    echo   4. 将生成的字符串填入 JwtKey 字段
    echo.
)

REM 提示用户
echo =========================================
echo   配置完成！
echo =========================================
echo.
echo 请检查以下配置项：
echo   1. database.DbPassWord - 数据库密码（强烈建议修改）
echo   2. database.DbHost - Docker 部署为 db，本地部署为 localhost
echo   3. server.SiteUrl - 网站域名（可选）
echo.
echo 下一步：
echo   - Docker 部署：运行 .\docker.ps1
echo   - 本地部署：运行 go run main.go
echo.
echo 详细使用说明请查看：QUICK_START.md
echo =========================================
echo.
pause
