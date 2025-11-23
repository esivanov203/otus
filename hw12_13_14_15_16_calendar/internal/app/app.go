package app

import (
	"context"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
)

type Application interface {
	CreateEvent(ctx context.Context, event model.Event) error
	UpdateEvent(ctx context.Context, event model.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (model.Event, error)

	ListEventsForDay(ctx context.Context, userID string, date time.Time) ([]model.Event, error)
	ListEventsForWeek(ctx context.Context, userID string, weekStart time.Time) ([]model.Event, error)
	ListEventsForMonth(ctx context.Context, userID string, monthStart time.Time) ([]model.Event, error)
}

type App struct {
	storage storage.Storage
	logger  logger.Logger
}

func New(logg logger.Logger, storage storage.Storage) *App {
	return &App{storage: storage, logger: logg}
}

func (a *App) CreateEvent(ctx context.Context, event model.Event) error {
	if err := event.ValidateCreate(); err != nil {
		return err
	}

	err := a.storage.CreateEvent(ctx, event)
	if err != nil {
		a.logger.Error(err.Error(), logger.Fields{"app": "create event", "title": event.Title})
		return err
	}

	a.logger.Info("success created", logger.Fields{"app": "create event", "title": event.Title})
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event model.Event) error {
	if err := event.ValidateUpdate(); err != nil {
		return err
	}

	err := a.storage.UpdateEvent(ctx, event)
	if err != nil {
		a.logger.Error(err.Error(), logger.Fields{"app": "update event", "id": event.ID})
		return err
	}

	a.logger.Info("success updated", logger.Fields{"app": "update event", "id": event.ID})
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	event := model.Event{ID: id}
	if err := event.ValidateOne(); err != nil {
		return err
	}

	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEvent(ctx context.Context, id string) (model.Event, error) {
	event := model.Event{ID: id}
	if err := event.ValidateOne(); err != nil {
		return event, err
	}

	return a.storage.GetEvent(ctx, id)
}

func (a *App) ListEventsForDay(ctx context.Context, userID string, date time.Time) ([]model.Event, error) {
	event := model.Event{
		UserID:    userID,
		DateStart: date,
	}
	if err := event.ValidateList(); err != nil {
		return nil, err
	}
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	return a.storage.ListEventsInRange(ctx, userID, start, end)
}

func (a *App) ListEventsForWeek(ctx context.Context, userID string, weekStart time.Time) ([]model.Event, error) {
	event := model.Event{
		UserID:    userID,
		DateStart: weekStart,
	}
	if err := event.ValidateList(); err != nil {
		return nil, err
	}

	start := weekStart.Truncate(24 * time.Hour)
	end := start.Add(7 * 24 * time.Hour)
	return a.storage.ListEventsInRange(ctx, userID, start, end)
}

func (a *App) ListEventsForMonth(ctx context.Context, userID string, monthStart time.Time) ([]model.Event, error) {
	event := model.Event{
		UserID:    userID,
		DateStart: monthStart,
	}
	if err := event.ValidateList(); err != nil {
		return nil, err
	}

	start := monthStart.Truncate(24 * time.Hour)
	end := start.AddDate(0, 1, 0)
	return a.storage.ListEventsInRange(ctx, userID, start, end)
}
