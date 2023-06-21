package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ApiKeyDbDatasource interface {
	CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error)
}

type apiKeyDbDatasource struct {
	database *sqldatabase.Sql
	pgxConn  *pgxpool.Pool
}

func NewApiKey(database *sqldatabase.Sql, pgxConn *pgxpool.Pool) ApiKeyDbDatasource {
	return &apiKeyDbDatasource{
		database: database,
		pgxConn:  pgxConn,
	}
}

func (repo *apiKeyDbDatasource) CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error) {
	var res entity.ApiKey
	if apiKey.ExpirationTime.IsZero() {
		// Add 3 month to time.Now for set 3 month of expiration time.
		expirationTime := time.Now().AddDate(0, 3, 0).UTC()
		apiKey.ExpirationTime = expirationTime
	}
	err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"api_key\" (app_version, revoked, expiration_time, create_time) VALUES ($1, $2, $3, $4) RETURNING id, app_version, revoked, expiration_time, create_time", apiKey.AppVersion, apiKey.Revoked, apiKey.ExpirationTime, time.Now().UTC()).Scan(&res.ID, &res.AppVersion, &res.Revoked, &res.ExpirationTime, &res.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &res, nil
}
