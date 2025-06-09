@echo off
echo ======================================
echo  eden* Docker一键启动脚本
echo ======================================
echo.

:: 设置窗口标题
title eden* Docker一键启动

:: 设置颜色
color 0B

:: 检查Docker是否安装
docker --version > nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] Docker未安装或未运行，请先安装并启动Docker。
    echo 下载地址: https://www.docker.com/products/docker-desktop
    echo.
    pause
    exit /b 1
)

echo [信息] 正在使用Docker Compose启动eden*...
echo.

:: 切换到docker目录
cd /d %~dp0\..\..\docker

:: 启动服务
docker-compose up -d

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [错误] Docker Compose启动失败，请检查错误信息。
    echo.
    pause
    exit /b 1
)

echo.
echo [成功] eden*已在Docker中启动！
echo.
echo 访问地址: http://localhost:8080
echo.
echo 提示: 要停止服务，请运行: docker-compose down
echo.

pause