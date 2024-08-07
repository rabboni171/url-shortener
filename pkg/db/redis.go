package db

import (
	"github.com/rabboni171/url-shortener/configs"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisConfig configs.DBRedisParams) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
}
