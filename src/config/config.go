package config

import (
	"os"
)

type Config struct {
	DBURL    string
	APIPort  string
	WALPath  string
	LogLevel string
}

func LoadConfig() *Config {
	return &Config{
		DBURL:    getEnv("DB_URL", "postgres://avocato:avocato_password@localhost:5432/avocato_ledger"),
		APIPort:  getEnv("API_PORT", "8080"),
		WALPath:  getEnv("WAL_PATH", "data/ledger.wal"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
