package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sacurio/jb-challenge/internal/util"
	"github.com/sirupsen/logrus"
)

type (
	// DatabaseConfig handles all the Database info values needed to establish a connection.
	DBConfig struct {
		User   string
		Pwd    string
		Port   string
		Host   string
		DBName string
	}

	// AppConfig handles all the chatbot app needed configuration.
	AppConfig struct {
		Name                string
		HttpServerPort      string
		WebSocketServerPort string
		SecretKey           string
		DBConfig            DBConfig
		RabbitConfig        RabbitMQConfig
		BotConfig           BotConfig
		Logger              *logrus.Logger
	}

	// RabbitMQConfig handles allthe RabbitMQ configuration.
	RabbitMQConfig struct {
		User     string
		Password string
		Host     string
		Port     string
	}

	BotConfig struct {
		Name     string
		URL      string
		Wildcard string
	}
)

// LoadConfig loads all config values from specified sources.
func LoadConfig(logger *logrus.Logger) (*AppConfig, error) {
	logger.Info("starting to load environment variables...")
	file := ".env-local"
	variablesMsg := "local"
	if os.Getenv("ENVIRONMENT") != "" {
		variablesMsg = "production"
		file = ".env"
	}

	logger.Infof("Ready to load %s variables.", variablesMsg)

	if err := godotenv.Load(file); err != nil {
		logger.Error("error on reading environment variables")
		return nil, err
	}

	logger.Info("Environment variables loaded successfully")
	return &AppConfig{
		Name:                os.Getenv("APP_NAME"),
		HttpServerPort:      os.Getenv("SERVER_PORT"),
		SecretKey:           os.Getenv("APP_SECRET_KEY"),
		WebSocketServerPort: os.Getenv("WEBSOCKET_SERVER_PORT"),
		DBConfig: DBConfig{
			User:   os.Getenv("MYSQL_USER"),
			Pwd:    os.Getenv("MYSQL_PASSWORD"),
			Port:   os.Getenv("MYSQL_PORT"),
			Host:   util.ValidateStringNotEmpty(os.Getenv("MYSQL_HOST"), "localhost"),
			DBName: os.Getenv("MYSQL_DATABASE"),
		},
		RabbitConfig: RabbitMQConfig{
			User:     os.Getenv("RABBITMQ_DEFAULT_USER"),
			Password: os.Getenv("RABBITMQ_DEFAULT_PASS"),
			Host:     util.ValidateStringNotEmpty(os.Getenv("RABBITMQ_HOST"), "localhost"),
			Port:     os.Getenv("RABBITMQ_DEFAULT_PORT"),
		},
		BotConfig: BotConfig{
			Name:     os.Getenv("BOT_NAME"),
			URL:      os.Getenv("BOT_URL"),
			Wildcard: os.Getenv("BOT_WILDCARD"),
		},
		Logger: logger,
	}, nil
}
