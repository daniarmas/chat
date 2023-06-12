package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) (*entity.RefreshToken, error)
	GetRefreshTokenByUserId(ctx context.Context, id string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error
	DeleteRefreshTokenByUserId(ctx context.Context, userId string) error
}

type refreshToken struct {
	database *sqldatabase.Sql
}

func NewRefreshTokenRepository(database *sqldatabase.Sql) RefreshTokenRepository {
	return &refreshToken{
		database: database,
	}
}

func (repo refreshToken) CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) (*entity.RefreshToken, error) {
	refreshTokenModel := models.RefreshTokenOrm{}
	refreshTokenModel.MapToRefreshTokenGorm(&refreshToken)
	result := repo.database.Gorm.Create(&refreshTokenModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := refreshTokenModel.MapFromRefreshTokenGorm()
	return res, nil
}

func (repo refreshToken) GetRefreshTokenByUserId(ctx context.Context, id string) (*entity.RefreshToken, error) {
	var refreshTokenOrm models.RefreshTokenOrm
	result := repo.database.Gorm.Where("user_id = ?", id).Take(&refreshTokenOrm)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, myerror.NotFoundError{}
		} else {
			return nil, myerror.InternalServerError{}
		}
	}
	res := refreshTokenOrm.MapFromRefreshTokenGorm()
	return res, nil
}

func (repo refreshToken) DeleteRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error {
	refreshTokenGorm := models.RefreshTokenOrm{}
	refreshTokenGorm.MapToRefreshTokenGorm(&refreshToken)
	result := repo.database.Gorm.Delete(&refreshTokenGorm)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return myerror.NotFoundError{}
	}
	return nil
}

func (repo refreshToken) DeleteRefreshTokenByUserId(ctx context.Context, userId string) error {
	refreshTokenGorm := models.RefreshTokenOrm{}
	result := repo.database.Gorm.Where("user_id = ?", userId).Delete(&refreshTokenGorm)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return myerror.NotFoundError{}
	}
	return nil
}
