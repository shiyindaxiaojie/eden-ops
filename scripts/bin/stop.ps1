# eden* 停止服务脚本 (PowerShell版本)

# 设置窗口标题
$host.UI.RawUI.WindowTitle = "eden* 停止服务"

# 输出标题
Write-Host "
=====================================" -ForegroundColor Red
Write-Host "    eden* 停止服务脚本" -ForegroundColor Red
Write-Host "=====================================" -ForegroundColor Red
Write-Host ""

# 停止服务
Write-Host "[信息] 正在停止服务..." -ForegroundColor Cyan
Write-Host ""

$stoppedBackend = $false
$stoppedFrontend = $false

# 查找并终止后端进程
$backendProcesses = Get-Process | Where-Object { $_.MainWindowTitle -like "*eden* 后端服务*" }
if ($backendProcesses) {
    foreach ($process in $backendProcesses) {
        Write-Host "正在终止后端进程 PID: $($process.Id)" -ForegroundColor Gray
        try {
            Stop-Process -Id $process.Id -Force -ErrorAction Stop
            Write-Host "[✓] " -ForegroundColor Green -NoNewline
            Write-Host "后端服务已停止"
            $stoppedBackend = $true
        } catch {
            Write-Host "[✗] " -ForegroundColor Red -NoNewline
            Write-Host "无法停止后端服务: $($_.Exception.Message)"
        }
    }
} else {
    Write-Host "未找到运行中的后端服务进程" -ForegroundColor Gray
}

# 查找并终止前端进程
$frontendProcesses = Get-Process | Where-Object { $_.MainWindowTitle -like "*eden* 前端服务*" }
if ($frontendProcesses) {
    foreach ($process in $frontendProcesses) {
        Write-Host "正在终止前端进程 PID: $($process.Id)" -ForegroundColor Gray
        try {
            Stop-Process -Id $process.Id -Force -ErrorAction Stop
            Write-Host "[✓] " -ForegroundColor Green -NoNewline
            Write-Host "前端服务已停止"
            $stoppedFrontend = $true
        } catch {
            Write-Host "[✗] " -ForegroundColor Red -NoNewline
            Write-Host "无法停止前端服务: $($_.Exception.Message)"
        }
    }
} else {
    Write-Host "未找到运行中的前端服务进程" -ForegroundColor Gray
}

# 检查是否还有相关进程在运行
$remainingBackendProcesses = Get-Process | Where-Object { $_.MainWindowTitle -like "*eden* 后端服务*" }
if ($remainingBackendProcesses) {
    Write-Host ""
    Write-Host "[警告] " -ForegroundColor Yellow -NoNewline
    Write-Host "仍有后端服务进程在运行"
}

$remainingFrontendProcesses = Get-Process | Where-Object { $_.MainWindowTitle -like "*eden* 前端服务*" }
if ($remainingFrontendProcesses) {
    Write-Host ""
    Write-Host "[警告] " -ForegroundColor Yellow -NoNewline
    Write-Host "仍有前端服务进程在运行"
}

# 停止Docker服务（如果有）
Write-Host ""
Write-Host "[信息] 检查是否有Docker容器在运行..." -ForegroundColor Cyan

try {
    # 检查是否安装了Docker
    $dockerVersion = docker --version
    
    # 检查是否有eden*相关的容器在运行
    $dockerContainers = docker ps --filter "name=eden-ops" --format "{{.ID}} {{.Names}}"
    
    if ($dockerContainers) {
        Write-Host "发现eden*相关的Docker容器正在运行" -ForegroundColor Yellow
        Write-Host "是否停止这些容器? (Y/N)" -ForegroundColor Yellow
        $response = Read-Host
        
        if ($response -eq "Y" -or $response -eq "y") {
            Write-Host "正在停止Docker容器..." -ForegroundColor Gray
            
            # 切换到docker目录
            $dockerDir = Join-Path -Path (Split-Path -Parent (Split-Path -Parent $PSScriptRoot)) -ChildPath "docker"
            if (Test-Path $dockerDir) {
                Set-Location -Path $dockerDir
                docker-compose down
                Write-Host "[✓] " -ForegroundColor Green -NoNewline
                Write-Host "Docker容器已停止"
            } else {
                Write-Host "[✗] " -ForegroundColor Red -NoNewline
                Write-Host "找不到Docker目录: $dockerDir"
            }
        } else {
            Write-Host "跳过停止Docker容器" -ForegroundColor Gray
        }
    } else {
        Write-Host "未发现eden*相关的Docker容器在运行" -ForegroundColor Gray
    }
} catch {
    Write-Host "Docker未安装或未运行" -ForegroundColor Gray
}

# 返回到原目录
Set-Location -Path $PSScriptRoot

Write-Host ""
if ($stoppedBackend -or $stoppedFrontend) {
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "服务停止操作完成"
} else {
    Write-Host "[信息] " -ForegroundColor Cyan -NoNewline
    Write-Host "没有找到需要停止的服务"
}

Write-Host ""
Write-Host "提示: 如果有服务无法停止，可以使用任务管理器手动结束进程" -ForegroundColor Gray
Write-Host ""

Write-Host "按任意键退出..." -ForegroundColor DarkGray
$null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")