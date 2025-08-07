-- +goose Up
-- +goose StatementBegin
-- id: a UUID that will serve as the primary key
-- created_at: a TIMESTAMP that can not be null
-- updated_at: a TIMESTAMP that can not be null
-- email: TEXT that can not be null and must be unique
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd