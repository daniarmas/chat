package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/datasource/dbdatasource"
	"github.com/daniarmas/chat/internal/entity"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) (*entity.RefreshToken, error)
	GetRefreshTokenByUserId(ctx context.Context, id string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error
	DeleteRefreshTokenByUserId(ctx context.Context, userId string) error
}

type refreshToken struct {
	refreshTokenDbDatasource dbdatasource.RefreshTokenDbDatasource
}

func NewRefreshTokenRepository(refreshTokenDbDatasource dbdatasource.RefreshTokenDbDatasource) RefreshTokenRepository {
	return &refreshToken{
		refreshTokenDbDatasource: refreshTokenDbDatasource,
	}
}

func (repo refreshToken) CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) (*entity.RefreshToken, error) {
	res, err := repo.refreshTokenDbDatasource.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo refreshToken) GetRefreshTokenByUserId(ctx context.Context, id string) (*entity.RefreshToken, error) {
	res, err := repo.refreshTokenDbDatasource.GetRefreshTokenByUserId(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo refreshToken) DeleteRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error {
	err := repo.refreshTokenDbDatasource.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (repo refreshToken) DeleteRefreshTokenByUserId(ctx context.Context, userId string) error {
	err := repo.refreshTokenDbDatasource.DeleteRefreshTokenByUserId(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}
