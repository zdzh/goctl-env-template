# goctl-env-template Examples

## Basic Example

### 1. Create config/config.go

```go
package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// Server name
	Name string `json:",env=SERVER_NAME"`

	// Server host (default: 0.0.0.0)
	Host string `json:",env=SERVER_HOST,default=0.0.0.0"`

	// Server port (default: 8888)
	Port int `json:",env=SERVER_PORT,default=8888"`

	// Server mode (optional, default: dev)
	Mode string `json:",env=SERVER_MODE,optional,default=dev"`
}
```

### 2. Generate .env template

```bash
./goctl-env-template -c config/config.go -o .env.template
```

### 3. Output (.env.template)

```bash
## Configuration
# Server name
SERVER_NAME=

# Server host (default: 0.0.0.0)
SERVER_HOST=0.0.0.0

# Server port (default: 8888)
SERVER_PORT=8888

# Server mode (optional, default: dev)
# SERVER_MODE=dev
```

## Nested Struct Example

### 1. Create config/config.go with nested structs

```go
package config

type Config struct {
	// Server name
	Name string `json:",env=SERVER_NAME"`

	Auth AuthConfig
	DB   DBConfig
}

type AuthConfig struct {
	// JWT secret key
	Secret string `json:",env=AUTH_SECRET"`

	// Auth timeout (default: 3600)
	Timeout int `json:",env=AUTH_TIMEOUT,default=3600"`

	// Token refresh interval (optional)
	RefreshInterval int `json:",env=AUTH_REFRESH_INTERVAL,optional"`
}

type DBConfig struct {
	// Database host
	Host string `json:",env=DB_HOST"`

	// Database port (default: 3306)
	Port int `json:",env=DB_PORT,default=3306"`

	// Database username
	Username string `json:",env=DB_USERNAME"`

	// Database password (optional)
	Password string `json:",env=DB_PASSWORD,optional"`
}
```

### 2. Generate .env template

```bash
./goctl-env-template -c config/config.go -o .env.template
```

### 3. Output (.env.template)

```bash
## Configuration
# Server name
SERVER_NAME=

## Auth Configuration
# JWT secret key
AUTH_SECRET=

# Auth timeout (default: 3600)
AUTH_TIMEOUT=3600

# Token refresh interval (optional)
# AUTH_REFRESH_INTERVAL=

## DB Configuration
# Database host
DB_HOST=

# Database port (default: 3306)
DB_PORT=3306

# Database username
DB_USERNAME=

# Database password (optional)
# DB_PASSWORD=
```


## Advanced Features

### Custom Grouping

Nested structs are automatically grouped by their names:

```go
type Config struct {
	Redis RedisConfig
	Cache CacheConfig
}
```

Will generate separate sections:
- `## Redis Configuration`
- `## Cache Configuration`

### Default Values

Specify default values in the struct tag:

```go
Port int `json:",env=PORT,default=8080"`
```

### Optional Fields

Mark fields as optional to comment them out in the template:

```go
Password string `json:",env=PASSWORD,optional"`
```

### Doc Comments

Use Go doc comments to add descriptions:

```go
// Database connection timeout in seconds
Timeout int `json:",env=DB_TIMEOUT,default=10"`
```

This will appear as:
```bash
# Database connection timeout in seconds
DB_TIMEOUT=10
```
