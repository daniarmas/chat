package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/gen"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ApiKeyDbDatasource interface {
	CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error)
}

type apiKeyDbDatasource struct {
	pgxConn *pgxpool.Pool
	queries *gen.Queries
}

func NewApiKey(pgxConn *pgxpool.Pool, queries *gen.Queries) ApiKeyDbDatasource {
	return &apiKeyDbDatasource{
		pgxConn: pgxConn,
		queries: queries,
	}
}

func (repo *apiKeyDbDatasource) CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error) {
	// var res entity.ApiKey
	// if apiKey.ExpirationTime.IsZero() {
	// 	// Add 3 month to time.Now for set 3 month of expiration time.
	// 	expirationTime := time.Now().AddDate(0, 3, 0).UTC()
	// 	apiKey.ExpirationTime = expirationTime
	// }
	// err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"api_key\" (app_version, revoked, expiration_time, create_time) VALUES ($1, $2, $3, $4) RETURNING id, app_version, revoked, expiration_time, create_time", apiKey.AppVersion, apiKey.Revoked, apiKey.ExpirationTime, time.Now().UTC()).Scan(&res.ID, &res.AppVersion, &res.Revoked, &res.ExpirationTime, &res.CreateTime)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// 	return nil, err
	// }
	// return &res, nil
	createTime := time.Now().UTC()
	res, err := repo.queries.CreateApiKey(ctx, gen.CreateApiKeyParams{AppVersion: apiKey.AppVersion, Revoked: apiKey.Revoked, ExpirationTime: apiKey.ExpirationTime, CreateTime: createTime})
	if err != nil {
		go log.Error().Msg(err.Error())
		return nil, err
	}
	return &entity.ApiKey{
		ID:             res.ID.String(),
		ExpirationTime: res.ExpirationTime,
		AppVersion:     res.AppVersion,
		Revoked:        res.Revoked,
		CreateTime:     res.CreateTime,
	}, nil
}
