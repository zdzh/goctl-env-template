#!/bin/bash

# 验证脚本

echo "==================================="
echo "goctl-env-template 功能验证"
echo "==================================="
echo ""

# 检查构建
echo "1. 检查构建..."
if [ -f "./goctl-env-template" ]; then
    echo "✅ CLI 工具已构建"
else
    echo "❌ CLI 工具未找到，正在构建..."
    make build
fi

if [ -f "./plugin/main" ]; then
    echo "✅ 插件已构建"
else
    echo "❌ 插件未找到，正在构建..."
    make build
fi

echo ""

# 检查文档
echo "2. 检查文档..."
docs=("README.md" "README_CN.md" "EXAMPLES.md" "QUICKSTART.md" "CHANGELOG.md")
for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        echo "✅ $doc 存在"
    else
        echo "❌ $doc 不存在"
    fi
done

echo ""

# 运行测试
echo "3. 运行测试..."
if go test -v > /dev/null 2>&1; then
    echo "✅ 所有测试通过"
else
    echo "❌ 测试失败"
    go test -v
fi

echo ""

# 生成 .env 模板
echo "4. 生成 .env 模板..."
./goctl-env-template -c config/config.go -o .env.template
if [ -f ".env.template" ]; then
    echo "✅ .env.template 已生成"
    
    # 验证格式
    echo ""
    echo "5. 验证格式..."
    
    # 检查分组注释
    if grep -q "^## " .env.template; then
        echo "✅ 分组注释使用 ##"
    else
        echo "❌ 分组注释未使用 ##"
    fi
    
    # 检查配置项注释
    if grep -q "^# " .env.template; then
        echo "✅ 配置项注释使用 #"
    else
        echo "❌ 配置项注释未使用 #"
    fi
    
    # 检查必需变量
    if grep -q "^AUTH_SECRET=" .env.template; then
        echo "✅ 必需变量未注释"
    else
        echo "❌ 必需变量被注释"
    fi
    
    # 检查可选变量
    if grep -q "^# AUTH_REFRESH_INTERVAL=" .env.template; then
        echo "✅ 可选变量已注释"
    else
        echo "❌ 可选变量未注释"
    fi
    
    # 检查注释位置
    comment_before=0
    comment_after=0
    while IFS= read -r line; do
        if [[ $line =~ ^#[[:space:]] && ! $line =~ ^## ]]; then
            # 读取下一行
            if IFS= read -r next_line; then
                if [[ $next_line =~ ^[A-Z_]+= ]]; then
                    ((comment_before++))
                fi
            fi
        fi
    done < .env.template
    
    if [ $comment_before -gt 0 ]; then
        echo "✅ 注释在配置项上方"
    else
        echo "⚠️  未检测到注释在配置项上方"
    fi
    
else
    echo "❌ .env.template 生成失败"
fi

echo ""
echo "==================================="
echo "验证完成"
echo "==================================="
echo ""
echo "生成的 .env.template 内容："
echo "-----------------------------------"
cat .env.template
echo "-----------------------------------"
