-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ADD COLUMN noticed BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
    DROP COLUMN IF EXISTS noticed;
-- +goose StatementEnd
