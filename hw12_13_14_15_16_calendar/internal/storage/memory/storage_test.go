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
	seedEvent := func(userID string, startTime, endTime time.Time) storage.Event {
		return storage.Event{
			ID:          uuid.NewString(),
			UserID:      userID,
			Title:       "Title",
			Description: "Description",
			DateStart:   startTime,
			DateEnd:     endTime,
		}
	}

	ctx := context.Background()
	store := New()

	now := time.Now()
	user1 := "user1"
	user2 := "user2"

	e1 := seedEvent(user1, now, now.Add(time.Hour))
	e2 := seedEvent(user1, now.Add(2*time.Hour), now.Add(3*time.Hour))
	e3 := seedEvent(user2, now.Add(30*time.Minute), now.Add(time.Hour+30*time.Minute))

	require.NoError(t, store.CreateEvent(ctx, e1))
	require.NoError(t, store.CreateEvent(ctx, e2))
	require.NoError(t, store.CreateEvent(ctx, e3))

	// проверяем GetEvent
	got, err := store.GetEvent(ctx, e1.ID)
	require.NoError(t, err)
	require.Equal(t, e1, got)

	// обновляем событие
	e1.Title = "Updated"
	require.NoError(t, store.UpdateEvent(ctx, e1))
	got, err = store.GetEvent(ctx, e1.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated", got.Title)

	// проверяем ListEventsInRange для user1
	from := now.Add(-time.Hour)
	to := now.Add(4 * time.Hour)
	events, err := store.ListEventsInRange(ctx, user1, from, to)
	require.NoError(t, err)
	require.Len(t, events, 2) // e1 и e2

	// удаление
	require.NoError(t, store.DeleteEvent(ctx, e1))
	_, err = store.GetEvent(ctx, e1.ID)
	require.ErrorIs(t, err, storage.ErrNotFound)
}

func TestMemoryStorageConcurrencySafety(t *testing.T) {
	store := New()
	ctx := context.Background()

	const goroutines = 50
	const eventsPerGoroutine = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			userID := uuid.NewString()
			for j := 0; j < eventsPerGoroutine; j++ {
				ev := storage.Event{
					ID:        uuid.NewString(),
					UserID:    userID,
					Title:     "Event",
					DateStart: time.Now().Add(time.Duration(j) * time.Hour),
					DateEnd:   time.Now().Add(time.Duration(j+1) * time.Hour),
				}

				// Создаем событие
				require.NoError(t, store.CreateEvent(ctx, ev))

				// Сразу обновляем
				ev.Title = "Updated Event"
				require.NoError(t, store.UpdateEvent(ctx, ev))

				// Получаем событие
				got, err := store.GetEvent(ctx, ev.ID)
				require.NoError(t, err)
				require.Equal(t, "Updated Event", got.Title)
			}
		}()
	}

	wg.Wait()

	// Проверяем количество всех событий
	store.mu.RLock()
	defer store.mu.RUnlock()
	count := len(store.events)
	require.Equal(t, goroutines*eventsPerGoroutine, count)
}
