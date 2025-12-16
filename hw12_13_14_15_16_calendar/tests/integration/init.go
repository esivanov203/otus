package integration

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type EventBody struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"dateStart"`
	DateEnd     time.Time `json:"dateEnd"`
	UserID      string    `json:"userId"`
}

func httpInit() (string, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return "", err
	}

	host := os.Getenv("CALENDAR_HOST")
	port := "8080"
	url := fmt.Sprintf("http://%s:%s/events", host, port)

	return url, nil
}
