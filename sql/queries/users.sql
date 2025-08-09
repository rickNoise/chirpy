-- name: CreateUser :one
-- id PK for users has a default UUID generated, so can leave out here
INSERT INTO
    users (
        created_at,
        updated_at,
        email,
        hashed_password
    )
VALUES (
        NOW(),
        NOW(),
        @email,
        @hashedPassword
    ) RETURNING *;

-- name: DeleteAllUsers :exec
-- deletes all users data in the users table
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT
    id,
    created_at,
    updated_at,
    email,
    hashed_password
FROM users
WHERE
    email = @useremail;