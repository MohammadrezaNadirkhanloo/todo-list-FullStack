package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type envConfig struct {
	DatabaseURL string
	Port        string
}

var Config envConfig

func (e *envConfig) LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, relying on system environment variables")
	}

	e.DatabaseURL = loadRequiredString("DATABASE_URL")
	e.Port = loadString("PORT", "3000")
}

func init() {
	Config.LoadConfig()
}

func loadString(key, fallback string) string { //برای متغیرهای اختیاری
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("🛑 %s env is not set, using fallback: %s\n", key, fallback)
		return fallback
	}
	return val
}

func loadRequiredString(key string) string { //برای متغیرهای الزامی
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("🛑 required env variable %s is not set", key)
	}
	return val
}
