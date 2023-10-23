// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package gen

import (
	"context"

	"github.com/google/uuid"
)

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
