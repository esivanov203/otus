package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCrudEvent(t *testing.T) {
	urlBase, err := httpInit()
	require.NoError(t, err)

	// create
	b := EventBody{
		Title:       "New Event",
		Description: "New Event Description",
		DateStart:   time.Now().Add(time.Hour * 24),
		DateEnd:     time.Now().Add(time.Hour * 24 * 7),
		UserID:      "u1-01",
	}

	bj, err := json.Marshal(b)
	require.NoError(t, err)

	//nolint:noctx
	resp, err := http.Post(urlBase, "application/json", bytes.NewBuffer(bj))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	id := string(respBody)
	require.NotEmpty(t, id)

	client := &http.Client{}
	// update
	b.ID = id
	b.Title = "Updated Event"

	bju, err := json.Marshal(b)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, urlBase, bytes.NewBuffer(bju))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusAccepted, resp.StatusCode)

	respBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	resID := string(respBody)
	require.Equal(t, id, resID)

	// getOne
	//nolint:noctx
	resp, err = http.Get(urlBase + "/" + id)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBodyGetOne, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodyGetOne := EventBody{}
	err = json.Unmarshal(respBodyGetOne, &bodyGetOne)
	require.NoError(t, err)
	require.Equal(t, b.Title, bodyGetOne.Title)

	// get list — день/неделя/месяц
	periods := map[string]time.Time{
		"day":   time.Now().Add(24 * time.Hour),  // сегодня
		"week":  time.Now().Add(-24 * time.Hour), // вчера
		"month": time.Now().Add(-24 * time.Hour), // вчера
	}

	for period, startDate := range periods {
		t.Run(fmt.Sprintf("list events period=%s", period), func(t *testing.T) {
			u, err := url.Parse(urlBase)
			require.NoError(t, err)

			q := u.Query()
			q.Set("startDate", startDate.Format("2006-01-02"))
			q.Set("userId", b.UserID)
			q.Set("period", period)
			u.RawQuery = q.Encode()

			//nolint:noctx
			resp, err := http.Get(u.String())
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var bodyGet []EventBody
			err = json.NewDecoder(resp.Body).Decode(&bodyGet)
			require.NoError(t, err)

			found := false
			for _, e := range bodyGet {
				if e.ID == b.ID {
					found = true
					break
				}
			}
			require.True(t, found, "event not found in period=%s", period)
		})
	}

	// delete
	deleteURL := fmt.Sprintf("%s/%s", urlBase, id)
	req, err = http.NewRequestWithContext(context.Background(), http.MethodDelete, deleteURL, nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// негативные тесты - по одному на CRUD.
// полный набор всех случаев - в юнит-тестах.
func TestCrudEventNegative(t *testing.T) {
	urlBase, err := httpInit()
	require.NoError(t, err)

	client := &http.Client{}

	// --------------------------------------------------------
	// CREATE — дата окончания раньше даты начала
	// --------------------------------------------------------
	t.Run("create invalid dates", func(t *testing.T) {
		bad := EventBody{
			Title:       "Bad Event",
			Description: "Wrong dates",
			DateStart:   time.Now().Add(5 * time.Hour),
			DateEnd:     time.Now().Add(1 * time.Hour), // <--- end < start !!!
			UserID:      "u1-neg",
		}

		j, _ := json.Marshal(bad)

		//nolint:noctx
		resp, err := http.Post(urlBase, "application/json", bytes.NewBuffer(j))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode) // 422
	})

	// --------------------------------------------------------
	// UPDATE — несуществующий ID
	// --------------------------------------------------------
	t.Run("update non-existing event", func(t *testing.T) {
		update := EventBody{
			ID:          "550e8400-e29b-41d4-a716-446655440000",
			Title:       "Update",
			Description: "none",
			DateStart:   time.Now().Add(time.Hour),
			DateEnd:     time.Now().Add(2 * time.Hour),
			UserID:      "u1-neg",
		}

		j, _ := json.Marshal(update)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, urlBase, bytes.NewBuffer(j))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	// --------------------------------------------------------
	// GET — несуществующий ID
	// --------------------------------------------------------
	t.Run("get not exists id", func(t *testing.T) {
		//nolint:noctx
		resp, err := http.Get(fmt.Sprintf("%s/550e8400-e29b-41d4-a716-446655440000", urlBase))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	// --------------------------------------------------------
	// GET LIST — неверный period
	// --------------------------------------------------------
	t.Run("get list invalid period", func(t *testing.T) {
		u, _ := url.Parse(urlBase)
		q := u.Query()
		q.Set("startDate", time.Now().Format("2006-01-02"))
		q.Set("userId", "u1-neg")
		q.Set("period", "WRONG_PERIOD") // <--- неверный период
		u.RawQuery = q.Encode()

		//nolint:noctx
		resp, err := http.Get(u.String())
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	// --------------------------------------------------------
	// GET LIST — отсутствие обязательных параметров
	// --------------------------------------------------------
	t.Run("get list missing params", func(t *testing.T) {
		//nolint:noctx
		resp, err := http.Get(urlBase)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}
