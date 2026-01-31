# goctl-env-template

A goctl plugin that extracts environment variables from `config/config.go` and generates a `.env` template file.

## Features

- Extract environment variables from struct tags (`json:"name,env=VAR_NAME"`)
- Support required and optional variables (marked with `optional` tag)
- Extract default values from struct tags (`default=value`)
- Parse doc comments as field descriptions
- Group configuration by struct nesting level
- Generate clean, well-commented `.env` template files

## Installation

```bash
go install github.com/zdzh/goctl-env-template@latest
```

Make sure the installed `goctl-env-template` is in your `$PATH`.

## Usage

### As a standalone CLI

```bash
goctl-env-template -c config/config.go -o .env.template
```

### As a goctl plugin

```bash
goctl plugin -p goctl-env-template --config config/config.go --output .env.template
```

## Config File Format

Define your configuration in `config/config.go` with environment variable tags:

```go
package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth  AuthConfig
	DB    DBConfig
}

type AuthConfig struct {
	Host     string `json:",env=AUTH_HOST"`
	Port     int    `json:",env=AUTH_PORT"`
	Secret   string `json:",env=AUTH_SECRET,optional"`
	Timeout  int    `json:",env=AUTH_TIMEOUT,default=30"`
}

type DBConfig struct {
	Host     string `json:",env=DB_HOST"`
	Port     int    `json:",env=DB_PORT,default=3306"`
	Username string `json:",env=DB_USERNAME"`
	Password string `json:",env=DB_PASSWORD,optional"`
}
```

## Generated .env Template

The plugin will generate a `.env.template` file with:

- Group headers marked with `##`
- Active lines for required variables
- Commented lines for optional variables (marked with `#`)
- Field comments (with `#`) placed above configuration items
- Default values as examples

Example output:

```bash
## Auth Configuration
# JWT secret key
AUTH_SECRET=

# Auth timeout in seconds
AUTH_TIMEOUT=30

# Token refresh interval (optional)
# AUTH_REFRESH_INTERVAL=

## DB Configuration
# Database host address
DB_HOST=

# Database port (default: 3306)
DB_PORT=3306

# Database username
DB_USERNAME=

# Database password (optional)
# DB_PASSWORD=
```

## Struct Tag Format

Use the following tag format in your config structs:

```go
FieldName FieldType `json:"name,env=ENV_VAR,optional,default=value"`
```

## Comment Format

- **Group headers**: Use `##` to mark configuration groups
- **Field comments**: Use `#` to mark field descriptions
- **Comment position**: Field comments are placed above configuration items

Example:

```bash
## Database Configuration
# Database host address
DB_HOST=

# Database port (default: 3306)
DB_PORT=3306
```

### Tag Options

- `env=VAR_NAME`: Environment variable name (required)
- `optional`: Mark as optional variable (optional)
- `default=value`: Default value (optional)

## CLI Options

```
Usage: goctl-env-template [flags]

Flags:
  -c, --config string   Config file path (default "config/config.go")
  -o, --output string   Output file path (default ".env.template")
  -h, --help            Show help
```

## Examples

### Basic usage with default paths

```bash
goctl-env-template
```

### Specify custom config and output paths

```bash
goctl-env-template -c internal/config/config.go -o .env.example
```

### Use with goctl plugin

```bash
goctl plugin -p goctl-env-template --config config/config.go --output .env
```

## Project Structure

```
goctl-env-template/
├── cmd/
│   └── root.go           # CLI command definition
├── internal/
│   ├── parser/
│   │   └── parser.go     # Go AST parser
│   ├── generator/
│   │   └── generator.go  # .env template generator
│   └── types/
│       └── config.go     # Internal data structures
├── plugin/
│   └── main.go           # goctl plugin entry point
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License
