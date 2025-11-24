package server

import (
	"context"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct{}

func (MockLogger) Debug(msg string, fields ...logger.Fields) { _ = msg; _ = fields }
func (MockLogger) Info(msg string, fields ...logger.Fields)  { _ = msg; _ = fields }
func (MockLogger) Warn(msg string, fields ...logger.Fields)  { _ = msg; _ = fields }
func (MockLogger) Error(msg string, fields ...logger.Fields) { _ = msg; _ = fields }

type MockApp struct {
	mock.Mock
}

func (m *MockApp) CreateEvent(ctx context.Context, event model.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockApp) UpdateEvent(ctx context.Context, event model.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockApp) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockApp) GetEvent(ctx context.Context, id string) (model.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Event), args.Error(1)
}

func (m *MockApp) ListEventsForDay(ctx context.Context, userID string, date time.Time) ([]model.Event, error) {
	args := m.Called(ctx, userID, date)
	return args.Get(0).([]model.Event), args.Error(1)
}

func (m *MockApp) ListEventsForWeek(ctx context.Context, userID string, weekStart time.Time) ([]model.Event, error) {
	args := m.Called(ctx, userID, weekStart)
	return args.Get(0).([]model.Event), args.Error(1)
}

func (m *MockApp) ListEventsForMonth(ctx context.Context, userID string, monthStart time.Time) ([]model.Event, error) {
	args := m.Called(ctx, userID, monthStart)
	return args.Get(0).([]model.Event), args.Error(1)
}
