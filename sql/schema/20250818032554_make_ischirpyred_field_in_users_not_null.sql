-- +goose Up
-- +goose StatementBegin
-- migration to ensure that is_chirpy_red field is NOT NULL

-- Step 1: Update all existing rows to set a default value for NULL values.
UPDATE users SET is_chirpy_red = FALSE WHERE is_chirpy_red IS NULL;

-- Step 2: Alter the table to make the is_chirpy_red column NOT NULL.
ALTER TABLE users ALTER COLUMN is_chirpy_red SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN is_chirpy_red DROP NOT NULL;
-- +goose StatementEnd