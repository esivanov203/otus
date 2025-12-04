package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/model"
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

func (s *Storage) CreateEvent(ctx context.Context, event model.Event) error {
	query := `
		INSERT INTO events 
		    (id, user_id, title, description, date_start, date_end, created_at, updated_at)
		VALUES (:id, :user_id, :title, :description, :date_start, :date_end, NOW(), NOW())
	`
	_, err := s.db.NamedExecContext(ctx, query, &event)

	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event model.Event) error {
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

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (model.Event, error) {
	var event model.Event
	query := `SELECT id, user_id, title, description, date_start, date_end FROM events WHERE id=$1`
	err := s.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Event{}, storage.ErrNotFound
		}
		return model.Event{}, err
	}

	return event, nil
}

func (s *Storage) ListEventsInRange(
	ctx context.Context,
	userID string,
	from, to time.Time,
) ([]model.Event, error) {
	var events []model.Event

	query := `
		SELECT id, user_id, title, description, date_start, date_end
		FROM events
		WHERE user_id=$1 AND date_start >= $2 AND date_start < $3
		ORDER BY date_start
	`

	err := s.db.SelectContext(ctx, &events, query, userID, from, to)
	return events, err
}

func (s *Storage) ListEventsTillNow(ctx context.Context) ([]model.Event, error) {
	var events []model.Event

	query := `
		SELECT id, user_id, title, description, date_start, date_end
		FROM events
		WHERE date_start <= $1 
			AND noticed = false
		ORDER BY date_start
		LIMIT 500
	`

	err := s.db.SelectContext(ctx, &events, query, time.Now())
	return events, err
}

func (s *Storage) UpdateNoticedEvent(ctx context.Context, id string) error {
	query := `
		UPDATE events
		SET noticed = true
		WHERE id = :id
	`
	event := model.Event{ID: id}
	res, err := s.db.NamedExecContext(ctx, query, &event)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return storage.ErrNotFound
	}

	return nil
}
