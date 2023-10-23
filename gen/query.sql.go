// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package gen

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO "refresh_token" (user_id, expiration_time, create_time) VALUES ($1, $2, $3) RETURNING id, user_id, expiration_time, create_time
`

type CreateRefreshTokenParams struct {
	UserID         uuid.UUID
	ExpirationTime time.Time
	CreateTime     time.Time
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.UserID, arg.ExpirationTime, arg.CreateTime)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExpirationTime,
		&i.CreateTime,
	)
	return i, err
}

const deleteAccessTokenByRefreshTokenId = `-- name: DeleteAccessTokenByRefreshTokenId :one
DELETE FROM "access_token" WHERE refresh_token_id = $1 RETURNING id, user_id, refresh_token_id, expiration_time, create_time
`

func (q *Queries) DeleteAccessTokenByRefreshTokenId(ctx context.Context, refreshTokenID uuid.UUID) (AccessToken, error) {
	row := q.db.QueryRowContext(ctx, deleteAccessTokenByRefreshTokenId, refreshTokenID)
	var i AccessToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshTokenID,
		&i.ExpirationTime,
		&i.CreateTime,
	)
	return i, err
}

const getRefreshTokenByUserId = `-- name: GetRefreshTokenByUserId :one
SELECT id, user_id, expiration_time, create_time FROM "refresh_token" WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetRefreshTokenByUserId(ctx context.Context, userID uuid.UUID) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getRefreshTokenByUserId, userID)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExpirationTime,
		&i.CreateTime,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, fullname, username, create_time FROM "user" WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Fullname,
		&i.Username,
		&i.CreateTime,
	)
	return i, err
}