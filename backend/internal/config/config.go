package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sacurio/jb-challenge/internal/util"
)

type (
	// DatabaseConfig handles all the Database info values needed to establish a connection.
	DatabaseConfig struct {
		User   string
		Pwd    string
		Port   string
		Host   string
		DbName string
	}

	// AppConfig handles all the chatbot app needed configuration.
	AppConfig struct {
		Name     string
		Port     string
		DBConfig DatabaseConfig
	}
)

// LoadConfig loads all config values from specified sources.
func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	return &AppConfig{
		Name: os.Getenv("APP_NAME"),
		Port: os.Getenv("SERVER_PORT"),
		DBConfig: DatabaseConfig{
			User:   os.Getenv("MYSQL_USER"),
			Pwd:    os.Getenv("MYSQL_PASSWORD"),
			Port:   os.Getenv("MYSQL_PORT"),
			Host:   util.ValidateStringNotEmpty(os.Getenv("MYSQL_HOST"), "localhost"),
			DbName: os.Getenv("MYSQL_DATABASE"),
		},
	}, nil
}
