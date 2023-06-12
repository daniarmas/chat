package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error)
}

type userRepository struct {
	database *sqldatabase.Sql
}

func NewUserRepository(database *sqldatabase.Sql) UserRepository {
	return &userRepository{
		database: database,
	}
}

func (repo *userRepository) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	var user *models.UserOrm
	result := repo.database.Gorm.Where("id = ?", id).Take(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, myerror.NotFoundError{}
		} else {
			return nil, myerror.InternalServerError{}
		}
	}
	res := user.MapFromUserGorm()
	return res, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *models.UserOrm
	result := repo.database.Gorm.Where("email = ?", email).Take(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, myerror.NotFoundError{}
		} else {
			return nil, myerror.InternalServerError{}
		}
	}
	res := user.MapFromUserGorm()
	return res, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error) {
	user := models.UserOrm{Email: email, Password: password, Username: username, Fullname: fullname}
	result := repo.database.Gorm.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	userEntity := user.MapFromUserGorm()
	return userEntity, nil
}
