package broker

import (
	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Connection{conn: conn, channel: ch}, nil
}

func (c *Connection) Channel() *amqp091.Channel {
	return c.channel
}

func (c *Connection) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
