package repository

import (
	"context"

	"github.com/daniarmas/chat/gen"
	"github.com/daniarmas/chat/internal/datasource/cacheds"
	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error)
}

type userRepository struct {
	userDbDatasource databaseds.UserDbDatasource
	userCacheDs      cacheds.UserCacheDatasource
	queries          *gen.Queries
}

func NewUser(userDbDatasource databaseds.UserDbDatasource, userCacheDs cacheds.UserCacheDatasource, queries *gen.Queries) UserRepository {
	return &userRepository{
		userDbDatasource: userDbDatasource,
		userCacheDs:      userCacheDs,
		queries:          queries,
	}
}

func (repo *userRepository) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	var user *models.User
	var res entity.User
	var err error
	user, err = repo.userCacheDs.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user.ID == "" {
		user, err = repo.userDbDatasource.GetUserById(ctx, id)
		if err != nil {
			return nil, err
		}
		err = repo.userCacheDs.CacheUserById(ctx, *user)
		if err != nil {
			return nil, err
		}
	}
	res = *user.MapFromUserModel()
	return &res, nil
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
