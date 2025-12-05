package storage

import (
	"context"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
)

type Storage interface {
	CreateEvent(ctx context.Context, event model.Event) error
	UpdateEvent(ctx context.Context, event model.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (model.Event, error)
	ListEventsInRange(ctx context.Context, userID string, from, to time.Time) ([]model.Event, error)

	ListEventsTillNow(ctx context.Context) ([]model.Event, error)
	UpdateNoticedEvent(ctx context.Context, id string) error

	Connect() error
	Close() error
}
