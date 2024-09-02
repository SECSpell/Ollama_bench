#!/bin/bash

# 设置程序名称
APP_NAME="ollama_bench"

# 创建一个输出目录
mkdir -p build

# 定义目标平台
platforms=(
    "windows/amd64"
    "windows/386"
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/386"
    "linux/arm"
    "linux/arm64"
)

# 遍历并编译每个平台
for platform in "${platforms[@]}"
do
    # 分割平台字符串
    IFS='/' read -r -a array <<< "$platform"
    GOOS=${array[0]}
    GOARCH=${array[1]}
    
    # 设置输出文件名
    if [ $GOOS = "windows" ]; then
        output_name=$APP_NAME'_'$GOOS'_'$GOARCH'.exe'
    else
        output_name=$APP_NAME'_'$GOOS'_'$GOARCH
    fi

    # 编译
    echo "编译 $output_name"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o build/$output_name main.go
    
    if [ $? -ne 0 ]; then
        echo "编译 $output_name 失败"
    else
        echo "编译 $output_name 成功"
    fi
done

echo "编译完成"
