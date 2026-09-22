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

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
)

type Config struct {
	Env        Environment      `mapstructure:"env"`
	Server     ServerConfig     `mapstructure:"server"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Cors       CorsConfig       `mapstructure:"cors"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Authz      AuthzConfig      `mapstructure:"authz"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Otp        OtpConfig        `mapstructure:"otp"`
	Password   PasswordConfig   `mapstructure:"password"`
	RateLimit  RateLimitConfig  `mapstructure:"ratelimit"`
	Pagination PaginationConfig `mapstructure:"pagination"`
}

func (c *Config) IsProduction() bool { return c.Env == EnvProduction } //?

type ServerConfig struct {
	Port              string        `mapstructure:"port"`
	Domain            string        `mapstructure:"domain"`
	RunMode           string        `mapstructure:"runMode"` //debug | release | test
	APIBasePath       string        `mapstructure:"apiBasePath"`
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdownTimeout"`
	RequestTimeout    time.Duration `mapstructure:"requestTimeout"`
	MaxBodyBytes      int64         `mapstructure:"maxBodyBytes"`
	TrustedProxies    []string      `mapstructure:"trustedProxies"`
	InsecureCookies   bool          `mapstructure:"insecureCookies"`
}

type RoutePosture string

const (
	PostureEnterprise RoutePosture = "enterprise"
	PosturePublic     RoutePosture = "public"
)

type AuthConfig struct {
	Posture RoutePosture `mapstructure:"posture"`
}

type AuthzMode string

const (
	AuthzRoles AuthzMode = "roles"
	AuthzRules AuthzMode = "rules"
	AuthzBoth  AuthzMode = "both"
)

type AuthzConfig struct {
	Mode    AuthzMode `mapstructure:"mode"`
	Enforce bool      `mapstructure:"enforce"`
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

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode, p.TimeZone,
	)
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

func (r RedisConfig) Addr() string { return fmt.Sprintf("%s:%s", r.Host, r.Port) }

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"accessSecret"`
	RefreshSecret string        `mapstructure:"refreshSecret"`
	AccessTTL     time.Duration `mapstructure:"accessTTL"`
	RefreshTTL    time.Duration `mapstructure:"refreshTTL"`
	Issuer        string        `mapstructure:"issuer"`
	Audience      string        `mapstructure:"audience"`
}

type CorsConfig struct {
	AllowOrigins     []string      `mapstructure:"allowOrigins"`
	AllowCredentials bool          `mapstructure:"allowCredentials"`
	MaxAge           time.Duration `mapstructure:"maxAge"`
}

type LoggerConfig struct {
	Backend        string `mapstructure:"backend"`
	Level          string `mapstructure:"level"`
	Encoding       string `mapstructure:"encoding"`
	FilePath       string `mapstructure:"filePath"`
	LogRequestBody bool   `mapstructure:"logRequestBody"`
}

type OtpConfig struct {
	Enabled    bool          `mapstructure:"enabled"`
	Digits     int           `mapstructure:"digits"`
	ExpireTime time.Duration `mapstructure:"expireTime"`

	ResendDelay time.Duration `mapstructure:"resendDelay"`

	MaxAttempts int `mapstructure:"maxAttempts"`
}

type PasswordConfig struct {
	MinLength        int  `mapstructure:"minLength"`
	MaxLength        int  `mapstructure:"maxLength"`
	IncludeDigits    bool `mapstructure:"includeDigits"`
	IncludeUppercase bool `mapstructure:"includeUppercase"`
	IncludeLowercase bool `mapstructure:"includeLowercase"`
	IncludeSymbols   bool `mapstructure:"includeSymbols"`

	Argon2Memory      uint32 `mapstructure:"argon2Memory"`
	Argon2Iterations  uint32 `mapstructure:"argon2Iterations"`
	Argon2Parallelism uint8  `mapstructure:"argon2Parallelism"`
}

type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`

	GlobalRPS   float64 `mapstructure:"globalRPS"`
	GlobalBurst int     `mapstructure:"globalBurst"`

	AuthRPS   float64 `mapstructure:"authRPS"`
	AuthBurst int     `mapstructure:"authBurst"`

	OtpRPS   float64 `mapstructure:"otpRPS"`
	OtpBurst int     `mapstructure:"otpBurst"`

	IdleTTL time.Duration `mapstructure:"idleTTL"`
	MaxKeys int           `mapstructure:"maxKeys"`
}

type PaginationConfig struct {
	DefaultPageSize int `mapstructure:"defaultPageSize"`

	MaxPageSize int `mapstructure:"maxPageSize"`
}

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
	if err := cfg.Validate(); err != nil {
		return nil, err
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
