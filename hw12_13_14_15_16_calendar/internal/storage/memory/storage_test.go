package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	seedEvent := func(userId string, startTime, endTime time.Time) storage.Event {
		return storage.Event{
			ID:          uuid.NewString(),
			Title:       "Test title",
			Description: "Test description",
			DateStart:   startTime,
			DateEnd:     endTime,
			UserID:      userId,
		}
	}

	ctx := context.Background()
	store := New()

	start := time.Now()
	end := start.Add(time.Hour)

	// create
	event := seedEvent("user1", start, end)
	err := store.CreateEvent(ctx, event)
	require.NoError(t, err)

	// create not fault with date for another user
	eventU2 := seedEvent("user2", start.Add(30*time.Minute), end.Add(30*time.Minute))
	err = store.CreateEvent(ctx, eventU2)
	require.NoError(t, err)

	// getOne
	getEvent, err := store.GetEvent(ctx, event.ID)
	require.NoError(t, err)
	require.Equal(t, event, getEvent)

	// update
	event.Title = "Updated title"
	err = store.UpdateEvent(ctx, event)
	require.NoError(t, err)

	getEvent, err = store.GetEvent(ctx, event.ID)
	require.NoError(t, err)
	require.Equal(t, event, getEvent)

	// get list
	list, err := store.GetEventsList(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)

	// delete
	err = store.DeleteEvent(ctx, event)
	require.NoError(t, err)

	// count
	count, err := store.GetEventsCount(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestStorageConcurrencySafety(t *testing.T) {
	store := New()
	ctx := context.Background()

	const goroutines = 50
	const eventsPerGoroutine = 20

	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			userID := uuid.New().String()
			for j := 0; j < eventsPerGoroutine; j++ {
				ev := storage.Event{
					ID:        uuid.NewString(),
					UserID:    userID,
					Title:     "Event",
					DateStart: time.Now().Add(time.Duration(j) * time.Hour),
					DateEnd:   time.Now().Add(time.Duration(j+1) * time.Hour),
				}
				err := store.CreateEvent(ctx, ev)
				require.NoError(t, err)

				// попробуем сразу обновить
				ev.Title = "Updated Event"
				err = store.UpdateEvent(ctx, ev)
				require.NoError(t, err)

				// прочитаем обратно
				got, err := store.GetEvent(ctx, ev.ID)
				require.NoError(t, err)
				require.Equal(t, "Updated Event", got.Title)
			}
		}()
	}

	wg.Wait()

	count, err := store.GetEventsCount(ctx)
	require.NoError(t, err)
	require.Equal(t, goroutines*eventsPerGoroutine, count)
}
