-- +goose Up
-- A new random id: A UUID
-- created_at: A non-null timestamp
-- updated_at: A non null timestamp
-- body: A non-null string
-- user_id: This should reference the id of the user who created the chirp, and ON DELETE CASCADE, which will cause a user's chirps to be deleted if the user is deleted.
-- +goose StatementBegin
CREATE TABLE chirps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chirps;
-- +goose StatementEnd