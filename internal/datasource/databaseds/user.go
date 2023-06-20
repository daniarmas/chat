package databaseds

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
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
}

func NewUser(database *sqldatabase.Sql, pgxConn *pgxpool.Pool) UserDbDatasource {
	return &userDbDatasource{
		database: database,
		pgxConn:  pgxConn,
	}
}

func (repo *userDbDatasource) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	// var user *models.UserOrm
	// result := repo.database.Gorm.Where("id = ?", id).Take(&user)
	// if result.Error != nil {
	// 	if result.Error.Error() == "record not found" {
	// 		return nil, myerror.NotFoundError{}
	// 	} else {
	// 		return nil, myerror.InternalServerError{}
	// 	}
	// }
	// res := user.MapFromUserGorm()
	// return res, nil

	row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, email, fullname, username, password, create_time FROM \"user\" WHERE id = $1;", id)

	// Scan the row into a User struct
	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &user, nil
}

func (repo *userDbDatasource) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, email, fullname, username, password, create_time FROM \"user\" WHERE email = $1;", email)

	// Scan the row into a User struct
	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &user, nil
}

func (repo *userDbDatasource) CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error) {
	user := models.UserOrm{Email: email, Password: password, Username: username, Fullname: fullname}
	result := repo.database.Gorm.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	userEntity := user.MapFromUserGorm()
	return userEntity, nil
}
