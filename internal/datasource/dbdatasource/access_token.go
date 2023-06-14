package dbdatasource

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type AccessTokenDbDatasource interface {
	CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error)
}

type accessTokenDbDatasource struct {
	database *sqldatabase.Sql
}

func NewAccessTokenDbDatasource(database *sqldatabase.Sql) AccessTokenDbDatasource {
	return &accessTokenDbDatasource{
		database: database,
	}
}

func (repo accessTokenDbDatasource) CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error) {
	accessTokenModel := models.AccessTokenOrm{}
	accessTokenModel.MapToAccessTokenGorm(&accessToken)
	result := repo.database.Gorm.Create(&accessTokenModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := accessTokenModel.MapFromAccessTokenGorm()
	return res, nil
}
