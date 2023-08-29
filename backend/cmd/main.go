package main

import (
	repository "github.com/sacurio/jb-challenge/internal/app/repository/user"
	"github.com/sacurio/jb-challenge/internal/app/service"
	"github.com/sacurio/jb-challenge/internal/config"
	"github.com/sacurio/jb-challenge/internal/util"
	db "github.com/sacurio/jb-challenge/pkg/db/mysql"
	"github.com/sacurio/jb-challenge/pkg/server"
	"github.com/sacurio/jb-challenge/pkg/websocket"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Panicf("App config values was not loaded: %s", err.Error())
	}
	log.Infof("Starting %s app...", util.ValidateStringNotEmpty(cfg.Name, util.DefaultAppName))

	jwtService := service.NewJWT([]byte(cfg.SecretKey), log)

	dbHandler := db.NewDB(cfg.DBConfig, log)

	userRepository := repository.NewUser(dbHandler.DB)
	userService := service.NewService(userRepository)

	webSocketService := websocket.NewWebSocketServer(cfg.WebSocketServerPort, jwtService, cfg.Logger)
	webSocketService.Start()

	srv := server.NewServer(
		util.ValidateStringNotEmpty(
			cfg.HttpServerPort,
			util.DefaultPort,
		),
		jwtService,
		userService,
		webSocketService,
		log)
	srv.StartHTTPServer()
}
