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
	Env    string       `mapstructure:"env"`
	Server ServerConfig `mapstructure:"server"`
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
