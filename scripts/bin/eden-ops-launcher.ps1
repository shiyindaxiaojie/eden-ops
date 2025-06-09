# eden* 启动器 (PowerShell版本)

# 函数: 显示菜单
function Show-Menu {
    Clear-Host
    Write-Host "
======================================="  -ForegroundColor Yellow
    Write-Host "       eden* 系统启动器" -ForegroundColor Yellow
    Write-Host "======================================="  -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  请选择启动方式:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  [1] 开发模式 - 启动后端和前端服务" -ForegroundColor White
    Write-Host "  [2] Docker模式 - 使用Docker Compose启动" -ForegroundColor White
    Write-Host "  [3] 仅启动后端服务" -ForegroundColor White
    Write-Host "  [4] 仅启动前端服务" -ForegroundColor White
    Write-Host "  [5] 重启服务" -ForegroundColor White
    Write-Host "  [6] 停止服务" -ForegroundColor White
    Write-Host "  [7] 退出" -ForegroundColor White
    Write-Host ""
    Write-Host "======================================="  -ForegroundColor Yellow
    Write-Host ""
}

# 函数: 检查环境
function Check-Environment {
    $envStatus = @{}
    
    # 检查Go
    try {
        $goVersion = go version
        $envStatus.Go = @{ Installed = $true; Version = $goVersion }
    } catch {
        $envStatus.Go = @{ Installed = $false; Version = $null }
    }
    
    # 检查Node.js和npm
    try {
        $nodeVersion = node -v
        $npmVersion = npm -v
        $envStatus.Node = @{ Installed = $true; Version = $nodeVersion; NpmVersion = $npmVersion }
    } catch {
        $envStatus.Node = @{ Installed = $false; Version = $null; NpmVersion = $null }
    }
    
    # 检查Docker
    try {
        $dockerVersion = docker --version
        $envStatus.Docker = @{ Installed = $true; Version = $dockerVersion }
    } catch {
        $envStatus.Docker = @{ Installed = $false; Version = $null }
    }
    
    return $envStatus
}

# 函数: 显示环境状态
function Show-Environment {
    param ($envStatus)
    
    Write-Host "环境检查:" -ForegroundColor Cyan
    
    # Go状态
    if ($envStatus.Go.Installed) {
        Write-Host "[✓] " -ForegroundColor Green -NoNewline
        Write-Host "Go: $($envStatus.Go.Version)"
    } else {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "Go未安装 (后端服务需要)"
    }
    
    # Node.js状态
    if ($envStatus.Node.Installed) {
        Write-Host "[✓] " -ForegroundColor Green -NoNewline
        Write-Host "Node.js: $($envStatus.Node.Version) (npm: $($envStatus.Node.NpmVersion))"
    } else {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "Node.js未安装 (前端服务需要)"
    }
    
    # Docker状态
    if ($envStatus.Docker.Installed) {
        Write-Host "[✓] " -ForegroundColor Green -NoNewline
        Write-Host "Docker: $($envStatus.Docker.Version)"
    } else {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "Docker未安装 (Docker模式需要)"
    }
    
    Write-Host ""
}

# 函数: 开发模式
function Start-DevMode {
    param ($envStatus)
    
    Write-Host ""
    Write-Host "正在启动开发模式..." -ForegroundColor Cyan
    Write-Host ""
    
    # 检查必要环境
    if (-not $envStatus.Go.Installed) {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "错误: Go未安装，无法启动后端服务"
        Write-Host "请安装Go: https://golang.org/dl/"
        return $false
    }
    
    if (-not $envStatus.Node.Installed) {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "错误: Node.js未安装，无法启动前端服务"
        Write-Host "请安装Node.js: https://nodejs.org/"
        return $false
    }
    
    # 启动后端服务
    Write-Host "[1/2] " -ForegroundColor Yellow -NoNewline
    Write-Host "启动后端服务..."
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..'; Write-Host '启动Go后端服务...' -ForegroundColor Cyan; go run cmd/server/main.go" -WindowStyle Normal
    
    # 启动前端服务
    Write-Host "[2/2] " -ForegroundColor Yellow -NoNewline
    Write-Host "启动前端服务..."
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..\web'; Write-Host '启动Vue前端服务...' -ForegroundColor Cyan; npm run dev" -WindowStyle Normal
    
    Write-Host ""
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "服务启动完成！"
    Write-Host ""
    Write-Host "后端服务: " -NoNewline
    Write-Host "http://localhost:8080" -ForegroundColor Cyan
    Write-Host "前端服务: " -NoNewline
    Write-Host "请查看前端服务窗口中显示的URL" -ForegroundColor Cyan
    
    return $true
}

