package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MPAccessToken string
	ServerPort    string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (c *Config) Load() {
	c.MPAccessToken = os.Getenv("MP_ACCESS_TOKEN")
	c.ServerPort = os.Getenv("SERVER_PORT")
}