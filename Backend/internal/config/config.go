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
	Env      string         `mapstructure:"env"`
	Server   ServerConfig   `mapstructure:"server"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type ServerConfig struct {
	Port              string        `mapstructure:"port"`
	RunMode           string        `mapstructure:"runMode"` // debug | release
	ReadTimeout       time.Duration `mapstructure:"readTimeout"`
	WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdownTimeout"`
	TrustedProxies    []string      `mapstructure:"trustedProxies"`
	IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
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

func mergeFileIfExists(v *viper.Viper, path string) error { // mergeFileIfExists فایل را در صورت وجود ادغام می‌کند.
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
