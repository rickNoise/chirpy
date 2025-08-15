-- +goose Up
-- +goose StatementBegin
-- Create a new database table with up/down migrations called refresh_tokens.
-- token: the primary key - it's just a string
-- created_at
-- updated_at
-- user_id: foreign key that deletes the row if the user is deleted
-- expires_at: the timestamp when the token expires
-- revoked_at: the timestamp when the token was revoked (null if not revoked)
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd