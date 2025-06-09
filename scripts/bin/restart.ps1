# eden* 重启脚本 (PowerShell版本)

# 设置窗口标题
$host.UI.RawUI.WindowTitle = "eden* 重启"

# 输出标题
Write-Host "
=====================================" -ForegroundColor Yellow
Write-Host "      eden* 重启脚本" -ForegroundColor Yellow
Write-Host "=====================================" -ForegroundColor Yellow
Write-Host ""

# 停止服务
Write-Host "[信息] 正在停止服务..." -ForegroundColor Cyan

# 查找并终止后端进程
$backendProcesses = Get-Process | Where-Object { $_.MainWindowTitle -like "*eden* 后端服务*" }
if ($backendProcesses) {
    foreach ($process in $backendProcesses) {
        Write-Host "正在终止后端进程 PID: $($process.Id)" -ForegroundColor Gray
        Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
    }
} else {
    Write-Host "未找到运行中的后端服务进程" -ForegroundColor Gray
}

# 查找并终止前端进程
$frontendProcesses = Get-Process | Where-Object { $_.MainWindowTitle -like "*eden* 前端服务*" }
if ($frontendProcesses) {
    foreach ($process in $frontendProcesses) {
        Write-Host "正在终止前端进程 PID: $($process.Id)" -ForegroundColor Gray
        Stop-Process -Id $process.Id -Force -ErrorAction SilentlyContinue
    }
} else {
    Write-Host "未找到运行中的前端服务进程" -ForegroundColor Gray
}

Write-Host "[信息] 服务已停止" -ForegroundColor Cyan
Write-Host ""

Write-Host "[信息] 等待3秒后重新启动服务..." -ForegroundColor Cyan
Start-Sleep -Seconds 3

# 检查是否安装了Go
try {
    $goVersion = go version
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "Go已安装: $goVersion"
} catch {
    Write-Host "[✗] " -ForegroundColor Red -NoNewline
    Write-Host "Go未安装，请先安装Go: https://golang.org/dl/"
    Write-Host "按任意键退出..."
    $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit
}

# 检查是否安装了Node.js和npm
try {
    $nodeVersion = node -v
    $npmVersion = npm -v
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "Node.js已安装: $nodeVersion (npm: $npmVersion)"
} catch {
    Write-Host "[✗] " -ForegroundColor Red -NoNewline
    Write-Host "Node.js未安装，请先安装Node.js: https://nodejs.org/"
    Write-Host "按任意键退出..."
    $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit
}

# 启动后端服务
Write-Host "[信息] 正在启动后端服务..." -ForegroundColor Cyan
Write-Host ""
Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..'; Write-Host '启动Go后端服务...' -ForegroundColor Cyan; go run cmd/server/main.go" -WindowStyle Normal

# 启动前端服务
Write-Host "[信息] 正在启动前端服务..." -ForegroundColor Cyan
Write-Host ""
Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..\web'; Write-Host '启动Vue前端服务...' -ForegroundColor Cyan; npm run dev" -WindowStyle Normal

Write-Host ""
Write-Host "[✓] " -ForegroundColor Green -NoNewline
Write-Host "服务重启完成！"
Write-Host ""
Write-Host "后端服务: " -NoNewline
Write-Host "http://localhost:8080" -ForegroundColor Cyan
Write-Host "前端服务: " -NoNewline
Write-Host "请查看前端服务窗口中显示的URL" -ForegroundColor Cyan
Write-Host ""
Write-Host "提示: 关闭此窗口不会停止服务，请手动关闭相应的服务窗口" -ForegroundColor Gray
Write-Host ""

Write-Host "按任意键退出..." -ForegroundColor DarkGray
$null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")