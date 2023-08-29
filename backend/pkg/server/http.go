package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"github.com/sacurio/jb-challenge/internal/app/handler"
	"github.com/sacurio/jb-challenge/internal/app/service"
	"github.com/sacurio/jb-challenge/pkg/websocket"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	server Server
}

// Server represents server struct needed to execute HTTP Server
type Server struct {
	port             string
	jwtService       service.JWTManager
	userService      service.User
	webSocketService websocket.WebSocketService
	log              *logrus.Logger
}

// NewServer returns a new instance of Server.
func NewServer(
	port string,
	jwtService service.JWTManager,
	userService service.User,
	webSocketService websocket.WebSocketService,
	log *logrus.Logger,
) *Server {
	return &Server{
		port:             port,
		jwtService:       jwtService,
		userService:      userService,
		webSocketService: webSocketService,
		log:              log,
	}
}

// StartHTTPServer set the HTTP Server configuration to exposes the services to clients.
func (s *Server) StartHTTPServer() {
	r := mux.NewRouter()
	r.Use(
		LoggingMiddleware(s.log),
		MetricsMiddleware(s.log),
		EnableCORSMiddleware,
	)

	public := r.NewRoute().Subrouter()
	protected := r.NewRoute().Subrouter()
	chat := r.NewRoute().Subrouter()

	public.HandleFunc("/user/validate/", handler.ValidateUser(s.userService, s.jwtService, s.log)).Methods("POST")

	protected.HandleFunc("/ws/", s.webSocketService.Serve())

	chat.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/chat/login", http.StatusOK)
	})
	chat.HandleFunc("/chat/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})
	chat.HandleFunc("/chat/index", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	staticFileDirectory := http.Dir("./static")
	staticFileHandler := http.StripPrefix("/chat/", http.FileServer(staticFileDirectory))
	chat.PathPrefix("/chat/").Handler(staticFileHandler)

	avatarDirectory := http.Dir("./static/assets/images/avatars")
	avatarHandler := http.StripPrefix("/chat/assets/images/avatars/", http.FileServer(avatarDirectory))
	chat.PathPrefix("/chat/assets/images/avatars/").Handler(avatarHandler)

	s.log.Infof("Starting HTTP Server on: %s", s.port)

	if err := http.ListenAndServe(":"+s.port, r); err != nil {
		log.Fatal(err)
	}
}
