-- +goose Up
-- +goose StatementBegin
-- Enable extension (if not done already)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Alter the users table
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4 ();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN id DROP DEFAULT;
-- +goose StatementEnd