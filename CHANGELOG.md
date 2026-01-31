# 更新日志

## 2025-01-31 - 注释格式更新

### 新增功能

1. **分组注释格式**
   - 使用 `##` 标记配置分组
   - 示例：`## Auth Configuration`

2. **配置项注释格式**
   - 使用 `#` 标记配置项说明
   - 示例：`# JWT secret key`

3. **注释位置**
   - 配置项说明注释现在位于配置项上方
   - 之前：配置项在上方，说明在下方
   - 现在：说明在上方，配置项在下方

### 格式对比

#### 之前
```bash
# Auth Configuration
AUTH_SECRET=
# JWT secret key

AUTH_TIMEOUT=3600
# Auth timeout (default: 3600)
```

#### 现在
```bash
## Auth Configuration
# JWT secret key
AUTH_SECRET=

# Auth timeout (default: 3600)
AUTH_TIMEOUT=3600
```

### 新增文档

- **README_CN.md**: 完整的中文使用文档
  - 功能特性介绍
  - 安装说明
  - 使用方法
  - 配置文件格式
  - 注释说明
  - CLI 选项
  - 使用示例
  - 高级用法
  - 开发指南

### 更新文档

- **README.md**: 更新注释格式说明
- **EXAMPLES.md**: 更新所有示例以反映新的注释格式
- **IMPLEMENTATION.md**: 更新功能说明

### 技术实现

修改文件：
- `internal/generator/generator.go`: 调整注释生成逻辑

主要变更：
1. `generateGroup()`: 分组标题从 `#` 改为 `##`
2. `generateField()`: 注释从配置项下方移到上方

### 使用方法

```bash
# 生成 .env 模板
./goctl-env-template -c config/config.go -o .env.template

# 或使用 Makefile
make demo
```

### 生成的示例

```bash
## Auth Configuration
# JWT secret key
AUTH_SECRET=

# Auth timeout in seconds
AUTH_TIMEOUT=3600

# Token refresh interval (optional)
# AUTH_REFRESH_INTERVAL=300

## DB Configuration
# Database host
DB_HOST=

# Database port
DB_PORT=3306

# Database username
DB_USERNAME=

# Database password (optional)
# DB_PASSWORD=
```

## 2025-01-31 - 初始版本

### 功能特性

✅ 从结构体标签提取环境变量
✅ 支持必填和可选变量
✅ 从结构体标签提取默认值
✅ 解析文档注释作为字段说明
✅ 按结构体嵌套级别分组配置
✅ 生成清晰、带注释的 .env 模板文件
✅ 全小写文件命名
