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
SELECT * FROM users WHERE email = @useremail;

-- name: UpdateEmailAndPasswordByUserId :one
-- updates a user record with a new hashed password and email address
UPDATE users
SET
    updated_at = NOW(),
    email = @newEmail,
    hashed_password = @newHashedPassword
WHERE
    id = @userId RETURNING *;

-- name: UpgradeUserToChirpyRedById :one
-- upgrades a user to chirpy red based on their ID by modifying the is_chirpy_field to true.
UPDATE users
SET
    updated_at = NOW(),
    is_chirpy_red = TRUE
WHERE
    id = @userid RETURNING *;