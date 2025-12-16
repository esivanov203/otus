package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/email_service"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/queue"
)

func Run(ctx context.Context, logg logger.Logger, q queue.CalendarQueueService, emlS email_service.EmailSender) {
	msgs, err := q.Consume(ctx, model.QueueName)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	for msg := range msgs {
		var n model.Notification
		err := json.Unmarshal(msg.Body, &n)
		if err != nil {
			logg.Error(err.Error())
			continue
		}
		eml := email_service.Email{
			UserId:  n.UserID,
			Subject: "Уведомление о событии: " + n.Title,
			Body:    fmt.Sprintf("Событие: %s\nНачало:%s", n.Title, n.Date),
			EventID: n.ID,
		}
		err = emlS.Send(ctx, eml)
		if err != nil {
			logg.Error(err.Error())
			continue
		}

		logg.Info("message has been received", logger.Fields{"id": n.ID})
	}
}
