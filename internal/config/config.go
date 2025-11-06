package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	Secret []byte
	TTL    time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if not exists)
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}

	return &Config{
		Server: ServerConfig{
			Port:            port,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 5 * time.Second,
		},
		Database: DatabaseConfig{
			Path: dbPath,
		},
		JWT: JWTConfig{
			Secret: []byte(jwtSecret),
			TTL:    24 * time.Hour,
		},
	}, nil
}
