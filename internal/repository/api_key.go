package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type ApiKeyRepository interface {
	CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error)
}

type apiKeyRepository struct {
	database *sqldatabase.Sql
}

func NewApiKeyRepository(database *sqldatabase.Sql) ApiKeyRepository {
	return &apiKeyRepository{
		database: database,
	}
}

func (repo *apiKeyRepository) CreateApiKey(ctx context.Context, apiKey *entity.ApiKey) (*entity.ApiKey, error) {
	apiKeyGorm := models.ApiKeyOrm{}
	apiKeyGorm.MapToApiKeyGorm(apiKey)
	result := repo.database.Gorm.Create(&apiKeyGorm)
	if result.Error != nil {
		return nil, result.Error
	}
	response := apiKeyGorm.MapFromApiKeyGorm()
	return response, nil
}
