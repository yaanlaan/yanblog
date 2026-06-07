@echo off
chcp 65001 >nul
echo =========================================
echo   配置检查工具
echo =========================================
echo.

REM 检查配置文件
if not exist "config\backend\config.yaml" (
    echo ❌ 配置文件不存在：config\backend\config.yaml
    echo    解决方法：运行 .\setup.bat 或复制 config\config_template.yaml
    pause
    exit /b 1
) else (
    echo ✓ 配置文件存在
)

REM 使用 PowerShell 检查配置
powershell -Command ^
    "$config = Get-Content 'config\backend\config.yaml' -Raw; ^
    $dbPass = if ($config -match 'DbPassWord:\s*(.+)') { $matches[1].Trim() } else { '' }; ^
    if ($dbPass -eq 'rootpassword') { ^
        Write-Host '️  警告：数据库密码仍为默认值，建议修改' -ForegroundColor Yellow ^
    } else { ^
        Write-Host '✓ 数据库密码已修改' -ForegroundColor Green ^
    }; ^
    $jwtKey = if ($config -match 'JwtKey:\s*(.+)') { $matches[1].Trim() } else { '' }; ^
    if ([string]::IsNullOrWhiteSpace($jwtKey)) { ^
        Write-Host '❌ JWT 密钥未设置' -ForegroundColor Red; ^
        Write-Host '   解决方法：运行 .\setup.bat 自动生成' -ForegroundColor Yellow; ^
        Write-Host '   或手动运行：openssl rand -hex 32' -ForegroundColor Yellow ^
    } elseif ($jwtKey.Length -lt 32) { ^
        Write-Host \"⚠️  警告：JWT 密钥长度不足（当前：$($jwtKey.Length) 位，建议：32 位）\" -ForegroundColor Yellow ^
    } else { ^
        Write-Host \"✓ JWT 密钥已设置（长度：$($jwtKey.Length) 位）\" -ForegroundColor Green ^
    }; ^
    $dbHost = if ($config -match 'DbHost:\s*(.+)') { $matches[1].Trim() } else { '未知' }; ^
    Write-Host \"ℹ️  数据库主机：$dbHost\"; ^
    $httpPort = if ($config -match 'HttpPort:\s*(.+)') { $matches[1].Trim() } else { '未知' }; ^
    Write-Host \"ℹ️  服务端口：$httpPort\"; ^
    Write-Host ''; ^
    Write-Host '=========================================' -ForegroundColor Cyan; ^
    Write-Host '  检查完成！' -ForegroundColor Green; ^
    Write-Host '=========================================' -ForegroundColor Cyan; ^
    Write-Host ''; ^
    if (Test-Path 'go.mod') { ^
        Write-Host '下一步：' -ForegroundColor Cyan; ^
        Write-Host '  - Docker 部署：.\docker.ps1' -ForegroundColor White; ^
        Write-Host '  - 本地部署：go run main.go' -ForegroundColor White ^
    } else { ^
        Write-Host '提示：请确保在项目根目录运行此脚本' -ForegroundColor Yellow ^
    }"

echo.
pause
