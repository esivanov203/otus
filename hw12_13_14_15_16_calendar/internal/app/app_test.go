package app

import (
	"context"
	"fmt"
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
	err := app.CreateEvent(ctx, EventDTO{
		ID:          "",
		UserID:      "",
		Title:       "Title",
		Description: "Desc",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Hour),
	})
	require.Error(t, err)

	// title пустой
	err = app.CreateEvent(ctx, EventDTO{
		ID:          "id",
		UserID:      "user",
		Title:       "",
		Description: "Desc",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Hour),
	})
	require.Error(t, err)

	// start после end
	err = app.CreateEvent(ctx, EventDTO{
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
	err := app.CreateEvent(ctx, EventDTO{
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
	err = app.UpdateEvent(ctx, EventDTO{
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
		err := app.CreateEvent(ctx, EventDTO{
			ID:          id,
			UserID:      "user1",
			Title:       "Title",
			Description: "Desc",
			DateStart:   time.Now().Add(time.Hour * time.Duration(i)),
			DateEnd:     time.Now().Add(time.Hour * time.Duration(i+1)),
		})
		require.NoError(t, err)
	}

	list, err := app.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, list, 5)

	count, err := app.CountEvents(ctx)
	require.NoError(t, err)
	require.Equal(t, 5, count)
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

				err := app.CreateEvent(ctx, EventDTO{
					ID:          id,
					UserID:      "user" + fmt.Sprint(i),
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

	count, err := app.CountEvents(ctx)
	require.NoError(t, err)
	require.Equal(t, goroutines*eventsPerGoroutine, count)
}
