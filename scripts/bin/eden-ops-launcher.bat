@echo off
setlocal enabledelayedexpansion

:: 设置窗口标题和颜色
title eden* 启动器
color 0E

:MENU
cls
echo =======================================
echo       eden* 系统启动器
echo =======================================
echo.
echo  请选择启动方式:
echo.
echo  [1] 开发模式 - 启动后端和前端服务
echo  [2] Docker模式 - 使用Docker Compose启动
echo  [3] 仅启动后端服务
echo  [4] 仅启动前端服务
echo  [5] 重启服务
echo  [6] 停止服务
echo  [7] 退出
echo.
echo =======================================
echo.

set /p choice=请输入选项 [1-7]: 

if "%choice%"=="1" goto DEV_MODE
if "%choice%"=="2" goto DOCKER_MODE
if "%choice%"=="3" goto BACKEND_ONLY
if "%choice%"=="4" goto FRONTEND_ONLY
if "%choice%"=="5" goto RESTART
if "%choice%"=="6" goto STOP
if "%choice%"=="7" goto EXIT

echo.
echo 无效的选项，请重新选择。
echo.
timeout /t 2 >nul
goto MENU

:DEV_MODE
echo.
echo 正在启动开发模式...
echo.

:: 启动后端服务（新窗口）
start "eden* 后端服务" cmd /c "cd /d %~dp0\..\.. && go run cmd/server/main.go"

:: 启动前端服务（新窗口）
start "eden* 前端服务" cmd /c "cd /d %~dp0\..\..\web && npm run dev"

echo 服务启动完成！
echo.
echo 后端服务: http://localhost:8080
echo 前端服务: 请查看前端服务窗口中显示的URL
echo.
echo 按任意键返回主菜单...
pause >nul
goto MENU

:DOCKER_MODE
echo.
echo 正在使用Docker Compose启动eden*...
echo.

:: 检查Docker是否安装
docker --version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] Docker未安装或未运行，请先安装并启动Docker。
    echo 下载地址: https://www.docker.com/products/docker-desktop
    echo.
    echo 按任意键返回主菜单...
    pause >nul
    goto MENU
)

:: 切换到docker目录并启动
cd /d %~dp0\..\..\docker
docker-compose up -d

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [错误] Docker Compose启动失败，请检查错误信息。
    echo.
) else (
    echo.
    echo [成功] eden*已在Docker中启动！
    echo.
    echo 访问地址: http://localhost:8080
    echo.
)

echo 按任意键返回主菜单...
pause >nul
goto MENU

:BACKEND_ONLY
echo.
echo 正在启动后端服务...
echo.

:: 启动后端服务（新窗口）
start "eden* 后端服务" cmd /c "cd /d %~dp0\..\.. && go run cmd/server/main.go"

echo 后端服务启动完成！
echo.
echo 访问地址: http://localhost:8080
echo.
echo 按任意键返回主菜单...
pause >nul
goto MENU

:FRONTEND_ONLY
echo.
echo 正在启动前端服务...
echo.

:: 启动前端服务（新窗口）
start "eden* 前端服务" cmd /c "cd /d %~dp0\..\..\web && npm run dev"

echo 前端服务启动完成！
echo.
echo 访问地址: 请查看前端服务窗口中显示的URL
echo.
echo 按任意键返回主菜单...
pause >nul
goto MENU

:RESTART
echo.
echo 正在重启服务...
echo.

:: 调用重启脚本
call %~dp0\restart.bat

echo.
echo 按任意键返回主菜单...
pause >nul
goto MENU

:STOP
echo.
echo 正在停止服务...
echo.

:: 调用停止脚本
call %~dp0\stop.bat

echo.
echo 按任意键返回主菜单...
pause >nul
goto MENU

:EXIT
echo.
echo 感谢使用 eden* 启动器！
echo.
timeout /t 2 >nul
exit /b 0