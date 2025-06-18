// rabbitmq/connection.go
package rabbitmq

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection holds the RabbitMQ connection and channel
type Connection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewConnection creates a new RabbitMQ connection
func NewConnection(defaultURL string) (*Connection, error) {
	// Read the RabbitMQ connection URL from an environment variable.
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = defaultURL
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("Successfully connected to RabbitMQ and opened a channel!")

	return &Connection{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close closes the RabbitMQ connection and channel
func (c *Connection) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}

// DeclareExchange declares a RabbitMQ exchange
func (c *Connection) DeclareExchange(name, exchangeType string) error {
	err := c.Channel.ExchangeDeclare(
		name,         // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	log.Printf("Declared exchange: %s (%s, durable)", name, exchangeType)
	return nil
}

// DeclareQueue declares a RabbitMQ queue
func (c *Connection) DeclareQueue(name string) (amqp.Queue, error) {
	queue, err := c.Channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	log.Printf("Declared queue: %s", name)
	return queue, nil
}

// BindQueue binds a queue to an exchange with a routing key
func (c *Connection) BindQueue(queueName, routingKey, exchangeName string) error {
	err := c.Channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	log.Printf("Bound queue '%s' to exchange '%s' with routing key '%s'", queueName, exchangeName, routingKey)
	return nil
}
