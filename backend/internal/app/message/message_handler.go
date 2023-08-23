package message

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sacurio/jb-challenge/internal/infraestructure/rabbitmq"
	"github.com/sacurio/jb-challenge/internal/util"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Post...
func PostMessage(messageService *MessageService, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message Message

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			logger.Errorf("error decoding message: %s", err.Error())
			util.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = messageService.PublishMessage(message)
		if err != nil {
			logger.Errorf("error posting message: %s", err.Error())
			util.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		data := util.MapResponse{
			"status": "ok",
		}

		dataJSON := util.CustomMarshall(w, data, logger)
		util.SendJSONResponse(w, http.StatusCreated, dataJSON)
	}
}

var wsClients map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)

// ConsumeMessages
func ConsumeMessages(messageService *MessageService, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
		}

		wsClients[ws] = true

		fmt.Println("Client connected")
		go func() {

			rabbitMQConfig := rabbitmq.NewRabbitMQConfig("rabbitmq_user", "&*dfgs33DaeaA!@", "localhost", "5672")
			rabbitMQHandler, err := rabbitmq.NewRabbitMQ(rabbitMQConfig, logger)
			if err != nil {
				logger.Panicf("Error on connecting to RabbitMQ service: %s", err.Error())
			}

			defer rabbitMQHandler.Close()

			ms := NewMessageService(rabbitMQHandler)
			for {
				writer("list", ms.RabbitMQ.Channel, ws)
			}

		}()

		reader(ws)
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func writer(queueName string, ch *amqp.Channel, conn *websocket.Conn) {
	ch.ExchangeDeclare("chatbot", "fanout", true, false, true, false, nil)
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("Error al consumir mensajes: %s", err)
	}

	fmt.Println("Esperando mensajes...")

	for msg := range msgs {
		for conn := range wsClients {
			conn.WriteMessage(websocket.TextMessage, msg.Body)
		}
	}
}