# 函数: Docker模式
function Start-DockerMode {
    param ($envStatus)
    
    Write-Host ""
    Write-Host "正在使用Docker Compose启动eden*..." -ForegroundColor Cyan
    Write-Host ""
    
    # 检查Docker是否安装
    if (-not $envStatus.Docker.Installed) {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "错误: Docker未安装或未运行"
        Write-Host "请安装Docker: https://www.docker.com/products/docker-desktop"
        return $false
    }
    
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
            Write-Host ""
            Write-Host "提示: 要停止服务，请运行: " -NoNewline
            Write-Host "docker-compose down" -ForegroundColor Yellow
        } else {
            Write-Host ""
            Write-Host "[✗] " -ForegroundColor Red -NoNewline
            Write-Host "Docker Compose启动失败，请检查错误信息。"
            return $false
        }
    } catch {
        Write-Host ""
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "执行Docker Compose时出错: $_"
        return $false
    }
    
    # 返回到原目录
    Set-Location -Path $PSScriptRoot
    return $true
}

# 函数: 仅启动后端
function Start-BackendOnly {
    param ($envStatus)
    
    Write-Host ""
    Write-Host "正在启动后端服务..." -ForegroundColor Cyan
    Write-Host ""
    
    # 检查Go是否安装
    if (-not $envStatus.Go.Installed) {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "错误: Go未安装，无法启动后端服务"
        Write-Host "请安装Go: https://golang.org/dl/"
        return $false
    }
    
    # 启动后端服务
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..'; Write-Host '启动Go后端服务...' -ForegroundColor Cyan; go run cmd/server/main.go" -WindowStyle Normal
    
    Write-Host ""
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "后端服务启动完成！"
    Write-Host ""
    Write-Host "访问地址: " -NoNewline
    Write-Host "http://localhost:8080" -ForegroundColor Cyan
    
    return $true
}

# 函数: 仅启动前端
function Start-FrontendOnly {
    param ($envStatus)
    
    Write-Host ""
    Write-Host "正在启动前端服务..." -ForegroundColor Cyan
    Write-Host ""
    
    # 检查Node.js是否安装
    if (-not $envStatus.Node.Installed) {
        Write-Host "[✗] " -ForegroundColor Red -NoNewline
        Write-Host "错误: Node.js未安装，无法启动前端服务"
        Write-Host "请安装Node.js: https://nodejs.org/"
        return $false
    }
    
    # 启动前端服务
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$PSScriptRoot\..\..\web'; Write-Host '启动Vue前端服务...' -ForegroundColor Cyan; npm run dev" -WindowStyle Normal
    
    Write-Host ""
    Write-Host "[✓] " -ForegroundColor Green -NoNewline
    Write-Host "前端服务启动完成！"
    Write-Host ""
    Write-Host "访问地址: " -NoNewline
    Write-Host "请查看前端服务窗口中显示的URL" -ForegroundColor Cyan
    
    return $true
}

# 函数: 重启服务
function Restart-Services {
    Write-Host ""
    Write-Host "正在重启服务..." -ForegroundColor Cyan
    Write-Host ""
    
    # 调用重启脚本
    & "$PSScriptRoot\restart.ps1"
    
    return $true
}

# 函数: 停止服务
function Stop-Services {
    Write-Host ""
    Write-Host "正在停止服务..." -ForegroundColor Cyan
    Write-Host ""
    
    # 调用停止脚本
    & "$PSScriptRoot\stop.ps1"
    
    return $true
}

# 设置窗口标题
$host.UI.RawUI.WindowTitle = "eden* 启动器"

# 主循环
while ($true) {
    # 显示菜单
    Show-Menu
    
    # 检查环境
    $envStatus = Check-Environment
    Show-Environment $envStatus
    
    # 获取用户选择
    $choice = Read-Host "请输入选项 [1-7]"
    
    switch ($choice) {
        "1" { 
            $result = Start-DevMode $envStatus
            if ($result) {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            } else {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            }
        }
        "2" { 
            $result = Start-DockerMode $envStatus
            if ($result) {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            } else {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            }
        }
        "3" { 
            $result = Start-BackendOnly $envStatus
            if ($result) {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            } else {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            }
        }
        "4" { 
            $result = Start-FrontendOnly $envStatus
            if ($result) {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            } else {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            }
        }
        "5" { 
            $result = Restart-Services
            if ($result) {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            }
        }
        "6" { 
            $result = Stop-Services
            if ($result) {
                Write-Host ""
                Write-Host "按任意键返回主菜单..." -ForegroundColor DarkGray
                $null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
            }
        }
        "7" { 
            Write-Host ""
            Write-Host "感谢使用 eden* 启动器！" -ForegroundColor Cyan
            Write-Host ""
            Start-Sleep -Seconds 1
            exit
        }
        default { 
            Write-Host ""
            Write-Host "无效的选项，请重新选择。" -ForegroundColor Red
            Write-Host ""
            Start-Sleep -Seconds 2
        }
    }
}