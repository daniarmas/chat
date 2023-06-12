package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type AccessTokenRepository interface {
	CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error)
}

type accessToken struct {
	database *sqldatabase.Sql
}

func NewAccessTokenRepository(database *sqldatabase.Sql) AccessTokenRepository {
	return &accessToken{
		database: database,
	}
}

func (repo accessToken) CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error) {
	accessTokenModel := models.AccessTokenOrm{}
	accessTokenModel.MapToAccessTokenGorm(&accessToken)
	result := repo.database.Gorm.Create(&accessTokenModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := accessTokenModel.MapFromAccessTokenGorm()
	return res, nil
}
