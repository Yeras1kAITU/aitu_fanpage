package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URI  string
	Name string
}

func Load() *Config {
	port := getEnv("PORT", "8080")

	return &Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			URI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Name: getEnv("DB_NAME", "aitu_fanpage"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
