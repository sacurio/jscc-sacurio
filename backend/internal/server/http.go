package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"github.com/sacurio/jb-challenge/internal/app/user"
	server "github.com/sacurio/jb-challenge/internal/server/user"
	"github.com/sirupsen/logrus"
)

// Server represents server struct needed to execute HTTP Server
type Server struct {
	port        string
	userService user.UserService
	log         *logrus.Logger
}

// NewServer returns a new instance of Server.
func NewServer(port string, userService user.UserService, log *logrus.Logger) *Server {
	return &Server{
		port:        port,
		userService: userService,
		log:         log,
	}
}

// StartHTTPServer set the HTTP Server configuration to exposes the services to clients.
func (s *Server) StartHTTPServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", server.DefaultHandler)
	r.HandleFunc("/user/validate/", func(w http.ResponseWriter, r *http.Request) {
		server.ValidateUser(w, r, s.userService)
	}).Methods("POST")

	s.log.Infof("Starting HTTP Server on: %s", s.port)

	if err := http.ListenAndServe(":"+s.port, r); err != nil {
		log.Fatal(err)
	}
}
