#!/bin/bash

# Tide SOLO模式演示脚本
# 用于快速测试SOLO模式功能

echo "🌊 Tide SOLO模式演示"
echo "===================="
echo

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 请先安装Go 1.19+"
    exit 1
fi

# 检查Node.js环境
if ! command -v node &> /dev/null; then
    echo "⚠️  未检测到Node.js，某些功能可能受限"
fi

# 检查Git环境
if ! command -v git &> /dev/null; then
    echo "⚠️  未检测到Git，部署功能可能受限"
fi

# 演示目录
demo_dir="solo-demo-$(date +%Y%m%d-%H%M%S)"
echo "📁 创建演示目录: $demo_dir"
mkdir -p "$demo_dir"
cd "$demo_dir"

# 测试需求列表
demos=(
    "创建一个响应式静态网站，包含导航栏和联系表单"
    "构建React应用，展示产品列表"
    "创建Express API，支持用户注册"
)

echo "🎯 可用演示:"
for i in "${!demos[@]}"; do
    echo "  $((i+1)). ${demos[$i]}"
done
echo

# 让用户选择演示
read -p "选择演示 (1-${#demos[@]}): " choice

# 验证选择
if [[ ! "$choice" =~ ^[1-3]$ ]]; then
    echo "❌ 无效选择，使用默认演示"
    choice=1
fi

# 获取选中的需求
selected_demo="${demos[$((choice-1))]}"
echo "🚀 启动SOLO模式: $selected_demo"
echo

# 运行SOLO模式
cd ..
go run solo/main.go "$selected_demo"

# 检查生成的项目
echo
echo "📊 演示结果:"
ls -la "$demo_dir" 2>/dev/null || echo "项目已生成在其他目录"

echo
echo "🎉 演示完成！"
echo "📂 项目文件位于: $(pwd)/$demo_dir"
echo "💡 按照上面的步骤完成部署"