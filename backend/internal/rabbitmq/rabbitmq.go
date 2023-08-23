package rabbitmq

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	// RabbitMQ defines the necessary methods to be implement
	RabbitMQ interface {
		Connect() (*amqp.Connection, error)
	}

	rabbitMQ struct {
		user     string
		password string
		host     string
		port     string
		logger   *logrus.Logger
	}
)

// NewRabbitMQ returns a new instance of RabbitMQ.
func NewRabbitMQ(user, password, host, port string, logger *logrus.Logger) RabbitMQ {
	return &rabbitMQ{
		user:     user,
		password: password,
		host:     host,
		port:     port,
		logger:   logger,
	}
}

func (r *rabbitMQ) buildURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", r.user, r.password, r.host, r.port)
}

// Connect establishes a connection to RabbitMQ.
func (r *rabbitMQ) Connect() (*amqp.Connection, error) {
	fmt.Println(r.buildURL())
	conn, err := amqp.Dial(r.buildURL())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
