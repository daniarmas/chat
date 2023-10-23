-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1 LIMIT 1;

-- name: GetRefreshTokenByUserId :one
SELECT * FROM "refresh_token" WHERE user_id = $1 LIMIT 1;