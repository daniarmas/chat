package ownredis

import (
	"github.com/daniarmas/chat/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDsn,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDb,
	})
	return rdb, nil
}
