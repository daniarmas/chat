package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/datasource/hashds"
	"github.com/daniarmas/chat/internal/entity"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type UserDbDatasource interface {
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error)
}

type userDbDatasource struct {
	database *sqldatabase.Sql
	pgxConn  *pgxpool.Pool
	hashDs   hashds.HashDatasource
}

func NewUser(database *sqldatabase.Sql, pgxConn *pgxpool.Pool, hashDs hashds.HashDatasource) UserDbDatasource {
	return &userDbDatasource{
		database: database,
		pgxConn:  pgxConn,
		hashDs:   hashDs,
	}
}

func (repo *userDbDatasource) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, email, fullname, username, password, create_time FROM \"user\" WHERE id = $1;", id)

	// Scan the row into a User struct
	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &user, nil
}

func (repo *userDbDatasource) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, email, fullname, username, password, create_time FROM \"user\" WHERE email = $1;", email)

	// Scan the row into a User struct
	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &user, nil
}

func (repo *userDbDatasource) CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error) {
	var user entity.User
	passwordHashed, _ := repo.hashDs.Hash(password)
	err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"user\" (email, fullname, username, password, create_time) VALUES ($1, $2, $3) RETURNING id, email, fullname, username, password, create_time", email, fullname, passwordHashed, time.Now().UTC()).Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return &user, nil
}
