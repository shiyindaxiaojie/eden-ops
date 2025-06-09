# eden* Docker一键启动脚本 (PowerShell版本)

# 设置窗口标题
$host.UI.RawUI.WindowTitle = "eden* Docker一键启动"

# 输出标题
Write-Host "
===================================="  -ForegroundColor Blue
Write-Host "  eden* Docker一键启动脚本" -ForegroundColor Blue
Write-Host "===================================="  -ForegroundColor Blue
Write-Host ""

# 检查Docker是否安装
try {
    $dockerVersion = docker --version
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "Docker已安装: $dockerVersion"
} catch {
    Write-Host "[✗] " -ForegroundColor Red -NoNewline
    Write-Host "Docker未安装或未运行，请先安装并启动Docker。"
    Write-Host "下载地址: https://www.docker.com/products/docker-desktop"
    Write-Host "按任意键退出..."
    $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit
}

# 检查Docker Compose是否可用
try {
    $dockerComposeVersion = docker-compose --version
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "Docker Compose已安装: $dockerComposeVersion"
} catch {
    Write-Host "[✗] " -ForegroundColor Red -NoNewline
    Write-Host "Docker Compose未安装或未运行。"
    Write-Host "Docker Desktop通常包含Docker Compose，请确保正确安装。"
    Write-Host "按任意键退出..."
    $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit
}

Write-Host ""
Write-Host "[1/1] " -ForegroundColor Yellow -NoNewline
Write-Host "正在使用Docker Compose启动eden*..."

# 切换到docker目录
Set-Location -Path "$PSScriptRoot\..\..\docker"

# 启动服务
try {
    docker-compose up -d
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "[✓] " -ForegroundColor Green -NoNewline
        Write-Host "eden*已在Docker中成功启动！"
        
        # 获取容器状态
        $containerStatus = docker-compose ps
        Write-Host ""
        Write-Host "容器状态:" -ForegroundColor Yellow
        Write-Host "$containerStatus"
        
        Write-Host ""
        Write-Host "访问地址: " -NoNewline
        Write-Host "http://localhost:8080" -ForegroundColor Cyan
    } else {
        Write-Host ""
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "Docker Compose启动失败，请检查错误信息。"
    }
} catch {
    Write-Host ""
    Write-Host "[✗] " -ForegroundColor Red -NoNewline
    Write-Host "执行Docker Compose时出错: $_"
}

Write-Host ""
Write-Host "提示: 要停止服务，请运行: " -NoNewline
Write-Host "docker-compose down" -ForegroundColor Yellow
Write-Host ""

Write-Host "按任意键退出此窗口..." -ForegroundColor DarkGray
$null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")