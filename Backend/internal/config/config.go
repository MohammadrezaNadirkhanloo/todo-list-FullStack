package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string           `mapstructure:"env"`
	Server     ServerConfig     `mapstructure:"server"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Authz      AuthzConfig      `mapstructure:"authz"`
	Password   PasswordConfig   `mapstructure:"password"`
	Pagination PaginationConfig `mapstructure:"pagination"`
	Cors       CORSConfig       `mapstructure:"cors"`
	CSRF       CSRFConfig       `mapstructure:"csrf"`
	RateLimit  RateLimitConfig  `mapstructure:"ratelimit"`
}

type ServerConfig struct {
	Port              string        `mapstructure:"port"`
	Domain            string        `mapstructure:"domain"`
	RunMode           string        `mapstructure:"runMode"`
	APIBasePath       string        `mapstructure:"apiBasePath"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdownTimeout"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
	RequestTimeout    time.Duration `mapstructure:"requestTimeout"`
	MaxBodyBytes      int64         `mapstructure:"maxBodyBytes"`
	InsecureCookies   bool          `mapstructure:"insecureCookies"`
	TrustedProxies    []string      `mapstructure:"trustedProxies"`
}

type PostgresConfig struct {
	Host               string        `mapstructure:"host"`
	Port               string        `mapstructure:"port"`
	User               string        `mapstructure:"user"`
	Password           string        `mapstructure:"password"`
	DBName             string        `mapstructure:"dbName"`
	SSLMode            string        `mapstructure:"sslMode"`
	TimeZone           string        `mapstructure:"timeZone"`
	MaxIdleConns       int           `mapstructure:"maxIdleConns"`
	MaxOpenConns       int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime    time.Duration `mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime    time.Duration `mapstructure:"connMaxIdleTime"`
	SlowQueryThreshold time.Duration `mapstructure:"slowQueryThreshold"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	PoolSize     int           `mapstructure:"poolSize"`
	MinIdleConns int           `mapstructure:"minIdleConns"`
}

type JWTConfig struct {
	AccessSecret string        `mapstructure:"accessSecret"`
	AccessTTL    time.Duration `mapstructure:"accessTTL"`
	RefreshTTL   time.Duration `mapstructure:"refreshTTL"`
	Issuer       string        `mapstructure:"issuer"`
	Audience     string        `mapstructure:"audience"`
}

type AuthzMode string

const (
	AuthzRoles AuthzMode = "roles"
	AuthzRules AuthzMode = "rules"
	AuthzBoth AuthzMode = "both"
)

type AuthzConfig struct {
	Enforce bool      `mapstructure:"enforce"`
	Mode    AuthzMode `mapstructure:"mode"`
}

type PasswordConfig struct {
	MinLength         int    `mapstructure:"minLength"`
	MaxLength         int    `mapstructure:"maxLength"`
	IncludeDigits     bool   `mapstructure:"includeDigits"`
	IncludeUppercase  bool   `mapstructure:"includeUppercase"`
	IncludeLowercase  bool   `mapstructure:"includeLowercase"`
	IncludeSymbols    bool   `mapstructure:"includeSymbols"`
	Argon2Memory      uint32 `mapstructure:"argon2Memory"`
	Argon2Iterations  uint32 `mapstructure:"argon2Iterations"`
	Argon2Parallelism uint8  `mapstructure:"argon2Parallelism"`
}

type PaginationConfig struct {
	DefaultPageSize int `mapstructure:"defaultPageSize"`
	MaxPageSize     int `mapstructure:"maxPageSize"`
}

type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allowOrigins"`
	AllowCredentials bool          `mapstructure:"allowCredentials"`
	MaxAge           time.Duration `mapstructure:"maxAge"`
}

type CSRFConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type RateLimitConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	GlobalRPS   float64       `mapstructure:"globalRPS"`
	GlobalBurst int           `mapstructure:"globalBurst"`
	AuthRPS     float64       `mapstructure:"authRPS"`
	AuthBurst   int           `mapstructure:"authBurst"`
	IdleTTL     time.Duration `mapstructure:"idleTTL"`
	MaxKeys     int           `mapstructure:"maxKeys"`
}

func (c *Config) IsProduction() bool { return c.Env == "production" }

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode, p.TimeZone,
	)
}

func (r RedisConfig) Addr() string { return r.Host + ":" + r.Port }

const envPrefix = "APP"

func Load(configDir string) (*Config, error) {
	v := viper.New()
	setDefaults(v)
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnvKeys(v)

	if configDir != "" {
		if err := mergeFileIfExists(v, filepath.Join(configDir, "config.yml")); err != nil {
			return nil, err
		}
		if env := v.GetString("env"); env != "" {
			path := filepath.Join(configDir, "config."+env+".yml")
			if err := mergeFileIfExists(v, path); err != nil {
				return nil, err
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

func mergeFileIfExists(v *viper.Viper, path string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("config: could not access %q: %w", path, err)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config: failed to open %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	v.SetConfigType("yml")
	if err := v.MergeConfig(f); err != nil {
		return fmt.Errorf("config: failed to merge %q: %w", path, err)
	}
	return nil
}
