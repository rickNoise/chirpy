-- name: CreateRefreshToken :one
-- creates a record in the refresh_tokens table
-- revoked_at field set to NULL for new tokens
INSERT INTO
    refresh_tokens (
        token,
        created_at,
        updated_at,
        user_id,
        expires_at,
        revoked_at
    )
VALUES (
        @token,
        NOW(),
        NOW(),
        @user_id,
        @expires_at,
        NULL
    ) RETURNING *;