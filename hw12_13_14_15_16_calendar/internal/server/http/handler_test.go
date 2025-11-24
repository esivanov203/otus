package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/server"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateHandler(t *testing.T) {
	mockApp := new(server.MockApp)
	mockLogger := server.MockLogger{}
	server := &Server{app: mockApp, logger: mockLogger}

	event := model.Event{
		ID:          "1",
		Title:       "Title",
		Description: "Desc",
		UserID:      "user1",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Hour),
	}

	mockApp.On("CreateEvent", mock.Anything, mock.AnythingOfType("model.Event")).Return(nil)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.createHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	mockApp.AssertCalled(t, "CreateEvent", mock.Anything, mock.AnythingOfType("model.Event"))
}

func TestUpdateHandler(t *testing.T) {
	mockApp := new(server.MockApp)
	mockLogger := server.MockLogger{}
	server := &Server{app: mockApp, logger: mockLogger}

	event := model.Event{
		ID:          "1",
		Title:       "Updated",
		Description: "Updated Desc",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Hour),
	}

	mockApp.On("UpdateEvent", mock.Anything, mock.AnythingOfType("model.Event")).Return(nil)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPut, "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.updateHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusAccepted, resp.StatusCode)
	mockApp.AssertCalled(t, "UpdateEvent", mock.Anything, mock.AnythingOfType("model.Event"))
}

func TestDeleteHandler(t *testing.T) {
	mockApp := new(server.MockApp)
	mockLogger := server.MockLogger{}
	server := &Server{app: mockApp, logger: mockLogger}

	mockApp.On("DeleteEvent", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	server.deleteHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	mockApp.AssertCalled(t, "DeleteEvent", mock.Anything, "1")
}

func TestGetOneHandler(t *testing.T) {
	mockApp := new(server.MockApp)
	mockLogger := server.MockLogger{}
	server := &Server{app: mockApp, logger: mockLogger}

	event := model.Event{
		ID:          "1",
		Title:       "Title",
		Description: "Desc",
	}

	mockApp.On("GetEvent", mock.Anything, "1").Return(event, nil)

	req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	server.getOneHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	var gotEvent model.Event
	err := json.NewDecoder(resp.Body).Decode(&gotEvent)
	require.NoError(t, err)
	require.Equal(t, "Title", gotEvent.Title)
	mockApp.AssertCalled(t, "GetEvent", mock.Anything, "1")
}

func TestListEventsForDayHandler(t *testing.T) {
	mockApp := new(server.MockApp)
	mockLogger := server.MockLogger{}
	server := &Server{app: mockApp, logger: mockLogger}

	userID := "user1"
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	events := []model.Event{
		{
			ID:        "1",
			Title:     "E1",
			UserID:    userID,
			DateStart: startDate,
			DateEnd:   startDate.Add(time.Hour),
		},
	}

	mockApp.
		On("ListEventsForDay", mock.Anything, userID, startDate).
		Return(events, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/events?period=day&userId=user1&startDate=2025-01-01",
		nil,
	)
	w := httptest.NewRecorder()

	server.getPeriodHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got []model.Event
	err := json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "1", got[0].ID)

	mockApp.AssertCalled(t, "ListEventsForDay", mock.Anything, userID, startDate)
}
