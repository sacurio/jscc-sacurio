package main

import (
	"github.com/sacurio/jb-challenge/internal/app/message"
	"github.com/sacurio/jb-challenge/internal/app/user"
	"github.com/sacurio/jb-challenge/internal/config"
	"github.com/sacurio/jb-challenge/internal/infraestructure/db"
	"github.com/sacurio/jb-challenge/internal/infraestructure/rabbitmq"
	"github.com/sacurio/jb-challenge/internal/infraestructure/server"
	"github.com/sacurio/jb-challenge/internal/util"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	app, err := config.LoadConfig()
	if err != nil {
		log.Panicf("App config values was not loaded: %s", err.Error())
	}
	log.Infof("Starting %s app...", util.ValidateStringNotEmpty(app.Name, util.DefaultAppName))

	dbHandler := db.NewDatabase(app.DBConfig.User, app.DBConfig.Pwd, app.DBConfig.Port, app.DBConfig.Host, app.DBConfig.DbName, log)
	dbHandler.SetupDB()

	rabbitMQConfig := rabbitmq.NewRabbitMQConfig(app.RMQConfig.User, app.RMQConfig.Password, app.RMQConfig.Host, app.RMQConfig.Port)
	rabbitMQHandler, err := rabbitmq.NewRabbitMQ(rabbitMQConfig, log)
	if err != nil {
		log.Panicf("Error on connecting to RabbitMQ service: %s", err.Error())
	}

	defer rabbitMQHandler.Close()

	messageService := message.NewMessageService(rabbitMQHandler)

	userRepo := user.NewUserRepository(dbHandler.DB)
	userValidator := user.NewUserValidator()
	userUserService := user.NewUserService(userRepo, userValidator)

	srv := server.NewServer(util.ValidateStringNotEmpty(app.Port, util.DefaultPort), userUserService, messageService, log)
	srv.StartHTTPServer()
}
