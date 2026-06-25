# YanBlog 部署验证脚本
# 用途：每次部署后自动验证网站可用性

$baseUrl = "http://localhost:3002"
$allPassed = $true

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  YanBlog 部署验证" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 1. 核心页面测试
Write-Host "[1] 核心页面测试" -ForegroundColor Yellow
$pages = @(
    @{url="$baseUrl/"; name="前台首页"},
    @{url="$baseUrl/admin/"; name="后台管理"},
    @{url="$baseUrl/api/v1/health"; name="API健康检查"}
)

foreach($p in $pages) {
    try {
        $response = Invoke-WebRequest -Uri $p.url -UseBasicParsing -TimeoutSec 5
        $code = $response.StatusCode
        if($code -eq 200) {
            Write-Host "  ✓ $($p.name): HTTP $code" -ForegroundColor Green
        } else {
            Write-Host "  ✗ $($p.name): HTTP $code" -ForegroundColor Red
            $allPassed = $false
        }
    } catch {
        Write-Host "  ✗ $($p.name): $($_.Exception.Message)" -ForegroundColor Red
        $allPassed = $false
    }
}
Write-Host ""

# 2. 静态资源测试
Write-Host "[2] 静态资源测试" -ForegroundColor Yellow
$resources = @(
    @{url="$baseUrl/admin/assets/index-CMJ96UUM.js"; name="后台JS"},
    @{url="$baseUrl/admin/assets/index-CGkevFw1.css"; name="后台CSS"},
    @{url="$baseUrl/admin/favicon.svg"; name="后台Favicon"},
    @{url="$baseUrl/uploads/defaults/hero.jpg"; name="默认背景图"},
    @{url="$baseUrl/uploads/defaults/avatar.jpg"; name="默认头像"}
)

foreach($r in $resources) {
    try {
        $response = Invoke-WebRequest -Uri $r.url -UseBasicParsing -TimeoutSec 5
        $code = $response.StatusCode
        if($code -eq 200) {
            Write-Host "  ✓ $($r.name): HTTP $code" -ForegroundColor Green
        } else {
            Write-Host "  ✗ $($r.name): HTTP $code" -ForegroundColor Red
            $allPassed = $false
        }
    } catch {
        Write-Host "  ✗ $($r.name): $($_.Exception.Message)" -ForegroundColor Red
        $allPassed = $false
    }
}
Write-Host ""

# 3. Gzip 压缩测试
Write-Host "[3] Gzip 压缩测试" -ForegroundColor Yellow
try {
    # 使用 curl.exe 避免 PowerShell 自动解压
    $tempFile = [System.IO.Path]::GetTempFileName()
    curl.exe -H "Accept-Encoding: gzip, deflate, br" -s -D - "$baseUrl/admin/assets/index-CMJ96UUM.js" -o $tempFile 2>&1 | Out-String | Out-Null
    $output = curl.exe -H "Accept-Encoding: gzip, deflate, br" -s -D - "$baseUrl/admin/assets/index-CMJ96UUM.js" -o $tempFile 2>&1
    Remove-Item $tempFile -Force -ErrorAction SilentlyContinue
    
    $outputStr = $output | Out-String
    $match = [regex]::Match($outputStr, 'Content-Encoding:\s*(\S+)')
    
    if($match.Success -and $match.Groups[1].Value -eq "gzip") {
        Write-Host "  ✓ Gzip 压缩已启用" -ForegroundColor Green
    } else {
        Write-Host "  ⚠ Gzip 未检测到，可能文件已压缩或太小" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ✗ Gzip 测试失败: $($_.Exception.Message)" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# 4. 缓存策略测试
Write-Host "[4] 缓存策略测试" -ForegroundColor Yellow
try {
    # 使用 curl.exe 获取完整响应头
    $tempFile = [System.IO.Path]::GetTempFileName()
    $output = curl.exe -s -D - "$baseUrl/admin/assets/index-CMJ96UUM.js" -o $tempFile 2>&1
    Remove-Item $tempFile -Force -ErrorAction SilentlyContinue
    
    $outputStr = $output | Out-String
    $match = [regex]::Match($outputStr, 'Cache-Control:\s*(.+?)\r?\n')
    
    if($match.Success) {
        $cacheControl = $match.Groups[1].Value.Trim()
        # 检查是否有 max-age 或 immutable
        if($cacheControl.Contains("max-age") -or $cacheControl.Contains("immutable")) {
            Write-Host "  ✓ 静态资源缓存策略: $cacheControl" -ForegroundColor Green
        } else {
            Write-Host "  ⚠ 缓存策略可能不完整 (Cache-Control: $cacheControl)" -ForegroundColor Yellow
        }
    } else {
        Write-Host "  ✗ 未找到 Cache-Control 响应头" -ForegroundColor Red
        $allPassed = $false
    }
} catch {
    Write-Host "  ✗ 缓存测试失败: $($_.Exception.Message)" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# 5. API 功能测试
Write-Host "[5] API 功能测试" -ForegroundColor Yellow
$apis = @(
    @{url="$baseUrl/api/v1/article?page=1&size=1"; name="文章列表"},
    @{url="$baseUrl/api/v1/tags"; name="标签列表"},
    @{url="$baseUrl/api/v1/category"; name="分类列表"},
    @{url="$baseUrl/api/v1/article/archive"; name="归档列表"},
    @{url="$baseUrl/api/v1/article/info/1"; name="文章详情"}
)

foreach($a in $apis) {
    try {
        $response = Invoke-WebRequest -Uri $a.url -UseBasicParsing -TimeoutSec 5
        $body = $response.Content | ConvertFrom-Json
        if($response.StatusCode -eq 200 -and $body.status -eq 200) {
            Write-Host "  ✓ $($a.name): 正常" -ForegroundColor Green
        } else {
            Write-Host "  ✗ $($a.name): status=$($body.status)" -ForegroundColor Red
            $allPassed = $false
        }
    } catch {
        Write-Host "  ✗ $($a.name): $($_.Exception.Message)" -ForegroundColor Red
        $allPassed = $false
    }
}
Write-Host ""

# 总结
Write-Host "========================================" -ForegroundColor Cyan
if($allPassed) {
    Write-Host "  ✓ 所有测试通过！" -ForegroundColor Green
} else {
    Write-Host "  ✗ 部分测试失败，请检查上方日志" -ForegroundColor Red
}
Write-Host "========================================" -ForegroundColor Cyan

exit $(if($allPassed) { 0 } else { 1 })
