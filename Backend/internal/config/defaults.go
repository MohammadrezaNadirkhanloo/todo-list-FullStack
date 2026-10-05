package config

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) {
	v.SetDefault("env", "development")
 
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.runMode", "debug")
	v.SetDefault("server.apiBasePath", "/api/v1")
	v.SetDefault("server.readTimeout", "10s")
	v.SetDefault("server.writeTimeout", "15s")
	v.SetDefault("server.shutdownTimeout", "15s")
	v.SetDefault("server.idleTimeout", "120s")
	v.SetDefault("server.readHeaderTimeout", "5s")
	v.SetDefault("server.requestTimeout", "25s")
	v.SetDefault("server.maxBodyBytes", 1048576)
	v.SetDefault("server.insecureCookies", true)
	v.SetDefault("server.trustedProxies", []string{})
 
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.dbName", "tododb")
	v.SetDefault("postgres.sslMode", "disable")
	v.SetDefault("postgres.timeZone", "UTC")
	v.SetDefault("postgres.maxIdleConns", 10)
	v.SetDefault("postgres.maxOpenConns", 50)
	v.SetDefault("postgres.connMaxLifetime", "30m")
	v.SetDefault("postgres.connMaxIdleTime", "5m")
	v.SetDefault("postgres.slowQueryThreshold", "200ms")
 
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.dialTimeout", "5s")
	v.SetDefault("redis.readTimeout", "3s")
	v.SetDefault("redis.writeTimeout", "3s")
	v.SetDefault("redis.poolSize", 10)
	v.SetDefault("redis.minIdleConns", 2)
 
	v.SetDefault("jwt.accessTTL", "15m")
	v.SetDefault("jwt.refreshTTL", "168h")
	v.SetDefault("jwt.issuer", "todo-api")
	v.SetDefault("jwt.audience", "todo-api")
 
	v.SetDefault("authz.enforce", true)
 
	v.SetDefault("password.minLength", 12)
	v.SetDefault("password.maxLength", 128)
	v.SetDefault("password.includeDigits", true)
	v.SetDefault("password.includeUppercase", true)
	v.SetDefault("password.includeLowercase", true)
	v.SetDefault("password.includeSymbols", false)
	v.SetDefault("password.argon2Memory", 65536)
	v.SetDefault("password.argon2Iterations", 3)
	v.SetDefault("password.argon2Parallelism", 2)
 
	v.SetDefault("pagination.defaultPageSize", 20)
	v.SetDefault("pagination.maxPageSize", 100)
}
 

var secretKeys = []string{
	"postgres.password",
	"redis.password",
	"jwt.accessSecret",
	"jwt.refreshSecret",
}

func bindEnvKeys(v *viper.Viper) {
	for _, key := range secretKeys {
		_ = v.BindEnv(key)
	}
}
