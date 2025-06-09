# eden* 一键启动脚本 (PowerShell版本)

# 设置窗口标题
$host.UI.RawUI.WindowTitle = "eden* 一键启动"

# 输出标题
Write-Host "
===================================="  -ForegroundColor Cyan
Write-Host "  eden*一键启动脚本" -ForegroundColor Cyan
Write-Host "===================================="  -ForegroundColor Cyan
Write-Host ""

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
    Write-Host "Node.js或npm未安装，请先安装Node.js: https://nodejs.org/"
    Write-Host "按任意键退出..."
    $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit
}

Write-Host ""
Write-Host "[1/2] " -ForegroundColor Yellow -NoNewline
Write-Host "正在启动后端服务..."

# 启动后端服务（新窗口）
Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..'; Write-Host '启动Go后端服务...' -ForegroundColor Cyan; go run cmd/server/main.go" -WindowStyle Normal

Write-Host "[2/2] " -ForegroundColor Yellow -NoNewline
Write-Host "正在启动前端服务..."

# 启动前端服务（新窗口）
Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..\web'; Write-Host '启动Vue前端服务...' -ForegroundColor Cyan; npm run dev" -WindowStyle Normal

Write-Host ""
Write-Host "[✓] " -ForegroundColor Green -NoNewline
Write-Host "服务启动完成！"
Write-Host ""
Write-Host "后端服务: " -NoNewline
Write-Host "http://localhost:8080" -ForegroundColor Cyan
Write-Host "前端服务: " -NoNewline
Write-Host "请查看前端服务窗口中显示的URL" -ForegroundColor Cyan
Write-Host ""
Write-Host "提示: 关闭此窗口不会停止服务，请手动关闭相应的服务窗口" -ForegroundColor Yellow
Write-Host ""

Write-Host "按任意键退出此窗口..." -ForegroundColor DarkGray
$null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")