// package rabbitmq

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	"github.com/streadway/amqp"
// )

// const (
// 	QueueName = "package_status"
// )

// type rabbitmqClient struct {
// 	conn          *amqp.Connection
// 	ch            *amqp.Channel
// 	connString    string
// 	packageStatus <-chan amqp.Delivery
// }

// func NewRabbitMQClient(connectionString string) (*rabbitmqClient, error) {
// 	c := &rabbitmqClient{}
// 	var err error

// 	c.conn, err = amqp.Dial(connectionString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c.ch, err = c.conn.Channel()
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = c.configureQueue()

// 	return c, err
// }

// func (c *rabbitmqClient) ConsumeByVehicleID(ctx context.Context, vehicleID string) ([]byte, error) {
// 	for msg := range c.packageStatus {
// 		if msg.MessageId == vehicleID {
// 			_ = msg.Ack(false)
// 			return msg.Body, nil
// 		}
// 	}
// 	return nil, errors.New("err when getting package status on channel")
// }

// func (c *rabbitmqClient) Publish(message message.Message) {
// 	jsonStr := fmt.Sprintf(`{ "content": %s }`, message.Content)

// 	_ = c.ch.Publish("", message.Queue, true, false, amqp.Publishing{
// 		ContentType: "application/json",
// 		MessageId:   message.ID,
// 		Body:        []byte(jsonStr),
// 	})
// }

// func (c *rabbitmqClient) Close() {
// 	c.ch.Close()
// 	c.conn.Close()
// }

// func (c *rabbitmqClient) configureQueue() error {
// 	_, err := c.ch.QueueDeclare(
// 		QueueName,
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	c.packageStatus, err = c.ch.Consume(
// 		QueueName,
// 		"",
// 		false,
// 		false,
// 		false,
// 		true,
// 		nil,
// 	)
// 	return err
// }
