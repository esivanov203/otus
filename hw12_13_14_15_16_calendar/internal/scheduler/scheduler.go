package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
)

func Run(
	ctx context.Context,
	logger logger.Logger,
	storage storage.Storage,
	q queue.CalendarQueueService,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	for {
		select {
		case <-ctx.Done():
			logger.Error(ctx.Err().Error())
			return

		case <-ticker.C:
			events, err := storage.ListEventsTillNow(ctx)
			if err != nil {
				logger.Error(ctx.Err().Error())
				return
			}
			for _, event := range events {
				n := model.NewNotificationFromEvent(event)
				data, err := json.Marshal(n)
				if err != nil {
					logger.Error(ctx.Err().Error())
					continue
				}
				err = q.Publish(ctx, model.QueueName, queue.Message{
					Body: data,
				})
				if err != nil {
					logger.Error(ctx.Err().Error())
					continue
				}
				logger.Info("published event: " + event.ID)

				if err := storage.UpdateNoticedEvent(ctx, event.ID); err != nil {
					logger.Error(ctx.Err().Error())
					continue
				}
				if event.DateStart.Before(oneYearAgo) {
					if err := storage.DeleteEvent(ctx, event.ID); err != nil {
						logger.Error(fmt.Sprintf("failed to delete event %s: %v", event.ID, err))
					}
				}
			}
		}
	}
}
