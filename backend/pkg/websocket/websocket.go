package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sacurio/jb-challenge/internal/app/model"
	"github.com/sirupsen/logrus"
)

type (
	WebSocketService interface {
		Read(client *Client) error
		Broadcast() error
		Start()
		ServeWs() http.HandlerFunc
	}

	webSocketServer struct {
		port      string
		clients   map[*Client]bool
		broadcast chan *model.Chat
		upgrader  websocket.Upgrader
		logger    *logrus.Logger
	}

	Client struct {
		Conn     *websocket.Conn
		Username string
	}

	Message struct {
		Type string     `json:"type"`
		User string     `json:"user,omitempty"`
		Chat model.Chat `json:"chat,omitempty"`
	}
)

func NewWebSocketServer(port string, logger *logrus.Logger) WebSocketService {
	return &webSocketServer{
		port:      port,
		clients:   make(map[*Client]bool),
		broadcast: make(chan *model.Chat),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		logger: logger,
	}
}

// Read define a receiver which will listen for new messages being sent to our WebSocket endpoint
func (wss *webSocketServer) Read(client *Client) error {
	for {
		// read in a message
		// readMessage returns messageType, message, err
		// messageType: 1-> Text Message, 2 -> Binary Message
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			wss.logger.Infof("error trying to read message from client: %s", err.Error())
			return err
		}

		m := &Message{}
		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("error while unmarshaling chat", err)
			continue
		}

		if m.Type == "bootup" {
			client.Username = m.User
			wss.logger.WithFields(logrus.Fields{
				"client info": client,
				"username":    client.Username,
			}).Info("client successfully mapped")
		} else {
			wss.logger.WithField(
				"type", m.Type).Info("received message")
			c := m.Chat
			c.Timestamp = time.Now().Unix()

			// // save in redis
			// id, err := redisrepo.CreateChat(&c)
			// if err != nil {
			// 	log.Println("error while saving chat in redis", err)
			// 	return
			// }

			wss.broadcast <- &c
		}
	}
}

func (wss *webSocketServer) Broadcast() error {
	for {
		message := <-wss.broadcast
		wss.logger.Println("new message to be broadcasted", message)

		for client := range wss.clients {
			err := client.Conn.WriteJSON(message)
			if err != nil {
				wss.logger.Errorf("Websocket error: %s", err)
				client.Conn.Close()
				delete(wss.clients, client)
			}
		}
	}
}

// define our WebSocket endpoint
func (wss *webSocketServer) ServeWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wss.logger.WithFields(logrus.Fields{
			"host":                r.Host,
			"URL query variables": r.URL.Query(),
		}).Info("Listening WebSocket request")

		ws, err := wss.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		client := &Client{Conn: ws}

		wss.clients[client] = true
		wss.logger.WithFields(logrus.Fields{
			"quantity":       len(wss.clients),
			"remote address": ws.RemoteAddr(),
		}).Info("connected clients")

		wss.Read(client)

		wss.logger.WithFields(logrus.Fields{
			"remote address": ws.RemoteAddr().String(),
		}).Info("closing connection")
		delete(wss.clients, client)
	}
}

func (wss *webSocketServer) Start() {
	// RabbitMQ
	wss.logger.Info("Staring websocket server...")
	go wss.Broadcast()
}
