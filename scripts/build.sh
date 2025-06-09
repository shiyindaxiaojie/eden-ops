#!/bin/bash

# 设置版本号
VERSION=$(git describe --tags --always --dirty)
if [ -z "$VERSION" ]; then
    VERSION="dev"
fi

# 设置构建时间
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# 设置构建参数
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME"

# 清理旧的构建文件
rm -rf build/package/*

# 构建多平台二进制文件
PLATFORMS=("windows/amd64" "linux/amd64" "darwin/amd64")

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BIN_NAME="eden-ops"
    
    if [ $GOOS = "windows" ]; then
        BIN_NAME="eden-ops.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "build/package/eden-ops_${GOOS}_${GOARCH}/$BIN_NAME" cmd/eden-ops/main.go
    
    if [ $? -ne 0 ]; then
        echo "Error building for $GOOS/$GOARCH"
        exit 1
    fi
    
    # 复制配置文件模板
    cp configs/config.yaml.example "build/package/eden-ops_${GOOS}_${GOARCH}/config.yaml.example"
    
    # 创建压缩包
    cd build/package
    if [ $GOOS = "windows" ]; then
        zip -r "eden-ops_${GOOS}_${GOARCH}.zip" "eden-ops_${GOOS}_${GOARCH}"
    else
        tar -czf "eden-ops_${GOOS}_${GOARCH}.tar.gz" "eden-ops_${GOOS}_${GOARCH}"
    fi
    cd ../..
done

echo "Build complete!" 