package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sacurio/jb-challenge/internal/app/model"
	"github.com/sacurio/jb-challenge/internal/app/service"
	"github.com/sacurio/jb-challenge/internal/config"
	"github.com/sirupsen/logrus"
)

type (
	WebSocketService interface {
		Read(client *Client) error
		Broadcast() error
		Start()
		Serve() http.HandlerFunc
	}

	webSocketServer struct {
		currentUser      string
		maxHistoryMsg    string
		port             string
		clients          map[*Client]bool
		broadcastChannel chan *model.Chat
		errorChannel     chan *model.ChatError
		historyChannel   chan []model.Message
		jwtService       service.JWTManager
		botService       service.Bot
		userService      service.User
		messageService   service.Message
		upgrader         websocket.Upgrader
		logger           *logrus.Logger
	}

	Client struct {
		Conn     *websocket.Conn
		Username string
	}

	Message struct {
		Type string     `json:"type"`
		User string     `json:"user,omitempty"`
		JWT  string     `json:"jwt,omitempty"`
		Chat model.Chat `json:"chat,omitempty"`
	}
)

func NewWebSocketServer(
	jwtService service.JWTManager,
	botService service.Bot,
	userService service.User,
	messageService service.Message,
	cfg *config.AppConfig) WebSocketService {
	return &webSocketServer{
		maxHistoryMsg:    cfg.MaxHistoryMsgs,
		port:             cfg.WebSocketServerPort,
		clients:          make(map[*Client]bool),
		broadcastChannel: make(chan *model.Chat),
		errorChannel:     make(chan *model.ChatError),
		historyChannel:   make(chan []model.Message),
		jwtService:       jwtService,
		botService:       botService,
		userService:      userService,
		messageService:   messageService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		logger: cfg.Logger,
	}
}

// Read define a receiver which will listen for new messages being sent to our WebSocket endpoint
func (wss *webSocketServer) Read(client *Client) error {
	for {
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

		wss.currentUser = m.User
		_, err = wss.jwtService.ValidateToken(m.JWT)
		if err != nil {
			return errors.New("invalid token")
		}

		if m.Type == "bootup" {
			wss.logger.WithFields(logrus.Fields{
				"client info": client,
				"username":    client.Username,
			}).Info("client successfully mapped")
		} else {
			wss.logger.WithField(
				"type", m.Type).Info("received message")
			c := m.Chat
			c.Timestamp = time.Now().Unix()
			c.From = m.User

			// // save in redis
			// id, err := redisrepo.CreateChat(&c)
			// if err != nil {
			// 	log.Println("error while saving chat in redis", err)
			// 	return
			// }

			if errCmd := wss.handleCommands(&c); errCmd != nil {
				return errCmd
			}
			wss.broadcastChannel <- &c
		}
	}
}

func (wss *webSocketServer) Broadcast() error {
	for {
		select {
		case message := <-wss.broadcastChannel:
			wss.logger.Info("new message to be broadcasted")

			if message.ToBePersisted {
				currentUser, err := wss.userService.GetByUsername(message.From)
				if err != nil {
					wss.logger.Errorf("User couldn't be retrieved from database, %s", err.Error())
					return err
				}
				if err := wss.messageService.Register(currentUser.ID, message.Msg, message.Timestamp); err != nil {
					wss.logger.Infof("Message was not saved, %s", err.Error())
				}
				wss.logger.Info("Message was saved successfully.")
			}

			for client := range wss.clients {
				wss.handleIfErrorExists(client, message)
			}

		case history := <-wss.historyChannel:
			wss.logger.Info("history messages ready to be send...")
			for client := range wss.clients {
				if client.Username == wss.currentUser {
					fmt.Printf("------------ %s\n", client.Username)
					wss.handleIfErrorExists(client, history)
				}
			}

		case errChannel := <-wss.errorChannel:
			errMsg := "Unauthorized. Token provided is not valid."
			wss.logger.Warn(errMsg)
			for client := range wss.clients {
				if client.Username == errChannel.CausedBy {
					wss.handleIfErrorExists(client, errChannel)
				}
			}
		}
	}
}

// define our WebSocket endpoint
func (wss *webSocketServer) Serve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wss.logger.WithFields(logrus.Fields{
			"host":                r.Host,
			"URL query variables": r.URL.Query(),
		}).Info("Listening WebSocket request")

		ws, err := wss.upgrader.Upgrade(w, r, nil)
		if err != nil {
			wss.logger.Error(err)
		}

		client := &Client{Conn: ws}

		wss.clients[client] = true
		wss.logger.WithFields(logrus.Fields{
			"quantity":       len(wss.clients),
			"remote address": ws.RemoteAddr(),
		}).Info("connected clients")

		var currentUser model.CurrentUser
	_:
		ws.ReadJSON(&currentUser)
		wss.currentUser = currentUser.Username
		wss.handleHistoryMessages()

		client.Username = currentUser.Username
		if err := wss.Read(client); err != nil {
			chatErr := model.ChatError{
				Msg:      err.Error(),
				Status:   http.StatusUnauthorized,
				CausedBy: client.Username,
			}
			wss.errorChannel <- &chatErr
			return
		}

		wss.handleCloseConnection(ws, client)
	}
}

func (wss *webSocketServer) Start() {
	// RabbitMQ
	wss.logger.Info("Staring websocket server...")
	go wss.Broadcast()
}

func (wss *webSocketServer) handleIfErrorExists(client *Client, data interface{}) {
	err := client.Conn.WriteJSON(data)
	if err != nil {
		wss.logger.Errorf("Websocket error: %s", err)
		client.Conn.Close()
		delete(wss.clients, client)
	}
}

func (wss *webSocketServer) handleCommands(c *model.Chat) error {
	wss.logger.Infof("command received: %s", c.Msg)
	if wss.botService.IsValidCommand(c.Msg) {
		botUsr, err := wss.botService.GetBotUser(wss.userService)
		if err != nil {
			msg := fmt.Sprintf("Bot user was not found, %s", err.Error())
			wss.logger.Info(msg)
			return errors.New(msg)
		}
		originalCmd := c.Msg
		c.From = botUsr.Username
		c.Msg = fmt.Sprintf("I'm working on %s. Please, hold on.", c.Msg)
		wss.broadcastChannel <- c

		channelResp := make(chan string)
		go wss.botService.ProcessCommandAsync(originalCmd, channelResp)

		result := <-channelResp
		c.Msg = result
		c.ToBePersisted = false
		return nil
	}

	c.ToBePersisted = true
	return nil
}

func (wss *webSocketServer) handleHistoryMessages() {
	maxHistory, err := strconv.Atoi(wss.maxHistoryMsg)
	if err != nil {
		chatErr := model.ChatError{
			Msg:    err.Error(),
			Status: http.StatusInternalServerError,
		}
		wss.errorChannel <- &chatErr
		return
	}

	lastMsgs := wss.messageService.GetLastMessages(maxHistory)
	wss.historyChannel <- lastMsgs
}

func (wss *webSocketServer) handleCloseConnection(ws *websocket.Conn, client *Client) {
	wss.logger.WithFields(logrus.Fields{
		"remote address": ws.RemoteAddr().String(),
	}).Info("closing connection")
	delete(wss.clients, client)
}
