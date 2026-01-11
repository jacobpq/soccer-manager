package config

import (
	"os"
)

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBDSN:     getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/soccer_manager?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
