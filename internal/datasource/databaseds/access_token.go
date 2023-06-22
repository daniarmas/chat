package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type AccessTokenDbDatasource interface {
	CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*models.AccessToken, error)
	GetAccessTokenById(ctx context.Context, id string) (*models.AccessToken, error)
	DeleteAccessTokenByRefreshTokenId(ctx context.Context, refreshTokenId string) (*models.AccessToken, error)
	DeleteAccessTokenByUserId(ctx context.Context, userId string) (*models.AccessToken, error)
}

type accessTokenDbDatasource struct {
	pgxConn  *pgxpool.Pool
}

func NewAccessToken(pgxConn *pgxpool.Pool) AccessTokenDbDatasource {
	return &accessTokenDbDatasource{
		pgxConn:  pgxConn,
	}
}

func (repo accessTokenDbDatasource) DeleteAccessTokenByUserId(ctx context.Context, userId string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	row := repo.pgxConn.QueryRow(context.Background(), "DELETE FROM \"access_token\" WHERE user_id = $1 RETURNING id, user_id, refresh_token_id, expiration_time, create_time;", userId)
	err := row.Scan(&accessToken.ID, &accessToken.UserId, &accessToken.RefreshTokenId, &accessToken.ExpirationTime, &accessToken.CreateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &accessToken, nil
}

func (repo accessTokenDbDatasource) DeleteAccessTokenByRefreshTokenId(ctx context.Context, refreshTokenId string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	row := repo.pgxConn.QueryRow(context.Background(), "DELETE FROM \"access_token\" WHERE refresh_token_id = $1 RETURNING id, user_id, refresh_token_id, expiration_time, create_time;", refreshTokenId)
	err := row.Scan(&accessToken.ID, &accessToken.UserId, &accessToken.RefreshTokenId, &accessToken.ExpirationTime, &accessToken.CreateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &accessToken, nil
}

func (repo accessTokenDbDatasource) CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*models.AccessToken, error) {
	var res models.AccessToken
	err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"access_token\" (refresh_token_id, user_id, expiration_time, create_time) VALUES ($1, $2, $3, $4) RETURNING id, user_id, refresh_token_id, expiration_time, create_time", accessToken.RefreshTokenId, accessToken.UserId, accessToken.ExpirationTime, time.Now().UTC()).Scan(&res.ID, &res.UserId, &res.RefreshTokenId, &res.ExpirationTime, &res.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &res, nil
}

func (repo accessTokenDbDatasource) GetAccessTokenById(ctx context.Context, id string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, user_id, refresh_token_id, expiration_time, create_time FROM \"access_token\" WHERE id = $1;", id)
	err := row.Scan(&accessToken.ID, &accessToken.UserId, &accessToken.RefreshTokenId, &accessToken.ExpirationTime, &accessToken.CreateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &accessToken, nil
}
