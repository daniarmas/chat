package cacheds

import (
	"context"
	"fmt"

	"github.com/daniarmas/chat/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type UserCacheDatasource interface {
	GetUser(ctx context.Context, keyValue string) (*models.UserOrm, error)
	CreateUser(ctx context.Context, user models.UserOrm) error
}

type userRedisDatasource struct {
	redis *redis.Client
}

func NewUserCacheDatasource(redis *redis.Client) UserCacheDatasource {
	return &userRedisDatasource{
		redis: redis,
	}
}

func (repo userRedisDatasource) GetUser(ctx context.Context, userId string) (*models.UserOrm, error) {
	var user models.UserOrm
	cacheKey := fmt.Sprintf("user:%s", userId)
	err := repo.redis.HGetAll(ctx, cacheKey).Scan(&user)
	if err != nil {
		go log.Error().Msg(err.Error())
		return nil, err
	}

	return &user, nil
}

func (repo userRedisDatasource) CreateUser(ctx context.Context, user models.UserOrm) error {
	cacheKey := fmt.Sprintf("user:%s", user.ID)
	if err := repo.redis.HSet(ctx, cacheKey, user).Err(); err != nil {
		go log.Error().Msg(err.Error())
		return err
	}

	return nil
}
