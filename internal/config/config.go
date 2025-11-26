package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port          string
	DBPath        string
	MigrationPath string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBPath:        getEnv("DB_PATH", "./data/cronnor.db"),
		MigrationPath: getEnv("MIGRATION_PATH", "./migrations/001_initial_schema.sql"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
