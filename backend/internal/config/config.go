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
		Name      string
		Port      string
		DBConfig  DatabaseConfig
		RMQConfig RabbitMQConfig
	}

	// RabbitMQConfig handles allthe RabbitMQ configuration.
	RabbitMQConfig struct {
		User     string
		Password string
		Host     string
		Port     string
	}
)

// LoadConfig loads all config values from specified sources.
func LoadConfig() (*AppConfig, error) {
	file := ".env-local"
	if os.Getenv("ENVIRONMENT") != "" {
		file = ".env"
	}

	if err := godotenv.Load(file); err != nil {
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
		RMQConfig: RabbitMQConfig{
			User:     os.Getenv("RABBITMQ_DEFAULT_USER"),
			Password: os.Getenv("RABBITMQ_DEFAULT_PASS"),
			Host:     util.ValidateStringNotEmpty(os.Getenv("RABBITMQ_HOST"), "localhost"),
			Port:     os.Getenv("RABBITMQ_DEFAULT_PORT"),
		},
	}, nil
}
