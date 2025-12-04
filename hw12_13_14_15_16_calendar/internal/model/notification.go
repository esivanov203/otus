package model

import "time"

const QueueName = "calendar_notifications"

type Notification struct {
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"dateStart"`
	UserID string    `json:"userId"`
}

func NewNotificationFromEvent(event Event) Notification {
	return Notification{
		ID:     event.ID,
		Title:  event.Title,
		UserID: event.UserID,
		Date:   event.DateStart,
	}
}
