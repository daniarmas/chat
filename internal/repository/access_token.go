package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/datasource/cacheds"
	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
)

type AccessTokenRepository interface {
	CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error)
	GetAccessTokenById(ctx context.Context, id string) (*entity.AccessToken, error)
	DeleteAccessTokenByRefreshTokenId(ctx context.Context, refreshTokenId string) error
	DeleteAccessTokenByUserId(ctx context.Context, userId string) error
}

type accessToken struct {
	accessTokenDbDatasource databaseds.AccessTokenDbDatasource
	accessTokenCacheDs      cacheds.AccessTokenCacheDatasource
}

func NewAccessToken(accessTokenDbDatasource databaseds.AccessTokenDbDatasource, accessTokenCacheDs cacheds.AccessTokenCacheDatasource) AccessTokenRepository {
	return &accessToken{
		accessTokenDbDatasource: accessTokenDbDatasource,
		accessTokenCacheDs:      accessTokenCacheDs,
	}
}

func (repo accessToken) DeleteAccessTokenByUserId(ctx context.Context, userId string) error {
	accessToken, err := repo.accessTokenDbDatasource.DeleteAccessTokenByUserId(ctx, userId)
	if err != nil {
		return err
	}
	err = repo.accessTokenCacheDs.DeleteAccessTokenCache(ctx, accessToken)
	if err != nil {
		return err
	}
	return nil
}

func (repo accessToken) DeleteAccessTokenByRefreshTokenId(ctx context.Context, refreshTokenId string) error {
	accessToken, err := repo.accessTokenDbDatasource.DeleteAccessTokenByRefreshTokenId(ctx, refreshTokenId)
	if err != nil {
		return err
	}
	err = repo.accessTokenCacheDs.DeleteAccessTokenCache(ctx, accessToken)
	if err != nil {
		return err
	}
	return nil
}

func (repo accessToken) CreateAccessToken(ctx context.Context, accessToken entity.AccessToken) (*entity.AccessToken, error) {
	res, err := repo.accessTokenDbDatasource.CreateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	err = repo.accessTokenCacheDs.CacheAccessTokenById(ctx, *res)
	if err != nil {
		return nil, err
	}
	return res.MapFromAccessTokenModel(), nil
}

func (repo accessToken) GetAccessTokenById(ctx context.Context, id string) (*entity.AccessToken, error) {
	var accessToken *models.AccessToken
	var res entity.AccessToken
	var err error
	accessToken, err = repo.accessTokenCacheDs.GetAccessToken(ctx, id)
	if err != nil {
		return nil, err
	}
	if accessToken.ID == "" {
		accessToken, err = repo.accessTokenDbDatasource.GetAccessTokenById(ctx, id)
		if err != nil {
			return nil, err
		}
		err = repo.accessTokenCacheDs.CacheAccessTokenById(ctx, *accessToken)
		if err != nil {
			return nil, err
		}
	}
	res = *accessToken.MapFromAccessTokenModel()
	return &res, nil
}
