-- +goose Up
-- +goose StatementBegin
-- Adds a non-null TEXT column to the users table called hashed_password. It should default to "unset" for existing users.
ALTER TABLE users
ADD COLUMN hashed_password TEXT NOT NULL DEFAULT 'unset';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN hashed_password;
-- +goose StatementEnd