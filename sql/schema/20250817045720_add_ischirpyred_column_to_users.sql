-- +goose Up
-- +goose StatementBegin
-- Include a new column on the users table called is_chirpy_red. This column should be a boolean, and it should default to false.
ALTER TABLE users ADD COLUMN is_chirpy_red BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_chirpy_red;
-- +goose StatementEnd