package memorystorage

import (
	"context"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"sync"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]model.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]model.Event),
	}
}

func (s *Storage) Connect() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event model.Event) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event model.Event) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrNotFound
	}
	event.UserID = s.events[event.ID].UserID
	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (model.Event, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.events[id]
	if !ok {
		return model.Event{}, storage.ErrNotFound
	}

	return e, nil
}

func (s *Storage) ListEventsInRange(
	ctx context.Context,
	userID string,
	from, to time.Time,
) ([]model.Event, error) {

	_ = ctx

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Event, 0)

	for _, e := range s.events {
		if e.UserID != userID {
			continue
		}

		if !e.DateStart.Before(to) {
			continue
		}
		if e.DateStart.Before(from) {
			continue
		}

		result = append(result, e)
	}

	return result, nil
}
