@echo off
echo ======================================
echo  eden* 一键启动脚本
echo ======================================
echo.

:: 设置窗口标题
title eden* 一键启动

:: 设置颜色
color 0A

echo [信息] 正在启动后端服务...
echo.

:: 启动后端服务（新窗口）
start "eden* 后端服务" cmd /c "cd /d %~dp0\..\.. && go run cmd/server/main.go"

echo [信息] 正在启动前端服务...
echo.

:: 启动前端服务（新窗口）
start "eden* 前端服务" cmd /c "cd /d %~dp0\..\..\web && npm run dev"

echo [成功] 服务启动完成！
echo.
echo 后端服务: http://localhost:8080
echo 前端服务: 请查看前端服务窗口中显示的URL
echo.
echo 提示: 关闭此窗口不会停止服务，请手动关闭相应的服务窗口
echo.

pause