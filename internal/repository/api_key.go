package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/entity"
)

type ApiKeyRepository interface {
	CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error)
}

type apiKeyRepository struct {
	apiKeyDbDatasource databaseds.ApiKeyDbDatasource
}

func NewApiKey(apiKeyDbDatasource databaseds.ApiKeyDbDatasource) ApiKeyRepository {
	return &apiKeyRepository{
		apiKeyDbDatasource: apiKeyDbDatasource,
	}
}

func (repo *apiKeyRepository) CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error) {
	res, err := repo.apiKeyDbDatasource.CreateApiKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return res, nil
}
