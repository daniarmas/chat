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

type AccessTokenCacheDatasource interface {
	GetAccessToken(ctx context.Context, keyValue string) (*models.AccessTokenOrm, error)
	CacheAccessTokenById(ctx context.Context, accessToken models.AccessTokenOrm) error
	DeleteAccessTokenCache(ctx context.Context, accessToken *models.AccessTokenOrm) error
}

type accessTokenRedisDatasource struct {
	redis *redis.Client
	cfg   *config.Config
}

func NewAccessTokenCacheDatasource(redis *redis.Client, cfg *config.Config) AccessTokenCacheDatasource {
	return &accessTokenRedisDatasource{
		redis: redis,
		cfg:   cfg,
	}
}

func (repo accessTokenRedisDatasource) DeleteAccessTokenCache(ctx context.Context, accessToken *models.AccessTokenOrm) error {
	cacheKeyById := fmt.Sprintf("access_token:%s", accessToken.ID)
	err := repo.redis.Del(ctx, cacheKeyById).Err()
	if err != nil {
		go log.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (repo accessTokenRedisDatasource) GetAccessToken(ctx context.Context, keyValue string) (*models.AccessTokenOrm, error) {
	var accessToken models.AccessTokenOrm
	cacheKey := fmt.Sprintf("access_token:%s", keyValue)
	err := repo.redis.HGetAll(ctx, cacheKey).Scan(&accessToken)
	if err != nil {
		go log.Error().Msg(err.Error())
		return nil, err
	}

	return &accessToken, nil
}

func (repo accessTokenRedisDatasource) CacheAccessTokenById(ctx context.Context, accessToken models.AccessTokenOrm) error {
	cacheKey := fmt.Sprintf("access_token:%s", accessToken.ID)
	if err := repo.redis.HSet(ctx, cacheKey, accessToken).Err(); err != nil {
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
