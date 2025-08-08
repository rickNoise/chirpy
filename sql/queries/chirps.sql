-- name: CreateChirp :one
-- creates a new chirp in the db tied to the creating user
INSERT INTO
    chirps (
        created_at,
        updated_at,
        body,
        user_id
    )
VALUES (NOW(), NOW(), @body, @user_id) RETURNING *;

-- name: GetAllChirps :many
-- Retrieves all chirps in ascending order by created_at.
SELECT (
        id, created_at, updated_at, body, user_id
    )
FROM chirps
ORDER BY created_at ASC;