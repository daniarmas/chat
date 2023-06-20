package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type AccessTokenDbDatasource interface {
	CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error)
}

type accessTokenDbDatasource struct {
	database *sqldatabase.Sql
	pgxConn  *pgxpool.Pool
}

func NewAccessToken(database *sqldatabase.Sql, pgxConn *pgxpool.Pool) AccessTokenDbDatasource {
	return &accessTokenDbDatasource{
		database: database,
		pgxConn:  pgxConn,
	}
}

func (repo accessTokenDbDatasource) CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error) {
	var res entity.AccessToken
	err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"access_token\" (refresh_token_id, user_id, expiration_time, create_time) VALUES ($1, $2, $3, $4) RETURNING id, user_id, refresh_token_id, expiration_time, create_time", accessToken.RefreshTokenId, accessToken.UserId, accessToken.ExpirationTime, time.Now().UTC()).Scan(&res.ID, &res.UserId, &res.RefreshTokenId, &res.ExpirationTime, &res.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &res, nil
}
