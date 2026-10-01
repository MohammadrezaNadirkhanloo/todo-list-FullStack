package config

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.runMode", "debug")
	v.SetDefault("server.readTimeout", "10s")
	v.SetDefault("server.writeTimeout", "15s")
	v.SetDefault("server.shutdownTimeout", "15s")
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
