package message

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sacurio/jb-challenge/internal/infraestructure/rabbitmq"
)

type MessageService struct {
	RabbitMQ  *rabbitmq.RabbitMQ
	WebSocket *websocket.Conn
}

func NewMessageService(rabbitMQ *rabbitmq.RabbitMQ) *MessageService {
	return &MessageService{
		RabbitMQ: rabbitMQ,
	}
}

func (s *MessageService) PublishMessage(msg Message) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return s.RabbitMQ.Publish("", msg.Queue, messageBytes)
}

func (s *MessageService) ConsumeMessages(queueName string, handler func(msg Message) error) error {
	err := s.RabbitMQ.Channel.ExchangeDeclare(
		"excha1",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = s.RabbitMQ.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = s.RabbitMQ.Channel.QueueBind(
		queueName,
		"your-routing-key",
		"excha1",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return s.RabbitMQ.Consume(queueName, func(body []byte) {
		var message Message
		err := json.Unmarshal(body, &message)
		if err != nil {

			return
		}

		err = handler(message)
		if err != nil {

		}
		fmt.Println(message)
	})
}

func (s *MessageService) SendWebSocketMessage(message Message) error {
	if s.WebSocket != nil {
		return s.WebSocket.WriteJSON(message)
	}
	return nil
}
