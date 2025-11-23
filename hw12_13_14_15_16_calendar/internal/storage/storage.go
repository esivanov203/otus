package storage

import (
	"context"
	"time"
)

type Storage interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, event Event) error
	GetEvent(ctx context.Context, id string) (Event, error)
	ListEventsInRange(ctx context.Context, userID string, from, to time.Time) ([]Event, error)

	Connect() error
	Close() error
}
