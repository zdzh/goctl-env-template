package config

// Auth配置包含JWT认证相关的配置项
type AuthConfig struct {
	// JWT secret key
	Secret string `json:",env=AUTH_SECRET"`

	// Auth timeout in seconds
	Timeout int `json:",env=AUTH_TIMEOUT,default=3600"`

	// Token refresh interval (optional)
	RefreshInterval int `json:",env=AUTH_REFRESH_INTERVAL,optional,default=300"`

	// Auth service endpoint (optional)
	Endpoint string `json:",env=AUTH_ENDPOINT,optional"`
}

// 数据库配置包含MySQL数据库连接相关参数
type DBConfig struct {
	// Database host
	Host string `json:",env=DB_HOST"`

	// Database port
	Port int `json:",env=DB_PORT,default=3306"`

	// Database username
	Username string `json:",env=DB_USERNAME"`

	// Database password (optional)
	Password string `json:",env=DB_PASSWORD,optional"`

	// Database name
	Database string `json:",env=DB_DATABASE"`

	// Database connection timeout (optional)
	Timeout int `json:",env=DB_TIMEOUT,optional,default=10"`

	// Max idle connections (optional)
	MaxIdleConns int `json:",env=DB_MAX_IDLE_CONNS,optional,default=10"`

	// Max open connections (optional)
	MaxOpenConns int `json:",env=DB_MAX_OPEN_CONNS,optional,default=100"`
}

// Redis配置包含Redis缓存连接相关参数
type RedisConfig struct {
	// Redis host
	Host string `json:",env=REDIS_HOST"`

	// Redis port
	Port int `json:",env=REDIS_PORT,default=6379"`

	// Redis password (optional)
	Password string `json:",env=REDIS_PASSWORD,optional"`

	// Redis database index (optional)
	DB int `json:",env=REDIS_DB,optional,default=0"`

	// Redis connection pool size (optional)
	PoolSize int `json:",env=REDIS_POOL_SIZE,optional,default=10"`

	// Redis minimum idle connections (optional)
	MinIdleConns int `json:",env=REDIS_MIN_IDLE_CONNS,optional,default=5"`
}
