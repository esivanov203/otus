package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	"time"
)

// EventDTO — входная структура для создания/обновления события
type EventDTO struct {
	ID          string
	UserID      string
	Title       string
	Description string
	DateStart   time.Time
	DateEnd     time.Time
}

// Application — интерфейс бизнес-логики
type Application interface {
	CreateEvent(ctx context.Context, dto EventDTO) error
	UpdateEvent(ctx context.Context, dto EventDTO) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	ListEvents(ctx context.Context) ([]storage.Event, error)
	CountEvents(ctx context.Context) (int, error)
}

// App — реализация Application
type App struct {
	storage storage.Storage
	logger  logger.Logger
}

func New(logg logger.Logger, storage storage.Storage) *App {
	return &App{storage: storage, logger: logg}
}

// CreateEvent создаёт новое событие
func (a *App) CreateEvent(ctx context.Context, dto EventDTO) error {
	if dto.ID == "" || dto.UserID == "" {
		return errors.New("id and userID are required")
	}
	if dto.Title == "" {
		return errors.New("title is required")
	}
	if dto.DateStart.After(dto.DateEnd) {
		return fmt.Errorf("start date must be before end date")
	}

	event := storage.Event{
		ID:          dto.ID,
		UserID:      dto.UserID,
		Title:       dto.Title,
		Description: dto.Description,
		DateStart:   dto.DateStart,
		DateEnd:     dto.DateEnd,
	}

	return a.storage.CreateEvent(ctx, event)
}

// UpdateEvent обновляет событие
func (a *App) UpdateEvent(ctx context.Context, dto EventDTO) error {
	if dto.ID == "" {
		return errors.New("id is required")
	}
	if dto.DateStart.After(dto.DateEnd) {
		return fmt.Errorf("start date must be before end date")
	}

	event := storage.Event{
		ID:          dto.ID,
		Title:       dto.Title,
		Description: dto.Description,
		DateStart:   dto.DateStart,
		DateEnd:     dto.DateEnd,
	}

	return a.storage.UpdateEvent(ctx, event)
}

// DeleteEvent удаляет событие по ID
func (a *App) DeleteEvent(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	event := storage.Event{ID: id}
	return a.storage.DeleteEvent(ctx, event)
}

// GetEvent получает событие по ID
func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	if id == "" {
		return storage.Event{}, errors.New("id is required")
	}
	return a.storage.GetEvent(ctx, id)
}

// ListEvents возвращает все события
func (a *App) ListEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetEventsList(ctx)
}

// CountEvents возвращает количество событий
func (a *App) CountEvents(ctx context.Context) (int, error) {
	return a.storage.GetEventsCount(ctx)
}
