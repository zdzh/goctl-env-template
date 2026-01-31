# goctl-env-template

一个 CLI 工具，用于从 `config/config.go` 中提取环境变量并生成 `.env` 模板文件。

## 功能特性

- 从结构体标签中提取环境变量 (`json:"name,env=VAR_NAME"`)
- 支持必填和可选变量（使用 `optional` 标记）
- 从结构体标签中提取默认值 (`default=value`)
- 解析文档注释作为字段说明
- 按结构体嵌套级别分组配置
- 生成清晰、带注释的 `.env` 模板文件

## 安装

```bash
go install github.com/zdzh/goctl-env-template@latest
```

确保安装的 `goctl-env-template` 在你的 `$PATH` 环境变量中。

## 使用方法

```bash
goctl-env-template -c config/config.go -o .env.template
```

## 配置文件格式

在 `config/config.go` 中使用环境变量标签定义配置：

```go
package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth  AuthConfig
	DB    DBConfig
}

type AuthConfig struct {
	// JWT 密钥
	Secret string `json:",env=AUTH_SECRET"`

	// 超时时间（单位：秒）
	Timeout int `json:",env=AUTH_TIMEOUT,default=30"`

	// 刷新间隔（可选）
	RefreshInterval int `json:",env=AUTH_REFRESH_INTERVAL,optional"`
}

type DBConfig struct {
	// 数据库主机地址
	Host string `json:",env=DB_HOST"`

	// 数据库端口（默认：3306）
	Port int `json:",env=DB_PORT,default=3306"`

	// 数据库用户名
	Username string `json:",env=DB_USERNAME"`

	// 数据库密码（可选）
	Password string `json:",env=DB_PASSWORD,optional"`
}
```

## 生成的 .env 模板

插件会生成一个 `.env.template` 文件，包含：

- 必填变量为激活状态（未注释）
- 可选变量为注释状态（使用 `#` 标记）
- 从文档注释提取的说明
- 默认值作为示例

示例输出：

```bash
## Auth 配置
# JWT 密钥
AUTH_SECRET=

# 超时时间（单位：秒）
AUTH_TIMEOUT=30

# 刷新间隔（可选）
# AUTH_REFRESH_INTERVAL=

## DB 配置
# 数据库主机地址
DB_HOST=

# 数据库端口（默认：3306）
# DB_PORT=3306

# 数据库用户名
DB_USERNAME=

# 数据库密码（可选）
# DB_PASSWORD=
```

## 结构体标签格式

在配置结构体中使用以下标签格式：

```go
FieldName FieldType `json:"name,env=ENV_VAR,optional,default=value"`
```

### 标签选项

- `env=VAR_NAME`: 环境变量名称（必填）
- `optional`: 标记为可选变量（可选）
- `default=value`: 默认值（可选）

## 注释说明

- **分组注释**: 使用 `##` 标记配置分组
- **字段注释**: 使用 `#` 标记字段说明
- **字段说明位置**: 位于配置项上方

示例：

```bash
## 数据库配置
# 数据库主机地址
DB_HOST=

# 数据库端口（默认：3306）
DB_PORT=3306
```

## CLI 选项

```
Usage: goctl-env-template [flags]

Flags:
  -c, --config string   配置文件路径（默认 "config/config.go"）
  -o, --output string   输出文件路径（默认 ".env.template"）
  -h, --help            显示帮助信息
```

## 使用示例

### 使用默认路径

```bash
goctl-env-template
```

### 指定自定义配置和输出路径

```bash
goctl-env-template -c internal/config/config.go -o .env.example
```

### 使用 Makefile

```bash
make build      # 构建工具
make test       # 测试生成
make demo       # 生成并显示 .env 模板
make clean      # 清理构建文件
```

## 项目结构

```
goctl-env-template/
├── main.go               # 主入口
├── internal/
│   ├── parser/
│   │   └── parser.go     # Go AST 解析器
│   ├── generator/
│   │   └── generator.go  # .env 模板生成器
│   └── types/
│       └── config.go     # 内部数据结构
├── config/
│   └── config.go         # 示例配置文件
├── go.mod
├── go.sum
├── Makefile
├── README.md             # 英文文档
├── README_CN.md          # 中文文档
└── EXAMPLES.md           # 使用示例
```

## 高级用法

### 自动分组

嵌套结构体根据名称自动分组：

```go
type Config struct {
	Redis RedisConfig
	Cache CacheConfig
}
```

将生成独立的分组：
- `## Redis 配置`
- `## Cache 配置`

### 默认值

在结构体标签中指定默认值：

```go
Port int `json:",env=PORT,default=8080"`
```

### 可选字段

标记字段为可选，在模板中注释掉：

```go
Password string `json:",env=PASSWORD,optional"`
```

### 文档注释

使用 Go 文档注释添加说明：

```go
// 数据库连接超时时间（秒）
Timeout int `json:",env=DB_TIMEOUT,default=10"`
```

这将在生成的模板中显示为：
```bash
# 数据库连接超时时间（秒）
DB_TIMEOUT=10
```

## 开发指南

### 构建项目

```bash
make build
```

### 运行测试

```bash
make test
```

### 清理构建文件

```bash
make clean
```

### 生成演示

```bash
make demo
```

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

MIT License
