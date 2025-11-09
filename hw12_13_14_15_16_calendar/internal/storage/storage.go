package storage

import "context"

type Storage interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, event Event) error
	GetEvent(ctx context.Context, id string) (Event, error)
	GetEventsList(ctx context.Context) ([]Event, error)
	GetEventsCount(ctx context.Context) (int, error)

	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}
