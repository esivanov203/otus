package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type EventBody struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"dateStart"`
	DateEnd     time.Time `json:"dateEnd"`
	UserID      string    `json:"userId"`
}

func TestMain(m *testing.M) {
	start := time.Now()

	exitCode := m.Run()

	finish := time.Now()

	if err := clearDB(start, finish); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "clearDB failed:", err)
		exitCode = 1
	}

	os.Exit(exitCode)
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

func clearDB(from time.Time, to time.Time) error {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	tables := []string{"events"}

	for _, table := range tables {
		query := fmt.Sprintf(
			"DELETE FROM %s WHERE created_at BETWEEN $1 AND $2",
			table,
		)
		if _, err := db.ExecContext(context.Background(), query, from, to); err != nil {
			return err
		}
	}

	return nil
}
