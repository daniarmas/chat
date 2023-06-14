package dbdatasource

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type ApiKeyDbDatasource interface {
	CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error)
}

type apiKeyDbDatasource struct {
	database *sqldatabase.Sql
}

func NewApiKeyDbDatasource(database *sqldatabase.Sql) ApiKeyDbDatasource {
	return &apiKeyDbDatasource{
		database: database,
	}
}

func (repo *apiKeyDbDatasource) CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error) {
	apiKeyGorm := models.ApiKeyOrm{}
	apiKeyGorm.MapToApiKeyGorm(apiKey)
	result := repo.database.Gorm.Create(&apiKeyGorm)
	if result.Error != nil {
		return nil, result.Error
	}
	response := apiKeyGorm.MapFromApiKeyGorm()
	return response, nil
}
