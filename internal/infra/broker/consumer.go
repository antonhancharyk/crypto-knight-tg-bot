package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch      *amqp.Channel
	queue   string
	handler func(msg []byte)
}

func NewConsumer(ch *amqp.Channel, queue string, handler func(msg []byte)) (*Consumer, error) {
	_, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{ch: ch, queue: queue, handler: handler}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		c.queue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-msgs:
				c.handler(m.Body)
			}
		}
	}()

	return nil
}
