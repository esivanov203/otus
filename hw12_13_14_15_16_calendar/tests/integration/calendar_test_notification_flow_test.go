package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNotificationFlow(t *testing.T) {
	urlBase, err := httpInit()
	require.NoError(t, err)

	// create event
	event := EventBody{
		Title:       "New Event",
		Description: "New Event Description",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Hour * 24),
		UserID:      "550e8400-e29b-41d4-a716-446655440099",
	}

	jsonEvent, err := json.Marshal(event)
	require.NoError(t, err)

	//nolint:noctx
	resp, err := http.Post(urlBase, "application/json", bytes.NewBuffer(jsonEvent))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	id := string(respBody)
	require.NotEmpty(t, id)

	path := os.Getenv("INTEGRATION_TESTS_LOG_PATH")
	if path == "" {
		t.Fatal("INTEGRATION_TESTS_LOG_PATH is not set")
	}

	require.Eventually(t, func() bool {
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			if line == id {
				return true
			}
		}

		return false
	}, 10*time.Second, 500*time.Millisecond)

	err = os.Remove(path)
	require.NoError(t, err)
}
