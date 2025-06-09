@echo off
echo ======================================
echo  eden* 重启脚本
echo ======================================
echo.

:: 设置窗口标题
title eden* 重启

:: 设置颜色
color 0E

echo [信息] 正在停止服务...
echo.

:: 查找并终止后端进程
for /f "tokens=2" %%a in ('tasklist /fi "WINDOWTITLE eq eden* 后端服务" /fo list ^| findstr "PID:"') do (
    echo 正在终止后端进程 PID: %%a
    taskkill /pid %%a /f >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        echo [成功] 后端服务已停止
    ) else (
        echo [警告] 无法停止后端服务
    )
)

:: 查找并终止前端进程
for /f "tokens=2" %%a in ('tasklist /fi "WINDOWTITLE eq eden* 前端服务" /fo list ^| findstr "PID:"') do (
    echo 正在终止前端进程 PID: %%a
    taskkill /pid %%a /f >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        echo [成功] 前端服务已停止
    ) else (
        echo [警告] 无法停止前端服务
    )
)

echo [信息] 服务已停止
echo.

echo [信息] 等待3秒后重新启动服务...
timeout /t 3 /nobreak >nul

:: 检查是否安装了Go
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] Go未安装，请先安装Go: https://golang.org/dl/
    echo.
    pause
    exit /b
) else (
    for /f "tokens=3" %%i in ('go version') do set GOVERSION=%%i
    echo [成功] Go已安装: %GOVERSION%
)

:: 检查是否安装了Node.js和npm
where node >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] Node.js未安装，请先安装Node.js: https://nodejs.org/
    echo.
    pause
    exit /b
) else (
    for /f "tokens=*" %%i in ('node -v') do set NODEVERSION=%%i
    for /f "tokens=*" %%i in ('npm -v') do set NPMVERSION=%%i
    echo [成功] Node.js已安装: %NODEVERSION% (npm: %NPMVERSION%)
)

echo [信息] 正在启动后端服务...
echo.

:: 启动后端服务（新窗口）
start "eden* 后端服务" cmd /c "cd /d %~dp0\..\.. && go run cmd/server/main.go"

echo [信息] 正在启动前端服务...
echo.

:: 启动前端服务（新窗口）
start "eden* 前端服务" cmd /c "cd /d %~dp0\..\..\web && npm run dev"

echo [成功] 服务重启完成！
echo.
echo 后端服务: http://localhost:8080
echo 前端服务: 请查看前端服务窗口中显示的URL
echo.
echo 提示: 关闭此窗口不会停止服务，请手动关闭相应的服务窗口
echo.

pause