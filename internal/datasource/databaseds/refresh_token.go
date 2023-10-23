package databaseds

import (
	"context"
	"strings"

	"github.com/daniarmas/chat/gen"
	"github.com/daniarmas/chat/internal/entity"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/google/uuid"
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
	pgxConn *pgxpool.Pool
	queries *gen.Queries
}

func NewRefreshToken(pgxConn *pgxpool.Pool, queries *gen.Queries) RefreshTokenDbDatasource {
	return &refreshTokenDbDatasource{
		pgxConn: pgxConn,
		queries: queries,
	}
}

func (repo refreshTokenDbDatasource) CreateRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) (*entity.RefreshToken, error) {
	// var res entity.RefreshToken
	// err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"refresh_token\" (user_id, expiration_time, create_time) VALUES ($1, $2, $3) RETURNING id, user_id, expiration_time, create_time", refreshToken.UserId, refreshToken.ExpirationTime, time.Now().UTC()).Scan(&res.ID, &res.UserId, &res.ExpirationTime, &res.CreateTime)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// 	return nil, err
	// }
	// return &res, nil
	res, err := repo.queries.CreateRefreshToken(ctx, gen.CreateRefreshTokenParams{UserID: uuid.MustParse(refreshToken.UserId), ExpirationTime: refreshToken.ExpirationTime, CreateTime: refreshToken.CreateTime})
	if err != nil {
		go log.Error().Msg(err.Error())
		return nil, err
	}
	return &entity.RefreshToken{
		ID:             res.ID.String(),
		UserId:         res.UserID.String(),
		ExpirationTime: res.ExpirationTime,
		CreateTime:     res.CreateTime,
	}, nil
}

func (repo refreshTokenDbDatasource) GetRefreshTokenByUserId(ctx context.Context, id string) (*entity.RefreshToken, error) {
	res, err := repo.queries.GetRefreshTokenByUserId(ctx, uuid.MustParse(id))
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			go log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			go log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &entity.RefreshToken{
		ID:             res.ID.String(),
		UserId:         res.UserID.String(),
		ExpirationTime: res.ExpirationTime,
		CreateTime:     res.CreateTime,
	}, nil
	// row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, user_id, expiration_time, create_time FROM \"refresh_token\" WHERE user_id = $1;", id)

	// // Scan the row into a User struct
	// var refreshToken entity.RefreshToken
	// err := row.Scan(&refreshToken.ID, &refreshToken.UserId, &refreshToken.ExpirationTime, &refreshToken.CreateTime)
	// if err != nil {
	// 	if err.Error() == "no rows in result set" {
	// 		go log.Error().Msg(err.Error())
	// 		return nil, myerror.NotFoundError{}
	// 	} else {
	// 		go log.Error().Msg(err.Error())
	// 		return nil, err
	// 	}
	// }
	// return &refreshToken, nil
}

func (repo refreshTokenDbDatasource) DeleteRefreshToken(ctx context.Context, refreshToken entity.RefreshToken) error {
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
