package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/server"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	mockApp := new(server.MockApp)
	serv := &GRPCServer{app: mockApp, logger: server.MockLogger{}}

	event := model.Event{
		ID:        "1",
		Title:     "Title",
		UserID:    "u1",
		DateStart: time.Now(),
		DateEnd:   time.Now().Add(time.Hour),
	}

	mockApp.On("CreateEvent", mock.Anything, mock.AnythingOfType("model.Event")).Return(nil)

	req := &proto.CreateEventRequest{
		Event: &proto.Event{
			Id:        event.ID,
			Title:     event.Title,
			UserId:    event.UserID,
			DateStart: event.DateStart.Unix(),
			DateEnd:   event.DateEnd.Unix(),
		},
	}

	resp, err := serv.CreateEvent(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, event.ID, resp.Event.Id)
	mockApp.AssertCalled(t, "CreateEvent", mock.Anything, mock.AnythingOfType("model.Event"))
}

func TestGetEvent(t *testing.T) {
	mockApp := new(server.MockApp)
	serv := &GRPCServer{app: mockApp, logger: server.MockLogger{}}

	event := model.Event{
		ID:    "1",
		Title: "Test",
	}

	mockApp.On("GetEvent", mock.Anything, "1").Return(event, nil)

	req := &proto.EventRequest{Id: "1"}
	resp, err := serv.GetEvent(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "1", resp.Event.Id)
	mockApp.AssertCalled(t, "GetEvent", mock.Anything, "1")
}

func TestDeleteEvent(t *testing.T) {
	mockApp := new(server.MockApp)
	serv := &GRPCServer{app: mockApp, logger: server.MockLogger{}}

	mockApp.On("DeleteEvent", mock.Anything, "1").Return(nil)

	req := &proto.DeleteEventRequest{Id: "1"}
	resp, err := serv.DeleteEvent(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	mockApp.AssertCalled(t, "DeleteEvent", mock.Anything, "1")
}

func TestListEventsPeriod(t *testing.T) {
	mockApp := new(server.MockApp)
	serv := &GRPCServer{app: mockApp, logger: server.MockLogger{}}

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	userID := "u1"

	events := []model.Event{
		{ID: "1", Title: "E1", UserID: userID, DateStart: startDate, DateEnd: startDate.Add(time.Hour)},
	}

	mockApp.On("ListEventsForDay", mock.Anything, userID, mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == 2025 && t.Month() == 1 && t.Day() == 1
	})).Return(events, nil)

	req := &proto.ListRequest{
		UserId: userID,
		Date:   startDate.Unix(),
		Period: "day",
	}

	resp, err := serv.ListEventsPeriod(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, resp.Event, 1)
	require.Equal(t, "1", resp.Event[0].Id)
	mockApp.AssertCalled(t, "ListEventsForDay", mock.Anything, userID, mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == 2025 && t.Month() == 1 && t.Day() == 1
	}))
}
