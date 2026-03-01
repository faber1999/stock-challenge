package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	StocksAPIURL       string
	StocksAPIToken     string
	SyncTimeout        time.Duration
	SyncMaxPages       int
	AutoSyncOnStartup  bool
	CORSAllowedOrigins string
}

var loadEnvOnce sync.Once

func Load() Config {
	loadEnvOnce.Do(func() {
		_ = godotenv.Load()
	})

	return Config{
		Port:               envOrDefault("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		StocksAPIURL:       envOrDefault("STOCKS_API_URL", "https://api.karenai.click/swechallenge/list"),
		StocksAPIToken:     envOrDefault("STOCKS_API_TOKEN", "1"),
		SyncTimeout:        time.Duration(envIntOrDefault("SYNC_TIMEOUT_SECONDS", 20)) * time.Second,
		SyncMaxPages:       envIntOrDefault("SYNC_MAX_PAGES", 50),
		AutoSyncOnStartup:  envBoolOrDefault("AUTO_SYNC_ON_STARTUP", false),
		CORSAllowedOrigins: envOrDefault("CORS_ALLOWED_ORIGINS", "*"),
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func envIntOrDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return n
}

func envBoolOrDefault(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return v
}
