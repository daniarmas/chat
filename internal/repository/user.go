package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/datasource/dbdatasource"
	"github.com/daniarmas/chat/internal/entity"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error)
}

type userRepository struct {
	userDbDatasource dbdatasource.UserDbDatasource
}

func NewUserRepository(userDbDatasource dbdatasource.UserDbDatasource) UserRepository {
	return &userRepository{
		userDbDatasource: userDbDatasource,
	}
}

func (repo *userRepository) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	user, err := repo.userDbDatasource.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := repo.userDbDatasource.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error) {
	user, err := repo.userDbDatasource.CreateUser(ctx, email, password, username, fullname)
	if err != nil {
		return nil, err
	}
	return user, nil
}
