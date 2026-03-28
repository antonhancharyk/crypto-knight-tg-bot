package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// HandlerFunc processes a single message body; a non-nil error causes a negative ack and redelivery.
type HandlerFunc func(msg []byte) error

// Consumer subscribes to a queue and dispatches deliveries to a HandlerFunc.
type Consumer struct {
	ch      *amqp.Channel
	queue   string
	handler HandlerFunc
}

// NewConsumer builds a Consumer for queue on channel ch.
func NewConsumer(ch *amqp.Channel, queue string, handler HandlerFunc) *Consumer {
	return &Consumer{ch: ch, queue: queue, handler: handler}
}

// Run starts a goroutine that consumes messages until ctx is canceled.
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
