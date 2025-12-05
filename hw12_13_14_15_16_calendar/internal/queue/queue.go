package queue

import (
	"context"
)

type CalendarQueueConfig struct {
	Dsn      string `yaml:"dsn"`
	Interval string `yaml:"interval"`
}

type Message struct {
	Body []byte
}

type CalendarQueueService interface {
	Publish(ctx context.Context, queue string, msg Message) error
	Consume(ctx context.Context, queue string) (<-chan Message, error)
	Close() error
}
