# go install 支持更新说明

## 问题

使用 `go install github.com/zdzh/goctl-env-template@latest` 安装时遇到错误：
```
github.com/zdzh/goctl-env-template: no non-test Go files in /Users/pikaka/go/pkg/mod/github.com/zdzh/goctl-env-template@v1.0.0
```

## 原因

`go install` 默认在以下位置查找 main 包：
1. 模块根目录
2. `cmd/<module-name>` 目录

之前的项目结构将 main 包放在 `cmd/root.go` 中，而根目录没有 main 包。

## 解决方案

### 1. 创建根目录的 main.go

在项目根目录创建了 `main.go`，内容从 `cmd/root.go` 复制并调整。

### 2. 更新项目结构

**修改前：**
```
goctl-env-template/
├── cmd/
│   └── root.go               # main 包
├── internal/
└── ...
```

**修改后：**
```
goctl-env-template/
├── main.go                   # main 包（根目录）
├── internal/
└── ...
```

### 3. 删除 cmd 目录

删除了 `cmd/` 目录，因为 main.go 已经在根目录。

### 4. 更新 Makefile

```makefile
# 修改前
build-cli:
	go build -o goctl-env-template cmd/root.go

# 修改后
build:
	go build -o goctl-env-template .
```

### 5. 更新模块名称

所有导入路径从 `github.com/yourusername/goctl-env-template` 更新为 `github.com/zdzh/goctl-env-template`。

## 验证

### 1. 本地构建

```bash
make build
./goctl-env-template -c config/config.go -o .env.template
```

✅ 构建成功
✅ 生成正确

### 2. 测试

```bash
make test
```

✅ 所有测试通过

### 3. go install（模拟）

```bash
go build -o goctl-env-template .
```

✅ 可以从根目录构建

## 使用方式

### 方式1：go install（推荐）

```bash
go install github.com/zdzh/goctl-env-template@latest
goctl-env-template -c config/config.go -o .env.template
```

### 方式2：本地构建

```bash
make build
./goctl-env-template -c config/config.go -o .env.template
```

## 项目结构更新

```bash
goctl-env-template/
├── main.go                   # 主入口（新增）
├── config/
│   └── config.go
├── internal/
│   ├── generator/
│   │   └── generator.go
│   ├── parser/
│   │   └── parser.go
│   └── types/
│       └── config.go
├── plugin/
│   └── main.go
├── 文档文件（8个）
├── 测试文件
└── Makefile
```

## 兼容性

✅ 所有现有功能保持不变
✅ 测试全部通过
✅ 文档已更新
✅ Makefile 已更新

## 总结

通过在根目录创建 `main.go` 并删除 `cmd/` 目录，项目现在可以通过 `go install` 正常安装。这是 Go 语言项目中 main 包的标准布局方式。
