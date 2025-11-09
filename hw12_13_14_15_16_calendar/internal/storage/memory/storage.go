package memorystorage

import (
	"context"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type Storage struct {
	mu     sync.RWMutex //nolint:unused
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// проверка на пересечение дат
	for _, e := range s.events {
		if e.UserID == event.UserID &&
			event.DateStart.Before(e.DateEnd) && e.DateStart.Before(event.DateEnd) {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrNotFound
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrNotFound
	}

	delete(s.events, event.ID)

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.events[id]
	if !ok {
		return storage.Event{}, storage.ErrNotFound
	}

	return e, nil
}

func (s *Storage) GetEventsList(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]storage.Event, 0, len(s.events))
	for _, e := range s.events {
		list = append(list, e)
	}

	return list, nil
}

func (s *Storage) GetEventsCount(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.events), nil
}
