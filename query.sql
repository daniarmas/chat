-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1 LIMIT 1;

-- name: GetRefreshTokenByUserId :one
SELECT * FROM "refresh_token" WHERE user_id = $1 LIMIT 1;

-- name: DeleteAccessTokenByRefreshTokenId :one
DELETE FROM "access_token" WHERE refresh_token_id = $1 RETURNING *;

-- name: CreateRefreshToken :one
INSERT INTO "refresh_token" (user_id, expiration_time, create_time) VALUES ($1, $2, $3) RETURNING *;