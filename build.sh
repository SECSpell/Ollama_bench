#!/bin/bash

# 设置程序名称 / Set the application name
APP_NAME="ollama_bench"

# 创建一个输出目录 / Create an output directory
mkdir -p build

# 定义目标平台 / Define target platforms
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

# 遍历并编译每个平台 / Iterate and compile for each platform
for platform in "${platforms[@]}"
do
    # 分割平台字符串 / Split the platform string
    IFS='/' read -r -a array <<< "$platform"
    GOOS=${array[0]}
    GOARCH=${array[1]}
    
    # 设置输出文件名 / Set the output file name
    if [ $GOOS = "windows" ]; then
        output_name=$APP_NAME'_'$GOOS'_'$GOARCH'.exe'
    else
        output_name=$APP_NAME'_'$GOOS'_'$GOARCH
    fi

    # 编译 / Compile
    echo "Compiling $output_name"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o build/$output_name main.go
    
    if [ $? -ne 0 ]; then
        echo "Failed to compile $output_name"
    else
        echo "Successfully compiled $output_name"
    fi
done

echo "Compilation completed"
