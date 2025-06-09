@echo off
echo ======================================
echo  eden* 停止服务脚本
echo ======================================
echo.

:: 设置窗口标题
title eden* 停止服务

:: 设置颜色
color 0C

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

echo.
echo [信息] 服务停止操作完成
echo.

:: 检查是否还有相关进程在运行
tasklist /fi "WINDOWTITLE eq eden* 后端服务" /fo list | findstr "PID:" >nul
if %ERRORLEVEL% EQU 0 (
    echo [警告] 仍有后端服务进程在运行
)

tasklist /fi "WINDOWTITLE eq eden* 前端服务" /fo list | findstr "PID:" >nul
if %ERRORLEVEL% EQU 0 (
    echo [警告] 仍有前端服务进程在运行
)

echo.
echo 提示: 如果有服务无法停止，可以使用任务管理器手动结束进程
echo.

pause