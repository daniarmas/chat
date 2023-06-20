package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type RefreshTokenDbDatasource interface {
	CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) (*entity.RefreshToken, error)
	GetRefreshTokenByUserId(ctx context.Context, id string) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error
	DeleteRefreshTokenByUserId(ctx context.Context, userId string) error
}

type refreshTokenDbDatasource struct {
	database *sqldatabase.Sql
	pgxConn  *pgxpool.Pool
}

func NewRefreshToken(database *sqldatabase.Sql, pgxConn *pgxpool.Pool) RefreshTokenDbDatasource {
	return &refreshTokenDbDatasource{
		database: database,
		pgxConn:  pgxConn,
	}
}

func (repo refreshTokenDbDatasource) CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) (*entity.RefreshToken, error) {
	// refreshTokenModel := models.RefreshTokenOrm{}
	// refreshTokenModel.MapToRefreshTokenGorm(&refreshToken)
	// result := repo.database.Gorm.Create(&refreshTokenModel)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	// res := refreshTokenModel.MapFromRefreshTokenGorm()
	// return res, nil

	var res entity.RefreshToken
	err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"refresh_token\" (user_id, expiration_time, create_time) VALUES ($1, $2, $3) RETURNING id, user_id, expiration_time, create_time", refreshToken.UserId, refreshToken.ExpirationTime, time.Now().UTC()).Scan(&res.ID, &res.UserId, &res.ExpirationTime, &res.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &res, nil
}

func (repo refreshTokenDbDatasource) GetRefreshTokenByUserId(ctx context.Context, id string) (*entity.RefreshToken, error) {
	// var refreshTokenOrm models.RefreshTokenOrm
	// result := repo.database.Gorm.Where("user_id = ?", id).Take(&refreshTokenOrm)
	// if result.Error != nil {
	// 	if result.Error.Error() == "record not found" {
	// 		return nil, myerror.NotFoundError{}
	// 	} else {
	// 		return nil, myerror.InternalServerError{}
	// 	}
	// }
	// res := refreshTokenOrm.MapFromRefreshTokenGorm()
	// return res, nil

	row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, user_id, expiration_time, create_time FROM \"refresh_token\" WHERE user_id = $1;", id)

	// Scan the row into a User struct
	var refreshToken entity.RefreshToken
	err := row.Scan(&refreshToken.ID, &refreshToken.UserId, &refreshToken.ExpirationTime, &refreshToken.CreateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			go log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			go log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &refreshToken, nil
}

func (repo refreshTokenDbDatasource) DeleteRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error {
	// refreshTokenGorm := models.RefreshTokenOrm{}
	// refreshTokenGorm.MapToRefreshTokenGorm(&refreshToken)
	// result := repo.database.Gorm.Delete(&refreshTokenGorm)
	// if result.Error != nil {
	// 	return result.Error
	// } else if result.RowsAffected == 0 {
	// 	return myerror.NotFoundError{}
	// }
	// return nil

	result, err := repo.pgxConn.Exec(context.Background(), "DELETE FROM \"refresh_token\" WHERE id = $1", refreshToken.ID)
	if err != nil {
		go log.Error().Msg(err.Error())
		return err
	}

	// Check if the record was actually deleted
	if result.RowsAffected() == 0 {
		return myerror.NotFoundError{}
	} else {
		return nil
	}
}

func (repo refreshTokenDbDatasource) DeleteRefreshTokenByUserId(ctx context.Context, userId string) error {
	// refreshTokenGorm := models.RefreshTokenOrm{}
	// result := repo.database.Gorm.Where("user_id = ?", userId).Delete(&refreshTokenGorm)
	// if result.Error != nil {
	// 	return result.Error
	// } else if result.RowsAffected == 0 {
	// 	return myerror.NotFoundError{}
	// }
	// return nil

	result, err := repo.pgxConn.Exec(context.Background(), "DELETE FROM \"refresh_token\" WHERE user_id = $1", userId)
	if err != nil {
		go log.Error().Msg(err.Error())
		return err
	}

	// Check if the record was actually deleted
	if result.RowsAffected() == 0 {
		return myerror.NotFoundError{}
	} else {
		return nil
	}
}
