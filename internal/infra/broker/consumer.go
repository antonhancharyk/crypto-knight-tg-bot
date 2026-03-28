package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFunc func(msg []byte) error

type Consumer struct {
	ch      *amqp.Channel
	queue   string
	handler HandlerFunc
}

func NewConsumer(ch *amqp.Channel, queue string, handler HandlerFunc) *Consumer {
	return &Consumer{ch: ch, queue: queue, handler: handler}
}

func (c *Consumer) Run(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		c.queue,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("amqp consume %q: %w", c.queue, err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-msgs:
				if m.Body == nil {
					continue
				}

				if err := c.handler(m.Body); err != nil {
					_ = m.Nack(false, true) //nolint:errcheck // broker will redeliver
					continue
				}

				_ = m.Ack(false) //nolint:errcheck // ack after successful handler
			}
		}
	}()

	return nil
}
