# 安装指南

## 方法1：go install（推荐）

这是最简单的安装方式，会将二进制文件安装到 `$GOPATH/bin` 目录。

```bash
go install github.com/zdzh/goctl-env-template@latest
```

确保 `$GOPATH/bin` 在你的 `$PATH` 环境变量中：

```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
```

验证安装：

```bash
goctl-env-template -h
```

## 方法2：克隆仓库构建

```bash
# 克隆仓库
git clone https://github.com/zdzh/goctl-env-template.git
cd goctl-env-template

# 构建
make build

# 安装（可选）
make install
```

## 方法3：使用 Makefile

如果你已经克隆了仓库：

```bash
# 构建
make build

# 或者直接安装
make install
```

## 验证安装

运行以下命令验证安装是否成功：

```bash
# 查看帮助
goctl-env-template -h

# 使用示例配置生成模板
cd goctl-env-template
goctl-env-template -c config/config.go -o .env.template

# 查看生成的文件
cat .env.template
```

## 作为 goctl 插件使用

如果你需要作为 goctl 插件使用：

```bash
# 构建插件
make build-plugin

# 使用插件
goctl plugin -p plugin/main --config config/config.go --output .env
```

## 常见问题

### Q: 找不到命令 goctl-env-template

A: 确保 `$GOPATH/bin` 在 `$PATH` 中：

```bash
# 查看 GOPATH
go env GOPATH

# 添加到 PATH（macOS/Linux）
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Q: go install 失败

A: 确保：
1. 网络连接正常
2. Go 版本 >= 1.21
3. 可以访问 GitHub

```bash
# 检查 Go 版本
go version

# 测试网络连接
ping github.com
```

### Q: 如何更新到最新版本

```bash
go install github.com/zdzh/goctl-env-template@latest
```

### Q: 如何卸载

```bash
# 删除二进制文件
rm $(go env GOPATH)/bin/goctl-env-template

# 或删除克隆的目录
rm -rf goctl-env-template
```

## 快速开始

安装完成后，快速生成 .env 模板：

```bash
# 1. 创建配置文件 config/config.go
# 2. 生成模板
goctl-env-template -c config/config.go -o .env.template

# 3. 查看生成的模板
cat .env.template
```

## 获取帮助

```bash
# 查看帮助
goctl-env-template -h

# 或
goctl-env-template --help
```

## 下一步

安装完成后，查看以下文档了解更多：

- [快速入门](QUICKSTART.md) - 5分钟快速入门
- [中文文档](README_CN.md) - 完整中文文档
- [英文文档](README.md) - 完整英文文档
- [使用示例](EXAMPLES.md) - 详细使用示例
