package app

import (
	"context"
	"fmt"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
	"sync"
	"testing"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestApp(store storage.Storage) *App {
	logg, _ := logger.New(logger.Conf{Level: "debug", Type: "console"})
	return New(logg, store)
}

func TestApp_CreateEventValidation(t *testing.T) {
	store := memorystorage.New()
	app := newTestApp(store)
	ctx := context.Background()

	// id/userID пустые
	err := app.CreateEvent(ctx, model.Event{
		ID:          "",
		UserID:      "",
		Title:       "Title",
		Description: "Desc",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Hour),
	})
	require.Error(t, err)

	// title пустой
	err = app.CreateEvent(ctx, model.Event{
		ID:          "id",
		UserID:      "user",
		Title:       "",
		Description: "Desc",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Hour),
	})
	require.Error(t, err)

	// start после end
	err = app.CreateEvent(ctx, model.Event{
		ID:          "id",
		UserID:      "user",
		Title:       "Title",
		Description: "Desc",
		DateStart:   time.Now().Add(time.Hour),
		DateEnd:     time.Now(),
	})
	require.Error(t, err)
}

func TestApp_CreateUpdateDeleteEvent(t *testing.T) {
	store := memorystorage.New()
	app := newTestApp(store)
	ctx := context.Background()

	id := uuid.NewString()
	userID := "user1"
	start := time.Now()
	end := start.Add(time.Hour)

	// Create
	err := app.CreateEvent(ctx, model.Event{
		ID:          id,
		UserID:      userID,
		Title:       "Title",
		Description: "Desc",
		DateStart:   start,
		DateEnd:     end,
	})
	require.NoError(t, err)

	// Get
	ev, err := app.GetEvent(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "Title", ev.Title)
	require.Equal(t, "Desc", ev.Description)

	// Update
	err = app.UpdateEvent(ctx, model.Event{
		ID:          id,
		Title:       "Updated",
		Description: "Updated Desc",
		DateStart:   start,
		DateEnd:     end,
	})
	require.NoError(t, err)

	ev, err = app.GetEvent(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "Updated", ev.Title)

	// Delete
	err = app.DeleteEvent(ctx, id)
	require.NoError(t, err)

	_, err = app.GetEvent(ctx, id)
	require.ErrorIs(t, err, storage.ErrNotFound)
}

func TestApp_ListAndCountEvents(t *testing.T) {
	store := memorystorage.New()
	app := newTestApp(store)
	ctx := context.Background()

	// Создадим несколько событий
	for i := 0; i < 5; i++ {
		id := uuid.NewString()
		err := app.CreateEvent(ctx, model.Event{
			ID:          id,
			UserID:      "user1",
			Title:       "Title",
			Description: "Desc",
			DateStart:   time.Now().Add(time.Hour * time.Duration(i)),
			DateEnd:     time.Now().Add(time.Hour * time.Duration(i+1)),
		})
		require.NoError(t, err)
	}
}

func TestListEventsForDayWeekMonth(t *testing.T) {
	store := memorystorage.New()
	app := newTestApp(store)
	ctx := context.Background()

	userID := "user1"

	baseDate := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	// Создаём события:
	events := []model.Event{
		{
			ID:          uuid.NewString(),
			UserID:      userID,
			Title:       "Day event",
			Description: "event for same day",
			DateStart:   baseDate,
			DateEnd:     baseDate.Add(time.Hour),
		},
		{
			ID:          uuid.NewString(),
			UserID:      userID,
			Title:       "Same week event",
			Description: "event for same week",
			DateStart:   baseDate.AddDate(0, 0, 2), // +2 дня
			DateEnd:     baseDate.AddDate(0, 0, 2).Add(time.Hour),
		},
		{
			ID:          uuid.NewString(),
			UserID:      userID,
			Title:       "Same month event",
			Description: "event for same month",
			DateStart:   baseDate.AddDate(0, 0, 10), // +10 дней
			DateEnd:     baseDate.AddDate(0, 0, 10).Add(time.Hour),
		},
	}

	for _, e := range events {
		require.NoError(t, app.CreateEvent(ctx, e))
	}

	listDay, err := app.ListEventsForDay(ctx, userID, baseDate)
	require.NoError(t, err)
	require.Len(t, listDay, 1)
	require.Equal(t, "Day event", listDay[0].Title)

	listWeek, err := app.ListEventsForWeek(ctx, userID, baseDate)
	require.NoError(t, err)
	require.Len(t, listWeek, 2) // day + week
	titlesWeek := []string{listWeek[0].Title, listWeek[1].Title}
	require.Contains(t, titlesWeek, "Day event")
	require.Contains(t, titlesWeek, "Same week event")

	listMonth, err := app.ListEventsForMonth(ctx, userID, baseDate)
	require.NoError(t, err)
	require.Len(t, listMonth, 3)
	titlesMonth := []string{
		listMonth[0].Title,
		listMonth[1].Title,
		listMonth[2].Title,
	}
	require.Contains(t, titlesMonth, "Day event")
	require.Contains(t, titlesMonth, "Same week event")
	require.Contains(t, titlesMonth, "Same month event")
}

func TestApp_ConcurrencySafety(t *testing.T) {
	store := memorystorage.New()
	app := newTestApp(store)
	ctx := context.Background()

	const goroutines = 20
	const eventsPerGoroutine = 10

	errCh := make(chan error, goroutines*eventsPerGoroutine)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				id := uuid.NewString()
				start := time.Now().Add(time.Minute * time.Duration(j))
				end := start.Add(time.Minute * 30)

				err := app.CreateEvent(ctx, model.Event{
					ID:          id,
					UserID:      fmt.Sprintf("user%d", i),
					Title:       "Title",
					Description: "Desc",
					DateStart:   start,
					DateEnd:     end,
				})
				if err != nil {
					errCh <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err)
	}

	total := 0
	now := time.Now()
	for i := 0; i < goroutines; i++ {
		userID := fmt.Sprintf("user%d", i)
		events, err := app.ListEventsForMonth(ctx, userID, now)
		require.NoError(t, err)
		total += len(events)
	}

	require.Equal(t, goroutines*eventsPerGoroutine, total)
}
