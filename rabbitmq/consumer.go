// rabbitmq/consumer.go
package rabbitmq

import (
	"log"
)

// Consumer represents a RabbitMQ message consumer
type Consumer struct {
	conn    *Connection
	handler func([]byte) error
}

// NewConsumer creates a new RabbitMQ message consumer
func NewConsumer(conn *Connection, handler func([]byte) error) *Consumer {
	return &Consumer{
		conn:    conn,
		handler: handler,
	}
}

// StartConsuming starts consuming messages from a queue
func (c *Consumer) StartConsuming(queueName string) error {
	msgs, err := c.conn.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message from queue '%s' (Routing Key: %s)", queueName, d.RoutingKey)

			if err := c.handler(d.Body); err != nil {
				log.Printf("Error processing message: %v. Re-queueing.", err)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	log.Printf("Started consuming messages from queue: %s", queueName)
	return nil
}
