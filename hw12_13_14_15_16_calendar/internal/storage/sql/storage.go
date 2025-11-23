package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgresql driver for sqlx
)

type Storage struct {
	db  *sqlx.DB
	dsn string
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect() error {
	db, err := sqlx.Connect("postgres", s.dsn)
	if err != nil {
		return err
	}
	s.db = db

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	query := `
		INSERT INTO events 
		    (id, user_id, title, description, date_start, date_end, created_at, updated_at)
		VALUES (:id, :user_id, :title, :description, :date_start, :date_end, NOW(), NOW())
	`
	_, err := s.db.NamedExecContext(ctx, query, &event)

	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := `
		UPDATE events
		SET title=:title,
		    description=:description,
		    date_start=:date_start,
		    date_end=:date_end,
		    updated_at=NOW()
		WHERE id=:id
	`
	res, err := s.db.NamedExecContext(ctx, query, &event)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, event storage.Event) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id = $1", event.ID)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	var event storage.Event
	query := `SELECT id, user_id, title, description, date_start, date_end FROM events WHERE id=$1`
	err := s.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrNotFound
		}
		return storage.Event{}, err
	}

	return event, nil
}

func (s *Storage) ListEventsInRange(
	ctx context.Context,
	userID string,
	from, to time.Time,
) ([]storage.Event, error) {
	var events []storage.Event

	query := `
		SELECT id, user_id, title, description, date_start, date_end
		FROM events
		WHERE user_id=$1 AND date_start >= $2 AND date_start < $3
		ORDER BY date_start
	`

	err := s.db.SelectContext(ctx, &events, query, userID, from, to)
	return events, err
}
