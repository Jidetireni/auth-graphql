package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DatabaseURL string `env:"DB_OPEN"`
}

type ServerConfig struct {
	Port string `env:"PORT"`
	Host string `env:"HOST"`
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using default values")
	}
	config := &Config{
		Database: DatabaseConfig{
			DatabaseURL: os.Getenv("DB_OPEN"),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
	}

	return config, nil

}

func (c *Config) GetDatabaseURL() string {
	cfg, err := New()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return cfg.Database.DatabaseURL
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
