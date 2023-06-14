package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/datasource/dbdatasource"
	"github.com/daniarmas/chat/internal/entity"
)

type AccessTokenRepository interface {
	CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error)
}

type accessToken struct {
	accessTokenDbDatasource dbdatasource.AccessTokenDbDatasource
}

func NewAccessTokenRepository(accessTokenDbDatasource dbdatasource.AccessTokenDbDatasource) AccessTokenRepository {
	return &accessToken{
		accessTokenDbDatasource: accessTokenDbDatasource,
	}
}

func (repo accessToken) CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error) {
	chat, err := repo.accessTokenDbDatasource.CreateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
