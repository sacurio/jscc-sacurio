package config

import (
	"os"

	"github.com/joho/godotenv"
)

// AppConfig handles all the chatbot app needed configuration.
type AppConfig struct {
	Name string
	Port string
}

// LoadConfig loads all config values from specified sources.
func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	return &AppConfig{
		Name: os.Getenv("APP_NAME"),
		Port: os.Getenv("SERVER_PORT"),
	}, nil
}
