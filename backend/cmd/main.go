package main

import (
	"github.com/sacurio/jb-challenge/internal/app/user"
	"github.com/sacurio/jb-challenge/internal/config"
	"github.com/sacurio/jb-challenge/internal/db"
	"github.com/sacurio/jb-challenge/internal/server"
	"github.com/sacurio/jb-challenge/internal/util"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	app, err := config.LoadConfig()
	if err != nil {
		log.Panicf("app config values was not loaded: %s", err.Error())
	}
	log.Infof("starting %s app...", util.ValidateStringNotEmpty(app.Name, util.DefaultAppName))
	log.Infof("%v", app)

	dbConn, err := db.ConnectMySQL()
	if err != nil {
		return
	}

	userRepo := user.NewUserRepository(dbConn)
	userValidator := user.NewUserValidator()
	userUserService := user.NewUserService(userRepo, userValidator)
	srv := server.NewServer(util.ValidateStringNotEmpty(app.Port, util.DefaultPort), userUserService, log)

	srv.StartHTTPServer()
}
