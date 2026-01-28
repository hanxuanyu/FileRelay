#!/bin/bash

# 获取脚本所在目录并切换到项目根目录
cd "$(dirname "$0")/.."

# 设置变量
APP_NAME="filerelay"
OUTPUT_DIR="output"

# 定义目标平台
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

echo "开始构建 $APP_NAME 多平台二进制文件..."

# 清理 output 目录
if [ -d "$OUTPUT_DIR" ]; then
    echo "正在清理 $OUTPUT_DIR 目录..."
    rm -rf "$OUTPUT_DIR"
fi

mkdir -p "$OUTPUT_DIR"

# 前端构建
echo "正在构建前端项目..."
cd webapp
npm install
npm run build
if [ $? -ne 0 ]; then
    echo "前端构建失败，停止编译。"
    exit 1
fi
cd ..

# 循环编译各平台
for PLATFORM in "${PLATFORMS[@]}"; do
    # 分离 OS 和 ARCH
    OS=$(echo $PLATFORM | cut -d'/' -f1)
    ARCH=$(echo $PLATFORM | cut -d'/' -f2)
    
    # 设置输出名称
    OUTPUT_NAME="${APP_NAME}-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "正在编译 ${OS}/${ARCH}..."
    GOOS=$OS GOARCH=$ARCH go build -o "${OUTPUT_DIR}/${OUTPUT_NAME}" main.go
    
    if [ $? -eq 0 ]; then
        echo "  ${OS}/${ARCH} 编译成功"
        # 压缩为 tar.gz
        tar -czf "${OUTPUT_DIR}/${OUTPUT_NAME}.tar.gz" -C "${OUTPUT_DIR}" "${OUTPUT_NAME}"
        # 删除原始二进制文件
        rm "${OUTPUT_DIR}/${OUTPUT_NAME}"
    else
        echo "  ${OS}/${ARCH} 编译失败"
    fi
done

echo "----------------------------------------"
echo "多平台打包完成！输出目录: $OUTPUT_DIR"
ls -R "$OUTPUT_DIR"
