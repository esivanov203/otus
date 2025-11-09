-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
     id UUID PRIMARY KEY,
     title TEXT NOT NULL,
     description TEXT,
     date_start TIMESTAMP NOT NULL,
     date_end TIMESTAMP NOT NULL,
     user_id TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS ids_events_date_start ON events(date_start);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
