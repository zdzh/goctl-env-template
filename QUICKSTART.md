# 快速入门

## 安装

```bash
go install github.com/zdzh/goctl-env-template@latest
```

## 5分钟快速使用

### 1. 创建配置文件

创建 `config/config.go`：

```go
package config

type Config struct {
	// 服务名称
	Name string `json:",env=SERVICE_NAME"`

	// 服务端口（默认：8080）
	Port int `json:",env=SERVICE_PORT,default=8080"`

	// 数据库主机
	DBHost string `json:",env=DB_HOST"`

	// 数据库端口
	DBPort int `json:",env=DB_PORT,default=3306"`

	// 数据库用户名
	DBUser string `json:",env=DB_USER"`

	// 数据库密码（可选）
	DBPassword string `json:",env=DB_PASSWORD,optional"`
}
```

### 2. 生成 .env 模板

```bash
goctl-env-template -c config/config.go -o .env.template
```

### 3. 生成的 .env.template

```bash
## Configuration
# 服务名称
SERVICE_NAME=

# 服务端口（默认：8080）
SERVICE_PORT=8080

## DB Configuration
# 数据库主机
DB_HOST=

# 数据库端口（默认：3306）
DB_PORT=3306

# 数据库用户名
DB_USER=

# 数据库密码（可选）
# DB_PASSWORD=
```

### 4. 复制并配置

```bash
cp .env.template .env
# 编辑 .env 文件，填写实际值
```

## 常用命令

```bash
# 使用默认路径生成
goctl-env-template

# 指定配置和输出文件
goctl-env-template -c internal/config/config.go -o .env.example

# 使用 goctl 插件
goctl plugin -p goctl-env-template --config config/config.go --output .env
```

## 标签格式

```go
FieldName Type `json:",env=ENV_NAME,optional,default=value"`
```

- `env=ENV_NAME`: 环境变量名（必填）
- `optional`: 可选字段（可选）
- `default=value`: 默认值（可选）

## 注释格式

- `## 分组名`: 配置分组
- `# 说明`: 配置项说明（位于配置项上方）

## 更多文档

- [完整文档 (中文)](README_CN.md)
- [Full Documentation (English)](README.md)
- [使用示例](EXAMPLES.md)
- [实现文档](IMPLEMENTATION.md)
