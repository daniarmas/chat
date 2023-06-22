package cacheds

import (
	"context"
	"fmt"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type UserCacheDatasource interface {
	GetUser(ctx context.Context, keyValue string) (*models.User, error)
	CacheUserById(ctx context.Context, user models.User) error
	CacheUserByEmail(ctx context.Context, user models.User) error
}

type userRedisDatasource struct {
	redis *redis.Client
	cfg   *config.Config
}

func NewUserCacheDatasource(redis *redis.Client, cfg *config.Config) UserCacheDatasource {
	return &userRedisDatasource{
		redis: redis,
		cfg:   cfg,
	}
}

func (repo userRedisDatasource) GetUser(ctx context.Context, userId string) (*models.User, error) {
	var user models.User
	cacheKey := fmt.Sprintf("user:%s", userId)
	err := repo.redis.HGetAll(ctx, cacheKey).Scan(&user)
	if err != nil {
		go log.Error().Msg(err.Error())
		return nil, err
	}

	return &user, nil
}

func (repo userRedisDatasource) CacheUserById(ctx context.Context, user models.User) error {
	cacheKey := fmt.Sprintf("user:%s", user.ID)
	if err := repo.redis.HSet(ctx, cacheKey, user).Err(); err != nil {
		go log.Error().Msg(err.Error())
		return err
	}
	err := repo.redis.Expire(ctx, cacheKey, time.Duration(repo.cfg.RedisExpirationTimeSeconds)*time.Second).Err()
	if err != nil {
		go log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (repo userRedisDatasource) CacheUserByEmail(ctx context.Context, user models.User) error {
	cacheKey := fmt.Sprintf("user:%s", user.Email)
	if err := repo.redis.HSet(ctx, cacheKey, user).Err(); err != nil {
		go log.Error().Msg(err.Error())
		return err
	}
	err := repo.redis.Expire(ctx, cacheKey, time.Duration(repo.cfg.RedisExpirationTimeSeconds)*time.Second).Err()
	if err != nil {
		go log.Error().Msg(err.Error())
		return err
	}

	return nil
}
