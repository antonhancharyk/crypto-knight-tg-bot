package broker

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("amqp dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close() //nolint:errcheck // rollback after failed channel
		return nil, fmt.Errorf("amqp channel: %w", err)
	}

	if err := ch.Qos(1, 0, false); err != nil {
		_ = ch.Close()   //nolint:errcheck // rollback after qos failure
		_ = conn.Close() //nolint:errcheck
		return nil, fmt.Errorf("amqp qos: %w", err)
	}

	return &Connection{conn: conn, channel: ch}, nil
}

func (c *Connection) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("queue declare %q: %w", name, err)
	}
	return nil
}

func (c *Connection) Channel() *amqp091.Channel {
	return c.channel
}

func (c *Connection) Close() {
	if c.channel != nil {
		_ = c.channel.Close() //nolint:errcheck // best-effort shutdown
	}
	if c.conn != nil {
		_ = c.conn.Close() //nolint:errcheck // best-effort shutdown
	}
}
