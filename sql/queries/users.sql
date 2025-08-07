-- name: CreateUser :one
-- id PK for users has a default UUID generated, so can leave out here
INSERT INTO
    users (created_at, updated_at, email)
VALUES (NOW(), NOW(), @email) RETURNING *;

-- name: DeleteAllUsers :exec
-- deletes all users data in the users table
DELETE FROM users;