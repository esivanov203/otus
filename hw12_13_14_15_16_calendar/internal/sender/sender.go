package sender

import (
	"context"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/queue"
)

func Run(ctx context.Context, logger logger.Logger, q queue.CalendarQueueService) {
	msgs, err := q.Consume(ctx, model.QueueName)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	for msg := range msgs {
		logger.Info(string(msg.Body))
	}
}
