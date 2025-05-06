package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MPAccessToken string
	ServerPort    string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}

	config := Config{
		MPAccessToken: os.Getenv("MP_ACCESS_TOKEN"),
		ServerPort:    os.Getenv("SERVER_PORT"),
	}

	if config.MPAccessToken == "" {
		return Config{}, fmt.Errorf("MP_ACCESS_TOKEN not set in environment")
	}

	if config.ServerPort == "" {
		config.ServerPort = "8080" // Default server port
	}

	return config, nil
}
