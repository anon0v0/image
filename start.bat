@echo off
echo ========================================
echo OneImg 图床系统启动脚本
echo ========================================

echo.
echo 正在检查环境...

:: 检查Go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未检测到Go环境，请先安装Go语言
    pause
    exit /b 1
)

:: 检查Node.js是否安装
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未检测到Node.js环境，请先安装Node.js
    pause
    exit /b 1
)

echo [成功] 环境检查通过

echo.
echo 正在构建前端...
cd frontend
if not exist "node_modules" (
    echo 正在安装前端依赖...
    npm install
    if %errorlevel% neq 0 (
        echo [错误] 前端依赖安装失败
        pause
        exit /b 1
    )
)

echo 正在构建前端项目...
npm run build
if %errorlevel% neq 0 (
    echo [错误] 前端构建失败
    pause
    exit /b 1
)

cd ..
echo [成功] 前端构建完成

echo.
echo 正在启动后端服务...
cd backend

:: 检查go.mod是否存在
if not exist "go.mod" (
    echo 正在初始化Go模块...
    go mod init oneimg/backend
    go mod tidy
)

:: 下载依赖
echo 正在下载Go依赖...
go mod tidy

:: 启动服务
echo 正在启动服务器...
echo.
echo ========================================
echo 服务启动成功！
echo 访问地址: http://localhost:8080
echo 按 Ctrl+C 停止服务
echo ========================================
echo.

go run main.go

pause