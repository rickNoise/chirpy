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