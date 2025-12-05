package queue

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewQueueService(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{Connection: conn, Channel: ch}, nil
}

func (r *RabbitMQService) Close() error {
	err := r.Channel.Close()
	errC := r.Connection.Close()
	if errC != nil {
		err = fmt.Errorf("close rabbit connection: %w %w", errC, err)
	}

	return err
}

func (r *RabbitMQService) Publish(ctx context.Context, queue string, message Message) error {
	_, err := r.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)

	go func() {
		err := r.Channel.Publish("", queue, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        message.Body,
		})
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (r *RabbitMQService) Consume(ctx context.Context, queue string) (<-chan Message, error) {
	_, err := r.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := r.Channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	out := make(chan Message)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return // канал закрыт брокером
				}
				out <- Message{Body: msg.Body}
			}
		}
	}()

	return out, nil
}
