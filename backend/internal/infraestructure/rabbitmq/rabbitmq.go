package rabbitmq

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	// RabbitMQ defines the structure to interact with RabbitMQ.
	RabbitMQ struct {
		Connection *amqp.Connection
		Channel    *amqp.Channel
		Logger     *logrus.Logger
	}

	// RabbitMQConfig holds the configuration values for connecting to RabbitMQ.
	RabbitMQConfig struct {
		User     string
		Password string
		Host     string
		Port     string
	}
)

// NewRabbitMQConfig returns a new instance of RabbitMQConfig.
func NewRabbitMQConfig(user, password, host, port string) RabbitMQConfig {
	return RabbitMQConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (rc RabbitMQConfig) buildURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", rc.User, rc.Password, rc.Host, rc.Port)
}

// NewRabbitMQ returns a new instance of RabbitMQ.
func NewRabbitMQ(config RabbitMQConfig, logger *logrus.Logger) (*RabbitMQ, error) {
	conn, err := amqp.Dial(config.buildURL())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close closes in a proper way a opened RabbitMQ connection.
func (rmq *RabbitMQ) Close() {
	if rmq.Channel != nil {
		rmq.Channel.Close()
	}

	if rmq.Connection != nil {
		rmq.Connection.Close()
	}
}

// Publish implements the needed functionality to send a message through a RabbitMQ channel.
func (rmq *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
	err := rmq.Channel.Publish(
		exchange,
		routingKey,
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

// Consume retrieves all the messages attached on a RabbitMQ channel.
func (rmq *RabbitMQ) Consume(queueName string, handler func([]byte)) error {
	msgs, err := rmq.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
		}
	}()

	return nil
}
