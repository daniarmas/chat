package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/entity"
)

type AccessTokenRepository interface {
	CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error)
}

type accessToken struct {
	accessTokenDbDatasource databaseds.AccessTokenDbDatasource
}

func NewAccessToken(accessTokenDbDatasource databaseds.AccessTokenDbDatasource) AccessTokenRepository {
	return &accessToken{
		accessTokenDbDatasource: accessTokenDbDatasource,
	}
}

func (repo accessToken) CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error) {
	res, err := repo.accessTokenDbDatasource.CreateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return res, nil
}
