package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func clear(t *testing.T) {
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
	require.NoError(t, err)
	defer db.Close()

	// Очистка базы перед тестом
	tables := []string{
		"events",
		// если будут другие таблицы
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s", table)
		_, err := db.ExecContext(context.Background(), query)
		require.NoError(t, err)
	}
}
