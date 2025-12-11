#!/bin/bash

echo "========================================"
echo "OneImg 图床系统启动脚本"
echo "========================================"

echo ""
echo "正在检查环境..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "[错误] 未检测到Go环境，请先安装Go语言"
    exit 1
fi

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "[错误] 未检测到Node.js环境，请先安装Node.js"
    exit 1
fi

echo "[成功] 环境检查通过"

echo ""
echo "正在构建前端..."
cd frontend

if [ ! -d "node_modules" ]; then
    echo "正在安装前端依赖..."
    npm install
    if [ $? -ne 0 ]; then
        echo "[错误] 前端依赖安装失败"
        exit 1
    fi
fi

echo "正在构建前端项目..."
npm run build
if [ $? -ne 0 ]; then
    echo "[错误] 前端构建失败"
    exit 1
fi

cd ..
echo "[成功] 前端构建完成"

echo ""
echo "正在启动后端服务..."
cd backend

# 检查go.mod是否存在
if [ ! -f "go.mod" ]; then
    echo "正在初始化Go模块..."
    go mod init oneimg/backend
    go mod tidy
fi

# 下载依赖
echo "正在下载Go依赖..."
go mod tidy

# 启动服务
echo "正在启动服务器..."
echo ""
echo "========================================"
echo "服务启动成功！"
echo "访问地址: http://localhost:8080"
echo "按 Ctrl+C 停止服务"
echo "========================================"
echo ""

go run main.go