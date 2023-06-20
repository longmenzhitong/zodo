#!/bin/bash

# 设置输出目录
output_dir="./bin"
output_name="zodo"

# 创建输出目录
mkdir -p "$output_dir"

# 构建Go服务
go build -o "$output_dir/$output_name" .

# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "构建成功！可执行文件已输出到 $output_dir"

    if [ "$1" = "run" ]; then
        # 切换到输出目录
        cd "$output_dir"

        # 运行可执行文件
        ./$output_name
    fi
else
    echo "构建失败！"
fi
